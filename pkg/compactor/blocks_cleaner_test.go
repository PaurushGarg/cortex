package compactor

import (
	"context"
	"crypto/rand"
	"fmt"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/oklog/ulid/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	prom_testutil "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thanos-io/thanos/pkg/block"
	"github.com/thanos-io/thanos/pkg/block/metadata"

	"github.com/cortexproject/cortex/pkg/storage/bucket"
	"github.com/cortexproject/cortex/pkg/storage/parquet"
	"github.com/cortexproject/cortex/pkg/storage/tsdb"
	"github.com/cortexproject/cortex/pkg/storage/tsdb/bucketindex"
	cortex_testutil "github.com/cortexproject/cortex/pkg/storage/tsdb/testutil"
	"github.com/cortexproject/cortex/pkg/storage/tsdb/users"
	"github.com/cortexproject/cortex/pkg/util"
	util_log "github.com/cortexproject/cortex/pkg/util/log"
	"github.com/cortexproject/cortex/pkg/util/services"
)

type testBlocksCleanerOptions struct {
	concurrency             int
	markersMigrationEnabled bool
	tenantDeletionDelay     time.Duration
	user4FilesExist         bool // User 4 has "FinishedTime" in tenant deletion marker set to "1h" ago.
}

func (o testBlocksCleanerOptions) String() string {
	return fmt.Sprintf("concurrency=%d, markers migration enabled=%v, tenant deletion delay=%v",
		o.concurrency, o.markersMigrationEnabled, o.tenantDeletionDelay)
}

func TestBlocksCleaner(t *testing.T) {
	for _, options := range []testBlocksCleanerOptions{
		{concurrency: 1, tenantDeletionDelay: 0, user4FilesExist: false},
		{concurrency: 1, tenantDeletionDelay: 2 * time.Hour, user4FilesExist: true},
		{concurrency: 1, markersMigrationEnabled: true},
		{concurrency: 2},
		{concurrency: 10},
	} {
		options := options

		t.Run(options.String(), func(t *testing.T) {
			t.Parallel()
			testBlocksCleanerWithOptions(t, options)
		})
	}
}

func TestBlockCleaner_KeyPermissionDenied(t *testing.T) {
	const userID = "user-1"

	bkt, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bkt = bucketindex.BucketWithGlobalMarkers(bkt)

	// Create blocks.
	ctx := context.Background()
	deletionDelay := 12 * time.Hour
	mbucket := &cortex_testutil.MockBucketFailure{
		Bucket: bkt,
		GetFailures: map[string]error{
			path.Join(userID, "bucket-index.json.gz"): cortex_testutil.ErrKeyAccessDeniedError,
		},
	}
	createTSDBBlock(t, bkt, userID, 10, 20, nil)

	cfg := BlocksCleanerConfig{
		DeletionDelay:      deletionDelay,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		BlockRanges:        (&tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}).ToMilliseconds(),
	}

	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, mbucket, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	blocksMarkedForDeletion := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, mbucket, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", reg, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)

	// Clean User with no error
	cleaner.bucketClient = bkt
	userLogger := util_log.WithUserID(userID, cleaner.logger)
	userBucket := bucket.NewUserBucketClient(userID, cleaner.bucketClient, cleaner.cfgProvider)
	err = cleaner.cleanUser(ctx, userLogger, userBucket, userID, false)
	require.NoError(t, err)
	s, err := bucketindex.ReadSyncStatus(ctx, bkt, userID, logger)
	require.NoError(t, err)
	require.Equal(t, bucketindex.Ok, s.Status)
	require.Equal(t, int64(0), s.NonQueryableUntil)

	// Clean with cmk error
	cleaner.bucketClient = mbucket
	userLogger = util_log.WithUserID(userID, cleaner.logger)
	userBucket = bucket.NewUserBucketClient(userID, cleaner.bucketClient, cleaner.cfgProvider)
	err = cleaner.cleanUser(ctx, userLogger, userBucket, userID, false)
	require.NoError(t, err)
	s, err = bucketindex.ReadSyncStatus(ctx, bkt, userID, logger)
	require.NoError(t, err)
	require.Equal(t, bucketindex.CustomerManagedKeyError, s.Status)
	require.Less(t, int64(0), s.NonQueryableUntil)

	// Re grant access to the key
	cleaner.bucketClient = bkt
	userLogger = util_log.WithUserID(userID, cleaner.logger)
	userBucket = bucket.NewUserBucketClient(userID, cleaner.bucketClient, cleaner.cfgProvider)
	err = cleaner.cleanUser(ctx, userLogger, userBucket, userID, false)
	require.NoError(t, err)
	s, err = bucketindex.ReadSyncStatus(ctx, bkt, userID, logger)
	require.NoError(t, err)
	require.Equal(t, bucketindex.Ok, s.Status)
	require.Less(t, int64(0), s.NonQueryableUntil)
}

