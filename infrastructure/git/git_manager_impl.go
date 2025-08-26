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
	commits, err := parseCommits(commitLog)
	if err != nil {
		return nil, err
	}

	slog.Debug("get commitIds from git", "commitIds", commits)
	return commits, nil
}

// git log --oneline の形式をパースして、Commit構造体のスライスを返す
func parseCommits(log string) ([]*model.Commit, error) {
	var commits []*model.Commit

	lines := strings.Split(strings.TrimSpace(log), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// 最初の空白で分割（commit id と残り）
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		id := parts[0]
		message := parts[1]

		commit := model.NewCommit(id, message, line)
		commits = append(commits, commit)
	}
	return commits, nil
}
