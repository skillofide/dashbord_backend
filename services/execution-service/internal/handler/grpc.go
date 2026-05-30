// Package handler implements the gRPC ExecutionService server.
package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/skillofide/execution-service/internal/judge"
	"github.com/skillofide/execution-service/internal/sandbox"
	"github.com/skillofide/execution-service/internal/worker"
	executionv1 "github.com/skillofide/proto/execution/v1"
	problemv1 "github.com/skillofide/proto/problem/v1"
)

// ExecutionHandler implements executionv1.ExecutionServiceServer.
type ExecutionHandler struct {
	executionv1.UnimplementedExecutionServiceServer
	sb      *sandbox.DockerSandbox
	judge   *judge.Judge
	worker  *worker.Worker
	probCli problemv1.ProblemServiceClient
	log     *zap.Logger
}

// New constructs an ExecutionHandler.
func New(sb *sandbox.DockerSandbox, j *judge.Judge, w *worker.Worker, probCli problemv1.ProblemServiceClient, log *zap.Logger) *ExecutionHandler {
	return &ExecutionHandler{sb: sb, judge: j, worker: w, probCli: probCli, log: log}
}

// RunCode executes code synchronously against visible (non-hidden) test cases.
// This maps to the "Run" button on the frontend.
func (h *ExecutionHandler) RunCode(ctx context.Context, req *executionv1.RunCodeRequest) (*executionv1.RunCodeResponse, error) {
	if err := validateRunRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Fetch visible test cases only
	tcResp, err := h.probCli.GetTestCases(ctx, &problemv1.GetTestCasesRequest{
		ProblemId:     req.ProblemId,
		IncludeHidden: false,
	})
	if err != nil {
		h.log.Error("fetch test cases for run", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "fetch test cases: %v", err)
	}

	jobId := newJobID()
	var testResults []*executionv1.TestResult
	var maxRuntime int64
	var compileError string

	runCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	for _, tc := range tcResp.TestCases {
		sbResult, err := h.sb.Run(runCtx, &sandbox.RunRequest{
			Language:      req.Language,
			Code:          req.Code,
			Input:         tc.Input,
			TimeLimitMs:   tc.TimeLimitMs,
			MemoryLimitMb: tc.MemoryLimitMb,
		})
		if err != nil {
			h.log.Error("sandbox run error", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "sandbox error: %v", err)
		}

		tr := h.judge.EvaluateTestCase(tc, sbResult)
		testResults = append(testResults, tr)

		if sbResult.ExecutionMs > maxRuntime {
			maxRuntime = sbResult.ExecutionMs
		}

		// Capture compile error from first failed case
		if tr.Status == judge.StatusCompileError && compileError == "" {
			compileError = tr.Error
		}
	}

	overallStatus := judge.OverallStatus(testResults)

	h.log.Info("run code completed",
		zap.String("job_id", jobId),
		zap.String("language", req.Language),
		zap.String("status", overallStatus),
		zap.Int("cases", len(testResults)),
	)

	return &executionv1.RunCodeResponse{
		JobId:         jobId,
		OverallStatus: overallStatus,
		TestResults:   testResults,
		CompileError:  compileError,
		Runtime:       maxRuntime,
	}, nil
}

// SubmitCode enqueues an async execution job against ALL test cases (including hidden).
// Returns immediately with a job ID. Results are published to NATS execution.result.
func (h *ExecutionHandler) SubmitCode(ctx context.Context, req *executionv1.SubmitCodeRequest) (*executionv1.SubmitCodeResponse, error) {
	if req.SubmissionId == "" || req.ProblemId == "" || req.Language == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "submission_id, problem_id, language, and code are required")
	}

	if err := h.worker.PublishRunJob(req); err != nil {
		h.log.Error("publish execution job", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "enqueue job: %v", err)
	}

	h.log.Info("submit code enqueued",
		zap.String("submission_id", req.SubmissionId),
		zap.String("language", req.Language),
	)

	return &executionv1.SubmitCodeResponse{JobId: req.SubmissionId}, nil
}

func validateRunRequest(req *executionv1.RunCodeRequest) error {
	if req.ProblemId == "" {
		return fmt.Errorf("problem_id is required")
	}
	if req.Language == "" {
		return fmt.Errorf("language is required")
	}
	if req.Code == "" {
		return fmt.Errorf("code is required")
	}
	supported := map[string]bool{"python": true, "javascript": true, "java": true, "cpp": true, "go": true}
	if !supported[strings.ToLower(req.Language)] {
		return fmt.Errorf("unsupported language: %s (supported: python, javascript, java, cpp, go)", req.Language)
	}
	return nil
}

func newJobID() string {
	return uuid.New().String()
}
