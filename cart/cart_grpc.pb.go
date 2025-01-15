// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.20.3
// source: api/cart/cart.proto

package cart

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	CartService_AddItem_FullMethodName   = "/cart.CartService/AddItem"
	CartService_GetCart_FullMethodName   = "/cart.CartService/GetCart"
	CartService_EmptyCart_FullMethodName = "/cart.CartService/EmptyCart"
)

// CartServiceClient is the client API for CartService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CartServiceClient interface {
	AddItem(ctx context.Context, in *AddItemReq, opts ...grpc.CallOption) (*AddItemResp, error)
	GetCart(ctx context.Context, in *GetCartReq, opts ...grpc.CallOption) (*GetCartResp, error)
	EmptyCart(ctx context.Context, in *EmptyCartReq, opts ...grpc.CallOption) (*EmptyCartResp, error)
}

type cartServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCartServiceClient(cc grpc.ClientConnInterface) CartServiceClient {
	return &cartServiceClient{cc}
}

func (c *cartServiceClient) AddItem(ctx context.Context, in *AddItemReq, opts ...grpc.CallOption) (*AddItemResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddItemResp)
	err := c.cc.Invoke(ctx, CartService_AddItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartServiceClient) GetCart(ctx context.Context, in *GetCartReq, opts ...grpc.CallOption) (*GetCartResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCartResp)
	err := c.cc.Invoke(ctx, CartService_GetCart_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartServiceClient) EmptyCart(ctx context.Context, in *EmptyCartReq, opts ...grpc.CallOption) (*EmptyCartResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyCartResp)
	err := c.cc.Invoke(ctx, CartService_EmptyCart_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CartServiceServer is the server API for CartService service.
// All implementations must embed UnimplementedCartServiceServer
// for forward compatibility.
type CartServiceServer interface {
	AddItem(context.Context, *AddItemReq) (*AddItemResp, error)
	GetCart(context.Context, *GetCartReq) (*GetCartResp, error)
	EmptyCart(context.Context, *EmptyCartReq) (*EmptyCartResp, error)
	mustEmbedUnimplementedCartServiceServer()
}

// UnimplementedCartServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCartServiceServer struct{}

func (UnimplementedCartServiceServer) AddItem(context.Context, *AddItemReq) (*AddItemResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedCartServiceServer) GetCart(context.Context, *GetCartReq) (*GetCartResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCart not implemented")
}
func (UnimplementedCartServiceServer) EmptyCart(context.Context, *EmptyCartReq) (*EmptyCartResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmptyCart not implemented")
}
func (UnimplementedCartServiceServer) mustEmbedUnimplementedCartServiceServer() {}
func (UnimplementedCartServiceServer) testEmbeddedByValue()                     {}

// UnsafeCartServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CartServiceServer will
// result in compilation errors.
type UnsafeCartServiceServer interface {
	mustEmbedUnimplementedCartServiceServer()
}

func RegisterCartServiceServer(s grpc.ServiceRegistrar, srv CartServiceServer) {
	// If the following call pancis, it indicates UnimplementedCartServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CartService_ServiceDesc, srv)
}

func _CartService_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServiceServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CartService_AddItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServiceServer).AddItem(ctx, req.(*AddItemReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CartService_GetCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCartReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServiceServer).GetCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CartService_GetCart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServiceServer).GetCart(ctx, req.(*GetCartReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CartService_EmptyCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyCartReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServiceServer).EmptyCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CartService_EmptyCart_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServiceServer).EmptyCart(ctx, req.(*EmptyCartReq))
	}
	return interceptor(ctx, in, info, handler)
}

// CartService_ServiceDesc is the grpc.ServiceDesc for CartService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CartService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cart.CartService",
	HandlerType: (*CartServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddItem",
			Handler:    _CartService_AddItem_Handler,
		},
		{
			MethodName: "GetCart",
			Handler:    _CartService_GetCart_Handler,
		},
		{
			MethodName: "EmptyCart",
			Handler:    _CartService_EmptyCart_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/cart/cart.proto",
}
