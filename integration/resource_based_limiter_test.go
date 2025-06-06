//go:build requires_docker
// +build requires_docker

package integration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cortexproject/cortex/integration/e2e"
	e2edb "github.com/cortexproject/cortex/integration/e2e/db"
	"github.com/cortexproject/cortex/integration/e2ecortex"
)

func Test_ResourceBasedLimiter_shouldStartWithoutError(t *testing.T) {
	s, err := e2e.NewScenario(networkName)
	require.NoError(t, err)
	defer s.Close()

	flags := mergeFlags(BlocksStorageFlags(), map[string]string{
		"-resource-monitor.resources": "cpu,heap",
	})

	// Start dependencies.
	consul := e2edb.NewConsul()
	minio := e2edb.NewMinio(9000, flags["-blocks-storage.s3.bucket-name"])
	require.NoError(t, s.StartAndWaitReady(consul, minio))

	// Start Cortex components.
	ingester := e2ecortex.NewIngester("ingester", e2ecortex.RingStoreConsul, consul.NetworkHTTPEndpoint(), mergeFlags(flags, map[string]string{
		"-ingester.query-protection.rejection.threshold.cpu-utilization":  "0.8",
		"-ingester.query-protection.rejection.threshold.heap-utilization": "0.8",
	}), "")
	storeGateway := e2ecortex.NewStoreGateway("store-gateway", e2ecortex.RingStoreConsul, consul.NetworkHTTPEndpoint(), mergeFlags(flags, map[string]string{
		"-store-gateway.query-protection.rejection.threshold.cpu-utilization":  "0.8",
		"-store-gateway.query-protection.rejection.threshold.heap-utilization": "0.8",
	}), "")
	require.NoError(t, s.StartAndWaitReady(ingester, storeGateway))
}
