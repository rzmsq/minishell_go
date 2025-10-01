package parser

import "strings"

type Command struct {
	Name string
	Args []string
}

type Pipeline Command

func Parse(argStr string) [][]Pipeline {
	pipelines := make([][]Pipeline, 0)
	pipeline := make([]Pipeline, 0)

	commandLines := strings.Split(argStr, "&&")
	for _, commandLine := range commandLines {
		pipes := strings.Split(commandLine, "|")
		for _, pipe := range pipes {
			fields := strings.Split(strings.TrimSpace(pipe), " ")
			pipeline = append(pipeline, Pipeline{Name: fields[0], Args: fields[1:]})
		}
		pipelines = append(pipelines, pipeline)
	}

	return pipelines
}
