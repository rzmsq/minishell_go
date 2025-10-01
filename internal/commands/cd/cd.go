package cd

import (
	"io"
	msErr "minishell_go/internal/miniShell_errors"
	"os"
)

type Cd struct {
	Path string
}

func (c *Cd) Run(stdout io.Writer) error {
	return os.Chdir(c.Path)
}

func (c *Cd) SetArguments(arg interface{}) error {
	if arg == nil {
		return msErr.ErrInvalidArg
	}

	switch arg.(type) {
	case []string:
		args := arg.([]string)
		if len(args) == 0 {
			return msErr.ErrInvalidArg
		}
		c.Path = arg.([]string)[0]
	case string:
		c.Path = arg.(string)
	default:
		return msErr.ErrInvalidArg
	}
	return nil
}
