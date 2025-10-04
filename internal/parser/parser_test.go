package parser

import (
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "simple command",
			input: "echo hello",
			want:  1,
		},
		{
			name:  "pipe command",
			input: "echo test | cat",
			want:  1,
		},
		{
			name:  "AND operator",
			input: "echo a && echo b",
			want:  1,
		},
		{
			name:  "OR operator",
			input: "echo a || echo b",
			want:  2,
		},
		{
			name:  "complex command",
			input: "echo a | cat && echo b || echo c",
			want:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.input)
			if len(got) != tt.want {
				t.Errorf("Parse() returned %d groups, want %d", len(got), tt.want)
			}
		})
	}
}

func TestParseEnvVars(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple env var",
			input: "$TEST_VAR",
			want:  "test_value",
		},
		{
			name:  "env var with braces",
			input: "${TEST_VAR}",
			want:  "test_value",
		},
		{
			name:  "text with env var",
			input: "hello $TEST_VAR world",
			want:  "hello test_value world",
		},
		{
			name:  "non-existent env var",
			input: "$NONEXISTENT",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseEnvVars(tt.input)
			if got != tt.want {
				t.Errorf("ParseEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePipelineWithRedirects(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Pipeline
	}{
		{
			name:  "output redirect",
			input: "echo test > file",
			want: Pipeline{
				Name:       "echo",
				Args:       []string{"test"},
				OutputFile: "file",
				InputFile:  "",
				AppendFile: "",
			},
		},
		{
			name:  "append redirect",
			input: "echo test >> file",
			want: Pipeline{
				Name:       "echo",
				Args:       []string{"test"},
				AppendFile: "file",
				InputFile:  "",
				OutputFile: "",
			},
		},
		{
			name:  "input redirect",
			input: "cat < input",
			want: Pipeline{
				Name:       "cat",
				Args:       []string{},
				InputFile:  "input",
				OutputFile: "",
				AppendFile: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePipelineWithRedirects(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePipelineWithRedirects() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
