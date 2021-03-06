// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/message.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Text                 string   `protobuf:"bytes,2,opt,name=text" json:"text,omitempty"`
	Timestamp            int64    `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_6f913b6ef9a380ad, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Message) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Message) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func init() {
	proto.RegisterType((*Message)(nil), "proto.Message")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeClient interface {
	RouteChat(ctx context.Context, opts ...grpc.CallOption) (Node_RouteChatClient, error)
}

type nodeClient struct {
	cc *grpc.ClientConn
}

func NewNodeClient(cc *grpc.ClientConn) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) RouteChat(ctx context.Context, opts ...grpc.CallOption) (Node_RouteChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Node_serviceDesc.Streams[0], "/proto.Node/RouteChat", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeRouteChatClient{stream}
	return x, nil
}

type Node_RouteChatClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type nodeRouteChatClient struct {
	grpc.ClientStream
}

func (x *nodeRouteChatClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *nodeRouteChatClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Node service

type NodeServer interface {
	RouteChat(Node_RouteChatServer) error
}

func RegisterNodeServer(s *grpc.Server, srv NodeServer) {
	s.RegisterService(&_Node_serviceDesc, srv)
}

func _Node_RouteChat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NodeServer).RouteChat(&nodeRouteChatServer{stream})
}

type Node_RouteChatServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type nodeRouteChatServer struct {
	grpc.ServerStream
}

func (x *nodeRouteChatServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *nodeRouteChatServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Node_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Node",
	HandlerType: (*NodeServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RouteChat",
			Handler:       _Node_RouteChat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/message.proto",
}

func init() { proto.RegisterFile("proto/message.proto", fileDescriptor_message_6f913b6ef9a380ad) }

var fileDescriptor_message_6f913b6ef9a380ad = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x03, 0xf3, 0x84, 0x58, 0xc1, 0x94,
	0x92, 0x3f, 0x17, 0xbb, 0x2f, 0x44, 0x5c, 0x48, 0x88, 0x8b, 0x25, 0x2f, 0x31, 0x37, 0x55, 0x82,
	0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xcc, 0x06, 0x89, 0x95, 0xa4, 0x56, 0x94, 0x48, 0x30, 0x41,
	0xc4, 0x40, 0x6c, 0x21, 0x19, 0x2e, 0xce, 0x92, 0xcc, 0xdc, 0xd4, 0xe2, 0x92, 0xc4, 0xdc, 0x02,
	0x09, 0x66, 0x05, 0x46, 0x0d, 0xe6, 0x20, 0x84, 0x80, 0x91, 0x25, 0x17, 0x8b, 0x5f, 0x7e, 0x4a,
	0xaa, 0x90, 0x21, 0x17, 0x67, 0x50, 0x7e, 0x69, 0x49, 0xaa, 0x73, 0x46, 0x62, 0x89, 0x10, 0x1f,
	0xc4, 0x52, 0x3d, 0xa8, 0x55, 0x52, 0x68, 0x7c, 0x25, 0x06, 0x0d, 0x46, 0x03, 0x46, 0x27, 0x79,
	0x2e, 0x81, 0xe4, 0xfc, 0x5c, 0xbd, 0xe2, 0x82, 0xa4, 0xc4, 0x52, 0xbd, 0xe2, 0x94, 0xe4, 0x8c,
	0xc4, 0x12, 0x27, 0xee, 0x60, 0x30, 0x1d, 0x00, 0x52, 0x9e, 0xc4, 0x06, 0xd6, 0x65, 0x0c, 0x08,
	0x00, 0x00, 0xff, 0xff, 0x33, 0x09, 0xae, 0x69, 0xd1, 0x00, 0x00, 0x00,
}
