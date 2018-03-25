// Code generated by protoc-gen-go.
// source: clientcompat.proto
// DO NOT EDIT!

/*
Package clientcompat is a generated protocol buffer package.

It is generated from these files:
	clientcompat.proto

It has these top-level messages:
	Empty
	Req
	Resp
	ClientCompatMessage
*/
package clientcompat

import proto "github.com/shutej/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ClientCompatMessage_CompatServiceMethod int32

const (
	ClientCompatMessage_NOOP   ClientCompatMessage_CompatServiceMethod = 0
	ClientCompatMessage_METHOD ClientCompatMessage_CompatServiceMethod = 1
)

var ClientCompatMessage_CompatServiceMethod_name = map[int32]string{
	0: "NOOP",
	1: "METHOD",
}
var ClientCompatMessage_CompatServiceMethod_value = map[string]int32{
	"NOOP":   0,
	"METHOD": 1,
}

func (x ClientCompatMessage_CompatServiceMethod) String() string {
	return proto.EnumName(ClientCompatMessage_CompatServiceMethod_name, int32(x))
}
func (ClientCompatMessage_CompatServiceMethod) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{3, 0}
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Req struct {
	V string `protobuf:"bytes,1,opt,name=v" json:"v,omitempty"`
}

func (m *Req) Reset()                    { *m = Req{} }
func (m *Req) String() string            { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()               {}
func (*Req) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Req) GetV() string {
	if m != nil {
		return m.V
	}
	return ""
}

type Resp struct {
	V int32 `protobuf:"varint,1,opt,name=v" json:"v,omitempty"`
}

func (m *Resp) Reset()                    { *m = Resp{} }
func (m *Resp) String() string            { return proto.CompactTextString(m) }
func (*Resp) ProtoMessage()               {}
func (*Resp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Resp) GetV() int32 {
	if m != nil {
		return m.V
	}
	return 0
}

type ClientCompatMessage struct {
	ServiceAddress string                                  `protobuf:"bytes,1,opt,name=service_address,json=serviceAddress" json:"service_address,omitempty"`
	Method         ClientCompatMessage_CompatServiceMethod `protobuf:"varint,2,opt,name=method,enum=twirp.clientcompat.ClientCompatMessage_CompatServiceMethod" json:"method,omitempty"`
	Request        []byte                                  `protobuf:"bytes,3,opt,name=request,proto3" json:"request,omitempty"`
}

func (m *ClientCompatMessage) Reset()                    { *m = ClientCompatMessage{} }
func (m *ClientCompatMessage) String() string            { return proto.CompactTextString(m) }
func (*ClientCompatMessage) ProtoMessage()               {}
func (*ClientCompatMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ClientCompatMessage) GetServiceAddress() string {
	if m != nil {
		return m.ServiceAddress
	}
	return ""
}

func (m *ClientCompatMessage) GetMethod() ClientCompatMessage_CompatServiceMethod {
	if m != nil {
		return m.Method
	}
	return ClientCompatMessage_NOOP
}

func (m *ClientCompatMessage) GetRequest() []byte {
	if m != nil {
		return m.Request
	}
	return nil
}

func init() {
	proto.RegisterType((*Empty)(nil), "twirp.clientcompat.Empty")
	proto.RegisterType((*Req)(nil), "twirp.clientcompat.Req")
	proto.RegisterType((*Resp)(nil), "twirp.clientcompat.Resp")
	proto.RegisterType((*ClientCompatMessage)(nil), "twirp.clientcompat.ClientCompatMessage")
	proto.RegisterEnum("twirp.clientcompat.ClientCompatMessage_CompatServiceMethod", ClientCompatMessage_CompatServiceMethod_name, ClientCompatMessage_CompatServiceMethod_value)
}

func init() { proto.RegisterFile("clientcompat.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 276 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4e, 0x83, 0x40,
	0x10, 0x86, 0x5d, 0xdb, 0x52, 0x9d, 0x20, 0x36, 0x83, 0x89, 0xd8, 0x13, 0xe1, 0x22, 0x89, 0x09,
	0x87, 0x7a, 0xec, 0xc9, 0xd6, 0x26, 0x5e, 0x28, 0x66, 0xeb, 0xc9, 0x8b, 0x41, 0x98, 0x28, 0x89,
	0xb8, 0x0b, 0xbb, 0x62, 0x7c, 0x0b, 0x9f, 0xcf, 0xa7, 0x31, 0x02, 0xc6, 0x1a, 0xd7, 0xe3, 0x7c,
	0xff, 0xcc, 0xb7, 0x93, 0x1d, 0xc0, 0xec, 0xa9, 0xa0, 0x67, 0x9d, 0x89, 0x52, 0xa6, 0x3a, 0x92,
	0xb5, 0xd0, 0x02, 0x51, 0xbf, 0x16, 0xb5, 0x8c, 0xb6, 0x93, 0x60, 0x0c, 0xa3, 0x55, 0x29, 0xf5,
	0x5b, 0xe0, 0xc2, 0x80, 0x53, 0x85, 0x36, 0xb0, 0xc6, 0x63, 0x3e, 0x0b, 0xf7, 0x39, 0x6b, 0x82,
	0x23, 0x18, 0x72, 0x52, 0xf2, 0x87, 0x8e, 0xbe, 0xe8, 0x07, 0x03, 0x77, 0xd9, 0x4a, 0x96, 0xad,
	0x24, 0x26, 0xa5, 0xd2, 0x07, 0xc2, 0x53, 0x38, 0x54, 0x54, 0x37, 0x45, 0x46, 0x77, 0x69, 0x9e,
	0xd7, 0xa4, 0x54, 0x6f, 0x72, 0x7a, 0x7c, 0xd1, 0x51, 0xdc, 0x80, 0x55, 0x92, 0x7e, 0x14, 0xb9,
	0xb7, 0xeb, 0xb3, 0xd0, 0x99, 0xcd, 0xa3, 0xbf, 0x9b, 0x45, 0x86, 0x17, 0xa2, 0xae, 0xda, 0x74,
	0xb6, 0xb8, 0x55, 0xf0, 0x5e, 0x85, 0x1e, 0x8c, 0x6b, 0xaa, 0x5e, 0x48, 0x69, 0x6f, 0xe0, 0xb3,
	0xd0, 0xe6, 0xdf, 0x65, 0x70, 0x06, 0xae, 0x61, 0x10, 0xf7, 0x60, 0xb8, 0x4e, 0x92, 0xeb, 0xc9,
	0x0e, 0x02, 0x58, 0xf1, 0xea, 0xe6, 0x2a, 0xb9, 0x9c, 0xb0, 0xd9, 0x3b, 0x83, 0x83, 0x5f, 0xdd,
	0x38, 0x07, 0xab, 0x9f, 0x38, 0x36, 0xed, 0xc9, 0xa9, 0x9a, 0x7a, 0xe6, 0x40, 0x49, 0x5c, 0x00,
	0xac, 0x85, 0x90, 0xbd, 0xe0, 0xc4, 0xd4, 0xd7, 0xfe, 0xff, 0xf4, 0xff, 0x68, 0xe1, 0xdc, 0xda,
	0xdb, 0xf4, 0xde, 0x6a, 0xcf, 0x79, 0xfe, 0x19, 0x00, 0x00, 0xff, 0xff, 0x16, 0xee, 0xa9, 0xa1,
	0xe4, 0x01, 0x00, 0x00,
}