func testBlocksCleanerWithOptions(t *testing.T, options testBlocksCleanerOptions) {
	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)

	// If the markers migration is enabled, then we create the fixture blocks without
	// writing the deletion marks in the global location, because they will be migrated
	// at startup.
	if !options.markersMigrationEnabled {
		bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)
	}

	// Create blocks.
	ctx := context.Background()
	now := time.Now()
	deletionDelay := 12 * time.Hour
	block1 := createTSDBBlock(t, bucketClient, "user-1", 10, 20, nil)
	block2 := createTSDBBlock(t, bucketClient, "user-1", 20, 30, nil)
	block3 := createTSDBBlock(t, bucketClient, "user-1", 30, 40, nil)
	block4 := ulid.MustNew(4, rand.Reader)
	block5 := ulid.MustNew(5, rand.Reader)
	block6 := createTSDBBlock(t, bucketClient, "user-1", 40, 50, nil)
	block7 := createTSDBBlock(t, bucketClient, "user-2", 10, 20, nil)
	block8 := createTSDBBlock(t, bucketClient, "user-2", 40, 50, nil)
	block11 := ulid.MustNew(11, rand.Reader)
	createDeletionMark(t, bucketClient, "user-1", block2, now.Add(-deletionDelay).Add(time.Hour))             // Block hasn't reached the deletion threshold yet.
	createDeletionMark(t, bucketClient, "user-1", block3, now.Add(-deletionDelay).Add(-time.Hour))            // Block reached the deletion threshold.
	createDeletionMark(t, bucketClient, "user-1", block4, now.Add(-deletionDelay).Add(time.Hour))             // Partial block hasn't reached the deletion threshold yet.
	createDeletionMark(t, bucketClient, "user-1", block5, now.Add(-deletionDelay).Add(-time.Hour))            // Partial block reached the deletion threshold.
	require.NoError(t, bucketClient.Delete(ctx, path.Join("user-1", block6.String(), metadata.MetaFilename))) // Partial block without deletion mark.
	createBlockVisitMarker(t, bucketClient, "user-1", block11)                                                // Partial block only has visit marker.
	createDeletionMark(t, bucketClient, "user-2", block7, now.Add(-deletionDelay).Add(-time.Hour))            // Block reached the deletion threshold.

	// Blocks for user-3, tenant marked for deletion.
	require.NoError(t, tsdb.WriteTenantDeletionMark(context.Background(), bucketClient, "user-3", tsdb.NewTenantDeletionMark(time.Now())))
	block9 := createTSDBBlock(t, bucketClient, "user-3", 10, 30, nil)
	block10 := createTSDBBlock(t, bucketClient, "user-3", 30, 50, nil)
	createParquetMarker(t, bucketClient, "user-3", block10)

	// User-4 with no more blocks, but couple of mark and debug files. Should be fully deleted.
	user4Mark := tsdb.NewTenantDeletionMark(time.Now())
	user4Mark.FinishedTime = time.Now().Unix() - 60 // Set to check final user cleanup.
	require.NoError(t, tsdb.WriteTenantDeletionMark(context.Background(), bucketClient, "user-4", user4Mark))
	user4DebugMetaFile := path.Join("user-4", block.DebugMetas, "meta.json")
	require.NoError(t, bucketClient.Upload(context.Background(), user4DebugMetaFile, strings.NewReader("some random content here")))

	// No Compact blocks marker
	createTSDBBlock(t, bucketClient, "user-5", 10, 30, nil)
	block12 := createTSDBBlock(t, bucketClient, "user-5", 30, 50, nil)
	createNoCompactionMark(t, bucketClient, "user-5", block12)

	// Create Parquet marker
	block13 := createTSDBBlock(t, bucketClient, "user-6", 30, 50, nil)
	// This block should be converted to Parquet format so counted as remaining.
	block14 := createTSDBBlock(t, bucketClient, "user-6", 30, 50, nil)
	createParquetMarker(t, bucketClient, "user-6", block13)

	// The fixtures have been created. If the bucket client wasn't wrapped to write
	// deletion marks to the global location too, then this is the right time to do it.
	if options.markersMigrationEnabled {
		bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)
	}

	cfg := BlocksCleanerConfig{
		DeletionDelay:                      deletionDelay,
		CleanupInterval:                    time.Minute,
		CleanupConcurrency:                 options.concurrency,
		BlockDeletionMarksMigrationEnabled: options.markersMigrationEnabled,
		TenantCleanupDelay:                 options.tenantDeletionDelay,
		BlockRanges:                        (&tsdb.DurationList{2 * time.Hour}).ToMilliseconds(),
	}

	reg := prometheus.NewPedanticRegistry()
	logger := log.NewNopLogger()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	cfgProvider.parquetConverterEnabled = map[string]bool{
		"user-3": true,
		"user-5": true,
		"user-6": true,
	}
	blocksMarkedForDeletion := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", reg, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)
	require.NoError(t, services.StartAndAwaitRunning(ctx, cleaner))
	defer services.StopAndAwaitTerminated(ctx, cleaner) //nolint:errcheck

	for _, tc := range []struct {
		path           string
		expectedExists bool
	}{
		// Check the storage to ensure only the block which has reached the deletion threshold
		// has been effectively deleted.
		{path: path.Join("user-1", block1.String(), metadata.MetaFilename), expectedExists: true},
		{path: path.Join("user-1", block3.String(), metadata.MetaFilename), expectedExists: false},
		{path: path.Join("user-2", block7.String(), metadata.MetaFilename), expectedExists: false},
		{path: path.Join("user-2", block8.String(), metadata.MetaFilename), expectedExists: true},
		// Should not delete a block with deletion mark who hasn't reached the deletion threshold yet.
		{path: path.Join("user-1", block2.String(), metadata.MetaFilename), expectedExists: true},
		{path: path.Join("user-1", bucketindex.BlockDeletionMarkFilepath(block2)), expectedExists: true},
		// Should delete a partial block with deletion mark who hasn't reached the deletion threshold yet.
		{path: path.Join("user-1", block4.String(), metadata.DeletionMarkFilename), expectedExists: false},
		{path: path.Join("user-1", bucketindex.BlockDeletionMarkFilepath(block4)), expectedExists: false},
		// Should delete a partial block with deletion mark who has reached the deletion threshold.
		{path: path.Join("user-1", block5.String(), metadata.DeletionMarkFilename), expectedExists: false},
		{path: path.Join("user-1", bucketindex.BlockDeletionMarkFilepath(block5)), expectedExists: false},
		// Should not delete a partial block without deletion mark.
		{path: path.Join("user-1", block6.String(), "index"), expectedExists: true},
		// Should delete a partial block with only visit marker.
		{path: path.Join("user-1", block11.String(), BlockVisitMarkerFile), expectedExists: false},
		// Should completely delete blocks for user-3, marked for deletion
		{path: path.Join("user-3", block9.String(), metadata.MetaFilename), expectedExists: false},
		{path: path.Join("user-3", block9.String(), "index"), expectedExists: false},
		{path: path.Join("user-3", block10.String(), metadata.MetaFilename), expectedExists: false},
		{path: path.Join("user-3", block10.String(), "index"), expectedExists: false},
		{path: path.Join("user-3", block10.String(), parquet.ConverterMarkerFileName), expectedExists: false},
		{path: path.Join("user-4", block.DebugMetas, "meta.json"), expectedExists: options.user4FilesExist},
		{path: path.Join("user-6", block13.String(), parquet.ConverterMarkerFileName), expectedExists: true},
		{path: path.Join("user-6", block14.String(), parquet.ConverterMarkerFileName), expectedExists: false},
	} {
		exists, err := bucketClient.Exists(ctx, tc.path)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedExists, exists, tc.path)
	}

	// Check if tenant deletion mark exists
	for _, tc := range []struct {
		user           string
		expectedExists bool
	}{
		{"user-3", true},
		{"user-4", options.user4FilesExist},
	} {
		exists, err := tsdb.TenantDeletionMarkExists(ctx, bucketClient, tc.user)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedExists, exists, tc.user)
	}

	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.runsStarted.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.runsCompleted.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(0), prom_testutil.ToFloat64(cleaner.runsFailed.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(7), prom_testutil.ToFloat64(cleaner.blocksCleanedTotal))
	assert.Equal(t, float64(0), prom_testutil.ToFloat64(cleaner.blocksFailedTotal))

	// Check the updated bucket index.
	for _, tc := range []struct {
		userID         string
		expectedIndex  bool
		expectedBlocks []ulid.ULID
		expectedMarks  []ulid.ULID
	}{
		{
			userID:         "user-1",
			expectedIndex:  true,
			expectedBlocks: []ulid.ULID{block1, block2 /* deleted: block3, block4, block5, block11, partial: block6 */},
			expectedMarks:  []ulid.ULID{block2},
		}, {
			userID:         "user-2",
			expectedIndex:  true,
			expectedBlocks: []ulid.ULID{block8},
			expectedMarks:  []ulid.ULID{},
		}, {
			userID:        "user-3",
			expectedIndex: false,
		}, {
			userID:         "user-6",
			expectedIndex:  true,
			expectedBlocks: []ulid.ULID{block13, block14},
			expectedMarks:  []ulid.ULID{},
		},
	} {
		idx, err := bucketindex.ReadIndex(ctx, bucketClient, tc.userID, nil, logger)
		if !tc.expectedIndex {
			assert.Equal(t, bucketindex.ErrIndexNotFound, err)
			continue
		}

		require.NoError(t, err)
		assert.ElementsMatch(t, tc.expectedBlocks, idx.Blocks.GetULIDs())
		assert.ElementsMatch(t, tc.expectedMarks, idx.BlockDeletionMarks.GetULIDs())
		s, err := bucketindex.ReadSyncStatus(ctx, bucketClient, tc.userID, logger)
		require.NoError(t, err)
		require.Equal(t, bucketindex.Ok, s.Status)
	}

	assert.NoError(t, prom_testutil.GatherAndCompare(reg, strings.NewReader(`
		# HELP cortex_bucket_blocks_count Total number of blocks in the bucket. Includes blocks marked for deletion, but not partial blocks.
		# TYPE cortex_bucket_blocks_count gauge
		cortex_bucket_blocks_count{user="user-1"} 2
		cortex_bucket_blocks_count{user="user-2"} 1
		cortex_bucket_blocks_count{user="user-5"} 2
		cortex_bucket_blocks_count{user="user-6"} 2
		# HELP cortex_bucket_blocks_marked_for_deletion_count Total number of blocks marked for deletion in the bucket.
		# TYPE cortex_bucket_blocks_marked_for_deletion_count gauge
		cortex_bucket_blocks_marked_for_deletion_count{user="user-1"} 1
		cortex_bucket_blocks_marked_for_deletion_count{user="user-2"} 0
		cortex_bucket_blocks_marked_for_deletion_count{user="user-5"} 0
		cortex_bucket_blocks_marked_for_deletion_count{user="user-6"} 0
		# HELP cortex_bucket_blocks_marked_for_no_compaction_count Total number of blocks marked for no compaction in the bucket.
		# TYPE cortex_bucket_blocks_marked_for_no_compaction_count gauge
		cortex_bucket_blocks_marked_for_no_compaction_count{user="user-1"} 0
		cortex_bucket_blocks_marked_for_no_compaction_count{user="user-2"} 0
		cortex_bucket_blocks_marked_for_no_compaction_count{user="user-5"} 1
		cortex_bucket_blocks_marked_for_no_compaction_count{user="user-6"} 0
		# HELP cortex_bucket_blocks_partials_count Total number of partial blocks.
		# TYPE cortex_bucket_blocks_partials_count gauge
		cortex_bucket_blocks_partials_count{user="user-1"} 2
		cortex_bucket_blocks_partials_count{user="user-2"} 0
		cortex_bucket_blocks_partials_count{user="user-5"} 0
		cortex_bucket_blocks_partials_count{user="user-6"} 0
		# HELP cortex_bucket_parquet_blocks_count Total number of parquet blocks in the bucket. Blocks marked for deletion are included.
		# TYPE cortex_bucket_parquet_blocks_count gauge
		cortex_bucket_parquet_blocks_count{user="user-5"} 0
		cortex_bucket_parquet_blocks_count{user="user-6"} 1
		# HELP cortex_bucket_parquet_unconverted_blocks_count Total number of unconverted parquet blocks in the bucket. Blocks marked for deletion are included.
		# TYPE cortex_bucket_parquet_unconverted_blocks_count gauge
		cortex_bucket_parquet_unconverted_blocks_count{user="user-5"} 0
		cortex_bucket_parquet_unconverted_blocks_count{user="user-6"} 0
	`),
		"cortex_bucket_blocks_count",
		"cortex_bucket_parquet_blocks_count",
		"cortex_bucket_parquet_unconverted_blocks_count",
		"cortex_bucket_blocks_marked_for_deletion_count",
		"cortex_bucket_blocks_marked_for_no_compaction_count",
		"cortex_bucket_blocks_partials_count",
	))
}

