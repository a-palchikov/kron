// Code generated by protoc-gen-go.
// source: servicepb/scheduler.proto
// DO NOT EDIT!

/*
Package servicepb is a generated protocol buffer package.

It is generated from these files:
	servicepb/scheduler.proto
	servicepb/feedback.proto

It has these top-level messages:
	ScheduleRequest
	ScheduleResponse
	Job
	Quota
	GetJobRequest
	GetJobResponse
	JobRef
	JobStartedRequest
	JobStartedResponse
	JobProgressStep
	JobProgressResponse
	JobStoppedRequest
	JobStoppedResponse
*/
package servicepb

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

type ScheduleRequest struct {
	Job   *Job     `protobuf:"bytes,1,opt,name=job" json:"job,omitempty"`
	Quota *Quota   `protobuf:"bytes,2,opt,name=quota" json:"quota,omitempty"`
	Tags  []string `protobuf:"bytes,3,rep,name=tags" json:"tags,omitempty"`
}

func (m *ScheduleRequest) Reset()         { *m = ScheduleRequest{} }
func (m *ScheduleRequest) String() string { return proto.CompactTextString(m) }
func (*ScheduleRequest) ProtoMessage()    {}

func (m *ScheduleRequest) GetJob() *Job {
	if m != nil {
		return m.Job
	}
	return nil
}

func (m *ScheduleRequest) GetQuota() *Quota {
	if m != nil {
		return m.Quota
	}
	return nil
}

type ScheduleResponse struct {
	Ok    bool   `protobuf:"varint,1,opt,name=ok" json:"ok,omitempty"`
	JobId uint64 `protobuf:"varint,2,opt,name=job_id" json:"job_id,omitempty"`
}

func (m *ScheduleResponse) Reset()         { *m = ScheduleResponse{} }
func (m *ScheduleResponse) String() string { return proto.CompactTextString(m) }
func (*ScheduleResponse) ProtoMessage()    {}

type Job struct {
	Name        string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	CommandLine string `protobuf:"bytes,2,opt,name=command_line" json:"command_line,omitempty"`
	When        string `protobuf:"bytes,3,opt,name=when" json:"when,omitempty"`
	Id          uint64 `protobuf:"varint,4,opt,name=id" json:"id,omitempty"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}

type Quota struct {
	CpuQuota    *Quota_CPU    `protobuf:"bytes,1,opt,name=cpu_quota" json:"cpu_quota,omitempty"`
	MemoryQuota *Quota_Memory `protobuf:"bytes,2,opt,name=memory_quota" json:"memory_quota,omitempty"`
	IoQuota     *Quota_IO     `protobuf:"bytes,3,opt,name=io_quota" json:"io_quota,omitempty"`
}

func (m *Quota) Reset()         { *m = Quota{} }
func (m *Quota) String() string { return proto.CompactTextString(m) }
func (*Quota) ProtoMessage()    {}

func (m *Quota) GetCpuQuota() *Quota_CPU {
	if m != nil {
		return m.CpuQuota
	}
	return nil
}

func (m *Quota) GetMemoryQuota() *Quota_Memory {
	if m != nil {
		return m.MemoryQuota
	}
	return nil
}

func (m *Quota) GetIoQuota() *Quota_IO {
	if m != nil {
		return m.IoQuota
	}
	return nil
}

type Quota_CPU struct {
	// Units: shares (with 1024 shares by default)
	Limit int64 `protobuf:"varint,1,opt,name=limit" json:"limit,omitempty"`
}

func (m *Quota_CPU) Reset()         { *m = Quota_CPU{} }
func (m *Quota_CPU) String() string { return proto.CompactTextString(m) }
func (*Quota_CPU) ProtoMessage()    {}

type Quota_Memory struct {
	// Units: bytes
	Limit int64 `protobuf:"varint,1,opt,name=limit" json:"limit,omitempty"`
}

func (m *Quota_Memory) Reset()         { *m = Quota_Memory{} }
func (m *Quota_Memory) String() string { return proto.CompactTextString(m) }
func (*Quota_Memory) ProtoMessage()    {}

type Quota_IO struct {
	// IO limits (weight division/throttling):
	//  1. read/write speed (blkio.throttle.read_bps_device, blkio.throttle.write_bps_device)
	//  2. (blkio.weight)
	// FIXME: Units: bytes
	Limit int64 `protobuf:"varint,1,opt,name=limit" json:"limit,omitempty"`
}

func (m *Quota_IO) Reset()         { *m = Quota_IO{} }
func (m *Quota_IO) String() string { return proto.CompactTextString(m) }
func (*Quota_IO) ProtoMessage()    {}

type GetJobRequest struct {
	JobId uint64 `protobuf:"varint,1,opt,name=job_id" json:"job_id,omitempty"`
}

func (m *GetJobRequest) Reset()         { *m = GetJobRequest{} }
func (m *GetJobRequest) String() string { return proto.CompactTextString(m) }
func (*GetJobRequest) ProtoMessage()    {}

type GetJobResponse struct {
	Job *Job `protobuf:"bytes,1,opt,name=job" json:"job,omitempty"`
}

func (m *GetJobResponse) Reset()         { *m = GetJobResponse{} }
func (m *GetJobResponse) String() string { return proto.CompactTextString(m) }
func (*GetJobResponse) ProtoMessage()    {}

func (m *GetJobResponse) GetJob() *Job {
	if m != nil {
		return m.Job
	}
	return nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for SchedulerService service

type SchedulerServiceClient interface {
	Schedule(ctx context.Context, in *ScheduleRequest, opts ...grpc.CallOption) (*ScheduleResponse, error)
	GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error)
}

type schedulerServiceClient struct {
	cc *grpc.ClientConn
}

func NewSchedulerServiceClient(cc *grpc.ClientConn) SchedulerServiceClient {
	return &schedulerServiceClient{cc}
}

func (c *schedulerServiceClient) Schedule(ctx context.Context, in *ScheduleRequest, opts ...grpc.CallOption) (*ScheduleResponse, error) {
	out := new(ScheduleResponse)
	err := grpc.Invoke(ctx, "/servicepb.SchedulerService/Schedule", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedulerServiceClient) GetJob(ctx context.Context, in *GetJobRequest, opts ...grpc.CallOption) (*GetJobResponse, error) {
	out := new(GetJobResponse)
	err := grpc.Invoke(ctx, "/servicepb.SchedulerService/GetJob", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SchedulerService service

type SchedulerServiceServer interface {
	Schedule(context.Context, *ScheduleRequest) (*ScheduleResponse, error)
	GetJob(context.Context, *GetJobRequest) (*GetJobResponse, error)
}

func RegisterSchedulerServiceServer(s *grpc.Server, srv SchedulerServiceServer) {
	s.RegisterService(&_SchedulerService_serviceDesc, srv)
}

func _SchedulerService_Schedule_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(ScheduleRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(SchedulerServiceServer).Schedule(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _SchedulerService_GetJob_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(GetJobRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(SchedulerServiceServer).GetJob(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _SchedulerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "servicepb.SchedulerService",
	HandlerType: (*SchedulerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Schedule",
			Handler:    _SchedulerService_Schedule_Handler,
		},
		{
			MethodName: "GetJob",
			Handler:    _SchedulerService_GetJob_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
