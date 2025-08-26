package usecase

import (
	"gitman/domain/model"
	"gitman/infrastructure/fzf"
	"gitman/infrastructure/git"
)

type GitBranchUsecase struct {
	fzfManager fzf.FzfManager
	gitManager git.GitManager
}

func NewGitBranchUsecase(fm fzf.FzfManager, gm git.GitManager) GitBranchUsecase {
	return GitBranchUsecase{
		fzfManager: fm,
		gitManager: gm,
	}
}

func (gau GitBranchUsecase) InteractiveBranchAction() error {
	targeBranch, err := gau.getBranches()
	if err != nil {
		return err
	}
	// コミットIDの選択をキャンセルした等の理由でnilとなった場合は何もしない
	if targeBranch == nil {
		return nil
	}

	actionType, err := gau.fzfManager.SelectBranchAction(targeBranch)
	if err != nil {
		return err
	}
	if actionType.IsEqual(model.BranchActionTypes.Unknown) {
		return nil
	}

	gau.gitManager.ExecuteBranchActionCommand(actionType, targeBranch)

	return nil
}

func (gau GitBranchUsecase) getBranches() (*model.Branch, error) {
	branches, err := gau.gitManager.GetBranches()
	if err != nil {
		return nil, err
	}

	selectedBranch, err := gau.fzfManager.SelectBranch(branches)
	if err != nil {
		return nil, err
	}
	if selectedBranch == nil {
		return nil, nil
	}

	return selectedBranch, nil

}
