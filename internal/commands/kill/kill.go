package kill

import (
	"io"
	"os"
	"strconv"

	msErr "minishell_go/internal/miniShell_errors"
)

type Kill struct {
	Pid int
}

func (k *Kill) Run(stdout io.Writer) error {
	process, err := os.FindProcess(k.Pid)
	if err != nil {
		return err
	}
	err = process.Kill()
	if err != nil {
		return err
	}
	return nil
}

func (k *Kill) SetArguments(arg interface{}) error {
	if arg == nil {
		return msErr.ErrInvalidArg
	}

	switch v := arg.(type) {
	case []string:
		args := v
		if len(args) == 0 {
			return msErr.ErrInvalidArg
		}

		n, err := strconv.Atoi(v[0])
		if err != nil {
			return msErr.ErrInvalidArg
		}
		k.Pid = n
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return msErr.ErrInvalidArg
		}
		k.Pid = n
	case int:
		k.Pid = v
	default:
		return msErr.ErrInvalidArg
	}
	return nil
}
