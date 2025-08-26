package git

import "gitman/domain/model"

type GitManager interface {
	GetCommits() ([]*model.Commit, error)
	GetBranches() ([]*model.Branch, error)
	ExecuteCommitActionCommand(actionType model.ActionType, commit *model.Commit) error
	ExecuteBranchActionCommand(actionType model.ActionType, branch *model.Branch) error
}
