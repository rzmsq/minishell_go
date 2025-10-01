package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"minishell_go/internal/executor"
	"os"
	"os/signal"
	"syscall"

	"minishell_go/internal/parser"
)

type Shell struct {
}

func (sh *Shell) Run() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			<-sigChan
			fmt.Println()
			fmt.Print("> ")
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Bye!")
				return io.EOF
			}
			_, err2 := fmt.Fprintf(os.Stderr, "%v\n", err)
			if err2 != nil {
				return err2
			}
		}
		command := parser.Parse(input)
		err = executor.Execute(command)
		if err != nil {
			_, err2 := fmt.Fprintf(os.Stderr, "%v\n", err)
			if err2 != nil {
				return err2
			}
			continue
		}
	}
}
