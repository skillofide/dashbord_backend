package userv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── Message Types ────────────────────────────────────────────────────────────

type VerifyUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyUserResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type CreateOrUpdateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateOrUpdateUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ─── Server Interface ─────────────────────────────────────────────────────────

type UserServiceServer interface {
	VerifyUser(context.Context, *VerifyUserRequest) (*VerifyUserResponse, error)
	CreateOrUpdateUser(context.Context, *CreateOrUpdateUserRequest) (*CreateOrUpdateUserResponse, error)
}

type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) VerifyUser(context.Context, *VerifyUserRequest) (*VerifyUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyUser not implemented")
}

func (UnimplementedUserServiceServer) CreateOrUpdateUser(context.Context, *CreateOrUpdateUserRequest) (*CreateOrUpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrUpdateUser not implemented")
}

// ─── Client Interface & Implementation ───────────────────────────────────────

type UserServiceClient interface {
	VerifyUser(ctx context.Context, in *VerifyUserRequest, opts ...grpc.CallOption) (*VerifyUserResponse, error)
	CreateOrUpdateUser(ctx context.Context, in *CreateOrUpdateUserRequest, opts ...grpc.CallOption) (*CreateOrUpdateUserResponse, error)
}

type userServiceClient struct{ cc grpc.ClientConnInterface }

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) VerifyUser(ctx context.Context, in *VerifyUserRequest, opts ...grpc.CallOption) (*VerifyUserResponse, error) {
	out := new(VerifyUserResponse)
	if err := c.cc.Invoke(ctx, "/user.v1.UserService/VerifyUser", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateOrUpdateUser(ctx context.Context, in *CreateOrUpdateUserRequest, opts ...grpc.CallOption) (*CreateOrUpdateUserResponse, error) {
	out := new(CreateOrUpdateUserResponse)
	if err := c.cc.Invoke(ctx, "/user.v1.UserService/CreateOrUpdateUser", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Service Registration & Descriptor ───────────────────────────────────────

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "VerifyUser", Handler: _UserService_VerifyUser_Handler},
		{MethodName: "CreateOrUpdateUser", Handler: _UserService_CreateOrUpdateUser_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _UserService_VerifyUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).VerifyUser(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/VerifyUser"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(UserServiceServer).VerifyUser(ctx, req.(*VerifyUserRequest))
		})
}

func _UserService_CreateOrUpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrUpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateOrUpdateUser(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/CreateOrUpdateUser"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(UserServiceServer).CreateOrUpdateUser(ctx, req.(*CreateOrUpdateUserRequest))
		})
}
