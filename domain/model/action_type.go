package model

import "strings"

type ActionType struct {
	Name    string
	Command string
	Options []string
	Help    string
}

func (a ActionType) IsEqual(target ActionType) bool {
	// 対象が定数のため名前だけで一致を確認する
	return a.Name == target.Name
}

// optionsを半角スペース区切りで1つの文字列に変換する
func (a ActionType) GetOptions() string {
	ret := ""
	if a.Options != nil && len(a.Options) >= 0 {
		ret = strings.Join(a.Options, " ")
		return ret
	}
	return ret
}
