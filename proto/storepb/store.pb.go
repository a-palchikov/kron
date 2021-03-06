// Code generated by protoc-gen-go.
// source: storepb/store.proto
// DO NOT EDIT!

/*
Package storepb is a generated protocol buffer package.

It is generated from these files:
	storepb/store.proto

It has these top-level messages:
	SetRequest
	SetResponse
	GetRequest
	GetResponse
*/
package storepb

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

type GetResponse_Error int32

const (
	GetResponse_None     GetResponse_Error = 0
	GetResponse_NotFound GetResponse_Error = 1
)

var GetResponse_Error_name = map[int32]string{
	0: "None",
	1: "NotFound",
}
var GetResponse_Error_value = map[string]int32{
	"None":     0,
	"NotFound": 1,
}

func (x GetResponse_Error) String() string {
	return proto.EnumName(GetResponse_Error_name, int32(x))
}

type SetRequest struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *SetRequest) Reset()         { *m = SetRequest{} }
func (m *SetRequest) String() string { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()    {}

type SetResponse struct {
}

func (m *SetResponse) Reset()         { *m = SetResponse{} }
func (m *SetResponse) String() string { return proto.CompactTextString(m) }
func (*SetResponse) ProtoMessage()    {}

type GetRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}

type GetResponse struct {
	Value []byte            `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Error GetResponse_Error `protobuf:"varint,2,opt,name=error,enum=storepb.GetResponse_Error" json:"error,omitempty"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("storepb.GetResponse_Error", GetResponse_Error_name, GetResponse_Error_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for StoreService service

type StoreServiceClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
}

type storeServiceClient struct {
	cc *grpc.ClientConn
}

func NewStoreServiceClient(cc *grpc.ClientConn) StoreServiceClient {
	return &storeServiceClient{cc}
}

func (c *storeServiceClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := grpc.Invoke(ctx, "/storepb.StoreService/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := grpc.Invoke(ctx, "/storepb.StoreService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for StoreService service

type StoreServiceServer interface {
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
}

func RegisterStoreServiceServer(s *grpc.Server, srv StoreServiceServer) {
	s.RegisterService(&_StoreService_serviceDesc, srv)
}

func _StoreService_Set_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(SetRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServiceServer).Set(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _StoreService_Get_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(GetRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServiceServer).Get(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _StoreService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "storepb.StoreService",
	HandlerType: (*StoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _StoreService_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _StoreService_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
