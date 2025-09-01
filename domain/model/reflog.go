package model

import (
	"fmt"
	"log/slog"
	"strings"
)

type Reflog struct {
	Id          string
	HeadPoint   string
	Message     string
	RawReflog   string
	ActionTypes []ActionType
}

func NewReflog(id string, HeadPoint string, message string, rawReflog string) *Reflog {
	return &Reflog{
		Id:          id,
		HeadPoint:   HeadPoint,
		Message:     message,
		RawReflog:   rawReflog,
		ActionTypes: ReflogActionTypes.All(),
	}
}

func (r Reflog) String() string {
	return r.Id
}

func FindReflogByHeadPoint(reflogs []*Reflog, headPoint string) (*Reflog, error) {
	for _, reflog := range reflogs {
		if reflog.HeadPoint == headPoint {
			slog.Info("reflog found", "HeadPoint", reflog.HeadPoint)
			return reflog, nil
		}
	}
	return nil, fmt.Errorf("reflog not found: %s", headPoint)
}

func (r Reflog) GetFullCommand(actionType ActionType) string {
	options := r.GetOptionsWithHeadPoint(actionType)
	onelineOptions := strings.Join(options, " ")

	fullCommand := fmt.Sprintf("%s %s", actionType.Command, onelineOptions)
	slog.Debug("Command:", "Command", actionType.Name, "fullCommand", fullCommand)

	return fullCommand
}

func (r Reflog) GetOptionsWithHeadPoint(actionType ActionType) []string {
	ret := actionType.Options
	ret = append(ret, r.HeadPoint)
	return ret
}

func (r Reflog) GetFzfInputForSelectActionType(actionType ActionType) string {
	// fzfに渡す形式: "表示名\tフルコマンド\t説明文"
	return fmt.Sprintf("%s\tDescription : %s\tCommand     : %s\n", actionType.Name, actionType.Help, r.GetFullCommand(actionType))
}

func ParseReflogs(reflogs string) ([]*Reflog, error) {
	if strings.TrimSpace(reflogs) == "" {
		return []*Reflog{}, nil
	}

	lines := strings.Split(strings.TrimSpace(reflogs), "\n")
	result := make([]*Reflog, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 1. commit id を取得
		spaceIndex := strings.IndexByte(line, ' ')
		if spaceIndex == -1 {
			slog.Warn("invalid reflog line, missing space", "line", line)
			continue
		}
		id := line[:spaceIndex]
		rest := strings.TrimSpace(line[spaceIndex+1:])

		// 2. HEAD@{n} の位置を探す
		headStart := strings.Index(rest, "HEAD@{")
		if headStart == -1 {
			slog.Warn("invalid reflog line, missing HEAD@{}", "line", line)
			continue
		}

		// 3. HEAD@{n} の終了位置（"}"）を探す
		headEnd := strings.Index(rest[headStart:], "}")
		if headEnd == -1 {
			slog.Warn("invalid reflog line, incomplete HEAD@{}", "line", line)
			continue
		}
		headEnd += headStart

		headPoint := rest[headStart : headEnd+1]

		// 4. コロン以降を message 本体として扱う
		colonIndex := strings.Index(rest[headEnd+1:], ":")
		if colonIndex == -1 {
			slog.Warn("invalid reflog line, missing colon", "line", line)
			continue
		}
		message := strings.TrimSpace(rest[headEnd+1+colonIndex+1:])

		reflog := NewReflog(id, headPoint, message, line)
		result = append(result, reflog)
	}

	return result, nil
}
