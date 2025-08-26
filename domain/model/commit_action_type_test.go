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
