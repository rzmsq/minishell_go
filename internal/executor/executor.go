package executor

import (
	"errors"
	"fmt"
	"minishell_go/internal/commands/cd"
	"minishell_go/internal/commands/echo"
	"minishell_go/internal/commands/kill"
	"minishell_go/internal/commands/ps"
	"minishell_go/internal/commands/pwd"
	msErr "minishell_go/internal/miniShell_errors"
	"minishell_go/internal/parser"
	"os"
	"os/exec"
)

type Executor interface {
	Run() error
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

func Execute(pipes []parser.Pipeline) error {
	if len(pipes) == 0 {
		return nil
	}

	if len(pipes) == 1 {
		cmdName := pipes[0].Name
		cmdArgs := pipes[0].Args

		if action, ok := getBuiltIn(cmdName); ok {
			if err := action.SetArguments(cmdArgs); err != nil {
				if errors.Is(err, msErr.ErrInvalidArg) {
					if _, werr := fmt.Fprint(os.Stderr, msErr.ErrInvalidArg, "\n"); werr != nil {
						return werr
					}
					return nil
				}
				return err
			}
			return action.Run()
		}

		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}
	return nil
}
