package parser

import "strings"

type Command struct {
	Name string
	Args []string
}

type Pipeline Command

func Parse(argStr string) [][][]Pipeline {
	allCommands := strings.Split(argStr, "||")
	orPipelines := make([][][]Pipeline, 0, len(allCommands))
	for _, command := range allCommands {
		commandLines := strings.Split(command, "&&")
		pipelines := make([][]Pipeline, 0, len(commandLines))
		for _, commandLine := range commandLines {
			pipes := strings.Split(commandLine, "|")
			pipeline := make([]Pipeline, 0, len(pipes))
			for _, pipe := range pipes {
				fields := strings.Split(strings.TrimSpace(pipe), " ")
				pipeline = append(pipeline, Pipeline{Name: fields[0], Args: fields[1:]})
			}
			pipelines = append(pipelines, pipeline)
			pipeline = nil
		}
		orPipelines = append(orPipelines, pipelines)
		pipelines = nil
	}

	return orPipelines
}
