package model

import (
	"fmt"
	"log/slog"
	"strings"
)

// git logで対象となったコミットを表す構造体
type Commit struct {
	Id           string
	Message      string
	RawCommitLog string
	ActionTypes  []ActionType
}

func NewCommit(id string, message string, rawCommitLog string) *Commit {
	return &Commit{
		Id:           id,
		Message:      message,
		RawCommitLog: rawCommitLog,
		ActionTypes:  CommitActionTypes.All(),
	}
}

func (c Commit) String() string {
	return c.Id
}

func FindCommitById(commits []*Commit, id string) (*Commit, error) {
	for _, commit := range commits {
		if commit.Id == id {
			return commit, nil
		}
	}
	return nil, fmt.Errorf("commit %s not found", id)
}

func (c Commit) GetFullCommand(actionType ActionType) string {
	options := c.GetOptionsWithCommitId(actionType)
	onelineOptions := strings.Join(options, " ")

	fullCommand := fmt.Sprintf("%s %s", actionType.Command, onelineOptions)
	slog.Debug("Command:", "Command", actionType.Name, "fullCommand", fullCommand)

	return fullCommand
}

func (c Commit) GetOptionsWithCommitId(actionType ActionType) []string {
	ret := actionType.Options
	ret = append(ret, c.Id)
	return ret
}

func (c Commit) GetFzfInputForSelectActionType(actionType ActionType) string {
	// fzfに渡す形式: "表示名\tフルコマンド\t説明文"
	return fmt.Sprintf("%s\tDescription : %s\tCommand     : %s\n", actionType.Name, actionType.Help, c.GetFullCommand(actionType))
}

// git log --oneline の形式をパースして、Commit構造体のスライスを返す
func ParseCommits(log string) ([]*Commit, error) {
	var commits []*Commit

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

		commit := NewCommit(id, message, line)
		commits = append(commits, commit)
	}
	return commits, nil
}
