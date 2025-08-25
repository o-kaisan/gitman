package git

import "gitman/domain/model"

type GitManager interface {
	GetCommits() ([]*model.Commit, error)
	ExecuteCommand(actionType model.ActionType, commit *model.Commit) error
}
