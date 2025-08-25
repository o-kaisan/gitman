package model

import "fmt"

type CommitActionTypeMap struct {
	GetCommitId             ActionType
	Diff                    ActionType
	RebaseInteractive       ActionType
	Revert                  ActionType
	RevertWithoutCommit     ActionType
	CherryPick              ActionType
	CherryPickWithoutCommit ActionType
	Switch                  ActionType
	UNKNOWN                 ActionType
}

var CommitActionTypes = CommitActionTypeMap{
	GetCommitId: ActionType{
		Name:    "get commit id",
		Command: "echo",
		Options: nil,
		Help:    "print commit id",
	},
	Diff: ActionType{
		Name:    "diff",
		Command: "git",
		Options: []string{"diff"},
		Help:    "Show changes between commits",
	},
	RebaseInteractive: ActionType{
		Name:    "rebase interactive",
		Command: "git",
		Options: []string{"rebase", "-i"},
		Help:    "Interactive rebase",
	},
	Revert: ActionType{
		Name:    "revert",
		Command: "git",
		Options: []string{"revert", "--edit"},
		Help:    "Revert commit",
	},
	RevertWithoutCommit: ActionType{
		Name:    "revert no commit",
		Command: "git",
		Options: []string{"revert", "--no-commit"},
		Help:    "Revert without committing",
	},
	CherryPick: ActionType{
		Name:    "cherry-pick",
		Command: "git",
		Options: []string{"cherry-pick"},
		Help:    "Cherry-pick commit",
	},
	CherryPickWithoutCommit: ActionType{
		Name:    "cherry-pick without commit",
		Command: "git",
		Options: []string{"cherry-pick", "--no-commit"},
		Help:    "Cherry-pick without committing",
	},
	Switch: ActionType{
		Name:    "switch",
		Command: "git",
		Options: []string{"switch"},
		Help:    "Switch branch to commit",
	},
	UNKNOWN: ActionType{
		Name:    "unknown",
		Command: "unknown",
		Options: nil,
		Help:    "unknown",
	},
}

func (c CommitActionTypeMap) All() []ActionType {
	return []ActionType{
		c.GetCommitId,
		c.Diff,
		c.RebaseInteractive,
		c.Revert,
		c.RevertWithoutCommit,
		c.CherryPick,
		c.CherryPickWithoutCommit,
		c.Switch,
	}
}

func (c CommitActionTypeMap) GetCommitActionTypes(action string) (ActionType, error) {
	switch action {
	case "get commit id":
		return c.GetCommitId, nil
	case "diff":
		return c.Diff, nil
	case "rebase interactive":
		return c.RebaseInteractive, nil
	case "revert":
		return c.Revert, nil
	case "revert no commit":
		return c.RevertWithoutCommit, nil
	case "switch":
		return c.Switch, nil
	case "cherry-pick":
		return c.CherryPick, nil
	case "cherry-pick without commit":
		return c.CherryPickWithoutCommit, nil
	default:
		return c.UNKNOWN, fmt.Errorf("unknown action: %s", action)
	}
}
