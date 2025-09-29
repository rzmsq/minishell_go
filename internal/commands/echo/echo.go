package echo

import (
	"fmt"
	"strings"
)

type Echo struct {
	str []string
}

func (e Echo) Run() error {
	fmt.Println(strings.Join(e.str, " "))
	return nil
}
