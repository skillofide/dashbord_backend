// Package executionv1 contains the Execution service types and gRPC service definitions.
package executionv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─── Message Types ────────────────────────────────────────────────────────────

// RunCodeRequest is used for the "Run" button — executes against visible test cases only.
type RunCodeRequest struct {
	ProblemId string `json:"problem_id"`
	Language  string `json:"language"` // python | javascript | java | cpp | go
	Code      string `json:"code"`
	UserId    string `json:"user_id"`
}

// TestResult holds the outcome for a single test case.
type TestResult struct {
	TestCaseId     string `json:"test_case_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	ActualOutput   string `json:"actual_output"`
	// Status: Accepted | WrongAnswer | TimeLimitExceeded | MemoryLimitExceeded | RuntimeError | CompileError
	Status      string `json:"status"`
	ExecutionMs int64  `json:"execution_ms"`
	MemoryKb    int64  `json:"memory_kb"`
	Error       string `json:"error,omitempty"`
}

type RunCodeResponse struct {
	JobId        string        `json:"job_id"`
	OverallStatus string       `json:"overall_status"`
	TestResults  []*TestResult `json:"test_results"`
	CompileError string        `json:"compile_error,omitempty"`
	Runtime      int64         `json:"runtime_ms"`
	Memory       int64         `json:"memory_kb"`
}

// SubmitCodeRequest is used for the "Submit" button — executes against ALL test cases.
type SubmitCodeRequest struct {
	SubmissionId string `json:"submission_id"`
	ProblemId    string `json:"problem_id"`
	Language     string `json:"language"`
	Code         string `json:"code"`
	UserId       string `json:"user_id"`
}

type SubmitCodeResponse struct {
	JobId string `json:"job_id"`
}

// ExecutionResult is the async result published to NATS after SubmitCode.
type ExecutionResult struct {
	SubmissionId  string        `json:"submission_id"`
	UserId        string        `json:"user_id"`
	ProblemId     string        `json:"problem_id"`
	OverallStatus string        `json:"overall_status"`
	TestResults   []*TestResult `json:"test_results"`
	CompileError  string        `json:"compile_error,omitempty"`
	Runtime       int64         `json:"runtime_ms"`
	Memory        int64         `json:"memory_kb"`
}

// ─── Server Interface ─────────────────────────────────────────────────────────

type ExecutionServiceServer interface {
	RunCode(context.Context, *RunCodeRequest) (*RunCodeResponse, error)
	SubmitCode(context.Context, *SubmitCodeRequest) (*SubmitCodeResponse, error)
}

type UnimplementedExecutionServiceServer struct{}

func (UnimplementedExecutionServiceServer) RunCode(context.Context, *RunCodeRequest) (*RunCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunCode not implemented")
}
func (UnimplementedExecutionServiceServer) SubmitCode(context.Context, *SubmitCodeRequest) (*SubmitCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitCode not implemented")
}

// ─── Client Interface & Implementation ───────────────────────────────────────

type ExecutionServiceClient interface {
	RunCode(ctx context.Context, in *RunCodeRequest, opts ...grpc.CallOption) (*RunCodeResponse, error)
	SubmitCode(ctx context.Context, in *SubmitCodeRequest, opts ...grpc.CallOption) (*SubmitCodeResponse, error)
}

type executionServiceClient struct{ cc grpc.ClientConnInterface }

func NewExecutionServiceClient(cc grpc.ClientConnInterface) ExecutionServiceClient {
	return &executionServiceClient{cc}
}

func (c *executionServiceClient) RunCode(ctx context.Context, in *RunCodeRequest, opts ...grpc.CallOption) (*RunCodeResponse, error) {
	out := new(RunCodeResponse)
	if err := c.cc.Invoke(ctx, "/execution.v1.ExecutionService/RunCode", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionServiceClient) SubmitCode(ctx context.Context, in *SubmitCodeRequest, opts ...grpc.CallOption) (*SubmitCodeResponse, error) {
	out := new(SubmitCodeResponse)
	if err := c.cc.Invoke(ctx, "/execution.v1.ExecutionService/SubmitCode", in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ─── Service Registration & Descriptor ───────────────────────────────────────

func RegisterExecutionServiceServer(s grpc.ServiceRegistrar, srv ExecutionServiceServer) {
	s.RegisterService(&ExecutionService_ServiceDesc, srv)
}

var ExecutionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "execution.v1.ExecutionService",
	HandlerType: (*ExecutionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "RunCode", Handler: _ExecutionService_RunCode_Handler},
		{MethodName: "SubmitCode", Handler: _ExecutionService_SubmitCode_Handler},
	},
	Streams: []grpc.StreamDesc{},
}

func _ExecutionService_RunCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutionServiceServer).RunCode(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/execution.v1.ExecutionService/RunCode"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ExecutionServiceServer).RunCode(ctx, req.(*RunCodeRequest))
		})
}

func _ExecutionService_SubmitCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutionServiceServer).SubmitCode(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: "/execution.v1.ExecutionService/SubmitCode"},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.(ExecutionServiceServer).SubmitCode(ctx, req.(*SubmitCodeRequest))
		})
}
