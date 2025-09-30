package parser

import "strings"

type Command struct {
	Name string
	Args []string
}

type Pipeline Command

func Parse(argStr string) []Pipeline {
	commands := make([]Pipeline, 0)

	commandsStr := strings.Split(argStr, "|")
	for _, commandStr := range commandsStr {
		fields := strings.Split(strings.TrimSpace(commandStr), " ")
		commands = append(commands, Pipeline{Name: fields[0], Args: fields[1:]})
	}

	return commands
}
