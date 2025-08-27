package fzf

import "gitman/domain/model"

type FzfManager interface {
	SelectCommit(commits []*model.Commit) (*model.Commit, error)
	SelectCommitAction(commit *model.Commit) (model.ActionType, error)
	SelectBranch(branches []*model.Branch) (*model.Branch, error)
	SelectBranchAction(branch *model.Branch) (model.ActionType, error)
	SelectReflog(reflogs []*model.Reflog) (*model.Reflog, error)
	SelectReflogAction(reflog *model.Reflog) (model.ActionType, error)
}
