package cursorcli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Runner struct {
	Bin       string
	Mode      string
	Workspace string
	Timeout   time.Duration
}

type RunResult struct {
	Text   string
	Raw    string
	Stderr string
}

func (r *Runner) Run(ctx context.Context, prompt, model string) (*RunResult, error) {
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args := []string{"-p", "--output-format", "text"}
	if m := strings.TrimSpace(r.Mode); m != "" {
		args = append(args, "--mode", m)
	}
	if ws := strings.TrimSpace(r.Workspace); ws != "" {
		args = append(args, "--workspace", ws)
	}
	if model = strings.TrimSpace(model); model != "" && !strings.EqualFold(model, "auto") {
		args = append(args, "--model", model)
	}
	args = append(args, prompt)

	cmd, err := r.command(runCtx, args...)
	if err != nil {
		return nil, err
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cursor agent failed: %w\nstderr: %s", err, stderr.String())
	}
	text := strings.TrimSpace(stdout.String())
	return &RunResult{Text: text, Raw: stdout.String(), Stderr: stderr.String()}, nil
}

// Stream runs agent with stream-json and yields text deltas via onDelta.
func (r *Runner) Stream(ctx context.Context, prompt, model string, onDelta func(string) error) (*RunResult, error) {
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args := []string{"-p", "--output-format", "stream-json", "--stream-partial-output"}
	if m := strings.TrimSpace(r.Mode); m != "" {
		args = append(args, "--mode", m)
	}
	if ws := strings.TrimSpace(r.Workspace); ws != "" {
		args = append(args, "--workspace", ws)
	}
	if model = strings.TrimSpace(model); model != "" && !strings.EqualFold(model, "auto") {
		args = append(args, "--model", model)
	}
	args = append(args, prompt)

	cmd, err := r.command(runCtx, args...)
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var full strings.Builder
	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		delta := extractStreamDelta(line)
		if delta == "" {
			continue
		}
		full.WriteString(delta)
		if onDelta != nil {
			if err := onDelta(delta); err != nil {
				_ = cmd.Process.Kill()
				return nil, err
			}
		}
	}
	waitErr := cmd.Wait()
	if scanErr := sc.Err(); scanErr != nil && scanErr != io.EOF {
		return nil, scanErr
	}
	if waitErr != nil {
		return nil, fmt.Errorf("cursor agent stream failed: %w\nstderr: %s", waitErr, stderr.String())
	}
	text := full.String()
	return &RunResult{Text: text, Raw: text, Stderr: stderr.String()}, nil
}

func (r *Runner) ListModels(ctx context.Context) ([]string, error) {
	runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cmd, err := r.command(runCtx, "--list-models")
	if err != nil {
		return nil, err
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []string{"auto"}, nil
	}
	var models []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Available") || strings.HasPrefix(line, "Usage") {
			continue
		}
		// Typical lines look like "model-name ..." — take first token.
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		name := fields[0]
		if strings.Contains(name, "/") || strings.Contains(name, "-") || name == "auto" {
			models = append(models, name)
		}
	}
	if len(models) == 0 {
		models = []string{"auto"}
	}
	return models, nil
}

func (r *Runner) command(ctx context.Context, args ...string) (*exec.Cmd, error) {
	bin := strings.TrimSpace(r.Bin)
	if bin == "" {
		return nil, fmt.Errorf("CURSOR_AGENT_BIN is empty")
	}
	if runtime.GOOS == "windows" && strings.HasSuffix(strings.ToLower(bin), ".cmd") {
		cmdArgs := append([]string{"/c", bin}, args...)
		return exec.CommandContext(ctx, "cmd.exe", cmdArgs...), nil
	}
	return exec.CommandContext(ctx, bin, args...), nil
}

func extractStreamDelta(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	// Common shapes from stream-json: {type,text} or nested message/delta
	for _, key := range []string{"text", "delta", "content"} {
		if v, ok := obj[key].(string); ok && v != "" {
			return v
		}
	}
	if msg, ok := obj["message"].(map[string]any); ok {
		if v, ok := msg["content"].(string); ok {
			return v
		}
	}
	return ""
}
