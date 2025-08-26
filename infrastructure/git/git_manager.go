package git

import "gitman/domain/model"

type GitManager interface {
	GetCommits() ([]*model.Commit, error)
	ExecuteCommitActionCommand(actionType model.ActionType, commit *model.Commit) error
}
