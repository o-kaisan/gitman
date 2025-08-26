package usecase

import (
	"gitman/domain/model"
	"gitman/infrastructure/fzf"
	"gitman/infrastructure/git"
)

type GitCommitUsecase struct {
	fzfManager fzf.FzfManager
	gitManager git.GitManager
}

func NewGitCommitUsecase(fm fzf.FzfManager, gm git.GitManager) GitCommitUsecase {
	return GitCommitUsecase{
		fzfManager: fm,
		gitManager: gm,
	}
}

func (gciu GitCommitUsecase) InteractiveCommitAction() error {
	targetCommit, err := gciu.getCommit()
	if err != nil {
		return err
	}
	// コミットIDの選択をキャンセルした等の理由でnilとなった場合は何もしない
	if targetCommit == nil {
		return nil
	}

	actionType, err := gciu.fzfManager.SelectCommitAction(targetCommit)
	if err != nil {
		return err
	}
	if actionType.IsEqual(model.CommitActionTypes.Unknown) {
		return nil
	}

	gciu.gitManager.ExecuteCommitActionCommand(actionType, targetCommit)

	return nil
}

// ユーザに対象となるコミットと実行したいコマンドを選択させる
func (gciu GitCommitUsecase) getCommit() (*model.Commit, error) {
	commits, err := gciu.gitManager.GetCommits()
	if err != nil {
		return nil, err
	}

	selectedCommit, err := gciu.fzfManager.SelectCommit(commits)
	if err != nil {
		return nil, err
	}
	if selectedCommit == nil {
		return nil, nil
	}
	return selectedCommit, nil
}
