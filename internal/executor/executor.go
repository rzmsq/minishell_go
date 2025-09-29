package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var errInvalidArg = errors.New("invalid argument")

var builtInCmd = map[string]Executor{
	"cd":   cd{},
	"pwd":  pwd{},
	"echo": echo{},
	"kill": kill{},
	"ps":   ps{},
}

type Executor interface {
	run() error
}

type cd struct {
	path string
}

func (c cd) run() error {
	return os.Chdir(c.path)
}

type pwd struct {
}

func (p pwd) run() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(dir)

	return nil
}

type echo struct {
	str []string
}

func (e echo) run() error {
	fmt.Println(strings.Join(e.str, " "))
	return nil
}

type kill struct {
	arg string
}

func (k kill) run() error {
	pid, err := strconv.Atoi(k.arg)
	if err != nil {
		return errInvalidArg
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

type ps struct {
}

func (p ps) run() error {
	return nil
}

func Execute(command []string) error {
	cmdName := command[0]
	cmdArgs := command[1:]

	if _, ok := builtInCmd[cmdName]; ok {
		action := builtInCmd[cmdName]
		err := action.run()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
