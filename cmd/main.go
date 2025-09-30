package main

import (
	"errors"
	"fmt"
	"io"
	sh "minishell_go/internal/shell"
	"os"
)

func main() {

	shell := sh.Shell{}

	err := shell.Run()
	if err != nil {
		if errors.Is(err, io.EOF) {
			os.Exit(0)
		}
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
