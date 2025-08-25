package fzf

import "gitman/domain/model"

type FzfManager interface {
	SelectCommitId(commits []*model.Commit) (*model.Commit, error)
	SelectCommitAction(commit *model.Commit) (model.ActionType, error)
}
