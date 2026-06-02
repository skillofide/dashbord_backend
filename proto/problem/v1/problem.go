// Package problemv1 contains the Problem service types and gRPC service definitions.
// Messages are plain Go structs (JSON-tagged); the gRPC codec/codec package
// handles encoding so protoc is not required.
package problemv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── Message Types ────────────────────────────────────────────────────────────

type ListProblemsRequest struct {
	SetId      string `json:"set_id"`
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"`
	Page       int32  `json:"page"`
	PageSize   int32  `json:"page_size"`
	UserId     string `json:"user_id"`
}

type ListProblemsResponse struct {
	Problems []*Problem `json:"problems"`
	Total    int32      `json:"total"`
	Page     int32      `json:"page"`
	PageSize int32      `json:"page_size"`
}

type GetProblemRequest struct {
	Id string `json:"id"` // UUID or slug
}

type Problem struct {
	Id           string        `json:"id"`
	Slug         string        `json:"slug"`
	Title        string        `json:"title"`
	Difficulty   string        `json:"difficulty"`
	Topic        string        `json:"topic"`
	Xp           int32         `json:"xp"`
	Statement    string        `json:"statement"`
	Constraints  []string      `json:"constraints"`
	Tags         []string      `json:"tags"`
	Examples     []*Example    `json:"examples"`
	Hints        []*Hint       `json:"hints"`
	StarterCodes *StarterCodes `json:"starter_codes"`
	SetId        string        `json:"set_id"`
	UserStatus   string        `json:"user_status"` // Solved | InProgress | Unsolved
}

type Example struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation"`
}

type Hint struct {
	Order int32  `json:"order"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type StarterCodes struct {
	Javascript string `json:"javascript"`
	Python     string `json:"python"`
	Java       string `json:"java"`
	Cpp        string `json:"cpp"`
	Go         string `json:"go"`
}

type GetTestCasesRequest struct {
	ProblemId     string `json:"problem_id"`
	IncludeHidden bool   `json:"include_hidden"`
}

type TestCase struct {
	Id             string `json:"id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsHidden       bool   `json:"is_hidden"`
	TimeLimitMs    int32  `json:"time_limit_ms"`
	MemoryLimitMb  int32  `json:"memory_limit_mb"`
	OrderIndex     int32  `json:"order_index"`
}

type GetTestCasesResponse struct {
	TestCases []*TestCase `json:"test_cases"`
}

type ListPracticeSetsRequest struct {
	UserId string `json:"user_id"` // optional — for progress %
}

type PracticeSet struct {
	Id            string  `json:"id"`
	Title         string  `json:"title"`
	Level         string  `json:"level"`
	LevelColor    string  `json:"level_color"`
	BgColor       string  `json:"bg_color"`
	TotalProblems int32   `json:"total_problems"`
	Progress      float32 `json:"progress"` // 0-100, user-specific
}

type ListPracticeSetsResponse struct {
	PracticeSets []*PracticeSet `json:"practice_sets"`
}

type GetProblemUserStatusRequest struct {
	UserId    string `json:"user_id"`
	ProblemId string `json:"problem_id"`
}

type ProblemUserStatus struct {
	UserId    string `json:"user_id"`
	ProblemId string `json:"problem_id"`
	Status    string `json:"status"` // Solved | InProgress | Unsolved
	SolvedAt  string `json:"solved_at"`
	Attempts  int32  `json:"attempts"`
}

// ─── Server Interface ─────────────────────────────────────────────────────────

type ProblemServiceServer interface {
	ListProblems(context.Context, *ListProblemsRequest) (*ListProblemsResponse, error)
	GetProblem(context.Context, *GetProblemRequest) (*Problem, error)
	GetTestCases(context.Context, *GetTestCasesRequest) (*GetTestCasesResponse, error)
	ListPracticeSets(context.Context, *ListPracticeSetsRequest) (*ListPracticeSetsResponse, error)
	GetProblemUserStatus(context.Context, *GetProblemUserStatusRequest) (*ProblemUserStatus, error)
}

// UnimplementedProblemServiceServer returns codes.Unimplemented for every method.
type UnimplementedProblemServiceServer struct{}

func (UnimplementedProblemServiceServer) ListProblems(context.Context, *ListProblemsRequest) (*ListProblemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProblems not implemented")
}
func (UnimplementedProblemServiceServer) GetProblem(context.Context, *GetProblemRequest) (*Problem, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProblem not implemented")
}
func (UnimplementedProblemServiceServer) GetTestCases(context.Context, *GetTestCasesRequest) (*GetTestCasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestCases not implemented")
}
func (UnimplementedProblemServiceServer) ListPracticeSets(context.Context, *ListPracticeSetsRequest) (*ListPracticeSetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPracticeSets not implemented")
}
func (UnimplementedProblemServiceServer) GetProblemUserStatus(context.Context, *GetProblemUserStatusRequest) (*ProblemUserStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProblemUserStatus not implemented")
}

