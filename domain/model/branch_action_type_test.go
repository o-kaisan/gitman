package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBranchActionTypeMap_GetBranchActionTypes(t *testing.T) {
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
			name: "対応するブランチアクション(diff)を取得すること",
			args: args{
				action: "diff",
			},
			want:           BranchActionTypes.Diff,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するブランチアクション(Get last commit)を取得すること",
			args: args{
				action: "get last commit",
			},
			want:           BranchActionTypes.GetLastCommitId,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するブランチアクション(Switch)を取得すること",
			args: args{
				action: "switch",
			},
			want:           BranchActionTypes.Switch,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するブランチアクション(RebaseInteractive)を取得すること",
			args: args{
				action: "rebase interactive",
			},
			want:           BranchActionTypes.RebaseInteractive,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するブランチアクション(rebase)を取得すること",
			args: args{
				action: "rebase",
			},
			want:           BranchActionTypes.Rebase,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するブランチアクション(Merge)を取得すること",
			args: args{
				action: "merge",
			},
			want:           BranchActionTypes.Merge,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "対応するブランチアクション(Delete)を取得すること",
			args: args{
				action: "delete",
			},
			want:           BranchActionTypes.Delete,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "不明なアクションが指定された場合、errorを返却すること",
			args: args{
				action: "dummy",
			},
			want:           BranchActionTypes.Unknown,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("unknown action: %s", "dummy"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			got, err := BranchActionTypes.GetBranchActionTypes(tt.args.action)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("BranchActionTypeMap.GetBranchActionTypes() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("BranchActionTypeMap.GetBranchActionTypes() error message = %v, wantErrMessage %v", err, tt.wantErrMessage)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BranchActionTypeMap.GetBranchActionTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranchActionTypeMap_All(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want []ActionType
	}{
		{
			name: "Unknown以外の全てのブランチアクションを取得すること",
			want: []ActionType{
				BranchActionTypes.Switch,
				BranchActionTypes.Diff,
				BranchActionTypes.Delete,
				BranchActionTypes.RebaseInteractive,
				BranchActionTypes.Rebase,
				BranchActionTypes.Merge,
				BranchActionTypes.GetLastCommitId,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt := tt
			if got := BranchActionTypes.All(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BranchActionTypeMap.All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSelectedBranchActionType(t *testing.T) {
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
			name: "fzfの選択結果を元に、対応するブランチアクションを取得すること",
			args: args{
				selectedLine: "diff\tDescription : hogehoge\tCommand     : fugafuga\n",
			},
			want:           BranchActionTypes.Diff,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "不明な文字列が指定された場合、Unknownを返却すること",
			args: args{
				selectedLine: "dummy\tDescription : hogehoge\tCommand     : fugafuga\n",
			},
			want:           BranchActionTypes.Unknown,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("unknown action: %s", "dummy"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt := tt
			got, err := ParseSelectedBranchActionType(tt.args.selectedLine)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("ParseSelectedBranchActionType() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("ParseSelectedBranchActionType() error = %v, wantErrMessage %v", err, tt.wantErrMessage)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSelectedBranchActionType() = %v, want %v", got, tt.want)
			}
		})
	}
}
