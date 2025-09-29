package executor

import (
	"minishell_go/internal/commands/cd"
	"minishell_go/internal/commands/echo"
	"minishell_go/internal/commands/kill"
	"minishell_go/internal/commands/ps"
	"minishell_go/internal/commands/pwd"
	"os"
	"os/exec"
)

type Executor interface {
	Run() error
	SetArguments(interface{}) error
}

func Execute(command []string) error {
	cmdName := command[0]
	cmdArgs := command[1:]

	var builtInCmd = map[string]Executor{
		"cd":   &cd.Cd{},
		"pwd":  &pwd.Pwd{},
		"echo": &echo.Echo{},
		"kill": &kill.Kill{},
		"ps":   &ps.Ps{},
	}

	if _, ok := builtInCmd[cmdName]; ok {
		action := builtInCmd[cmdName]
		err := action.SetArguments(cmdArgs)
		if err != nil {
			return err
		}
		err = action.Run()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
