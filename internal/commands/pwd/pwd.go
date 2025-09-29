package pwd

import (
	"fmt"
	"os"
)

type Pwd struct {
}

func (p *Pwd) Run() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(dir)

	return nil
}

func (p *Pwd) SetArguments(interface{}) error {
	return nil
}
