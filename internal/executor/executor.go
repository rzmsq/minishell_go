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
	for i := range buffers {
		buffers[i] = &bytes.Buffer{}
	}

	for i := 0; i < len(pipes); i++ {
		pipe := pipes[i]

		args := make([]string, len(pipe.Args))
		for j, arg := range pipe.Args {
			args[j] = parser.ParseEnvVars(arg)
		}

		var input io.Reader
		var output io.Writer

		input, closeInput, err := getInput(pipe, i, buffers)
		if err != nil {
			return err
		}
		defer closeInput()

		output, closeOutput, err := getOutput(pipe, i, len(pipes), buffers)
		if err != nil {
			return err
		}
		defer closeOutput()

		if action, ok := getBuiltIn(pipe.Name); ok {
			if err = action.SetArguments(args); err != nil {
				return err
			}
			if err = action.Run(output); err != nil {
				return err
			}
		} else {
			cmd := exec.Command(pipe.Name, args...)
			cmd.Stdin = input
			cmd.Stdout = output
			cmd.Stderr = os.Stderr

			if err = cmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

func getInput(pipe parser.Pipeline, index int, buffers []*bytes.Buffer) (io.Reader, func(), error) {
	if pipe.InputFile != "" {
		file, err := os.Open(pipe.InputFile)
		if err != nil {
			return nil, func() {}, err
		}
		return file, func() {
			err = file.Close()
			if err != nil {
				return
			}
		}, nil
	}

	if index == 0 {
		return os.Stdin, func() {}, nil
	}

	return buffers[index-1], func() {}, nil
}

func getOutput(pipe parser.Pipeline, index, totalPipes int, buffers []*bytes.Buffer) (io.Writer, func(), error) {
	if pipe.OutputFile != "" {
		file, err := os.Create(pipe.OutputFile)
		if err != nil {
			return nil, func() {}, err
		}
		return file, func() {
			err = file.Close()
			if err != nil {
				return
			}
		}, nil
	}

	if pipe.AppendFile != "" {
		file, err := os.OpenFile(pipe.AppendFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, func() {}, err
		}
		return file, func() {
			err = file.Close()
			if err != nil {
				return
			}
		}, nil
	}

	if index == totalPipes-1 {
		return os.Stdout, func() {}, nil
	}

	return buffers[index], func() {}, nil
}
