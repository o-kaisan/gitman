package model

import (
	"fmt"
	"log/slog"
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
	fullCommand := fmt.Sprintf("%s %s", actionType.Command, c.Id)
	if len(actionType.Options) > 0 {
		// オプションがない場合
		fullCommand = fmt.Sprintf("%s %s %s", actionType.Command, actionType.GetOptions(), c.Id)
	}
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