func TestBlocksCleaner_ShouldContinueOnBlockDeletionFailure(t *testing.T) {
	const userID = "user-1"

	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)

	// Create blocks.
	ctx := context.Background()
	now := time.Now()
	deletionDelay := 12 * time.Hour
	block1 := createTSDBBlock(t, bucketClient, userID, 10, 20, nil)
	block2 := createTSDBBlock(t, bucketClient, userID, 20, 30, nil)
	block3 := createTSDBBlock(t, bucketClient, userID, 30, 40, nil)
	block4 := createTSDBBlock(t, bucketClient, userID, 40, 50, nil)
	createDeletionMark(t, bucketClient, userID, block2, now.Add(-deletionDelay).Add(-time.Hour))
	createDeletionMark(t, bucketClient, userID, block3, now.Add(-deletionDelay).Add(-time.Hour))
	createDeletionMark(t, bucketClient, userID, block4, now.Add(-deletionDelay).Add(-time.Hour))

	// To emulate a failure deleting a block, we wrap the bucket client in a mocked one.
	bucketClient = &cortex_testutil.MockBucketFailure{
		Bucket:         bucketClient,
		DeleteFailures: []string{path.Join(userID, block3.String(), metadata.MetaFilename)},
	}

	cfg := BlocksCleanerConfig{
		DeletionDelay:      deletionDelay,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		BlockRanges:        (&tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}).ToMilliseconds(),
	}

	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	blocksMarkedForDeletion := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", nil, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)
	require.NoError(t, services.StartAndAwaitRunning(ctx, cleaner))
	defer services.StopAndAwaitTerminated(ctx, cleaner) //nolint:errcheck

	for _, tc := range []struct {
		path           string
		expectedExists bool
	}{
		{path: path.Join(userID, block1.String(), metadata.MetaFilename), expectedExists: true},
		{path: path.Join(userID, block2.String(), metadata.MetaFilename), expectedExists: false},
		{path: path.Join(userID, block3.String(), metadata.MetaFilename), expectedExists: true},
		{path: path.Join(userID, block4.String(), metadata.MetaFilename), expectedExists: false},
	} {
		exists, err := bucketClient.Exists(ctx, tc.path)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedExists, exists, tc.path)
	}

	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.runsStarted.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.runsCompleted.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(0), prom_testutil.ToFloat64(cleaner.runsFailed.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(2), prom_testutil.ToFloat64(cleaner.blocksCleanedTotal))
	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.blocksFailedTotal))

	// Check the updated bucket index.
	idx, err := bucketindex.ReadIndex(ctx, bucketClient, userID, nil, logger)
	require.NoError(t, err)
	assert.ElementsMatch(t, []ulid.ULID{block1, block3}, idx.Blocks.GetULIDs())
	assert.ElementsMatch(t, []ulid.ULID{block3}, idx.BlockDeletionMarks.GetULIDs())
}

