package cd

import "os"

type Cd struct {
	path string
}

func (c Cd) Run() error {
	return os.Chdir(c.path)
}
