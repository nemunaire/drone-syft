package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		// If the output spec contains a file path (FORMAT=FILE), resolve it to an absolute path.
		if idx := strings.Index(o, "="); idx >= 0 {
			filePath := o[idx+1:]
			if !filepath.IsAbs(filePath) {
				if wd, err := os.Getwd(); err == nil {
					o = o[:idx+1] + filepath.Join(wd, filePath)
				}
			}
		}
		cmdArgs = append(cmdArgs, "--output", o)
	}

	if args.SourceName != "" {
		cmdArgs = append(cmdArgs, "--source-name", args.SourceName)
	}

	if sourceVersion != "" {
		cmdArgs = append(cmdArgs, "--source-version", sourceVersion)
	}

	syftBin := "syft"
	if _, err := os.Stat("/syft"); err == nil {
		syftBin = "/syft"
	}

	cmd := exec.CommandContext(ctx, syftBin, cmdArgs...)
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
