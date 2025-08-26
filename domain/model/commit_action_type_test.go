package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCommitActionTypeMap_GetCommitActionTypes(t *testing.T) {
	t.Parallel()
	type args struct {
		action string
	}
	tests := []struct {
		name           string
		args           args
		want           ActionType
		wantErr        bool
		wantErrMessage error
	}{
		{
			name: "対応するコミットアクション(get commit id)を取得すること",
			args: args{
				action: "get commit id",
			},
			want:           CommitActionTypes.GetCommitId,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するコミットアクション(diff)を取得すること",
			args: args{
				action: "diff",
			},
			want:           CommitActionTypes.Diff,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するコミットアクション(rebase -i)を取得すること",
			args: args{
				action: "rebase interactive",
			},
			want:           CommitActionTypes.RebaseInteractive,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するコミットアクション(revert)を取得すること",
			args: args{
				action: "revert",
			},
			want:           CommitActionTypes.Revert,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するコミットアクション(revert no commit)を取得すること",
			args: args{
				action: "revert no commit",
			},
			want:           CommitActionTypes.RevertWithoutCommit,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するコミットアクション(cherry-pick)を取得すること",
			args: args{
				action: "cherry-pick",
			},
			want:           CommitActionTypes.CherryPick,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するコミットアクション(cherry-pick without commit)を取得すること",
			args: args{
				action: "cherry-pick without commit",
			},
			want:           CommitActionTypes.CherryPickWithoutCommit,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "不明な文字列が来た場合にはUNKNOWNとエラーを返すこと",
			args: args{
				action: "unknown action",
			},
			want:           CommitActionTypes.Unknown,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("unknown action: unknown action"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := CommitActionTypes
			got, err := c.GetCommitActionTypes(tt.args.action)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("CommitActionTypeMap.GetCommitActionTypes() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("CommitActionTypeMap.GetCommitActionTypes() error = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommitActionTypeMap.GetCommitActionTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommitActionTypeMap_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want []ActionType
	}{
		{
			name: "全てのコミットアクションを取得すること",
			want: []ActionType{
				CommitActionTypes.GetCommitId,
				CommitActionTypes.Diff,
				CommitActionTypes.RebaseInteractive,
				CommitActionTypes.Revert,
				CommitActionTypes.RevertWithoutCommit,
				CommitActionTypes.CherryPick,
				CommitActionTypes.CherryPickWithoutCommit,
				CommitActionTypes.Switch,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommitActionTypes.All(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommitActionTypeMap.All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSelectedCommitActionType(t *testing.T) {
	t.Parallel()
	type args struct {
		selectedLine string
	}
	tests := []struct {
		name           string
		args           args
		want           ActionType
		wantErr        bool
		wantErrMessage error
	}{
		{
			name: "fzfの選択結果を元に、対応するコミットアクションを取得すること",
			args: args{
				selectedLine: "diff\tDescription : hogehoge\tCommand     : fugafuga\n",
			},
			want:           CommitActionTypes.Diff,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "不明な文字列が指定された場合、Unknownを返却すること",
			args: args{
				selectedLine: "dummy\tDescription : hogehoge\tCommand     : fugafuga\n",
			},
			want:           CommitActionTypes.Unknown,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("unknown action: %s", "dummy"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt := tt
			got, err := ParseSelectedCommitActionType(tt.args.selectedLine)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("ParseSelectedCommitActionType() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("ParseSelectedCommitActionType() error = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSelectedCommitActionType() = %v, want %v", got, tt.want)
			}
		})
	}
}
