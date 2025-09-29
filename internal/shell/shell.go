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

func Run() error {
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
		args := parser.Parse(input)
		err = executor.Execute(args)
		if err != nil {
			return err
		}
	}
}
