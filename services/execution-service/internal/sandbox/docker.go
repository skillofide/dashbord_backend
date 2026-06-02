// Package sandbox runs user code inside ephemeral Docker containers.
package sandbox

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
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
	ProblemId     string
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
	// Add compilation grace period for compiled languages
	langLower := strings.ToLower(req.Language)
	if langLower == "cpp" || langLower == "go" {
		timeLimitMs += 8000
	} else if langLower == "java" {
		// Java needs extra time: javac compilation (3-5s) + JVM startup (2-3s) in Alpine Docker
		timeLimitMs += 15000
	}
	memLimitMb := req.MemoryLimitMb
	if memLimitMb <= 0 {
		memLimitMb = 256
	}

	// Wrap user code with driver code to read input and print output
	wrappedCode := wrapUserCode(req.ProblemId, req.Language, req.Code)

	// Build environment variables to pass code + input into the container
	envVars := []string{
		"USER_CODE=" + wrappedCode,
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

	// Collect stdout/stderr.
	// Docker's ContainerLogs returns a multiplexed stream. Each frame has an
	// 8-byte header: [stream(1)] [pad(3)] [size(4 big-endian)], then payload.
	// We demultiplex it manually since stdcopy is not vendored.
	logReader, err := s.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	var stdout, stderr string
	if err == nil {
		defer logReader.Close()
		var stdoutBuf, stderrBuf bytes.Buffer
		hdr := make([]byte, 8)
		for {
			_, err := io.ReadFull(logReader, hdr)
			if err != nil {
				break
			}
			frameSize := binary.BigEndian.Uint32(hdr[4:8])
			if frameSize == 0 {
				continue
			}
			payload := make([]byte, frameSize)
			if _, err := io.ReadFull(logReader, payload); err != nil {
				break
			}
			switch hdr[0] {
			case 1: // stdout
				stdoutBuf.Write(payload)
			case 2: // stderr
				stderrBuf.Write(payload)
			}
		}
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

func wrapUserCode(problemId string, language string, code string) string {
	language = strings.ToLower(language)
	problemId = strings.ToLower(problemId)

	// Map database UUIDs to slugs for custom driver wrappers
	switch problemId {
	case "54574a34-9a68-4e65-ab9a-af05db4c0001":
		problemId = "op1"
	case "54574a34-9a68-4e65-ab9a-af05db4c0002":
		problemId = "arr5"
	case "54574a34-9a68-4e65-ab9a-af05db4c0003":
		problemId = "cond1"
	case "54574a34-9a68-4e65-ab9a-af05db4c0004":
		problemId = "loop1"
	case "54574a34-9a68-4e65-ab9a-af05db4c0005":
		problemId = "str2"
	}

	switch language {
	case "javascript":
		var driver string
		switch problemId {
		case "op1":
			driver = `
const fs = require('fs');
const input = fs.readFileSync('/dev/stdin', 'utf-8').trim().split('\n');
const a = parseInt(input[0]);
const b = parseInt(input[1]);
console.log(JSON.stringify(arithmeticOperations(a, b)));
`
		case "arr5":
			driver = `
const fs = require('fs');
const input = fs.readFileSync('/dev/stdin', 'utf-8').trim().split('\n');
const nums = JSON.parse(input[0]);
const target = parseInt(input[1]);
console.log(JSON.stringify(twoSum(nums, target)));
`
		case "cond1":
			driver = `
const fs = require('fs');
const n = parseInt(fs.readFileSync('/dev/stdin', 'utf-8').trim());
console.log(JSON.stringify(checkEvenOdd(n)));
`
		case "loop1":
			driver = `
const fs = require('fs');
const n = parseInt(fs.readFileSync('/dev/stdin', 'utf-8').trim());
console.log(sumOfN(n));
`
		case "str2":
			driver = `
const fs = require('fs');
const s = JSON.parse(fs.readFileSync('/dev/stdin', 'utf-8').trim());
console.log(isPalindrome(s));
`
		default:
			driver = `
const fs = require('fs');
const input = fs.readFileSync('/dev/stdin', 'utf-8').trim();
try {
    const parsed = JSON.parse(input);
    if (typeof solveChallenge === 'function') {
        console.log(JSON.stringify(solveChallenge(parsed)));
    }
} catch (e) {
    if (typeof solveChallenge === 'function') {
        console.log(solveChallenge(input));
    }
}
`
		}
		return code + "\n" + driver

	case "python":
		var driver string
		switch problemId {
		case "op1":
			driver = `
import sys
input_data = sys.stdin.read().splitlines()
a = int(input_data[0])
b = int(input_data[1])
print(arithmeticOperations(a, b))
`
		case "arr5":
			driver = `
import sys, json
input_data = sys.stdin.read().splitlines()
nums = json.loads(input_data[0])
target = int(input_data[1])
print(json.dumps(twoSum(nums, target)))
`
		case "cond1":
			driver = `
import sys, json
n = int(sys.stdin.read().trim())
print(json.dumps(checkEvenOdd(n)))
`
		case "loop1":
			driver = `
import sys
n = int(sys.stdin.read().trim())
print(sumOfN(n))
`
		case "str2":
			driver = `
import sys, json
s = json.loads(sys.stdin.read().trim())
print(str(isPalindrome(s)).lower())
`
		default:
			driver = `
import sys, json
input_data = sys.stdin.read().trim()
try:
    parsed = json.loads(input_data)
    if 'solveChallenge' in globals():
        print(json.dumps(solveChallenge(parsed)))
except:
    if 'solveChallenge' in globals():
        print(solveChallenge(input_data))
`
		}
		return code + "\n" + driver

	case "java":
		var javaMain string
		switch problemId {
		case "op1":
			javaMain = `
    public static void main(String[] args) throws Exception {
        java.io.BufferedReader br = new java.io.BufferedReader(new java.io.InputStreamReader(System.in));
        String line1 = br.readLine();
        String line2 = br.readLine();
        if (line1 == null || line2 == null) return;
        int a = Integer.parseInt(line1.trim());
        int b = Integer.parseInt(line2.trim());
        Solution sol = new Solution();
        int[] res = sol.arithmeticOperations(a, b);
        System.out.println(java.util.Arrays.toString(res).replace(" ", ""));
    }
`
		case "arr5":
			javaMain = `
    public static void main(String[] args) throws Exception {
        java.io.BufferedReader br = new java.io.BufferedReader(new java.io.InputStreamReader(System.in));
        String arrayStr = br.readLine();
        String targetStr = br.readLine();
        if (arrayStr == null || targetStr == null) return;
        arrayStr = arrayStr.trim();
        int target = Integer.parseInt(targetStr.trim());
        arrayStr = arrayStr.substring(1, arrayStr.length() - 1);
        String[] tokens = arrayStr.split(",");
        int[] nums = new int[tokens.length];
        for (int i = 0; i < tokens.length; i++) {
            nums[i] = Integer.parseInt(tokens[i].trim());
        }
        Solution sol = new Solution();
        int[] res = sol.twoSum(nums, target);
        System.out.println(java.util.Arrays.toString(res).replace(" ", ""));
    }
`
		case "cond1":
			javaMain = `
    public static void main(String[] args) throws Exception {
        java.io.BufferedReader br = new java.io.BufferedReader(new java.io.InputStreamReader(System.in));
        String line = br.readLine();
        if (line == null) return;
        int n = Integer.parseInt(line.trim());
        Solution sol = new Solution();
        String res = sol.checkEvenOdd(n);
        System.out.println("\"" + res + "\"");
    }
`
		case "loop1":
			javaMain = `
    public static void main(String[] args) throws Exception {
        java.io.BufferedReader br = new java.io.BufferedReader(new java.io.InputStreamReader(System.in));
        String line = br.readLine();
        if (line == null) return;
        int n = Integer.parseInt(line.trim());
        Solution sol = new Solution();
        System.out.println(sol.sumOfN(n));
    }
`
		case "str2":
			javaMain = `
    public static void main(String[] args) throws Exception {
        java.io.BufferedReader br = new java.io.BufferedReader(new java.io.InputStreamReader(System.in));
        String s = br.readLine();
        if (s == null) return;
        s = s.trim();
        if (s.startsWith("\"") && s.endsWith("\"")) {
            s = s.substring(1, s.length() - 1);
        }
        Solution sol = new Solution();
        System.out.println(sol.isPalindrome(s));
    }
`
		default:
			javaMain = `
    public static void main(String[] args) throws Exception {
        java.io.BufferedReader br = new java.io.BufferedReader(new java.io.InputStreamReader(System.in));
        String input = br.readLine();
        if (input == null) return;
        Solution sol = new Solution();
        System.out.println(sol.solveChallenge(input.trim()));
    }
`
		}
		// Inject inside public class Solution
		lastBrace := strings.LastIndex(code, "}")
		if lastBrace != -1 {
			return code[:lastBrace] + "\n" + javaMain + "\n}"
		}
		return code

	case "cpp":
		var driver string
		switch problemId {
		case "op1":
			driver = `
#include <iostream>
int main() {
    int a, b;
    if (std::cin >> a >> b) {
        Solution sol;
        std::vector<int> res = sol.arithmeticOperations(a, b);
        std::cout << "[";
        for (size_t i = 0; i < res.size(); ++i) {
            std::cout << res[i] << (i == res.size() - 1 ? "" : ",");
        }
        std::cout << "]" << std::endl;
    }
    return 0;
}
`
		case "arr5":
			driver = `
#include <iostream>
#include <string>
#include <sstream>
int main() {
    std::string arrayStr;
    int target;
    if (std::getline(std::cin, arrayStr) && std::cin >> target) {
        if (arrayStr.front() == '[') arrayStr = arrayStr.substr(1);
        if (arrayStr.back() == ']') arrayStr.pop_back();
        std::vector<int> nums;
        std::stringstream ss(arrayStr);
        std::string item;
        while (std::getline(ss, item, ',')) {
            nums.push_back(std::stoi(item));
        }
        Solution sol;
        std::vector<int> res = sol.twoSum(nums, target);
        std::cout << "[";
        for (size_t i = 0; i < res.size(); ++i) {
            std::cout << res[i] << (i == res.size() - 1 ? "" : ",");
        }
        std::cout << "]" << std::endl;
    }
    return 0;
}
`
		case "cond1":
			driver = `
#include <iostream>
int main() {
    int n;
    if (std::cin >> n) {
        Solution sol;
        std::cout << "\"" << sol.checkEvenOdd(n) << "\"" << std::endl;
    }
    return 0;
}
`
		case "loop1":
			driver = `
#include <iostream>
int main() {
    int n;
    if (std::cin >> n) {
        Solution sol;
        std::cout << sol.sumOfN(n) << std::endl;
    }
    return 0;
}
`
		case "str2":
			driver = `
#include <iostream>
#include <string>
int main() {
    std::string s;
    if (std::getline(std::cin, s)) {
        if (s.front() == '"' && s.back() == '"') {
            s = s.substr(1, s.length() - 2);
        }
        Solution sol;
        std::cout << (sol.isPalindrome(s) ? "true" : "false") << std::endl;
    }
    return 0;
}
`
		default:
			driver = `
#include <iostream>
#include <string>
int main() {
    std::string s;
    if (std::getline(std::cin, s)) {
        Solution sol;
        sol.solveChallenge();
    }
    return 0;
}
`
		}
		return code + "\n" + driver
	}

	return code
}
