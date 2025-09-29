package miniShell_errors

import "errors"

var (
	ErrInvalidArg   = errors.New("invalid argument")
	ErrGetProcesses = errors.New("error getting processes")
)
