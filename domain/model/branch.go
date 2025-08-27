package model

import (
	"fmt"
	"log/slog"
	"strings"
)

// git branchで対象となったブランチを表す構造体
type Branch struct {
	Current             bool
	Name                string
	LastCommitId        string
	LastCommitMessage   string
	RawGirBranchMessage string
	ActionTypes         []ActionType
}

func NewBranch(current bool, name string, lastCommitId string, lastCommitMessage string, rawGirBranchMessage string) *Branch {
	return &Branch{
		Current:             current,
		Name:                name,
		LastCommitId:        lastCommitId,
		LastCommitMessage:   lastCommitMessage,
		RawGirBranchMessage: rawGirBranchMessage,
		ActionTypes:         BranchActionTypes.All(),
	}
}

func (b Branch) String() string {
	return b.Name
}

func FindBranchByBranchName(branches []*Branch, branchName string) (*Branch, error) {
	for _, branch := range branches {
		if branch.Name == branchName {
			return branch, nil
		}
	}
	return nil, fmt.Errorf("branch %s not found", branchName)
}

func (b Branch) GetFullCommand(actionType ActionType) string {
	options := b.GetOptionsWithBranchInfo(actionType)
	onelineOptions := strings.Join(options, " ")

	fullCommand := fmt.Sprintf("%s %s", actionType.Command, onelineOptions)
	slog.Debug("Command:", "Command", actionType.Name, "fullCommand", fullCommand)

	return fullCommand
}

func (b Branch) GetOptionsWithBranchInfo(actionType ActionType) []string {
	ret := actionType.Options
	slog.Debug("actionType", "actionType", actionType.Name)

	if actionType.IsEqual(BranchActionTypes.GetLastCommitId) {
		ret = append(ret, b.LastCommitId)
	} else {
		ret = append(ret, b.Name)
	}
	return ret
}

func (b Branch) GetFzfInputForSelectActionType(actionType ActionType) string {
	// fzfに渡す形式: "アクション名\tフルコマンド\t説明文"
	// Commandの空白は、Descriptionとコロンの位置が合わないための調整用
	return fmt.Sprintf("%s\tDescription : %s\tCommand     : %s\n", actionType.Name, actionType.Help, b.GetFullCommand(actionType))
}

func ParseBranches(log string) []*Branch {
	var branches []*Branch

	lines := strings.Split(strings.TrimSpace(log), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		isCurrent := false
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") {
			isCurrent = true
			line = strings.TrimSpace(line[1:]) // 先頭の "*" を除去
		}

		// 特殊ケース: "remotes/origin/HEAD -> origin/main"
		if strings.Contains(line, "->") {
			parts := strings.SplitN(line, "->", 2)
			name := strings.TrimSpace(parts[0])
			msg := "-> " + strings.TrimSpace(parts[1])
			branch := NewBranch(isCurrent, name, "NOTHING", msg, line)
			branches = append(branches, branch)
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			// 不正な行はスキップ
			continue
		}

		name := fields[0]
		commit := fields[1]
		branchLog := strings.Join(fields[2:], " ")

		branch := NewBranch(isCurrent, name, commit, branchLog, line)
		branches = append(branches, branch)
	}

	slog.Debug("get branches from git", "branches", branches)
	return branches
}
