package executor

import (
	"bytes"
	"io"
	"minishell_go/internal/parser"
	"os"
	"strings"
	"testing"
)

func TestExecuteBuiltinCommands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, output string)
	}{
		{
			name:    "pwd command",
			input:   "pwd",
			wantErr: false,
			validate: func(t *testing.T, output string) {
				wd, _ := os.Getwd()
				if !strings.Contains(output, wd) {
					t.Errorf("expected pwd output to contain %s, got %s", wd, output)
				}
			},
		},
		{
			name:    "echo command",
			input:   "echo hello world",
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "hello,world") {
					t.Errorf("expected 'hello,world', got %s", output)
				}
			},
		},
		{
			name:    "cd command",
			input:   "cd /tmp",
			wantErr: false,
			validate: func(t *testing.T, output string) {
				wd, _ := os.Getwd()
				if wd != "/tmp" {
					t.Errorf("expected to be in /tmp, got %s", wd)
				}
				os.Chdir(t.TempDir())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipes := parser.Parse(tt.input)
			err := Execute(pipes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutePipeline(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple pipe echo to cat",
			input:    "echo test | cat",
			expected: "test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipes := parser.Parse(tt.input)

			// Перехватываем stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := Execute(pipes)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !strings.Contains(buf.String(), tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestExecuteWithRedirection(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		input    string
		validate func(t *testing.T)
		wantErr  bool
	}{
		{
			name:  "output redirection",
			input: "echo test > " + tmpDir + "/output.txt",
			validate: func(t *testing.T) {
				content, err := os.ReadFile(tmpDir + "/output.txt")
				if err != nil {
					t.Fatalf("failed to read output file: %v", err)
				}
				if !strings.Contains(string(content), "test") {
					t.Errorf("expected 'test' in file, got %s", content)
				}
			},
			wantErr: false,
		},
		{
			name:  "append redirection",
			input: "echo world >> " + tmpDir + "/append.txt",
			validate: func(t *testing.T) {
				content, err := os.ReadFile(tmpDir + "/append.txt")
				if err != nil {
					t.Fatalf("failed to read append file: %v", err)
				}
				if !strings.Contains(string(content), "world") {
					t.Errorf("expected 'world' in file, got %s", content)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipes := parser.Parse(tt.input)
			err := Execute(pipes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.validate != nil {
				tt.validate(t)
			}
		})
	}
}

func TestExecuteLogicalOperators(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "AND operator - both succeed",
			input:   "echo test && echo success",
			wantErr: false,
		},
		{
			name:    "OR operator - first fails",
			input:   "false || echo fallback",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipes := parser.Parse(tt.input)
			err := Execute(pipes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBuiltIn(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{"cd builtin", "cd", true},
		{"pwd builtin", "pwd", true},
		{"echo builtin", "echo", true},
		{"kill builtin", "kill", true},
		{"ps builtin", "ps", true},
		{"not a builtin", "ls", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := getBuiltIn(tt.command)
			if got != tt.want {
				t.Errorf("getBuiltIn(%s) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}
