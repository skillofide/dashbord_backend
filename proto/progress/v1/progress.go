// Package progressv1 contains the Progress service types and gRPC service definitions.
package progressv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── Message Types ────────────────────────────────────────────────────────────

type GetUserProgressRequest struct {
	UserId string `json:"user_id"`
}

type SetProgress struct {
	SetId    string  `json:"set_id"`
	Title    string  `json:"title"`
	Progress float32 `json:"progress"` // 0–100
	Solved   int32   `json:"solved"`
	Total    int32   `json:"total"`
}

type UserProgress struct {
	UserId         string         `json:"user_id"`
	TotalSolved    int32          `json:"total_solved"`
	TotalAttempted int32          `json:"total_attempted"`
	EasySolved     int32          `json:"easy_solved"`
	MediumSolved   int32          `json:"medium_solved"`
	HardSolved     int32          `json:"hard_solved"`
	CurrentStreak  int32          `json:"current_streak"`
	LongestStreak  int32          `json:"longest_streak"`
	TotalXp        int32          `json:"total_xp"`
	SetProgress    []*SetProgress `json:"set_progress"`
}

type GetProblemStatusRequest struct {
	UserId    string `json:"user_id"`
	ProblemId string `json:"problem_id"`
}

type ProblemStatus struct {
	UserId    string `json:"user_id"`
	ProblemId string `json:"problem_id"`
	// Status: Solved | InProgress | Unsolved
	Status   string `json:"status"`
	SolvedAt string `json:"solved_at,omitempty"`
	Attempts int32  `json:"attempts"`
}

type UpdateProblemStatusRequest struct {
	UserId    string `json:"user_id"`
	ProblemId string `json:"problem_id"`
	SetId     string `json:"set_id"`
	// Status: Solved | InProgress
	Status    string `json:"status"`
	IsCorrect bool   `json:"is_correct"`
	RuntimeMs int64  `json:"runtime_ms"`
	MemoryKb  int64  `json:"memory_kb"`
	Language  string `json:"language"`
	XpEarned  int32  `json:"xp_earned"`
}

type UpdateProblemStatusResponse struct {
	Success    bool  `json:"success"`
	NewTotalXp int32 `json:"new_total_xp"`
}

// ─── Server Interface ─────────────────────────────────────────────────────────

type ProgressServiceServer interface {
	GetUserProgress(context.Context, *GetUserProgressRequest) (*UserProgress, error)
	GetProblemStatus(context.Context, *GetProblemStatusRequest) (*ProblemStatus, error)
	UpdateProblemStatus(context.Context, *UpdateProblemStatusRequest) (*UpdateProblemStatusResponse, error)
}

type UnimplementedProgressServiceServer struct{}

func (UnimplementedProgressServiceServer) GetUserProgress(context.Context, *GetUserProgressRequest) (*UserProgress, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProgress not implemented")
}
func (UnimplementedProgressServiceServer) GetProblemStatus(context.Context, *GetProblemStatusRequest) (*ProblemStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProblemStatus not implemented")
}
func (UnimplementedProgressServiceServer) UpdateProblemStatus(context.Context, *UpdateProblemStatusRequest) (*UpdateProblemStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProblemStatus not implemented")
}

// ─── Client Interface & Implementation ───────────────────────────────────────

type ProgressServiceClient interface {
	GetUserProgress(ctx context.Context, in *GetUserProgressRequest, opts ...grpc.CallOption) (*UserProgress, error)
	GetProblemStatus(ctx context.Context, in *GetProblemStatusRequest, opts ...grpc.CallOption) (*ProblemStatus, error)
	UpdateProblemStatus(ctx context.Context, in *UpdateProblemStatusRequest, opts ...grpc.CallOption) (*UpdateProblemStatusResponse, error)
}

type progressServiceClient struct{ cc grpc.ClientConnInterface }

func NewProgressServiceClient(cc grpc.ClientConnInterface) ProgressServiceClient {
	return &progressServiceClient{cc}
}

func (c *progressServiceClient) GetUserProgress(ctx context.Context, in *GetUserProgressRequest, opts ...grpc.CallOption) (*UserProgress, error) {
	out := new(UserProgress)
	if err := c.cc.Invoke(ctx, "/progress.v1.ProgressService/GetUserProgress", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *progressServiceClient) GetProblemStatus(ctx context.Context, in *GetProblemStatusRequest, opts ...grpc.CallOption) (*ProblemStatus, error) {
	out := new(ProblemStatus)
	if err := c.cc.Invoke(ctx, "/progress.v1.ProgressService/GetProblemStatus", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *progressServiceClient) UpdateProblemStatus(ctx context.Context, in *UpdateProblemStatusRequest, opts ...grpc.CallOption) (*UpdateProblemStatusResponse, error) {
	out := new(UpdateProblemStatusResponse)
	if err := c.cc.Invoke(ctx, "/progress.v1.ProgressService/UpdateProblemStatus", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Service Registration & Descriptor ───────────────────────────────────────

func RegisterProgressServiceServer(s grpc.ServiceRegistrar, srv ProgressServiceServer) {
	s.RegisterService(&ProgressService_ServiceDesc, srv)
}

var ProgressService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "progress.v1.ProgressService",
	HandlerType: (*ProgressServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "GetUserProgress", Handler: _ProgressService_GetUserProgress_Handler},
		{MethodName: "GetProblemStatus", Handler: _ProgressService_GetProblemStatus_Handler},
		{MethodName: "UpdateProblemStatus", Handler: _ProgressService_UpdateProblemStatus_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _ProgressService_GetUserProgress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserProgressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProgressServiceServer).GetUserProgress(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/progress.v1.ProgressService/GetUserProgress"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProgressServiceServer).GetUserProgress(ctx, req.(*GetUserProgressRequest))
		})
}

func _ProgressService_GetProblemStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProblemStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProgressServiceServer).GetProblemStatus(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/progress.v1.ProgressService/GetProblemStatus"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProgressServiceServer).GetProblemStatus(ctx, req.(*GetProblemStatusRequest))
		})
}

func _ProgressService_UpdateProblemStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProblemStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProgressServiceServer).UpdateProblemStatus(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/progress.v1.ProgressService/UpdateProblemStatus"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProgressServiceServer).UpdateProblemStatus(ctx, req.(*UpdateProblemStatusRequest))
		})
}
