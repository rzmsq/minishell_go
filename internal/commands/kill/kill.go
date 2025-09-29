package kill

import (
	"os"
	"strconv"

	msErr "minishell_go/internal/miniShell_errors"
)

type Kill struct {
	arg string
}

func (k Kill) Run() error {
	pid, err := strconv.Atoi(k.arg)
	if err != nil {
		return msErr.ErrInvalidArg
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = process.Kill()
	if err != nil {
		return err
	}
	return nil
}
