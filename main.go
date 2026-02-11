package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
	"github.com/yashikota/koeda/cmd"
)

var Version string

func main() {
	app := &cli.App{
		Name:    "koeda",
		Usage:   "GitHub repository fuzzy finder",
		Version: getVersion(),
		Action:  cmd.RootCommand.Action,
		Flags:   cmd.RootCommand.Flags,
		Commands: []*cli.Command{
			cmd.UpdateCommand,
		},
		HideHelpCommand: true, // We want the root command to be the default action
	}

	if err := app.Run(os.Args); err != nil {
		if exitErr, ok := err.(cli.ExitCoder); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getVersion() string {
	if Version != "" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "(devel)" {
			return info.Main.Version
		}
		if v, ok := getVCSBuildVersion(info); ok {
			return v
		}
	}
	return "(unset)"
}

func getVCSBuildVersion(info *debug.BuildInfo) (string, bool) {
	var (
		revision string
		dirty    string
	)
	for _, v := range info.Settings {
		switch v.Key {
		case "vcs.revision":
			revision = v.Value
		case "vcs.modified":
			dirty = " (dirty)"
		}
	}
	if revision == "" {
		return "", false
	}
	return revision + dirty, true
}