func TestBlocksCleaner_ShouldRebuildBucketIndexOnCorruptedOne(t *testing.T) {
	const userID = "user-1"

	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)

	// Create blocks.
	ctx := context.Background()
	now := time.Now()
	deletionDelay := 12 * time.Hour
	block1 := createTSDBBlock(t, bucketClient, userID, 10, 20, nil)
	block2 := createTSDBBlock(t, bucketClient, userID, 20, 30, nil)
	block3 := createTSDBBlock(t, bucketClient, userID, 30, 40, nil)
	createDeletionMark(t, bucketClient, userID, block2, now.Add(-deletionDelay).Add(-time.Hour))
	createDeletionMark(t, bucketClient, userID, block3, now.Add(-deletionDelay).Add(time.Hour))

	// Write a corrupted bucket index.
	require.NoError(t, bucketClient.Upload(ctx, path.Join(userID, bucketindex.IndexCompressedFilename), strings.NewReader("invalid!}")))

	cfg := BlocksCleanerConfig{
		DeletionDelay:      deletionDelay,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		BlockRanges:        (&tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}).ToMilliseconds(),
	}

	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	blocksMarkedForDeletion := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", nil, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)
	require.NoError(t, services.StartAndAwaitRunning(ctx, cleaner))
	defer services.StopAndAwaitTerminated(ctx, cleaner) //nolint:errcheck

	for _, tc := range []struct {
		path           string
		expectedExists bool
	}{
		{path: path.Join(userID, block1.String(), metadata.MetaFilename), expectedExists: true},
		{path: path.Join(userID, block2.String(), metadata.MetaFilename), expectedExists: false},
		{path: path.Join(userID, block3.String(), metadata.MetaFilename), expectedExists: true},
	} {
		exists, err := bucketClient.Exists(ctx, tc.path)
		require.NoError(t, err)
		assert.Equal(t, tc.expectedExists, exists, tc.path)
	}

	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.runsStarted.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.runsCompleted.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(0), prom_testutil.ToFloat64(cleaner.runsFailed.WithLabelValues(activeStatus)))
	assert.Equal(t, float64(1), prom_testutil.ToFloat64(cleaner.blocksCleanedTotal))
	assert.Equal(t, float64(0), prom_testutil.ToFloat64(cleaner.blocksFailedTotal))

	// Check the updated bucket index.
	idx, err := bucketindex.ReadIndex(ctx, bucketClient, userID, nil, logger)
	require.NoError(t, err)
	assert.ElementsMatch(t, []ulid.ULID{block1, block3}, idx.Blocks.GetULIDs())
	assert.ElementsMatch(t, []ulid.ULID{block3}, idx.BlockDeletionMarks.GetULIDs())
	s, err := bucketindex.ReadSyncStatus(ctx, bucketClient, userID, logger)
	require.NoError(t, err)
	require.Equal(t, bucketindex.Ok, s.Status)
}

func TestBlocksCleaner_ShouldRemoveMetricsForTenantsNotBelongingAnymoreToTheShard(t *testing.T) {
	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)

	// Create blocks.
	createTSDBBlock(t, bucketClient, "user-1", 10, 20, nil)
	createTSDBBlock(t, bucketClient, "user-1", 20, 30, nil)
	createTSDBBlock(t, bucketClient, "user-2", 30, 40, nil)

	cfg := BlocksCleanerConfig{
		DeletionDelay:      time.Hour,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		BlockRanges:        (&tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}).ToMilliseconds(),
	}

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	blocksMarkedForDeletion := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", reg, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)
	activeUsers, deleteUsers, err := cleaner.scanUsers(ctx)
	require.NoError(t, err)
	require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, true))
	require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))

	assert.NoError(t, prom_testutil.GatherAndCompare(reg, strings.NewReader(`
		# HELP cortex_bucket_blocks_count Total number of blocks in the bucket. Includes blocks marked for deletion, but not partial blocks.
		# TYPE cortex_bucket_blocks_count gauge
		cortex_bucket_blocks_count{user="user-1"} 2
		cortex_bucket_blocks_count{user="user-2"} 1
		# HELP cortex_bucket_blocks_marked_for_deletion_count Total number of blocks marked for deletion in the bucket.
		# TYPE cortex_bucket_blocks_marked_for_deletion_count gauge
		cortex_bucket_blocks_marked_for_deletion_count{user="user-1"} 0
		cortex_bucket_blocks_marked_for_deletion_count{user="user-2"} 0
		# HELP cortex_bucket_blocks_partials_count Total number of partial blocks.
		# TYPE cortex_bucket_blocks_partials_count gauge
		cortex_bucket_blocks_partials_count{user="user-1"} 0
		cortex_bucket_blocks_partials_count{user="user-2"} 0
	`),
		"cortex_bucket_blocks_count",
		"cortex_bucket_blocks_marked_for_deletion_count",
		"cortex_bucket_blocks_partials_count",
	))

	// Override the users scanner to reconfigure it to only return a subset of users.
	cleaner.usersScanner, err = users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cleaner.usersScanner = users.NewShardedScanner(cleaner.usersScanner, func(userID string) (bool, error) { return userID == "user-1", nil }, logger)

	// Create new blocks, to double check expected metrics have changed.
	createTSDBBlock(t, bucketClient, "user-1", 40, 50, nil)
	createTSDBBlock(t, bucketClient, "user-2", 50, 60, nil)

	activeUsers, deleteUsers, err = cleaner.scanUsers(ctx)
	require.NoError(t, err)
	require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, false))
	require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))

	assert.NoError(t, prom_testutil.GatherAndCompare(reg, strings.NewReader(`
		# HELP cortex_bucket_blocks_count Total number of blocks in the bucket. Includes blocks marked for deletion, but not partial blocks.
		# TYPE cortex_bucket_blocks_count gauge
		cortex_bucket_blocks_count{user="user-1"} 3
		# HELP cortex_bucket_blocks_marked_for_deletion_count Total number of blocks marked for deletion in the bucket.
		# TYPE cortex_bucket_blocks_marked_for_deletion_count gauge
		cortex_bucket_blocks_marked_for_deletion_count{user="user-1"} 0
		# HELP cortex_bucket_blocks_partials_count Total number of partial blocks.
		# TYPE cortex_bucket_blocks_partials_count gauge
		cortex_bucket_blocks_partials_count{user="user-1"} 0
	`),
		"cortex_bucket_blocks_count",
		"cortex_bucket_blocks_marked_for_deletion_count",
		"cortex_bucket_blocks_partials_count",
	))
}

