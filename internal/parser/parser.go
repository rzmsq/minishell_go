package parser

import (
	"os"
	"regexp"
	"strings"
)

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

func ParseEnvVars(s string) string {
	re := regexp.MustCompile(`\$\{?([A-Za-z_][A-Za-z0-9_]*)\}?`)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		varName := strings.TrimPrefix(match, "$")
		varName = strings.TrimPrefix(varName, "{")
		varName = strings.TrimSuffix(varName, "}")

		if value, exists := os.LookupEnv(varName); exists {
			return value
		}
		return ""
	})
}
