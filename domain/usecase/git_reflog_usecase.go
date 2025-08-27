package usecase

import (
	"gitman/domain/model"
	"gitman/infrastructure/fzf"
	"gitman/infrastructure/git"
)

type GitReflogUsecase struct {
	fzfManager fzf.FzfManager
	gitManager git.GitManager
}

func NewGitReflogUsecase(fm fzf.FzfManager, gm git.GitManager) GitReflogUsecase {
	return GitReflogUsecase{
		fzfManager: fm,
		gitManager: gm,
	}
}

func (gru GitReflogUsecase) InteractiveReflogAction() error {
	targetReflog, err := gru.getReflog()
	if err != nil {
		return err
	}
	// コミットIDの選択をキャンセルした等の理由でnilとなった場合は何もしない
	if targetReflog == nil {
		return nil
	}

	actionType, err := gru.fzfManager.SelectReflogAction(targetReflog)
	if err != nil {
		return err
	}
	if actionType.IsEqual(model.ReflogActionTypes.Unknown) {
		return nil
	}

	gru.gitManager.ExecuteReflogActionCommand(actionType, targetReflog)

	return nil
}

// ユーザに対象となるコミットと実行したいコマンドを選択させる
func (gru GitReflogUsecase) getReflog() (*model.Reflog, error) {
	reflogs, err := gru.gitManager.GetReflogs()
	if err != nil {
		return nil, err
	}

	selectedReflog, err := gru.fzfManager.SelectReflog(reflogs)
	if err != nil {
		return nil, err
	}
	if selectedReflog == nil {
		return nil, nil
	}
	return selectedReflog, nil
}