func TestBlocksCleaner_ListBlocksOutsideRetentionPeriod(t *testing.T) {
	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)
	ctx := context.Background()
	logger := log.NewNopLogger()

	id1 := createTSDBBlock(t, bucketClient, "user-1", 5000, 6000, nil)
	id2 := createTSDBBlock(t, bucketClient, "user-1", 6000, 7000, nil)
	id3 := createTSDBBlock(t, bucketClient, "user-1", 7000, 8000, nil)

	w := bucketindex.NewUpdater(bucketClient, "user-1", nil, logger)
	idx, _, _, err := w.UpdateIndex(ctx, nil)
	require.NoError(t, err)

	assert.ElementsMatch(t, []ulid.ULID{id1, id2, id3}, idx.Blocks.GetULIDs())

	// Excessive retention period (wrapping epoch)
	result := listBlocksOutsideRetentionPeriod(idx, time.Unix(10, 0).Add(-time.Hour))
	assert.ElementsMatch(t, []ulid.ULID{}, result.GetULIDs())

	// Normal operation - varying retention period.
	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(6, 0))
	assert.ElementsMatch(t, []ulid.ULID{}, result.GetULIDs())

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(7, 0))
	assert.ElementsMatch(t, []ulid.ULID{id1}, result.GetULIDs())

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(8, 0))
	assert.ElementsMatch(t, []ulid.ULID{id1, id2}, result.GetULIDs())

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(9, 0))
	assert.ElementsMatch(t, []ulid.ULID{id1, id2, id3}, result.GetULIDs())

	// Avoiding redundant marking - blocks already marked for deletion.

	mark1 := &bucketindex.BlockDeletionMark{ID: id1}
	mark2 := &bucketindex.BlockDeletionMark{ID: id2}

	idx.BlockDeletionMarks = bucketindex.BlockDeletionMarks{mark1}

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(7, 0))
	assert.ElementsMatch(t, []ulid.ULID{}, result.GetULIDs())

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(8, 0))
	assert.ElementsMatch(t, []ulid.ULID{id2}, result.GetULIDs())

	idx.BlockDeletionMarks = bucketindex.BlockDeletionMarks{mark1, mark2}

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(7, 0))
	assert.ElementsMatch(t, []ulid.ULID{}, result.GetULIDs())

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(8, 0))
	assert.ElementsMatch(t, []ulid.ULID{}, result.GetULIDs())

	result = listBlocksOutsideRetentionPeriod(idx, time.Unix(9, 0))
	assert.ElementsMatch(t, []ulid.ULID{id3}, result.GetULIDs())
}

