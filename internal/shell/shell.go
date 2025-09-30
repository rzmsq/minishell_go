package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"minishell_go/internal/executor"
	"os"

	"minishell_go/internal/parser"
)

type Shell struct {
}

func (sh *Shell) Run() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Bye!")
				return io.EOF
			}
			return err
		}
		command := parser.Parse(input)
		err = executor.Execute(command)
		if err != nil {
			return err
		}
	}
}
