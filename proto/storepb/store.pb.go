// Code generated by protoc-gen-go.
// source: storepb/store.proto
// DO NOT EDIT!

/*
Package storepb is a generated protocol buffer package.

It is generated from these files:
	storepb/store.proto

It has these top-level messages:
	SetScheduleRequest
	SetScheduleResponse
	ScheduleJobRequest
	ScheduleJobResponse
*/
package storepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import servicepb "github.com/a-palchikov/kron/proto/servicepb"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type SetScheduleRequest struct {
	Jobs []*servicepb.Job `protobuf:"bytes,1,rep,name=jobs" json:"jobs,omitempty"`
}

func (m *SetScheduleRequest) Reset()         { *m = SetScheduleRequest{} }
func (m *SetScheduleRequest) String() string { return proto.CompactTextString(m) }
func (*SetScheduleRequest) ProtoMessage()    {}

func (m *SetScheduleRequest) GetJobs() []*servicepb.Job {
	if m != nil {
		return m.Jobs
	}
	return nil
}

type SetScheduleResponse struct {
}

func (m *SetScheduleResponse) Reset()         { *m = SetScheduleResponse{} }
func (m *SetScheduleResponse) String() string { return proto.CompactTextString(m) }
func (*SetScheduleResponse) ProtoMessage()    {}

type ScheduleJobRequest struct {
	Job *servicepb.Job `protobuf:"bytes,1,opt,name=job" json:"job,omitempty"`
}

func (m *ScheduleJobRequest) Reset()         { *m = ScheduleJobRequest{} }
func (m *ScheduleJobRequest) String() string { return proto.CompactTextString(m) }
func (*ScheduleJobRequest) ProtoMessage()    {}

func (m *ScheduleJobRequest) GetJob() *servicepb.Job {
	if m != nil {
		return m.Job
	}
	return nil
}

type ScheduleJobResponse struct {
}

func (m *ScheduleJobResponse) Reset()         { *m = ScheduleJobResponse{} }
func (m *ScheduleJobResponse) String() string { return proto.CompactTextString(m) }
func (*ScheduleJobResponse) ProtoMessage()    {}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for StoreService service

type StoreServiceClient interface {
	SetSchedule(ctx context.Context, in *SetScheduleRequest, opts ...grpc.CallOption) (*SetScheduleResponse, error)
	ScheduleJob(ctx context.Context, in *ScheduleJobRequest, opts ...grpc.CallOption) (*ScheduleJobResponse, error)
}

type storeServiceClient struct {
	cc *grpc.ClientConn
}

func NewStoreServiceClient(cc *grpc.ClientConn) StoreServiceClient {
	return &storeServiceClient{cc}
}

func (c *storeServiceClient) SetSchedule(ctx context.Context, in *SetScheduleRequest, opts ...grpc.CallOption) (*SetScheduleResponse, error) {
	out := new(SetScheduleResponse)
	err := grpc.Invoke(ctx, "/storepb.StoreService/SetSchedule", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeServiceClient) ScheduleJob(ctx context.Context, in *ScheduleJobRequest, opts ...grpc.CallOption) (*ScheduleJobResponse, error) {
	out := new(ScheduleJobResponse)
	err := grpc.Invoke(ctx, "/storepb.StoreService/ScheduleJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for StoreService service

type StoreServiceServer interface {
	SetSchedule(context.Context, *SetScheduleRequest) (*SetScheduleResponse, error)
	ScheduleJob(context.Context, *ScheduleJobRequest) (*ScheduleJobResponse, error)
}

func RegisterStoreServiceServer(s *grpc.Server, srv StoreServiceServer) {
	s.RegisterService(&_StoreService_serviceDesc, srv)
}

func _StoreService_SetSchedule_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(SetScheduleRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServiceServer).SetSchedule(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _StoreService_ScheduleJob_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(ScheduleJobRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServiceServer).ScheduleJob(ctx, in)
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
			MethodName: "SetSchedule",
			Handler:    _StoreService_SetSchedule_Handler,
		},
		{
			MethodName: "ScheduleJob",
			Handler:    _StoreService_ScheduleJob_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}