package echo

import (
	"fmt"
	"io"
	"strings"

	msErr "minishell_go/internal/miniShell_errors"
)

type Echo struct {
	Str []string
}

func (e *Echo) Run(stdout io.Writer) error {
	_, err := fmt.Fprintln(stdout, strings.Join(e.Str, ","))
	if err != nil {
		return err
	}
	return nil
}

func (e *Echo) SetArguments(arg interface{}) error {
	if arg == nil {
		return msErr.ErrInvalidArg
	}

	switch v := arg.(type) {
	case []string:
		e.Str = v
	case string:
		e.Str = []string{v}
	default:
		return msErr.ErrInvalidArg
	}
	return nil
}

func (e *Echo) RunWithIO(stdin io.Reader, stdout io.Writer) error {
	_, err := fmt.Fprintln(stdout, strings.Join(e.Str, " "))
	return err
}
