// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: rules.proto

package rulespb

import (
	fmt "fmt"
	_ "github.com/cortexproject/cortex/pkg/cortexpb"
	github_com_cortexproject_cortex_pkg_cortexpb "github.com/cortexproject/cortex/pkg/cortexpb"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	types "github.com/gogo/protobuf/types"
	_ "github.com/golang/protobuf/ptypes/duration"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// RuleGroupDesc is a proto representation of a cortex rule group.
type RuleGroupDesc struct {
	Name      string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace string        `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Interval  time.Duration `protobuf:"bytes,3,opt,name=interval,proto3,stdduration" json:"interval"`
	Rules     []*RuleDesc   `protobuf:"bytes,4,rep,name=rules,proto3" json:"rules,omitempty"`
	User      string        `protobuf:"bytes,6,opt,name=user,proto3" json:"user,omitempty"`
	// The options field can be used to extend Cortex Ruler functionality without
	// having to repeatedly redefine the proto description. It can also be leveraged
	// to create custom `ManagerOpts` based on rule configs which can then be passed
	// to the Prometheus Manager.
	Options     []*types.Any                                                `protobuf:"bytes,9,rep,name=options,proto3" json:"options,omitempty"`
	Limit       int64                                                       `protobuf:"varint,10,opt,name=limit,proto3" json:"limit,omitempty"`
	QueryOffset *time.Duration                                              `protobuf:"bytes,11,opt,name=queryOffset,proto3,stdduration" json:"queryOffset,omitempty"`
	Labels      []github_com_cortexproject_cortex_pkg_cortexpb.LabelAdapter `protobuf:"bytes,12,rep,name=labels,proto3,customtype=github.com/cortexproject/cortex/pkg/cortexpb.LabelAdapter" json:"labels"`
}

func (m *RuleGroupDesc) Reset()      { *m = RuleGroupDesc{} }
func (*RuleGroupDesc) ProtoMessage() {}
func (*RuleGroupDesc) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e722d3e922f0937, []int{0}
}
func (m *RuleGroupDesc) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RuleGroupDesc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RuleGroupDesc.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RuleGroupDesc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuleGroupDesc.Merge(m, src)
}
func (m *RuleGroupDesc) XXX_Size() int {
	return m.Size()
}
func (m *RuleGroupDesc) XXX_DiscardUnknown() {
	xxx_messageInfo_RuleGroupDesc.DiscardUnknown(m)
}

var xxx_messageInfo_RuleGroupDesc proto.InternalMessageInfo

func (m *RuleGroupDesc) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RuleGroupDesc) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *RuleGroupDesc) GetInterval() time.Duration {
	if m != nil {
		return m.Interval
	}
	return 0
}

func (m *RuleGroupDesc) GetRules() []*RuleDesc {
	if m != nil {
		return m.Rules
	}
	return nil
}

func (m *RuleGroupDesc) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *RuleGroupDesc) GetOptions() []*types.Any {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *RuleGroupDesc) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *RuleGroupDesc) GetQueryOffset() *time.Duration {
	if m != nil {
		return m.QueryOffset
	}
	return nil
}

// RuleDesc is a proto representation of a Prometheus Rule
type RuleDesc struct {
	Expr          string                                                      `protobuf:"bytes,1,opt,name=expr,proto3" json:"expr,omitempty"`
	Record        string                                                      `protobuf:"bytes,2,opt,name=record,proto3" json:"record,omitempty"`
	Alert         string                                                      `protobuf:"bytes,3,opt,name=alert,proto3" json:"alert,omitempty"`
	For           time.Duration                                               `protobuf:"bytes,4,opt,name=for,proto3,stdduration" json:"for"`
	Labels        []github_com_cortexproject_cortex_pkg_cortexpb.LabelAdapter `protobuf:"bytes,5,rep,name=labels,proto3,customtype=github.com/cortexproject/cortex/pkg/cortexpb.LabelAdapter" json:"labels"`
	Annotations   []github_com_cortexproject_cortex_pkg_cortexpb.LabelAdapter `protobuf:"bytes,6,rep,name=annotations,proto3,customtype=github.com/cortexproject/cortex/pkg/cortexpb.LabelAdapter" json:"annotations"`
	KeepFiringFor time.Duration                                               `protobuf:"bytes,13,opt,name=keepFiringFor,proto3,stdduration" json:"keep_firing_for"`
}

func (m *RuleDesc) Reset()      { *m = RuleDesc{} }
func (*RuleDesc) ProtoMessage() {}
func (*RuleDesc) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e722d3e922f0937, []int{1}
}
func (m *RuleDesc) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RuleDesc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RuleDesc.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RuleDesc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuleDesc.Merge(m, src)
}
func (m *RuleDesc) XXX_Size() int {
	return m.Size()
}
func (m *RuleDesc) XXX_DiscardUnknown() {
	xxx_messageInfo_RuleDesc.DiscardUnknown(m)
}

var xxx_messageInfo_RuleDesc proto.InternalMessageInfo

func (m *RuleDesc) GetExpr() string {
	if m != nil {
		return m.Expr
	}
	return ""
}

func (m *RuleDesc) GetRecord() string {
	if m != nil {
		return m.Record
	}
	return ""
}

func (m *RuleDesc) GetAlert() string {
	if m != nil {
		return m.Alert
	}
	return ""
}

func (m *RuleDesc) GetFor() time.Duration {
	if m != nil {
		return m.For
	}
	return 0
}

func (m *RuleDesc) GetKeepFiringFor() time.Duration {
	if m != nil {
		return m.KeepFiringFor
	}
	return 0
}

func init() {
	proto.RegisterType((*RuleGroupDesc)(nil), "rules.RuleGroupDesc")
	proto.RegisterType((*RuleDesc)(nil), "rules.RuleDesc")
}

func init() { proto.RegisterFile("rules.proto", fileDescriptor_8e722d3e922f0937) }

var fileDescriptor_8e722d3e922f0937 = []byte{
	// 554 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x53, 0x31, 0x6f, 0xd3, 0x40,
	0x18, 0xf5, 0x61, 0xc7, 0x75, 0xce, 0x44, 0xad, 0x8e, 0x08, 0xb9, 0x05, 0x5d, 0xa2, 0x4a, 0x48,
	0x99, 0x1c, 0xa9, 0x88, 0x81, 0x01, 0xa1, 0x44, 0xa5, 0x48, 0x11, 0x12, 0xc8, 0x23, 0x42, 0xaa,
	0x6c, 0xe7, 0x6c, 0x4c, 0x1d, 0x9f, 0x39, 0x9f, 0x51, 0xb3, 0xf1, 0x13, 0x18, 0xf9, 0x09, 0xfc,
	0x94, 0x0e, 0x0c, 0x61, 0xab, 0x18, 0x02, 0x71, 0x16, 0xc4, 0xd4, 0x9f, 0x80, 0xee, 0xce, 0x86,
	0x00, 0x03, 0x65, 0x80, 0x29, 0xdf, 0x77, 0xef, 0xde, 0x7d, 0xef, 0x7b, 0x2f, 0x86, 0x36, 0x2b,
	0x53, 0x52, 0xb8, 0x39, 0xa3, 0x9c, 0xa2, 0x96, 0x6c, 0xf6, 0xba, 0x31, 0x8d, 0xa9, 0x3c, 0x19,
	0x8a, 0x4a, 0x81, 0x7b, 0x38, 0xa6, 0x34, 0x4e, 0xc9, 0x50, 0x76, 0x41, 0x19, 0x0d, 0xa7, 0x25,
	0xf3, 0x79, 0x42, 0xb3, 0x1a, 0xdf, 0xfd, 0x15, 0xf7, 0xb3, 0x79, 0x0d, 0xdd, 0x8d, 0x13, 0xfe,
	0xbc, 0x0c, 0xdc, 0x90, 0xce, 0x86, 0x21, 0x65, 0x9c, 0x9c, 0xe6, 0x8c, 0xbe, 0x20, 0x21, 0xaf,
	0xbb, 0x61, 0x7e, 0x12, 0x37, 0x40, 0x50, 0x17, 0x8a, 0xba, 0xff, 0x5e, 0x87, 0x1d, 0xaf, 0x4c,
	0xc9, 0x43, 0x46, 0xcb, 0xfc, 0x90, 0x14, 0x21, 0x42, 0xd0, 0xc8, 0xfc, 0x19, 0x71, 0x40, 0x1f,
	0x0c, 0xda, 0x9e, 0xac, 0xd1, 0x4d, 0xd8, 0x16, 0xbf, 0x45, 0xee, 0x87, 0xc4, 0xb9, 0x22, 0x81,
	0x1f, 0x07, 0xe8, 0x3e, 0xb4, 0x92, 0x8c, 0x13, 0xf6, 0xca, 0x4f, 0x1d, 0xbd, 0x0f, 0x06, 0xf6,
	0xc1, 0xae, 0xab, 0xc4, 0xba, 0x8d, 0x58, 0xf7, 0xb0, 0x5e, 0x66, 0x6c, 0x9d, 0x2d, 0x7b, 0xda,
	0xdb, 0x4f, 0x3d, 0xe0, 0x7d, 0x27, 0xa1, 0x5b, 0x50, 0x39, 0xe3, 0x18, 0x7d, 0x7d, 0x60, 0x1f,
	0x6c, 0xbb, 0xca, 0x34, 0xa1, 0x4b, 0x48, 0xf2, 0x14, 0x2a, 0x94, 0x95, 0x05, 0x61, 0x8e, 0xa9,
	0x94, 0x89, 0x1a, 0xb9, 0x70, 0x8b, 0xe6, 0xe2, 0xe1, 0xc2, 0x69, 0x4b, 0x72, 0xf7, 0xb7, 0xd1,
	0xa3, 0x6c, 0xee, 0x35, 0x97, 0x50, 0x17, 0xb6, 0xd2, 0x64, 0x96, 0x70, 0x07, 0xf6, 0xc1, 0x40,
	0xf7, 0x54, 0x83, 0x1e, 0x40, 0xfb, 0x65, 0x49, 0xd8, 0xfc, 0x71, 0x14, 0x15, 0x84, 0x3b, 0xf6,
	0x65, 0x96, 0x00, 0x72, 0x89, 0x4d, 0x1e, 0xca, 0xa0, 0x99, 0xfa, 0x01, 0x49, 0x0b, 0xe7, 0xaa,
	0xd4, 0x72, 0xcd, 0x6d, 0x4c, 0x77, 0x1f, 0x89, 0xf3, 0x27, 0x7e, 0xc2, 0xc6, 0x23, 0x61, 0xc0,
	0xc7, 0x65, 0xef, 0xaf, 0x42, 0x53, 0xfc, 0xd1, 0xd4, 0xcf, 0x39, 0x61, 0x5e, 0x3d, 0x65, 0x62,
	0x58, 0xad, 0x1d, 0x73, 0x62, 0x58, 0x5b, 0x3b, 0xd6, 0xc4, 0xb0, 0xac, 0x9d, 0xf6, 0xfe, 0x07,
	0x1d, 0x5a, 0x8d, 0x6d, 0xc2, 0x2f, 0xf1, 0x68, 0x93, 0xa4, 0xa8, 0xd1, 0x75, 0x68, 0x32, 0x12,
	0x52, 0x36, 0xad, 0x63, 0xac, 0x3b, 0xe1, 0x8b, 0x9f, 0x12, 0xc6, 0x65, 0x80, 0x6d, 0x4f, 0x35,
	0xe8, 0x0e, 0xd4, 0x23, 0xca, 0x1c, 0xe3, 0xf2, 0xa1, 0x8a, 0xfb, 0x1b, 0x3e, 0xb4, 0xfe, 0x87,
	0x0f, 0xe8, 0x14, 0xda, 0x7e, 0x96, 0x51, 0xee, 0xab, 0x3f, 0x82, 0xf9, 0x4f, 0x87, 0x6e, 0x8e,
	0x42, 0xcf, 0x60, 0xe7, 0x84, 0x90, 0xfc, 0x28, 0x61, 0x49, 0x16, 0x1f, 0x51, 0xe6, 0x74, 0xfe,
	0x64, 0xd5, 0x0d, 0xa1, 0xe0, 0xeb, 0xb2, 0xb7, 0x2d, 0x78, 0xc7, 0x91, 0x24, 0x1e, 0x47, 0x94,
	0x49, 0xf7, 0x7e, 0x7e, 0x4c, 0x26, 0xdb, 0x19, 0xdf, 0x5b, 0xac, 0xb0, 0x76, 0xbe, 0xc2, 0xda,
	0xc5, 0x0a, 0x83, 0xd7, 0x15, 0x06, 0xef, 0x2a, 0x0c, 0xce, 0x2a, 0x0c, 0x16, 0x15, 0x06, 0x9f,
	0x2b, 0x0c, 0xbe, 0x54, 0x58, 0xbb, 0xa8, 0x30, 0x78, 0xb3, 0xc6, 0xda, 0x62, 0x8d, 0xb5, 0xf3,
	0x35, 0xd6, 0x9e, 0x6e, 0xc9, 0x6f, 0x26, 0x0f, 0x02, 0x53, 0x6a, 0xb8, 0xfd, 0x2d, 0x00, 0x00,
	0xff, 0xff, 0x39, 0x21, 0xc0, 0xef, 0x8a, 0x04, 0x00, 0x00,
}

func (this *RuleGroupDesc) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RuleGroupDesc)
	if !ok {
		that2, ok := that.(RuleGroupDesc)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Name != that1.Name {
		return false
	}
	if this.Namespace != that1.Namespace {
		return false
	}
	if this.Interval != that1.Interval {
		return false
	}
	if len(this.Rules) != len(that1.Rules) {
		return false
	}
	for i := range this.Rules {
		if !this.Rules[i].Equal(that1.Rules[i]) {
			return false
		}
	}
	if this.User != that1.User {
		return false
	}
	if len(this.Options) != len(that1.Options) {
		return false
	}
	for i := range this.Options {
		if !this.Options[i].Equal(that1.Options[i]) {
			return false
		}
	}
	if this.Limit != that1.Limit {
		return false
	}
	if this.QueryOffset != nil && that1.QueryOffset != nil {
		if *this.QueryOffset != *that1.QueryOffset {
			return false
		}
	} else if this.QueryOffset != nil {
		return false
	} else if that1.QueryOffset != nil {
		return false
	}
	if len(this.Labels) != len(that1.Labels) {
		return false
	}
	for i := range this.Labels {
		if !this.Labels[i].Equal(that1.Labels[i]) {
			return false
		}
	}
	return true
}
func (this *RuleDesc) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RuleDesc)
	if !ok {
		that2, ok := that.(RuleDesc)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Expr != that1.Expr {
		return false
	}
	if this.Record != that1.Record {
		return false
	}
	if this.Alert != that1.Alert {
		return false
	}
	if this.For != that1.For {
		return false
	}
	if len(this.Labels) != len(that1.Labels) {
		return false
	}
	for i := range this.Labels {
		if !this.Labels[i].Equal(that1.Labels[i]) {
			return false
		}
	}
	if len(this.Annotations) != len(that1.Annotations) {
		return false
	}
	for i := range this.Annotations {
		if !this.Annotations[i].Equal(that1.Annotations[i]) {
			return false
		}
	}
	if this.KeepFiringFor != that1.KeepFiringFor {
		return false
	}
	return true
}
func (this *RuleGroupDesc) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 13)
	s = append(s, "&rulespb.RuleGroupDesc{")
	s = append(s, "Name: "+fmt.Sprintf("%#v", this.Name)+",\n")
	s = append(s, "Namespace: "+fmt.Sprintf("%#v", this.Namespace)+",\n")
	s = append(s, "Interval: "+fmt.Sprintf("%#v", this.Interval)+",\n")
	if this.Rules != nil {
		s = append(s, "Rules: "+fmt.Sprintf("%#v", this.Rules)+",\n")
	}
	s = append(s, "User: "+fmt.Sprintf("%#v", this.User)+",\n")
	if this.Options != nil {
		s = append(s, "Options: "+fmt.Sprintf("%#v", this.Options)+",\n")
	}
	s = append(s, "Limit: "+fmt.Sprintf("%#v", this.Limit)+",\n")
	s = append(s, "QueryOffset: "+fmt.Sprintf("%#v", this.QueryOffset)+",\n")
	s = append(s, "Labels: "+fmt.Sprintf("%#v", this.Labels)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *RuleDesc) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 11)
	s = append(s, "&rulespb.RuleDesc{")
	s = append(s, "Expr: "+fmt.Sprintf("%#v", this.Expr)+",\n")
	s = append(s, "Record: "+fmt.Sprintf("%#v", this.Record)+",\n")
	s = append(s, "Alert: "+fmt.Sprintf("%#v", this.Alert)+",\n")
	s = append(s, "For: "+fmt.Sprintf("%#v", this.For)+",\n")
	s = append(s, "Labels: "+fmt.Sprintf("%#v", this.Labels)+",\n")
	s = append(s, "Annotations: "+fmt.Sprintf("%#v", this.Annotations)+",\n")
	s = append(s, "KeepFiringFor: "+fmt.Sprintf("%#v", this.KeepFiringFor)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringRules(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *RuleGroupDesc) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RuleGroupDesc) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RuleGroupDesc) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Labels) > 0 {
		for iNdEx := len(m.Labels) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.Labels[iNdEx].Size()
				i -= size
				if _, err := m.Labels[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintRules(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x62
		}
	}
	if m.QueryOffset != nil {
		n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(*m.QueryOffset, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(*m.QueryOffset):])
		if err1 != nil {
			return 0, err1
		}
		i -= n1
		i = encodeVarintRules(dAtA, i, uint64(n1))
		i--
		dAtA[i] = 0x5a
	}
	if m.Limit != 0 {
		i = encodeVarintRules(dAtA, i, uint64(m.Limit))
		i--
		dAtA[i] = 0x50
	}
	if len(m.Options) > 0 {
		for iNdEx := len(m.Options) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Options[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRules(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarintRules(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Rules) > 0 {
		for iNdEx := len(m.Rules) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Rules[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRules(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	n2, err2 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.Interval, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.Interval):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintRules(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x1a
	if len(m.Namespace) > 0 {
		i -= len(m.Namespace)
		copy(dAtA[i:], m.Namespace)
		i = encodeVarintRules(dAtA, i, uint64(len(m.Namespace)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintRules(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RuleDesc) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RuleDesc) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RuleDesc) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n3, err3 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.KeepFiringFor, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.KeepFiringFor):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintRules(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0x6a
	if len(m.Annotations) > 0 {
		for iNdEx := len(m.Annotations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.Annotations[iNdEx].Size()
				i -= size
				if _, err := m.Annotations[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintRules(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Labels) > 0 {
		for iNdEx := len(m.Labels) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.Labels[iNdEx].Size()
				i -= size
				if _, err := m.Labels[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintRules(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	n4, err4 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.For, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.For):])
	if err4 != nil {
		return 0, err4
	}
	i -= n4
	i = encodeVarintRules(dAtA, i, uint64(n4))
	i--
	dAtA[i] = 0x22
	if len(m.Alert) > 0 {
		i -= len(m.Alert)
		copy(dAtA[i:], m.Alert)
		i = encodeVarintRules(dAtA, i, uint64(len(m.Alert)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Record) > 0 {
		i -= len(m.Record)
		copy(dAtA[i:], m.Record)
		i = encodeVarintRules(dAtA, i, uint64(len(m.Record)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Expr) > 0 {
		i -= len(m.Expr)
		copy(dAtA[i:], m.Expr)
		i = encodeVarintRules(dAtA, i, uint64(len(m.Expr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRules(dAtA []byte, offset int, v uint64) int {
	offset -= sovRules(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RuleGroupDesc) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovRules(uint64(l))
	}
	l = len(m.Namespace)
	if l > 0 {
		n += 1 + l + sovRules(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.Interval)
	n += 1 + l + sovRules(uint64(l))
	if len(m.Rules) > 0 {
		for _, e := range m.Rules {
			l = e.Size()
			n += 1 + l + sovRules(uint64(l))
		}
	}
	l = len(m.User)
	if l > 0 {
		n += 1 + l + sovRules(uint64(l))
	}
	if len(m.Options) > 0 {
		for _, e := range m.Options {
			l = e.Size()
			n += 1 + l + sovRules(uint64(l))
		}
	}
	if m.Limit != 0 {
		n += 1 + sovRules(uint64(m.Limit))
	}
	if m.QueryOffset != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdDuration(*m.QueryOffset)
		n += 1 + l + sovRules(uint64(l))
	}
	if len(m.Labels) > 0 {
		for _, e := range m.Labels {
			l = e.Size()
			n += 1 + l + sovRules(uint64(l))
		}
	}
	return n
}

func (m *RuleDesc) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Expr)
	if l > 0 {
		n += 1 + l + sovRules(uint64(l))
	}
	l = len(m.Record)
	if l > 0 {
		n += 1 + l + sovRules(uint64(l))
	}
	l = len(m.Alert)
	if l > 0 {
		n += 1 + l + sovRules(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.For)
	n += 1 + l + sovRules(uint64(l))
	if len(m.Labels) > 0 {
		for _, e := range m.Labels {
			l = e.Size()
			n += 1 + l + sovRules(uint64(l))
		}
	}
	if len(m.Annotations) > 0 {
		for _, e := range m.Annotations {
			l = e.Size()
			n += 1 + l + sovRules(uint64(l))
		}
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.KeepFiringFor)
	n += 1 + l + sovRules(uint64(l))
	return n
}

func sovRules(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRules(x uint64) (n int) {
	return sovRules(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *RuleGroupDesc) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForRules := "[]*RuleDesc{"
	for _, f := range this.Rules {
		repeatedStringForRules += strings.Replace(f.String(), "RuleDesc", "RuleDesc", 1) + ","
	}
	repeatedStringForRules += "}"
	repeatedStringForOptions := "[]*Any{"
	for _, f := range this.Options {
		repeatedStringForOptions += strings.Replace(fmt.Sprintf("%v", f), "Any", "types.Any", 1) + ","
	}
	repeatedStringForOptions += "}"
	s := strings.Join([]string{`&RuleGroupDesc{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Namespace:` + fmt.Sprintf("%v", this.Namespace) + `,`,
		`Interval:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Interval), "Duration", "duration.Duration", 1), `&`, ``, 1) + `,`,
		`Rules:` + repeatedStringForRules + `,`,
		`User:` + fmt.Sprintf("%v", this.User) + `,`,
		`Options:` + repeatedStringForOptions + `,`,
		`Limit:` + fmt.Sprintf("%v", this.Limit) + `,`,
		`QueryOffset:` + strings.Replace(fmt.Sprintf("%v", this.QueryOffset), "Duration", "duration.Duration", 1) + `,`,
		`Labels:` + fmt.Sprintf("%v", this.Labels) + `,`,
		`}`,
	}, "")
	return s
}
func (this *RuleDesc) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RuleDesc{`,
		`Expr:` + fmt.Sprintf("%v", this.Expr) + `,`,
		`Record:` + fmt.Sprintf("%v", this.Record) + `,`,
		`Alert:` + fmt.Sprintf("%v", this.Alert) + `,`,
		`For:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.For), "Duration", "duration.Duration", 1), `&`, ``, 1) + `,`,
		`Labels:` + fmt.Sprintf("%v", this.Labels) + `,`,
		`Annotations:` + fmt.Sprintf("%v", this.Annotations) + `,`,
		`KeepFiringFor:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.KeepFiringFor), "Duration", "duration.Duration", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringRules(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *RuleGroupDesc) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRules
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RuleGroupDesc: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RuleGroupDesc: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Interval", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.Interval, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rules", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rules = append(m.Rules, &RuleDesc{})
			if err := m.Rules[len(m.Rules)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Options", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Options = append(m.Options, &types.Any{})
			if err := m.Options[len(m.Options)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			m.Limit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Limit |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueryOffset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.QueryOffset == nil {
				m.QueryOffset = new(time.Duration)
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(m.QueryOffset, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Labels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Labels = append(m.Labels, github_com_cortexproject_cortex_pkg_cortexpb.LabelAdapter{})
			if err := m.Labels[len(m.Labels)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRules(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRules
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRules
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RuleDesc) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRules
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RuleDesc: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RuleDesc: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Expr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Expr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Record", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Record = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Alert", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Alert = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field For", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.For, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Labels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Labels = append(m.Labels, github_com_cortexproject_cortex_pkg_cortexpb.LabelAdapter{})
			if err := m.Labels[len(m.Labels)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Annotations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Annotations = append(m.Annotations, github_com_cortexproject_cortex_pkg_cortexpb.LabelAdapter{})
			if err := m.Annotations[len(m.Annotations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeepFiringFor", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRules
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRules
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRules
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.KeepFiringFor, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRules(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRules
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRules
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRules(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRules
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRules
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRules
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthRules
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthRules
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRules
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipRules(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthRules
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthRules = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRules   = fmt.Errorf("proto: integer overflow")
)
