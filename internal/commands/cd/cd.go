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

	switch v := arg.(type) {
	case []string:
		if len(v) == 0 {
			return msErr.ErrInvalidArg
		}
		c.Path = v[0]
	case string:
		c.Path = v
	default:
		return msErr.ErrInvalidArg
	}
	return nil
}
