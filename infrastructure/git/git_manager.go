package git

import "gitman/domain/model"

type GitManager interface {
	GetCommits() ([]*model.Commit, error)
	GetBranches() ([]*model.Branch, error)
	GetReflogs() ([]*model.Reflog, error)
	ExecuteCommitActionCommand(actionType model.ActionType, commit *model.Commit) error
	ExecuteBranchActionCommand(actionType model.ActionType, branch *model.Branch) error
	ExecuteReflogActionCommand(actionType model.ActionType, reflog *model.Reflog) error
}
