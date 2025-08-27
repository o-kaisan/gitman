package git

import (
	"fmt"
	"gitman/common"
	"gitman/domain/model"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type GitManagerImpl struct{}

func NewGitManager() (GitManager, error) {
	isValid, err := validGit()
	if isValid {
		return nil, err
	}

	return &GitManagerImpl{}, nil
}

func validGit() (bool, error) {
	_, err := exec.Command("git", "--version").Output()
	if err != nil {
		return true, fmt.Errorf("failed to get git version: %w", err)
	}
	return false, nil
}

func (gm GitManagerImpl) ExecuteCommitActionCommand(actionType model.ActionType, commit *model.Commit) error {
	cmd := exec.Command(actionType.Command, commit.GetOptionsWithCommitId(actionType)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 実行
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

func (gm GitManagerImpl) GetCommits() ([]*model.Commit, error) {
	logDisplayLimit := common.GetEnvWithString("GITMAN_LOG_DISPLAY_LIMIT", "100")
	cmd := exec.Command("git", "log", "--oneline", "--decorate", "-n", logDisplayLimit)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git log command: %w", err)
	}

	commitLog := strings.TrimSpace(string(out))
	commits, err := model.ParseCommits(commitLog)
	if err != nil {
		return nil, err
	}

	slog.Debug("get commitIds from git", "commitIds", commits)
	return commits, nil
}

func (gm GitManagerImpl) GetBranches() ([]*model.Branch, error) {
	cmd := exec.Command("git", "branch", "--all", "--verbose")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git branch command: %w", err)
	}

	reflogs, err := model.ParseBranches(string(out))
	if err != nil {
		return nil, err
	}
	return reflogs, nil
}

func (gm GitManagerImpl) ExecuteBranchActionCommand(actionType model.ActionType, branch *model.Branch) error {
	cmd := exec.Command(actionType.Command, branch.GetOptionsWithBranchInfo(actionType)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 実行
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}

func (gm GitManagerImpl) GetReflogs() ([]*model.Reflog, error) {
	cmd := exec.Command("git", "reflog", "-n", "50")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git reflog command: %w", err)
	}

	slog.Debug("get reflogs from git", "reflogs", out)
	reflogs, err := model.ParseReflogs(string(out))
	if err != nil {
		return nil, err
	}
	slog.Debug("get reflogs from git", "reflogs", reflogs)
	return reflogs, nil
}

func (gm GitManagerImpl) ExecuteReflogActionCommand(actionType model.ActionType, reflog *model.Reflog) error {
	cmd := exec.Command(actionType.Command, reflog.GetOptionsWithReflogId(actionType)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 実行
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}
