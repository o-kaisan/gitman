package di

import (
	"gitman/domain/usecase"
	"gitman/infrastructure/fzf"
	"gitman/infrastructure/git"
)

type Container struct {
	GitCommitUsecase usecase.GitCommitUsecase
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
	gcu := usecase.NewGitCommitUsecase(fm, gm)

	return Container{
		GitCommitUsecase: gcu,
	}
}
