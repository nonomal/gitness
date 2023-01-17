// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: merge.proto

package rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MergeServiceClient is the client API for MergeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MergeServiceClient interface {
	MergeBranch(ctx context.Context, in *MergeBranchRequest, opts ...grpc.CallOption) (*MergeBranchResponse, error)
}

type mergeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMergeServiceClient(cc grpc.ClientConnInterface) MergeServiceClient {
	return &mergeServiceClient{cc}
}

func (c *mergeServiceClient) MergeBranch(ctx context.Context, in *MergeBranchRequest, opts ...grpc.CallOption) (*MergeBranchResponse, error) {
	out := new(MergeBranchResponse)
	err := c.cc.Invoke(ctx, "/rpc.MergeService/MergeBranch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MergeServiceServer is the server API for MergeService service.
// All implementations must embed UnimplementedMergeServiceServer
// for forward compatibility
type MergeServiceServer interface {
	MergeBranch(context.Context, *MergeBranchRequest) (*MergeBranchResponse, error)
	mustEmbedUnimplementedMergeServiceServer()
}

// UnimplementedMergeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMergeServiceServer struct {
}

func (UnimplementedMergeServiceServer) MergeBranch(context.Context, *MergeBranchRequest) (*MergeBranchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeBranch not implemented")
}
func (UnimplementedMergeServiceServer) mustEmbedUnimplementedMergeServiceServer() {}

// UnsafeMergeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MergeServiceServer will
// result in compilation errors.
type UnsafeMergeServiceServer interface {
	mustEmbedUnimplementedMergeServiceServer()
}

func RegisterMergeServiceServer(s grpc.ServiceRegistrar, srv MergeServiceServer) {
	s.RegisterService(&MergeService_ServiceDesc, srv)
}

func _MergeService_MergeBranch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeBranchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MergeServiceServer).MergeBranch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.MergeService/MergeBranch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MergeServiceServer).MergeBranch(ctx, req.(*MergeBranchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MergeService_ServiceDesc is the grpc.ServiceDesc for MergeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MergeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.MergeService",
	HandlerType: (*MergeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MergeBranch",
			Handler:    _MergeService_MergeBranch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "merge.proto",
}