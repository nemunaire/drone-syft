package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline
	Level            string   `envconfig:"PLUGIN_LOG_LEVEL"`
	SelectCatalogers []string `envconfig:"PLUGIN_SELECT_CATALOGERS"`
	Output           []string `envconfig:"PLUGIN_OUTPUT"`
	SourceName       string   `envconfig:"PLUGIN_SOURCE_NAME"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	// Derive source version: prefer tag, fall back to commit SHA.
	sourceVersion := args.Tag.Name
	if sourceVersion == "" {
		sourceVersion = args.Commit.Rev
	}

	cmdArgs := []string{"scan", "dir:."}

	for _, c := range args.SelectCatalogers {
		cmdArgs = append(cmdArgs, "--select-catalogers", c)
	}

	for _, o := range args.Output {
		cmdArgs = append(cmdArgs, "--output", o)
	}

	if args.SourceName != "" {
		cmdArgs = append(cmdArgs, "--source-name", args.SourceName)
	}

	if sourceVersion != "" {
		cmdArgs = append(cmdArgs, "--source-version", sourceVersion)
	}

	cmd := exec.CommandContext(ctx, "syft", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}
