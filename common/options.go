package common

import (
	"fmt"
)

// 起動オプション
func Usage() string {
	branchCmd := GetEnvWithString("GITMAN_BRANCH_ALIAS", "br")
	logCmd := GetEnvWithString("GITMAN_LOG_ALIAS", "l")

	return fmt.Sprintf(`usage: gitman [options] [command]

options:
  -h, --help       show this usage
  -v, --version    display the version
  -d, --debug      enable debug mode

commands:
  branch, %s       show current branch
  log, %s          show commit log

environment variables:
  GITMAN_DEBUG                debug mode (default: "false")
  GITMAN_LOG_ALIAS            change log command alias (default: "l")
  GITMAN_LOG_DISPLAY_LIMIT    change log display limit (default: "1000")
  GITMAN_BRANCH_ALIAS         change branch command alias (default: "br")`, branchCmd, logCmd)
}

type (
	Options struct {
		Help    bool
		Version bool
		Log     bool
		Debug   bool
		Branch  bool
	}
)

func newOptions() *Options {
	return &Options{
		Help:    false,
		Version: false,
		Debug:   false,
		Log:     false,
		Branch:  false,
	}
}

type Args []string

func ParseOptions(args Args) *Options {
	opts := newOptions()
	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			opts.Help = true
		case "-v", "--version":
			opts.Version = true
		case "-d", "--debug":
			opts.Debug = true
		case "log", GetEnvWithString("GITMAN_LOG_ALIAS", "l"):
			opts.Log = true
		case "branch", GetEnvWithString("GITMAN_BRANCH_ALIAS", "br"):
			opts.Branch = true
		default:
			fmt.Printf("unrecognized option %s", arg)
			// 不明なオプションがあった場合はヘルプを表示
			opts.Help = true
		}
	}
	return opts
}
