// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: game.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GameService_Move_FullMethodName             = "/game.GameService/Move"
	GameService_MoveServerStream_FullMethodName = "/game.GameService/MoveServerStream"
)

// GameServiceClient is the client API for GameService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameServiceClient interface {
	// 自身の移動をサーバーに通知する
	Move(ctx context.Context, opts ...grpc.CallOption) (GameService_MoveClient, error)
	// 他のキャラクタの移動をクライアントに通知する
	MoveServerStream(ctx context.Context, in *MoveServerStreamRequest, opts ...grpc.CallOption) (GameService_MoveServerStreamClient, error)
}

type gameServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGameServiceClient(cc grpc.ClientConnInterface) GameServiceClient {
	return &gameServiceClient{cc}
}

func (c *gameServiceClient) Move(ctx context.Context, opts ...grpc.CallOption) (GameService_MoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &GameService_ServiceDesc.Streams[0], GameService_Move_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gameServiceMoveClient{stream}
	return x, nil
}

type GameService_MoveClient interface {
	Send(*MoveRequest) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type gameServiceMoveClient struct {
	grpc.ClientStream
}

func (x *gameServiceMoveClient) Send(m *MoveRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gameServiceMoveClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gameServiceClient) MoveServerStream(ctx context.Context, in *MoveServerStreamRequest, opts ...grpc.CallOption) (GameService_MoveServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &GameService_ServiceDesc.Streams[1], GameService_MoveServerStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gameServiceMoveServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GameService_MoveServerStreamClient interface {
	Recv() (*MoveServerStreamResponse, error)
	grpc.ClientStream
}

type gameServiceMoveServerStreamClient struct {
	grpc.ClientStream
}

func (x *gameServiceMoveServerStreamClient) Recv() (*MoveServerStreamResponse, error) {
	m := new(MoveServerStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GameServiceServer is the server API for GameService service.
// All implementations must embed UnimplementedGameServiceServer
// for forward compatibility
type GameServiceServer interface {
	// 自身の移動をサーバーに通知する
	Move(GameService_MoveServer) error
	// 他のキャラクタの移動をクライアントに通知する
	MoveServerStream(*MoveServerStreamRequest, GameService_MoveServerStreamServer) error
	mustEmbedUnimplementedGameServiceServer()
}

// UnimplementedGameServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGameServiceServer struct {
}

func (UnimplementedGameServiceServer) Move(GameService_MoveServer) error {
	return status.Errorf(codes.Unimplemented, "method Move not implemented")
}
func (UnimplementedGameServiceServer) MoveServerStream(*MoveServerStreamRequest, GameService_MoveServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method MoveServerStream not implemented")
}
func (UnimplementedGameServiceServer) mustEmbedUnimplementedGameServiceServer() {}

// UnsafeGameServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameServiceServer will
// result in compilation errors.
type UnsafeGameServiceServer interface {
	mustEmbedUnimplementedGameServiceServer()
}

func RegisterGameServiceServer(s grpc.ServiceRegistrar, srv GameServiceServer) {
	s.RegisterService(&GameService_ServiceDesc, srv)
}

func _GameService_Move_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GameServiceServer).Move(&gameServiceMoveServer{stream})
}

type GameService_MoveServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*MoveRequest, error)
	grpc.ServerStream
}

type gameServiceMoveServer struct {
	grpc.ServerStream
}

func (x *gameServiceMoveServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gameServiceMoveServer) Recv() (*MoveRequest, error) {
	m := new(MoveRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GameService_MoveServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MoveServerStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GameServiceServer).MoveServerStream(m, &gameServiceMoveServerStreamServer{stream})
}

type GameService_MoveServerStreamServer interface {
	Send(*MoveServerStreamResponse) error
	grpc.ServerStream
}

type gameServiceMoveServerStreamServer struct {
	grpc.ServerStream
}

func (x *gameServiceMoveServerStreamServer) Send(m *MoveServerStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// GameService_ServiceDesc is the grpc.ServiceDesc for GameService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GameService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "game.GameService",
	HandlerType: (*GameServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Move",
			Handler:       _GameService_Move_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "MoveServerStream",
			Handler:       _GameService_MoveServerStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "game.proto",
}