func TestBlocksCleaner_ShouldRemoveBlocksOutsideRetentionPeriod(t *testing.T) {
	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)

	ts := func(hours int) int64 {
		return time.Now().Add(time.Duration(hours)*time.Hour).Unix() * 1000
	}

	block1 := createTSDBBlock(t, bucketClient, "user-1", ts(-10), ts(-8), nil)
	block2 := createTSDBBlock(t, bucketClient, "user-1", ts(-8), ts(-6), nil)
	block3 := createTSDBBlock(t, bucketClient, "user-2", ts(-10), ts(-8), nil)
	block4 := createTSDBBlock(t, bucketClient, "user-2", ts(-8), ts(-6), nil)

	cfg := BlocksCleanerConfig{
		DeletionDelay:      time.Hour,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		BlockRanges:        (&tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}).ToMilliseconds(),
	}

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewPedanticRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	blocksMarkedForDeletion := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", reg, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)

	assertBlockExists := func(user string, block ulid.ULID, expectExists bool) {
		exists, err := bucketClient.Exists(ctx, path.Join(user, block.String(), metadata.MetaFilename))
		require.NoError(t, err)
		assert.Equal(t, expectExists, exists)
	}

	// Existing behaviour - retention period disabled.
	{
		// clean up cleaner visit marker before running test
		bucketClient.Delete(ctx, path.Join("user-1", GetCleanerVisitMarkerFilePath())) //nolint:errcheck
		bucketClient.Delete(ctx, path.Join("user-2", GetCleanerVisitMarkerFilePath())) //nolint:errcheck

		cfgProvider.userRetentionPeriods["user-1"] = 0
		cfgProvider.userRetentionPeriods["user-2"] = 0

		activeUsers, deleteUsers, err := cleaner.scanUsers(ctx)
		require.NoError(t, err)
		require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, true))
		require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))
		assertBlockExists("user-1", block1, true)
		assertBlockExists("user-1", block2, true)
		assertBlockExists("user-2", block3, true)
		assertBlockExists("user-2", block4, true)

		assert.NoError(t, prom_testutil.GatherAndCompare(reg, strings.NewReader(`
			# HELP cortex_bucket_blocks_count Total number of blocks in the bucket. Includes blocks marked for deletion, but not partial blocks.
			# TYPE cortex_bucket_blocks_count gauge
			cortex_bucket_blocks_count{user="user-1"} 2
			cortex_bucket_blocks_count{user="user-2"} 2
			# HELP cortex_bucket_blocks_marked_for_deletion_count Total number of blocks marked for deletion in the bucket.
			# TYPE cortex_bucket_blocks_marked_for_deletion_count gauge
			cortex_bucket_blocks_marked_for_deletion_count{user="user-1"} 0
			cortex_bucket_blocks_marked_for_deletion_count{user="user-2"} 0
			# HELP cortex_compactor_blocks_marked_for_deletion_total Total number of blocks marked for deletion in compactor.
			# TYPE cortex_compactor_blocks_marked_for_deletion_total counter
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-1"} 0
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-2"} 0
			`),
			"cortex_bucket_blocks_count",
			"cortex_bucket_blocks_marked_for_deletion_count",
			"cortex_compactor_blocks_marked_for_deletion_total",
		))
	}

	// Retention enabled only for a single user, but does nothing.
	{
		// clean up cleaner visit marker before running test
		bucketClient.Delete(ctx, path.Join("user-1", GetCleanerVisitMarkerFilePath())) //nolint:errcheck
		bucketClient.Delete(ctx, path.Join("user-2", GetCleanerVisitMarkerFilePath())) //nolint:errcheck

		cfgProvider.userRetentionPeriods["user-1"] = 9 * time.Hour

		activeUsers, deleteUsers, err := cleaner.scanUsers(ctx)
		require.NoError(t, err)
		require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, false))
		require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))
		assertBlockExists("user-1", block1, true)
		assertBlockExists("user-1", block2, true)
		assertBlockExists("user-2", block3, true)
		assertBlockExists("user-2", block4, true)
	}

	// Retention enabled only for a single user, marking a single block.
	// Note the block won't be deleted yet due to deletion delay.
	{
		// clean up cleaner visit marker before running test
		bucketClient.Delete(ctx, path.Join("user-1", GetCleanerVisitMarkerFilePath())) //nolint:errcheck
		bucketClient.Delete(ctx, path.Join("user-2", GetCleanerVisitMarkerFilePath())) //nolint:errcheck

		cfgProvider.userRetentionPeriods["user-1"] = 7 * time.Hour

		activeUsers, deleteUsers, err := cleaner.scanUsers(ctx)
		require.NoError(t, err)
		require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, false))
		require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))
		assertBlockExists("user-1", block1, true)
		assertBlockExists("user-1", block2, true)
		assertBlockExists("user-2", block3, true)
		assertBlockExists("user-2", block4, true)

		assert.NoError(t, prom_testutil.GatherAndCompare(reg, strings.NewReader(`
			# HELP cortex_bucket_blocks_count Total number of blocks in the bucket. Includes blocks marked for deletion, but not partial blocks.
			# TYPE cortex_bucket_blocks_count gauge
			cortex_bucket_blocks_count{user="user-1"} 2
			cortex_bucket_blocks_count{user="user-2"} 2
			# HELP cortex_bucket_blocks_marked_for_deletion_count Total number of blocks marked for deletion in the bucket.
			# TYPE cortex_bucket_blocks_marked_for_deletion_count gauge
			cortex_bucket_blocks_marked_for_deletion_count{user="user-1"} 1
			cortex_bucket_blocks_marked_for_deletion_count{user="user-2"} 0
			# HELP cortex_compactor_blocks_marked_for_deletion_total Total number of blocks marked for deletion in compactor.
			# TYPE cortex_compactor_blocks_marked_for_deletion_total counter
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-1"} 1
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-2"} 0
			`),
			"cortex_bucket_blocks_count",
			"cortex_bucket_blocks_marked_for_deletion_count",
			"cortex_compactor_blocks_marked_for_deletion_total",
		))
	}

	// Marking the block again, before the deletion occurs, should not cause an error.
	{
		// clean up cleaner visit marker before running test
		bucketClient.Delete(ctx, path.Join("user-1", GetCleanerVisitMarkerFilePath())) //nolint:errcheck
		bucketClient.Delete(ctx, path.Join("user-2", GetCleanerVisitMarkerFilePath())) //nolint:errcheck

		activeUsers, deleteUsers, err := cleaner.scanUsers(ctx)
		require.NoError(t, err)
		require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, false))
		require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))
		assertBlockExists("user-1", block1, true)
		assertBlockExists("user-1", block2, true)
		assertBlockExists("user-2", block3, true)
		assertBlockExists("user-2", block4, true)
	}

	// Reduce the deletion delay. Now the block will be deleted.
	{
		// clean up cleaner visit marker before running test
		bucketClient.Delete(ctx, path.Join("user-1", GetCleanerVisitMarkerFilePath())) //nolint:errcheck
		bucketClient.Delete(ctx, path.Join("user-2", GetCleanerVisitMarkerFilePath())) //nolint:errcheck

		cleaner.cfg.DeletionDelay = 0

		activeUsers, deleteUsers, err := cleaner.scanUsers(ctx)
		require.NoError(t, err)
		require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, false))
		require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))
		assertBlockExists("user-1", block1, false)
		assertBlockExists("user-1", block2, true)
		assertBlockExists("user-2", block3, true)
		assertBlockExists("user-2", block4, true)

		assert.NoError(t, prom_testutil.GatherAndCompare(reg, strings.NewReader(`
			# HELP cortex_bucket_blocks_count Total number of blocks in the bucket. Includes blocks marked for deletion, but not partial blocks.
			# TYPE cortex_bucket_blocks_count gauge
			cortex_bucket_blocks_count{user="user-1"} 1
			cortex_bucket_blocks_count{user="user-2"} 2
			# HELP cortex_bucket_blocks_marked_for_deletion_count Total number of blocks marked for deletion in the bucket.
			# TYPE cortex_bucket_blocks_marked_for_deletion_count gauge
			cortex_bucket_blocks_marked_for_deletion_count{user="user-1"} 0
			cortex_bucket_blocks_marked_for_deletion_count{user="user-2"} 0
			# HELP cortex_compactor_blocks_marked_for_deletion_total Total number of blocks marked for deletion in compactor.
			# TYPE cortex_compactor_blocks_marked_for_deletion_total counter
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-1"} 1
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-2"} 0
			`),
			"cortex_bucket_blocks_count",
			"cortex_bucket_blocks_marked_for_deletion_count",
			"cortex_compactor_blocks_marked_for_deletion_total",
		))
	}

	// Retention enabled for other user; test deleting multiple blocks.
	{
		// clean up cleaner visit marker before running test
		bucketClient.Delete(ctx, path.Join("user-1", GetCleanerVisitMarkerFilePath())) //nolint:errcheck
		bucketClient.Delete(ctx, path.Join("user-2", GetCleanerVisitMarkerFilePath())) //nolint:errcheck

		cfgProvider.userRetentionPeriods["user-2"] = 5 * time.Hour

		activeUsers, deleteUsers, err := cleaner.scanUsers(ctx)
		require.NoError(t, err)
		require.NoError(t, cleaner.cleanUpActiveUsers(ctx, activeUsers, false))
		require.NoError(t, cleaner.cleanDeletedUsers(ctx, deleteUsers))
		assertBlockExists("user-1", block1, false)
		assertBlockExists("user-1", block2, true)
		assertBlockExists("user-2", block3, false)
		assertBlockExists("user-2", block4, false)

		assert.NoError(t, prom_testutil.GatherAndCompare(reg, strings.NewReader(`
			# HELP cortex_bucket_blocks_count Total number of blocks in the bucket. Includes blocks marked for deletion, but not partial blocks.
			# TYPE cortex_bucket_blocks_count gauge
			cortex_bucket_blocks_count{user="user-1"} 1
			cortex_bucket_blocks_count{user="user-2"} 0
			# HELP cortex_bucket_blocks_marked_for_deletion_count Total number of blocks marked for deletion in the bucket.
			# TYPE cortex_bucket_blocks_marked_for_deletion_count gauge
			cortex_bucket_blocks_marked_for_deletion_count{user="user-1"} 0
			cortex_bucket_blocks_marked_for_deletion_count{user="user-2"} 0
			# HELP cortex_compactor_blocks_marked_for_deletion_total Total number of blocks marked for deletion in compactor.
			# TYPE cortex_compactor_blocks_marked_for_deletion_total counter
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-1"} 1
			cortex_compactor_blocks_marked_for_deletion_total{reason="retention",user="user-2"} 2
			`),
			"cortex_bucket_blocks_count",
			"cortex_bucket_blocks_marked_for_deletion_count",
			"cortex_compactor_blocks_marked_for_deletion_total",
		))
	}
}

