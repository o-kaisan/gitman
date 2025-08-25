package cli

import (
	"fmt"
	"gitman/common"
	"gitman/di"
	"log/slog"
	"os"
)

type Cli struct {
	options   *common.Options
	container di.Container
}

// New is a constructor for Handler.
func New(opts *common.Options, container di.Container) Cli {
	return Cli{
		options:   opts,
		container: container,
	}
}

// handle executes the command based on the parsed options.
func (c Cli) Handle() error {
	slog.Debug("Original arguments.", "args", os.Args)
	slog.Debug("Parsed options.", "options", c.options)

	switch {
	case c.options.Help:
		fmt.Println(common.Usage())

	case c.options.Version:
		fmt.Printf("gitman version %s\n", common.GetVersionFromGit())

	case c.options.Log:
		err := c.container.GitCommitUsecase.InteractiveCommitAction()
		if err != nil {
			return err
		}

	default:
		fmt.Println("Oops! No arguments were given.")
		fmt.Println("Use 'gitman --help' to see available commands.")
	}

	return nil
}
