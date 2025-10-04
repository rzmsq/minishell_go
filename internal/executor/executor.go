package executor

import (
	"bytes"
	"io"
	"minishell_go/internal/commands/cd"
	"minishell_go/internal/commands/echo"
	"minishell_go/internal/commands/kill"
	"minishell_go/internal/commands/ps"
	"minishell_go/internal/commands/pwd"
	"minishell_go/internal/parser"
	"os"
	"os/exec"
)

type Executor interface {
	Run(stdout io.Writer) error
	SetArguments(interface{}) error
}

func getBuiltIn(name string) (Executor, bool) {
	var builtInCmd = map[string]Executor{
		"cd":   &cd.Cd{},
		"pwd":  &pwd.Pwd{},
		"echo": &echo.Echo{},
		"kill": &kill.Kill{},
		"ps":   &ps.Ps{},
	}
	action, ok := builtInCmd[name]
	return action, ok
}

func Execute(pipes [][][]parser.Pipeline) error {
	if len(pipes) == 0 {
		return nil
	}
	for _, andPipelines := range pipes {
		allSucceeded := true
		for _, pipeline := range andPipelines {
			err := executePipeline(pipeline)
			if err != nil {
				allSucceeded = false
				break
			}
		}
		if allSucceeded {
			return nil
		}
	}
	return nil
}

func executePipeline(pipes []parser.Pipeline) error {
	if len(pipes) == 0 {
		return nil
	}

	buffers := make([]*bytes.Buffer, len(pipes)-1)

	for i := 0; i < len(pipes); i++ {
		pipe := pipes[i]

		args := make([]string, len(pipe.Args))
		for j, arg := range pipe.Args {
			args[j] = parser.ParseEnvVars(arg)
		}

		var input io.Reader
		var output io.Writer

		if i == 0 {
			input = os.Stdin
		} else {
			input = buffers[i-1]
		}

		if i == len(pipes)-1 {
			output = os.Stdout
		} else {
			output = buffers[i]
		}

		if action, ok := getBuiltIn(pipe.Name); ok {
			if err := action.SetArguments(args); err != nil {
				return err
			}
			if err := action.Run(output); err != nil {
				return err
			}
		} else {
			cmd := exec.Command(pipe.Name, args...)
			cmd.Stdin = input
			cmd.Stdout = output
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}
