// Package judge compares sandbox output against expected output and produces verdicts.
package judge

import (
	"strings"

	"github.com/skillofide/execution-service/internal/sandbox"
	executionv1 "github.com/skillofide/proto/execution/v1"
	problemv1 "github.com/skillofide/proto/problem/v1"
)

// Status constants for test results.
const (
	StatusAccepted            = "Accepted"
	StatusWrongAnswer         = "WrongAnswer"
	StatusTimeLimitExceeded   = "TimeLimitExceeded"
	StatusMemoryLimitExceeded = "MemoryLimitExceeded"
	StatusRuntimeError        = "RuntimeError"
	StatusCompileError        = "CompileError"
)

// Judge evaluates sandbox results against expected outputs.
type Judge struct{}

// New creates a Judge.
func New() *Judge { return &Judge{} }

// EvaluateTestCase compares sandbox output to a test case and returns a TestResult.
func (j *Judge) EvaluateTestCase(tc *problemv1.TestCase, res *sandbox.RunResult) *executionv1.TestResult {
	tr := &executionv1.TestResult{
		TestCaseId:     tc.Id,
		Input:          tc.Input,
		ExpectedOutput: tc.ExpectedOutput,
		ActualOutput:   strings.TrimSpace(res.Stdout),
		ExecutionMs:    res.ExecutionMs,
	}

	switch {
	case res.TimedOut:
		tr.Status = StatusTimeLimitExceeded
		tr.Error = "Time limit exceeded"

	case res.OOMKilled:
		tr.Status = StatusMemoryLimitExceeded
		tr.Error = "Memory limit exceeded"

	case res.ExitCode != 0:
		tr.Status = StatusRuntimeError
		tr.Error = strings.TrimSpace(res.Stderr)
		if tr.Error == "" {
			tr.Error = "Non-zero exit code"
		}

	case normalize(res.Stdout) == normalize(tc.ExpectedOutput):
		tr.Status = StatusAccepted

	default:
		tr.Status = StatusWrongAnswer
	}

	return tr
}

// OverallStatus computes the aggregate status from individual test results.
// Priority: CompileError > TLE > MLE > RuntimeError > WrongAnswer > Accepted
func OverallStatus(results []*executionv1.TestResult) string {
	if len(results) == 0 {
		return StatusAccepted
	}

	priority := map[string]int{
		StatusCompileError:        10,
		StatusTimeLimitExceeded:   9,
		StatusMemoryLimitExceeded: 8,
		StatusRuntimeError:        7,
		StatusWrongAnswer:         6,
		StatusAccepted:            0,
	}

	worst := StatusAccepted
	for _, r := range results {
		if priority[r.Status] > priority[worst] {
			worst = r.Status
		}
	}
	return worst
}

// normalize trims trailing whitespace and normalizes line endings for comparison.
func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	var out []string
	for _, l := range lines {
		trimmed := strings.TrimRight(l, " \t\r")
		out = append(out, trimmed)
	}
	return strings.TrimRight(strings.Join(out, "\n"), "\n")
}