func TestBlocksCleaner_CleanPartitionedGroupInfo(t *testing.T) {
	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)

	ts := func(hours int) int64 {
		return time.Now().Add(time.Duration(hours)*time.Hour).Unix() * 1000
	}

	userID := "user-1"
	partitionedGroupID := uint32(123)
	partitionCount := 1
	startTime := ts(-10)
	endTime := ts(-8)
	block1 := createTSDBBlock(t, bucketClient, userID, startTime, endTime, nil)
	block2 := createTSDBBlock(t, bucketClient, userID, startTime, endTime, nil)
	createNoCompactionMark(t, bucketClient, userID, block2)

	cfg := BlocksCleanerConfig{
		DeletionDelay:      time.Hour,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		ShardingStrategy:   util.ShardingStrategyShuffle,
		CompactionStrategy: util.CompactionStrategyPartitioning,
		BlockRanges:        (&tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}).ToMilliseconds(),
	}

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewPedanticRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	blocksMarkedForDeletion := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", reg, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)

	userBucket := bucket.NewUserBucketClient(userID, bucketClient, cfgProvider)

	partitionedGroupInfo := PartitionedGroupInfo{
		PartitionedGroupID: partitionedGroupID,
		PartitionCount:     partitionCount,
		Partitions: []Partition{
			{
				PartitionID: 0,
				Blocks:      []ulid.ULID{block1, block2},
			},
		},
		RangeStart:   startTime,
		RangeEnd:     endTime,
		CreationTime: time.Now().Add(-5 * time.Minute).Unix(),
		Version:      PartitionedGroupInfoVersion1,
	}
	_, err = UpdatePartitionedGroupInfo(ctx, userBucket, logger, partitionedGroupInfo)
	require.NoError(t, err)

	visitMarker := &partitionVisitMarker{
		PartitionedGroupID: partitionedGroupID,
		PartitionID:        0,
		Status:             Completed,
		VisitTime:          time.Now().Add(-2 * time.Minute).Unix(),
	}
	visitMarkerManager := NewVisitMarkerManager(userBucket, logger, "dummy-cleaner", visitMarker)
	err = visitMarkerManager.updateVisitMarker(ctx)
	require.NoError(t, err)

	cleaner.cleanPartitionedGroupInfo(ctx, userBucket, logger, userID)

	partitionedGroupFileExists, err := userBucket.Exists(ctx, GetPartitionedGroupFile(partitionedGroupID))
	require.NoError(t, err)
	require.False(t, partitionedGroupFileExists)

	block1DeletionMarkerExists, err := userBucket.Exists(ctx, path.Join(block1.String(), metadata.DeletionMarkFilename))
	require.NoError(t, err)
	require.True(t, block1DeletionMarkerExists)

	block2DeletionMarkerExists, err := userBucket.Exists(ctx, path.Join(block2.String(), metadata.DeletionMarkFilename))
	require.NoError(t, err)
	require.False(t, block2DeletionMarkerExists)
}

func TestBlocksCleaner_DeleteEmptyBucketIndex(t *testing.T) {
	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)

	userID := "user-1"

	cfg := BlocksCleanerConfig{
		DeletionDelay:      time.Hour,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		ShardingStrategy:   util.ShardingStrategyShuffle,
		CompactionStrategy: util.CompactionStrategyPartitioning,
		BlockRanges:        (&tsdb.DurationList{2 * time.Hour, 12 * time.Hour, 24 * time.Hour}).ToMilliseconds(),
	}

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewPedanticRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, reg)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	blocksMarkedForDeletion := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name: blocksMarkedForDeletionName,
		Help: blocksMarkedForDeletionHelp,
	}, append(commonLabels, reasonLabelName))
	dummyGaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{}, []string{"test"})

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 60*time.Second, cfgProvider, logger, "test-cleaner", reg, time.Minute, 30*time.Second, blocksMarkedForDeletion, dummyGaugeVec)

	userBucket := bucket.NewUserBucketClient(userID, bucketClient, cfgProvider)

	debugMetaFile := path.Join(block.DebugMetas, "meta.json")
	require.NoError(t, userBucket.Upload(context.Background(), debugMetaFile, strings.NewReader("some random content here")))

	partitionedGroupInfo := PartitionedGroupInfo{
		PartitionedGroupID: 1234,
		PartitionCount:     1,
		Partitions: []Partition{
			{
				PartitionID: 0,
				Blocks:      []ulid.ULID{},
			},
		},
		RangeStart:   0,
		RangeEnd:     2,
		CreationTime: time.Now().Add(-5 * time.Minute).Unix(),
		Version:      PartitionedGroupInfoVersion1,
	}
	_, err = UpdatePartitionedGroupInfo(ctx, userBucket, logger, partitionedGroupInfo)
	require.NoError(t, err)
	partitionedGroupFile := GetPartitionedGroupFile(partitionedGroupInfo.PartitionedGroupID)

	err = cleaner.cleanUser(ctx, logger, userBucket, userID, false)
	require.NoError(t, err)

	_, err = bucketindex.ReadIndex(ctx, bucketClient, userID, cfgProvider, logger)
	require.ErrorIs(t, err, bucketindex.ErrIndexNotFound)

	_, err = userBucket.WithExpectedErrs(userBucket.IsObjNotFoundErr).Get(ctx, bucketindex.SyncStatusFile)
	require.True(t, userBucket.IsObjNotFoundErr(err))

	_, err = userBucket.WithExpectedErrs(userBucket.IsObjNotFoundErr).Get(ctx, debugMetaFile)
	require.True(t, userBucket.IsObjNotFoundErr(err))

	_, err = userBucket.WithExpectedErrs(userBucket.IsObjNotFoundErr).Get(ctx, partitionedGroupFile)
	require.True(t, userBucket.IsObjNotFoundErr(err))
}

func TestBlocksCleaner_ParquetMetrics(t *testing.T) {
	// Create metrics
	reg := prometheus.NewPedanticRegistry()
	blocksMarkedForDeletion := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "cortex_compactor_blocks_marked_for_deletion_total",
			Help: "Total number of blocks marked for deletion in compactor.",
		},
		[]string{"user", "reason"},
	)
	remainingPlannedCompactions := promauto.With(reg).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cortex_compactor_remaining_planned_compactions",
			Help: "Total number of remaining planned compactions.",
		},
		[]string{"user"},
	)

	// Create the blocks cleaner
	cleaner := NewBlocksCleaner(
		BlocksCleanerConfig{
			BlockRanges: (&tsdb.DurationList{
				2 * time.Hour,
				12 * time.Hour,
			}).ToMilliseconds(),
		},
		nil, // bucket not needed
		nil, // usersScanner not needed
		0,
		&mockConfigProvider{
			parquetConverterEnabled: map[string]bool{
				"user1": true,
			},
		},
		log.NewNopLogger(),
		"test",
		reg,
		0,
		0,
		blocksMarkedForDeletion,
		remainingPlannedCompactions,
	)

	// Create test blocks in the index
	now := time.Now()
	idx := &bucketindex.Index{
		Blocks: bucketindex.Blocks{
			{
				ID:      ulid.MustNew(ulid.Now(), rand.Reader),
				MinTime: now.Add(-3 * time.Hour).UnixMilli(),
				MaxTime: now.UnixMilli(),
				Parquet: &parquet.ConverterMarkMeta{},
			},
			{
				ID:      ulid.MustNew(ulid.Now(), rand.Reader),
				MinTime: now.Add(-3 * time.Hour).UnixMilli(),
				MaxTime: now.UnixMilli(),
				Parquet: nil,
			},
			{
				ID:      ulid.MustNew(ulid.Now(), rand.Reader),
				MinTime: now.Add(-5 * time.Hour).UnixMilli(),
				MaxTime: now.UnixMilli(),
				Parquet: nil,
			},
		},
	}

	// Update metrics
	cleaner.updateBucketMetrics("user1", true, idx, 0, 0)

	// Verify metrics
	require.NoError(t, prom_testutil.CollectAndCompare(cleaner.tenantParquetBlocks, strings.NewReader(`
		# HELP cortex_bucket_parquet_blocks_count Total number of parquet blocks in the bucket. Blocks marked for deletion are included.
		# TYPE cortex_bucket_parquet_blocks_count gauge
		cortex_bucket_parquet_blocks_count{user="user1"} 1
	`)))

	require.NoError(t, prom_testutil.CollectAndCompare(cleaner.tenantParquetUnConvertedBlocks, strings.NewReader(`
		# HELP cortex_bucket_parquet_unconverted_blocks_count Total number of unconverted parquet blocks in the bucket. Blocks marked for deletion are included.
		# TYPE cortex_bucket_parquet_unconverted_blocks_count gauge
		cortex_bucket_parquet_unconverted_blocks_count{user="user1"} 2
	`)))
}

