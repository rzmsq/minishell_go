package echo

import (
	"fmt"
	"strings"

	msErr "minishell_go/internal/miniShell_errors"
)

type Echo struct {
	Str []string
}

func (e *Echo) Run() error {
	fmt.Println(strings.Join(e.Str, " "))
	return nil
}

func (e *Echo) SetArguments(arg interface{}) error {
	if arg == nil {
		return msErr.ErrInvalidArg
	}

	switch arg.(type) {
	case []string:
		e.Str = arg.([]string)
	case string:
		e.Str = []string{arg.(string)}
	default:
		return msErr.ErrInvalidArg
	}
	return nil
}
