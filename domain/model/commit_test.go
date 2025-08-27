package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindCommitById(t *testing.T) {
	t.Parallel()
	type args struct {
		commits []*Commit
		id      string
	}
	tests := []struct {
		name           string
		args           args
		want           *Commit
		wantErr        bool
		wantErrMessage error
	}{
		{
			name: "コミットIDを指定してコミットを取得できること",
			args: args{
				commits: []*Commit{
					NewCommit("commit1", "commit1", "commit1"),
					NewCommit("commit2", "commit2", "commit2"),
					NewCommit("commit3", "commit3", "commit3"),
				},
				id: "commit2",
			},
			want:           NewCommit("commit2", "commit2", "commit2"),
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "存在しないコミットIDを指定した場合にエラーが返ること(検索対象に含まれない)",
			args: args{
				commits: []*Commit{
					NewCommit("commit1", "commit1", "commit1"),
					NewCommit("commit2", "commit2", "commit2"),
					NewCommit("commit3", "commit3", "commit3"),
				},
				id: "dummy",
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("commit dummy not found"),
		},
		{
			name: "存在しないコミットIDを指定した場合にエラーが返ること(検索対象が空)",
			args: args{
				commits: []*Commit{},
				id:      "dummy",
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("commit dummy not found"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindCommitById(tt.args.commits, tt.args.id)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("FindCommitById() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("FindCommitById() error = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindCommitById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit_GetFzfInputForSelectActionType(t *testing.T) {
	type fields struct {
		Id           string
		Message      string
		RawCommitLog string
		ActionTypes  []ActionType
	}
	type args struct {
		actionType ActionType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "アクションを選択するためにfzfに渡す文字列を生成できること",
			fields: fields{
				Id:           "dummy",
				Message:      "commit message",       // 使わない
				RawCommitLog: "dummy commit message", // 使わない
				ActionTypes:  []ActionType{CommitActionTypes.Diff},
			},
			args: args{
				actionType: CommitActionTypes.Diff,
			},
			want: "diff\tDescription : Show changes between commits\tCommand     : git diff dummy\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := Commit{
				Id:           tt.fields.Id,
				Message:      tt.fields.Message,
				RawCommitLog: tt.fields.RawCommitLog,
				ActionTypes:  tt.fields.ActionTypes,
			}
			if got := c.GetFzfInputForSelectActionType(tt.args.actionType); got != tt.want {
				t.Errorf("Commit.GetFzfInputForSelectActionType() = %v, want %v", got, tt.want)
			}
		})
	}
}
