package ps

import (
	"fmt"
	"io"
	msErr "minishell_go/internal/miniShell_errors"
	"sort"
	"text/tabwriter"

	"github.com/tklauser/ps"
)

type Ps struct {
}

func (p *Ps) Run(stdout io.Writer) error {
	processes, err := ps.Processes()
	if err != nil {
		return fmt.Errorf("%f: %v", msErr.ErrGetProcesses, err)
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].PID() < processes[j].PID()
	})

	w := tabwriter.NewWriter(stdout, 0, 0, 3, ' ', tabwriter.Debug)

	_, err = fmt.Fprintln(w, "PID\tPPID\tCOMMAND\t")
	if err != nil {
		return err
	}

	for _, proc := range processes {
		_, err = fmt.Fprintf(w, "%d\t%d\t%s\t\n", proc.PID(), proc.PPID(), proc.Command())
		if err != nil {
			return err
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (p *Ps) SetArguments(interface{}) error {
	return nil
}
