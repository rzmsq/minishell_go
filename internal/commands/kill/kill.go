package kill

import (
	"os"
	"strconv"

	msErr "minishell_go/internal/miniShell_errors"
)

type Kill struct {
	Pid int
}

func (k *Kill) Run() error {
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

	switch arg.(type) {
	case []string:
		args := arg.([]string)
		if len(args) == 0 {
			return msErr.ErrInvalidArg
		}

		n, err := strconv.Atoi(arg.([]string)[0])
		if err != nil {
			return msErr.ErrInvalidArg
		}
		k.Pid = n
	case string:
		n, err := strconv.Atoi(arg.(string))
		if err != nil {
			return msErr.ErrInvalidArg
		}
		k.Pid = n
	case int:
		n, ok := arg.(int)
		if !ok {
			return msErr.ErrInvalidArg
		}
		k.Pid = n
	default:
		return msErr.ErrInvalidArg
	}
	return nil
}
