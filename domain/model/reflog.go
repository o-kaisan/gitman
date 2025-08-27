package model

import (
	"fmt"
	"log/slog"
	"regexp"
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

func FindReflogById(reflogs []*Reflog, id string) (*Reflog, error) {
	for _, reflog := range reflogs {
		if reflog.Id == id {
			return reflog, nil
		}
	}
	return nil, fmt.Errorf("reflog not found: %s", id)
}

func (r Reflog) GetFullCommand(actionType ActionType) string {
	options := r.GetOptionsWithReflogId(actionType)
	onelineOptions := strings.Join(options, " ")

	fullCommand := fmt.Sprintf("%s %s", actionType.Command, onelineOptions)
	slog.Debug("Command:", "Command", actionType.Name, "fullCommand", fullCommand)

	return fullCommand
}

func (r Reflog) GetOptionsWithReflogId(actionType ActionType) []string {
	ret := actionType.Options
	ret = append(ret, r.Id)
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

	// Regular expression to parse reflog entries
	// Pattern: {hash} HEAD@{index}: {action}: {message}
	// HEAD@{index} が並ぶようにあえて(origin/main, origin/HEAD, main) のような情報を含めない
	reflogPattern := regexp.MustCompile(`^([a-f0-9]+)\s+(HEAD@\{[0-9]+\}):\s+([^:]+):\s*(.*)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := reflogPattern.FindStringSubmatch(line)
		if len(matches) != 5 {
			// If the pattern doesn't match, skip this line or handle as needed
			continue
		}

		slog.Debug("matches", "matches", matches)

		id := matches[1]
		headPoint := matches[2]
		message := strings.TrimSpace(matches[4])

		reflog := NewReflog(id, headPoint, message, line)

		result = append(result, reflog)
	}

	return result, nil
}
