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

		// 1. 最初の空白で commit id と残りを分離
		spaceIndex := strings.IndexByte(line, ' ')
		if spaceIndex == -1 {
			slog.Warn("invalid reflog line, missing space", "line", line)
			continue
		}

		id := line[:spaceIndex]
		rest := strings.TrimSpace(line[spaceIndex+1:])

		// 2. HEAD@{n} とメッセージ部分に分割
		colonIndex := strings.Index(rest, ":")
		if colonIndex == -1 {
			slog.Warn("invalid reflog line, missing colon", "line", line)
			continue
		}

		headPoint := strings.TrimSpace(rest[:colonIndex])
		message := strings.TrimSpace(rest[colonIndex+1:])

		// 3. さらに最初の action とメッセージに分割
		actionIndex := strings.Index(message, ":")
		if actionIndex != -1 {
			message = strings.TrimSpace(message[actionIndex+1:])
		}

		reflog := NewReflog(id, headPoint, message, line)
		result = append(result, reflog)

		slog.Info("debug", "HeadPoint", fmt.Sprintf("%q", headPoint), "Message", fmt.Sprintf("%q", message))
	}

	return result, nil
}
