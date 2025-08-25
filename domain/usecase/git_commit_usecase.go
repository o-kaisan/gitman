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
	if actionType.IsEqual(model.CommitActionTypes.UNKNOWN) {
		return nil
	}

	gciu.gitManager.ExecuteCommand(actionType, targetCommit)

	return nil
}

// ユーザに対象となるコミットと実行したいコマンドを選択させる
func (gciu GitCommitUsecase) getCommit() (*model.Commit, error) {
	commits, err := gciu.gitManager.GetCommits()
	if err != nil {
		return nil, err
	}

	commitId, err := gciu.fzfManager.SelectCommitId(commits)
	if err != nil {
		return nil, err
	}
	if commitId == "" {
		return nil, nil
	}
	targetCommit, err := model.FindCommitById(commits, commitId)
	if err != nil {
		return nil, err
	}
	return targetCommit, nil
}
