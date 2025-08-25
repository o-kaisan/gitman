package main

import (
	"gitman/common"
	"gitman/interface/cli"
	"log/slog"
	"os"

	"gitman/di"
)

func main() {
	// Parsing options
	opts := common.ParseOptions(os.Args)

	// Setting log level (all logging must be after this line)
	common.SetupGlobalLogger(opts.Debug)

	container := di.NewContainer()

	// executing the command
	err := cli.New(opts, container).Handle()
	if err != nil {
		slog.Error("failed to execute gitman", "error", err)
		os.Exit(1)
	}

	// exit
	os.Exit(0)
}