// ─── Client Interface & Implementation ───────────────────────────────────────

type ProblemServiceClient interface {
	ListProblems(ctx context.Context, in *ListProblemsRequest, opts ...grpc.CallOption) (*ListProblemsResponse, error)
	GetProblem(ctx context.Context, in *GetProblemRequest, opts ...grpc.CallOption) (*Problem, error)
	GetTestCases(ctx context.Context, in *GetTestCasesRequest, opts ...grpc.CallOption) (*GetTestCasesResponse, error)
	ListPracticeSets(ctx context.Context, in *ListPracticeSetsRequest, opts ...grpc.CallOption) (*ListPracticeSetsResponse, error)
	GetProblemUserStatus(ctx context.Context, in *GetProblemUserStatusRequest, opts ...grpc.CallOption) (*ProblemUserStatus, error)
}

type problemServiceClient struct{ cc grpc.ClientConnInterface }

func NewProblemServiceClient(cc grpc.ClientConnInterface) ProblemServiceClient {
	return &problemServiceClient{cc}
}

func (c *problemServiceClient) ListProblems(ctx context.Context, in *ListProblemsRequest, opts ...grpc.CallOption) (*ListProblemsResponse, error) {
	out := new(ListProblemsResponse)
	if err := c.cc.Invoke(ctx, "/problem.v1.ProblemService/ListProblems", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *problemServiceClient) GetProblem(ctx context.Context, in *GetProblemRequest, opts ...grpc.CallOption) (*Problem, error) {
	out := new(Problem)
	if err := c.cc.Invoke(ctx, "/problem.v1.ProblemService/GetProblem", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *problemServiceClient) GetTestCases(ctx context.Context, in *GetTestCasesRequest, opts ...grpc.CallOption) (*GetTestCasesResponse, error) {
	out := new(GetTestCasesResponse)
	if err := c.cc.Invoke(ctx, "/problem.v1.ProblemService/GetTestCases", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *problemServiceClient) ListPracticeSets(ctx context.Context, in *ListPracticeSetsRequest, opts ...grpc.CallOption) (*ListPracticeSetsResponse, error) {
	out := new(ListPracticeSetsResponse)
	if err := c.cc.Invoke(ctx, "/problem.v1.ProblemService/ListPracticeSets", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *problemServiceClient) GetProblemUserStatus(ctx context.Context, in *GetProblemUserStatusRequest, opts ...grpc.CallOption) (*ProblemUserStatus, error) {
	out := new(ProblemUserStatus)
	if err := c.cc.Invoke(ctx, "/problem.v1.ProblemService/GetProblemUserStatus", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Service Registration & Descriptor ───────────────────────────────────────

func RegisterProblemServiceServer(s grpc.ServiceRegistrar, srv ProblemServiceServer) {
	s.RegisterService(&ProblemService_ServiceDesc, srv)
}

var ProblemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "problem.v1.ProblemService",
	HandlerType: (*ProblemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "ListProblems", Handler: _ProblemService_ListProblems_Handler},
		{MethodName: "GetProblem", Handler: _ProblemService_GetProblem_Handler},
		{MethodName: "GetTestCases", Handler: _ProblemService_GetTestCases_Handler},
		{MethodName: "ListPracticeSets", Handler: _ProblemService_ListPracticeSets_Handler},
		{MethodName: "GetProblemUserStatus", Handler: _ProblemService_GetProblemUserStatus_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _ProblemService_ListProblems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProblemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProblemServiceServer).ListProblems(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/problem.v1.ProblemService/ListProblems"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProblemServiceServer).ListProblems(ctx, req.(*ListProblemsRequest))
		})
}

func _ProblemService_GetProblem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProblemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProblemServiceServer).GetProblem(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/problem.v1.ProblemService/GetProblem"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProblemServiceServer).GetProblem(ctx, req.(*GetProblemRequest))
		})
}

func _ProblemService_GetTestCases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTestCasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProblemServiceServer).GetTestCases(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/problem.v1.ProblemService/GetTestCases"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProblemServiceServer).GetTestCases(ctx, req.(*GetTestCasesRequest))
		})
}

func _ProblemService_ListPracticeSets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPracticeSetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProblemServiceServer).ListPracticeSets(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/problem.v1.ProblemService/ListPracticeSets"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProblemServiceServer).ListPracticeSets(ctx, req.(*ListPracticeSetsRequest))
		})
}

func _ProblemService_GetProblemUserStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProblemUserStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProblemServiceServer).GetProblemUserStatus(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/problem.v1.ProblemService/GetProblemUserStatus"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ProblemServiceServer).GetProblemUserStatus(ctx, req.(*GetProblemUserStatusRequest))
		})
}
