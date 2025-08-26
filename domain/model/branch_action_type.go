package model

import (
	"fmt"
	"log/slog"
	"strings"
)

type BranchActionTypeMap struct {
	Switch            ActionType
	GetLastCommitId   ActionType
	Diff              ActionType
	RebaseInteractive ActionType
	Rebase            ActionType
	Merge             ActionType
	Delete            ActionType
	Unknown           ActionType
}

var BranchActionTypes = BranchActionTypeMap{
	Switch: ActionType{
		Name:    "switch",
		Command: "git",
		Options: []string{"switch"},
		Help:    "Switch branch to selected branch",
	},
	GetLastCommitId: ActionType{
		Name:    "get last commit",
		Command: "echo",
		Options: nil,
		Help:    "print branch last commit id",
	},
	Diff: ActionType{
		Name:    "diff",
		Command: "git",
		Options: []string{"diff"},
		Help:    "Show changes between current branch and selected branch",
	},
	RebaseInteractive: ActionType{
		Name:    "rebase interactive",
		Command: "git",
		Options: []string{"rebase", "-i"},
		Help:    "Interactive rebase to selected branch",
	},
	Rebase: ActionType{
		Name:    "rebase",
		Command: "git",
		Options: []string{"rebase"},
		Help:    "Rebase to selected branch",
	},
	Merge: ActionType{
		Name:    "merge",
		Command: "git",
		Options: []string{"merge"},
		Help:    "Merge to selected branch",
	},
	Delete: ActionType{
		Name:    "delete",
		Command: "git",
		Options: []string{"branch", "-d"},
		Help:    "Delete branch",
	},
	Unknown: ActionType{
		Name:    "unknown",
		Command: "unknown",
		Options: nil,
		Help:    "unknown",
	},
}

func (b BranchActionTypeMap) All() []ActionType {
	return []ActionType{
		b.Switch,
		b.Diff,
		b.Delete,
		b.RebaseInteractive,
		b.Rebase,
		b.Merge,
		b.GetLastCommitId,
	}
}

func (b BranchActionTypeMap) GetBranchActionTypes(action string) (ActionType, error) {
	switch action {
	case "switch":
		return b.Switch, nil
	case "get last commit":
		return b.GetLastCommitId, nil
	case "diff":
		return b.Diff, nil
	case "rebase interactive":
		return b.RebaseInteractive, nil
	case "rebase":
		return b.Rebase, nil
	case "merge":
		return b.Merge, nil
	case "delete":
		return b.Delete, nil
	default:
		return b.Unknown, fmt.Errorf("unknown action: %s", action)
	}
}

func ParseSelectedBranchActionType(selectedLine string) (ActionType, error) {
	slog.Debug("Selected action from fzf", "selected", selectedLine)
	if selectedLine == "" {
		slog.Debug("No action selected")
		return BranchActionTypes.Unknown, nil
	}

	// タブで分割
	fields := strings.Split(selectedLine, "\t")

	// 最初のフィールドだけ取得
	selectedActionType := fields[0]

	result, err := BranchActionTypes.GetBranchActionTypes(selectedActionType)
	if err != nil {
		return BranchActionTypes.Unknown, err
	}
	return result, nil
}
