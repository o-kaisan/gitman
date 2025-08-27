package di

import (
	"gitman/domain/usecase"
	"gitman/infrastructure/fzf"
	"gitman/infrastructure/git"
)

type Container struct {
	GitBranchUsecase usecase.GitBranchUsecase
	GitCommitUsecase usecase.GitCommitUsecase
	GitReflogUsecase usecase.GitReflogUsecase
}

func NewContainer() Container {
	// infrastructureの初期化
	gm, err := git.NewGitManager()
	if err != nil {
		panic(err)
	}

	fm, err := fzf.NewFzfManager()
	if err != nil {
		panic(err)
	}

	// Usecaseの初期化
	gbu := usecase.NewGitBranchUsecase(fm, gm)
	gcu := usecase.NewGitCommitUsecase(fm, gm)
	gru := usecase.NewGitReflogUsecase(fm, gm)

	return Container{
		GitBranchUsecase: gbu,
		GitCommitUsecase: gcu,
		GitReflogUsecase: gru,
	}
}
