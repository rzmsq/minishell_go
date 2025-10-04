package parser

import (
	"os"
	"regexp"
	"strings"
)

type Command struct {
	Name       string
	Args       []string
	InputFile  string
	OutputFile string
	AppendFile string
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
				pipeline = append(pipeline, parsePipelineWithRedirects(pipe))
			}
			pipelines = append(pipelines, pipeline)
		}
		orPipelines = append(orPipelines, pipelines)
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

func parsePipelineWithRedirects(pipe string) Pipeline {
	pipe = strings.TrimSpace(pipe)
	p := Pipeline{}

	if idx := strings.Index(pipe, ">"); idx != -1 {
		parts := strings.SplitN(pipe, ">", 2)
		pipe = strings.TrimSpace(parts[0])
		p.OutputFile = strings.TrimSpace(parts[1])
	}

	if idx := strings.Index(pipe, ">>"); idx != -1 {
		parts := strings.SplitN(pipe, ">>", 2)
		pipe = strings.TrimSpace(parts[0])
		p.AppendFile = strings.TrimSpace(parts[1])
		p.OutputFile = "" // приоритет у >>
	}

	if idx := strings.Index(pipe, "<"); idx != -1 {
		parts := strings.SplitN(pipe, "<", 2)
		pipe = strings.TrimSpace(parts[0])
		p.InputFile = strings.TrimSpace(parts[1])
	}

	fields := strings.Fields(pipe)
	if len(fields) > 0 {
		p.Name = fields[0]
		p.Args = fields[1:]
	}

	return p
}
