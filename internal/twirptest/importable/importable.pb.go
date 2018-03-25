// Code generated by protoc-gen-go.
// source: importable.proto
// DO NOT EDIT!

/*
Package importable is a generated protocol buffer package.

Test to make sure that importing other packages doesnt break

It is generated from these files:
	importable.proto

It has these top-level messages:
	Msg
*/
package importable

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

type Msg struct {
}

func (m *Msg) Reset()                    { *m = Msg{} }
func (m *Msg) String() string            { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()               {}
func (*Msg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Msg)(nil), "twirp.internal.twirptest.importable.Msg")
}

func init() { proto.RegisterFile("importable.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 114 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0xcc, 0x2d, 0xc8,
	0x2f, 0x2a, 0x49, 0x4c, 0xca, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x52, 0x2e, 0x29,
	0xcf, 0x2c, 0x2a, 0xd0, 0xcb, 0xcc, 0x2b, 0x49, 0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0x03, 0x73, 0x4b,
	0x52, 0x8b, 0x4b, 0xf4, 0x10, 0x4a, 0x95, 0x58, 0xb9, 0x98, 0x7d, 0x8b, 0xd3, 0x8d, 0x12, 0xb9,
	0x98, 0x83, 0xcb, 0x92, 0x85, 0xa2, 0xb8, 0x58, 0x82, 0x53, 0xf3, 0x52, 0x84, 0x34, 0xf4, 0x88,
	0xd0, 0xab, 0xe7, 0x5b, 0x9c, 0x2e, 0x45, 0xb4, 0x4a, 0x27, 0x9e, 0x28, 0x2e, 0x84, 0x48, 0x12,
	0x1b, 0xd8, 0x8d, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x37, 0x1f, 0x5d, 0x0f, 0xb7, 0x00,
	0x00, 0x00,
}
