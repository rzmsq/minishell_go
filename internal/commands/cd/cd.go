package cd

import (
	msErr "minishell_go/internal/miniShell_errors"
	"os"
)

type Cd struct {
	Path string
}

func (c *Cd) Run() error {
	return os.Chdir(c.Path)
}

func (c *Cd) SetArguments(arg interface{}) error {
	if arg == nil {
		return msErr.ErrInvalidArg
	}

	switch arg.(type) {
	case string:
		c.Path = arg.(string)
	default:
		return msErr.ErrInvalidArg
	}
	return nil
}
