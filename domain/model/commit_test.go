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
	t.Parallel()
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
				Message:      "commit message",
				RawCommitLog: "dummy commit message",
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
			t.Parallel()
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

func TestParseCommits(t *testing.T) {
	t.Parallel()
	type args struct {
		log string
	}
	tests := []struct {
		name string
		args args
		want []*Commit
	}{
		{
			name: "コミットを解析できること",
			args: args{
				log: "e51c238 [from now] 2025/08/27 19:27:45\n" +
					"b721c76 gitman reflogのreset --hardアクションで指定する値が誤っていたのを修正した\n" +
					"82ed65c (origin/main, origin/HEAD, main) Merge pull request #8 from o-kaisan/add-reflog-view\n",
			},
			want: []*Commit{
				NewCommit("e51c238", "[from now] 2025/08/27 19:27:45", "e51c238 [from now] 2025/08/27 19:27:45"),
				NewCommit("b721c76", "gitman reflogのreset --hardアクションで指定する値が誤っていたのを修正した", "b721c76 gitman reflogのreset --hardアクションで指定する値が誤っていたのを修正した"),
				NewCommit("82ed65c", "(origin/main, origin/HEAD, main) Merge pull request #8 from o-kaisan/add-reflog-view", "82ed65c (origin/main, origin/HEAD, main) Merge pull request #8 from o-kaisan/add-reflog-view"),
			},
		},
		{
			name: "コミットを解析できない場合はエラーを返すこと",
			args: args{
				log: "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseCommits(tt.args.log)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit_GetFullCommand(t *testing.T) {
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
			name: "アクションを選択するためにfzfに渡す文字列を生成できること(単体オプション)",
			fields: fields{
				Id:           "dummy",
				Message:      "commit message",
				RawCommitLog: "dummy commit message",
				ActionTypes:  []ActionType{CommitActionTypes.Diff},
			},
			args: args{
				actionType: CommitActionTypes.Diff,
			},
			want: "git diff dummy",
		},
		{
			name: "アクションを選択するためにfzfに渡す文字列を生成できること(複数オプション)",
			fields: fields{
				Id:           "dummy",
				Message:      "commit message",
				RawCommitLog: "dummy commit message",
				ActionTypes:  []ActionType{CommitActionTypes.CherryPickWithoutCommit},
			},
			args: args{
				actionType: CommitActionTypes.CherryPickWithoutCommit,
			},
			want: "git cherry-pick --no-commit dummy",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := Commit{
				Id:           tt.fields.Id,
				Message:      tt.fields.Message,
				RawCommitLog: tt.fields.RawCommitLog,
				ActionTypes:  tt.fields.ActionTypes,
			}
			if got := c.GetFullCommand(tt.args.actionType); got != tt.want {
				t.Errorf("Commit.GetFullCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit_GetOptionsWithCommitId(t *testing.T) {
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
		want   []string
	}{
		{
			name: "アクションを選択するためにfzfに渡す文字列を生成できること(単体オプション)",
			fields: fields{
				Id:           "dummy",
				Message:      "commit message",
				RawCommitLog: "dummy commit message",
				ActionTypes:  []ActionType{CommitActionTypes.Diff},
			},
			args: args{
				actionType: CommitActionTypes.Diff,
			},
			want: []string{"diff", "dummy"},
		},
		{
			name: "アクションを選択するためにfzfに渡す文字列を生成できること(複数オプション)",
			fields: fields{
				Id:           "dummy",
				Message:      "commit message",
				RawCommitLog: "dummy commit message",
				ActionTypes:  []ActionType{CommitActionTypes.CherryPickWithoutCommit},
			},
			args: args{
				actionType: CommitActionTypes.CherryPickWithoutCommit,
			},
			want: []string{"cherry-pick", "--no-commit", "dummy"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := Commit{
				Id:           tt.fields.Id,
				Message:      tt.fields.Message,
				RawCommitLog: tt.fields.RawCommitLog,
				ActionTypes:  tt.fields.ActionTypes,
			}
			if got := c.GetOptionsWithCommitId(tt.args.actionType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Commit.GetOptionsWithCommitId() = %v, want %v", got, tt.want)
			}
		})
	}
}
