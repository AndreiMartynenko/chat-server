// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: chat_api.proto

package chat_v1

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

// ChatAPIServicesClient is the client API for ChatAPIServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatAPIServicesClient interface {
	Create(ctx context.Context, in *CreateNewChatRequest, opts ...grpc.CallOption) (*CreateNewChatResponse, error)
	Delete(ctx context.Context, in *DeleteChatRequest, opts ...grpc.CallOption) (*DeleteChatResponse, error)
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
}

type chatAPIServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewChatAPIServicesClient(cc grpc.ClientConnInterface) ChatAPIServicesClient {
	return &chatAPIServicesClient{cc}
}

func (c *chatAPIServicesClient) Create(ctx context.Context, in *CreateNewChatRequest, opts ...grpc.CallOption) (*CreateNewChatResponse, error) {
	out := new(CreateNewChatResponse)
	err := c.cc.Invoke(ctx, "/chat_v1.ChatAPIServices/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatAPIServicesClient) Delete(ctx context.Context, in *DeleteChatRequest, opts ...grpc.CallOption) (*DeleteChatResponse, error) {
	out := new(DeleteChatResponse)
	err := c.cc.Invoke(ctx, "/chat_v1.ChatAPIServices/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatAPIServicesClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, "/chat_v1.ChatAPIServices/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatAPIServicesServer is the server API for ChatAPIServices service.
// All implementations must embed UnimplementedChatAPIServicesServer
// for forward compatibility
type ChatAPIServicesServer interface {
	Create(context.Context, *CreateNewChatRequest) (*CreateNewChatResponse, error)
	Delete(context.Context, *DeleteChatRequest) (*DeleteChatResponse, error)
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	mustEmbedUnimplementedChatAPIServicesServer()
}

// UnimplementedChatAPIServicesServer must be embedded to have forward compatible implementations.
type UnimplementedChatAPIServicesServer struct {
}

func (UnimplementedChatAPIServicesServer) Create(context.Context, *CreateNewChatRequest) (*CreateNewChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedChatAPIServicesServer) Delete(context.Context, *DeleteChatRequest) (*DeleteChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedChatAPIServicesServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedChatAPIServicesServer) mustEmbedUnimplementedChatAPIServicesServer() {}

// UnsafeChatAPIServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatAPIServicesServer will
// result in compilation errors.
type UnsafeChatAPIServicesServer interface {
	mustEmbedUnimplementedChatAPIServicesServer()
}

func RegisterChatAPIServicesServer(s grpc.ServiceRegistrar, srv ChatAPIServicesServer) {
	s.RegisterService(&ChatAPIServices_ServiceDesc, srv)
}

func _ChatAPIServices_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNewChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatAPIServicesServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat_v1.ChatAPIServices/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatAPIServicesServer).Create(ctx, req.(*CreateNewChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatAPIServices_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatAPIServicesServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat_v1.ChatAPIServices/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatAPIServicesServer).Delete(ctx, req.(*DeleteChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatAPIServices_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatAPIServicesServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat_v1.ChatAPIServices/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatAPIServicesServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatAPIServices_ServiceDesc is the grpc.ServiceDesc for ChatAPIServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatAPIServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat_v1.ChatAPIServices",
	HandlerType: (*ChatAPIServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _ChatAPIServices_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ChatAPIServices_Delete_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _ChatAPIServices_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat_api.proto",
}
