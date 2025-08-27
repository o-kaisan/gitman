package model

import (
	"fmt"
	"log/slog"
	"strings"
)

type ReflogActionTypeMap struct {
	ResetHard ActionType
	Unknown   ActionType
}

var ReflogActionTypes = ReflogActionTypeMap{
	ResetHard: ActionType{
		Name:    "reset hard",
		Command: "git",
		Options: []string{"reset", "--hard"},
		Help:    "Hard reset to selected commit",
	},
	Unknown: ActionType{
		Name:    "unknown",
		Command: "unknown",
		Options: []string{},
		Help:    "unknown",
	},
}

func (r ReflogActionTypeMap) All() []ActionType {
	return []ActionType{
		r.ResetHard,
	}
}

func (r ReflogActionTypeMap) GetReflogActionTypes(action string) (ActionType, error) {
	switch action {
	case "reset hard":
		return r.ResetHard, nil
	default:
		return r.Unknown, fmt.Errorf("unknown action: %s", action)
	}
}

func ParseSelectedReflogActionType(selectedLine string) (ActionType, error) {
	slog.Debug("Selected action from fzf", "selected", selectedLine)
	if selectedLine == "" {
		slog.Debug("No action selected")
		return ReflogActionTypes.Unknown, nil
	}

	// タブで分割
	fields := strings.Split(selectedLine, "\t")

	// 最初のフィールドだけ取得
	selectedActionType := fields[0]

	result, err := ReflogActionTypes.GetReflogActionTypes(selectedActionType)
	if err != nil {
		return ReflogActionTypes.Unknown, err
	}
	return result, nil
}
