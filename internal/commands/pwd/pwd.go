package pwd

import (
	"fmt"
	"io"
	"os"
)

type Pwd struct {
}

func (p *Pwd) Run(stdout io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, dir)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pwd) SetArguments(interface{}) error {
	return nil
}