func TestBlocksCleaner_EmitUserMetrics(t *testing.T) {
	bucketClient, _ := cortex_testutil.PrepareFilesystemBucket(t)
	bucketClient = bucketindex.BucketWithGlobalMarkers(bucketClient)

	cfg := BlocksCleanerConfig{
		DeletionDelay:      time.Hour,
		CleanupInterval:    time.Minute,
		CleanupConcurrency: 1,
		ShardingStrategy:   util.ShardingStrategyShuffle,
		CompactionStrategy: util.CompactionStrategyPartitioning,
	}

	ctx := context.Background()
	logger := log.NewNopLogger()
	registry := prometheus.NewPedanticRegistry()
	scanner, err := users.NewScanner(tsdb.UsersScannerConfig{
		Strategy: tsdb.UserScanStrategyList,
	}, bucketClient, logger, registry)
	require.NoError(t, err)
	cfgProvider := newMockConfigProvider()
	dummyCounterVec := prometheus.NewCounterVec(prometheus.CounterOpts{}, []string{"test"})
	remainingPlannedCompactions := promauto.With(registry).NewGaugeVec(prometheus.GaugeOpts{
		Name: "cortex_compactor_remaining_planned_compactions",
		Help: "Total number of plans that remain to be compacted. Only available with shuffle-sharding strategy",
	}, commonLabels)

	cleaner := NewBlocksCleaner(cfg, bucketClient, scanner, 15*time.Minute, cfgProvider, logger, "test-cleaner", registry, time.Minute, 30*time.Second, dummyCounterVec, remainingPlannedCompactions)

	ts := func(hours int) int64 {
		return time.Now().Add(time.Duration(hours)*time.Hour).Unix() * 1000
	}

	userID := "user-1"
	partitionedGroupID := uint32(123)
	partitionCount := 5
	startTime := ts(-10)
	endTime := ts(-8)
	userBucket := bucket.NewUserBucketClient(userID, bucketClient, cfgProvider)
	partitionedGroupInfo := PartitionedGroupInfo{
		PartitionedGroupID: partitionedGroupID,
		PartitionCount:     partitionCount,
		Partitions: []Partition{
			{
				PartitionID: 0,
			},
			{
				PartitionID: 1,
			},
			{
				PartitionID: 2,
			},
			{
				PartitionID: 3,
			},
			{
				PartitionID: 4,
			},
		},
		RangeStart:   startTime,
		RangeEnd:     endTime,
		CreationTime: time.Now().Add(-1 * time.Hour).Unix(),
		Version:      PartitionedGroupInfoVersion1,
	}
	_, err = UpdatePartitionedGroupInfo(ctx, userBucket, logger, partitionedGroupInfo)
	require.NoError(t, err)

	//InProgress with valid VisitTime
	v0 := &partitionVisitMarker{
		PartitionedGroupID: partitionedGroupID,
		PartitionID:        0,
		Status:             InProgress,
		VisitTime:          time.Now().Add(-2 * time.Minute).Unix(),
	}
	v0Manager := NewVisitMarkerManager(userBucket, logger, "dummy-cleaner", v0)
	err = v0Manager.updateVisitMarker(ctx)
	require.NoError(t, err)

	//InProgress with expired VisitTime
	v1 := &partitionVisitMarker{
		PartitionedGroupID: partitionedGroupID,
		PartitionID:        1,
		Status:             InProgress,
		VisitTime:          time.Now().Add(-30 * time.Minute).Unix(),
	}
	v1Manager := NewVisitMarkerManager(userBucket, logger, "dummy-cleaner", v1)
	err = v1Manager.updateVisitMarker(ctx)
	require.NoError(t, err)

	//V2 and V3 are pending
	//V4 is completed
	v4 := &partitionVisitMarker{
		PartitionedGroupID: partitionedGroupID,
		PartitionID:        4,
		Status:             Completed,
		VisitTime:          time.Now().Add(-20 * time.Minute).Unix(),
	}
	v4Manager := NewVisitMarkerManager(userBucket, logger, "dummy-cleaner", v4)
	err = v4Manager.updateVisitMarker(ctx)
	require.NoError(t, err)

	cleaner.emitUserParititionMetrics(ctx, logger, userBucket, userID)

	metricNames := []string{
		"cortex_compactor_remaining_planned_compactions",
		"cortex_compactor_in_progress_compactions",
		"cortex_compactor_oldest_partition_offset",
	}

	// Check tracked Prometheus metrics
	expectedMetrics := `
        # HELP cortex_compactor_in_progress_compactions Total number of in progress compactions. Only available with shuffle-sharding strategy and partitioning compaction strategy
		# TYPE cortex_compactor_in_progress_compactions gauge
		cortex_compactor_in_progress_compactions{user="user-1"} 1
		# HELP cortex_compactor_oldest_partition_offset Time in seconds between now and the oldest created partition group not completed. Only available with shuffle-sharding strategy and partitioning compaction strategy
		# TYPE cortex_compactor_oldest_partition_offset gauge
		cortex_compactor_oldest_partition_offset{user="user-1"} 3600
        # HELP cortex_compactor_remaining_planned_compactions Total number of plans that remain to be compacted. Only available with shuffle-sharding strategy
		# TYPE cortex_compactor_remaining_planned_compactions gauge
		cortex_compactor_remaining_planned_compactions{user="user-1"} 3
	`

	assert.NoError(t, prom_testutil.GatherAndCompare(registry, strings.NewReader(expectedMetrics), metricNames...))
}

type mockConfigProvider struct {
	userRetentionPeriods    map[string]time.Duration
	parquetConverterEnabled map[string]bool
}

func (m *mockConfigProvider) ParquetConverterEnabled(userID string) bool {
	if result, ok := m.parquetConverterEnabled[userID]; ok {
		return result
	}
	return false
}

func newMockConfigProvider() *mockConfigProvider {
	return &mockConfigProvider{
		userRetentionPeriods:    make(map[string]time.Duration),
		parquetConverterEnabled: make(map[string]bool),
	}
}

func (m *mockConfigProvider) CompactorBlocksRetentionPeriod(user string) time.Duration {
	if result, ok := m.userRetentionPeriods[user]; ok {
		return result
	}
	return 0
}

func (m *mockConfigProvider) S3SSEType(user string) string {
	return ""
}

func (m *mockConfigProvider) S3SSEKMSKeyID(userID string) string {
	return ""
}

func (m *mockConfigProvider) S3SSEKMSEncryptionContext(userID string) string {
	return ""
}
