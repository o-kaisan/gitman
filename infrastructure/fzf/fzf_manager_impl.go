package fzf

import (
	"bytes"
	"fmt"
	"gitman/domain/model"
	"log/slog"
	"os/exec"
	"strings"
)

type FzfManagerImpl struct{}

// NewFzfManager は FzfManagerImpl を返す
func NewFzfManager() (FzfManager, error) {
	isValid, err := isValidFzf()
	if !isValid {
		return nil, err
	}

	return &FzfManagerImpl{}, nil
}

// isValidFzf は fzf のバージョンを検証し、fzf がインストールされているかどうかを返す
func isValidFzf() (bool, error) {
	v, err := exec.Command("fzf", "--version").Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ProcessState.ExitCode() == 127 { // 127 is "command not found"
			return false, fmt.Errorf("fzf is not installed. Please install fzf(0.65.x or later) \ndownload from https://github.com/junegunn/fzf")
		}
		return false, fmt.Errorf("failed to execute fzf. Please check whether fzf is installed: err=(%w)", err)
	}

	// バージョンが 0.65.x未満の場合はエラーを返す
	if ver := strings.Split(string(v), ".")[0]; ver != "0" || strings.Compare(strings.Split(string(v), ".")[1], "65") < 0 {
		return false, fmt.Errorf("fzf version is too old. Please check 'fzf --version'\nif fzf version is > 0.65.x upgrade fzf(0.65.x or later) \ndownload from https://github.com/junegunn/fzf")
	}

	return true, nil
}

func (fm FzfManagerImpl) SelectCommitId(commits []*model.Commit) (*model.Commit, error) {
	cmd := exec.Command("fzf",
		"--ansi",
		"--prompt=gitman-log> ",
		"--preview", "echo {} | awk '{print $1}' | xargs git show --color=always --stat -p",
		"--preview-window=right:60%:wrap",                       // 右側に60%、折り返し表示
		"--bind", "shift-down:preview-down,shift-up:preview-up", // ctrl+j / ctrl+k で移動
		"--bind", "pgdn:preview-page-down,pgup:preview-page-up",
		"--bind", "ctrl-s:toggle-preview",
	)

	var in bytes.Buffer
	for _, commit := range commits {
		in.WriteString(commit.RawCommitLog + "\n")
	}
	cmd.Stdin = &in

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// ユーザーがキャンセルした場合（ESCキーやCtrl+C）
			if exitErr.ExitCode() == 1 || exitErr.ExitCode() == 130 {
				slog.Debug("User cancelled commit id selection")
				return nil, nil
			}
		}
		return nil, fmt.Errorf("fzf failed: %w", err)
	}

	selected := strings.TrimSpace(out.String())
	if selected == "" {
		return nil, nil // 選択なしはエラーにせず空文字
	}

	commitId := strings.Fields(selected)[0]
	slog.Debug("selected commitId", "commitId", commitId)

	selectedCommit, err := model.FindCommitById(commits, commitId)
	if err != nil {
		return nil, err
	}

	return selectedCommit, nil
}

func (fm FzfManagerImpl) SelectCommitAction(commit *model.Commit) (model.ActionType, error) {
	if commit == nil {
		return model.ActionType{}, fmt.Errorf("commit cannot be nil")
	}

	// fzfコマンドの基本設定
	cmd := exec.Command("fzf",
		"--ansi",
		"--prompt=gitman-action> ",
		"--delimiter", "\t", // タブを区切りに指定
		"--with-nth=1",                           // 1列目 (ActionName) だけを候補リストに表示
		"--preview", "printf '%s\n%s\n' {2} {3}", // 2列目=fullCommand, 3列目=Help
		"--preview-window=right:70%:wrap",
		"--border",
	)

	// 入力データの準備
	var in bytes.Buffer
	slog.Debug("ActionTypes", "commit.ActionTypes", commit.ActionTypes)
	for _, actionType := range commit.ActionTypes {

		// fzfに渡す形式: "表示名\tフルコマンド\t説明文"
		in.WriteString(commit.GetFzfInputForSelectActionType(actionType))
	}

	slog.Debug("fzf input", "input", in.String())
	cmd.Stdin = &in

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	// コマンド実行
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// ユーザーがキャンセルした場合（ESCキーやCtrl+C）
			if exitErr.ExitCode() == 1 || exitErr.ExitCode() == 130 {
				slog.Debug("User cancelled action selection")
				return model.CommitActionTypes.UNKNOWN, nil
			}
		}
		return model.ActionType{}, fmt.Errorf("fzf failed: %w, stderr: %s", err, errOut.String())
	}

	selected := strings.TrimSpace(out.String())
	SelectedActionType, err := parseSelectedActionType(selected)
	if err != nil {
		return model.ActionType{}, fmt.Errorf("failed to parse selected action type: %w", err)
	}

	return SelectedActionType, nil
}

func parseSelectedActionType(selectedLine string) (model.ActionType, error) {
	slog.Debug("Selected action from fzf", "selected", selectedLine)
	if selectedLine == "" {
		slog.Debug("No action selected")
		return model.CommitActionTypes.UNKNOWN, nil
	}

	// タブで分割
	fields := strings.Split(selectedLine, "\t")

	// 最初のフィールドだけ取得
	selectedActionType := fields[0]

	result, err := model.CommitActionTypes.GetCommitActionTypes(selectedActionType)
	if err != nil {
		return model.CommitActionTypes.UNKNOWN, err
	}
	return result, nil
}
