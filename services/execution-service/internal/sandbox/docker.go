// Package sandbox runs user code inside ephemeral Docker containers.
package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

// Language-to-Docker-image mapping for the language runners.
var languageImages = map[string]string{
	"python":     "skillofide/runner-python:latest",
	"javascript": "skillofide/runner-javascript:latest",
	"java":       "skillofide/runner-java:latest",
	"cpp":        "skillofide/runner-cpp:latest",
	"go":         "golang:1.22-alpine",
}

// RunRequest contains everything needed to execute user code against one test case.
type RunRequest struct {
	Language      string
	Code          string
	Input         string
	TimeLimitMs   int32
	MemoryLimitMb int32
}

// RunResult is the output from executing one test case in the sandbox.
type RunResult struct {
	Stdout      string
	Stderr      string
	ExitCode    int64
	ExecutionMs int64
	MemoryKb    int64
	TimedOut    bool
	OOMKilled   bool
}

// DockerSandbox executes code inside Docker containers with resource limits.
type DockerSandbox struct {
	cli *client.Client
	log *zap.Logger
}

// New creates a DockerSandbox using the DOCKER_HOST environment variable (or default socket).
func New(log *zap.Logger) (*DockerSandbox, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}
	return &DockerSandbox{cli: cli, log: log}, nil
}

// Run executes code for a single test case and returns the raw output.
func (s *DockerSandbox) Run(ctx context.Context, req *RunRequest) (*RunResult, error) {
	image, ok := languageImages[req.Language]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", req.Language)
	}

	timeLimitMs := req.TimeLimitMs
	if timeLimitMs <= 0 {
		timeLimitMs = 2000
	}
	memLimitMb := req.MemoryLimitMb
	if memLimitMb <= 0 {
		memLimitMb = 256
	}

	// Build environment variables to pass code + input into the container
	envVars := []string{
		"USER_CODE=" + req.Code,
		"USER_INPUT=" + req.Input,
	}

	// Create container config
	containerCfg := &container.Config{
		Image: image,
		Env:   envVars,
	}
	if req.Language == "go" {
		containerCfg.Entrypoint = []string{"sh", "-c"}
		containerCfg.Cmd = []string{`echo "$USER_CODE" > /tmp/solution.go && echo "$USER_INPUT" | go run /tmp/solution.go`}
	}

	// Create container with strict resource limits
	resp, err := s.cli.ContainerCreate(ctx,
		containerCfg,
		&container.HostConfig{
			NetworkMode: "none", // no network access for security
			Resources: container.Resources{
				Memory:   int64(memLimitMb) * 1024 * 1024,
				NanoCPUs: 1_000_000_000, // 1 CPU core
			},
		},
		nil, nil, "",
	)
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}

	containerID := resp.ID
	defer func() {
		// Always clean up
		rmCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.cli.ContainerRemove(rmCtx, containerID, container.RemoveOptions{Force: true}) //nolint:errcheck
	}()

	start := time.Now()

	if err := s.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	// Wait with a strict time limit
	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, time.Duration(timeLimitMs)*time.Millisecond+500*time.Millisecond)
	defer timeoutCancel()

	statusCh, errCh := s.cli.ContainerWait(timeoutCtx, containerID, container.WaitConditionNotRunning)

	var exitCode int64
	var timedOut bool
	var oomKilled bool

	select {
	case waitResp := <-statusCh:
		exitCode = waitResp.StatusCode
		if waitResp.Error != nil {
			s.log.Warn("container wait error", zap.String("err", waitResp.Error.Message))
		}
	case err := <-errCh:
		if err != nil {
			timedOut = true
			s.cli.ContainerKill(context.Background(), containerID, "SIGKILL") //nolint:errcheck
		}
	case <-timeoutCtx.Done():
		timedOut = true
		s.cli.ContainerKill(context.Background(), containerID, "SIGKILL") //nolint:errcheck
	}

	executionMs := time.Since(start).Milliseconds()

	// Collect inspect info for memory stats
	inspect, inspectErr := s.cli.ContainerInspect(ctx, containerID)
	if inspectErr == nil && inspect.State != nil {
		oomKilled = inspect.State.OOMKilled
	}

	// Collect stdout/stderr
	logReader, err := s.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	var stdout, stderr string
	if err == nil {
		defer logReader.Close()
		var stdoutBuf, stderrBuf bytes.Buffer
		// Docker multiplexes stdout/stderr with 8-byte headers; use io.Copy for simplicity
		io.Copy(&stdoutBuf, logReader) //nolint:errcheck
		stdout = stdoutBuf.String()
		stderr = stderrBuf.String()
	}

	return &RunResult{
		Stdout:      stdout,
		Stderr:      stderr,
		ExitCode:    exitCode,
		ExecutionMs: executionMs,
		TimedOut:    timedOut,
		OOMKilled:   oomKilled,
	}, nil
}
