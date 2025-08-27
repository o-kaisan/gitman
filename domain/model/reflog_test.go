package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindReflogByHeadPoint(t *testing.T) {
	t.Parallel()
	type args struct {
		reflogs []*Reflog
		id      string
	}
	tests := []struct {
		name           string
		args           args
		want           *Reflog
		wantErr        bool
		wantErrMessage error
	}{
		{
			name: "reflogのIDを指定して対象のreflogが取得できること",
			args: args{
				reflogs: []*Reflog{
					NewReflog("5d5d5d", "HEAD@{0}", "test message", "5d5d5d HEAD@{0}: test action: test message"),
					NewReflog("1x1x1x", "HEAD@{1}", "test message", "5d5d5d HEAD@{1}: test action: test message"),
					NewReflog("2x2x2x", "HEAD@{2}", "test message", "5d5d5d HEAD@{2}: test action: test message"),
				},
				id: "HEAD@{2}",
			},
			want:           NewReflog("2x2x2x", "HEAD@{2}", "test message", "5d5d5d HEAD@{2}: test action: test message"),
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "存在しないreflogIDを指定した場合にエラーが返ること(検索対象に含まれない)",
			args: args{
				reflogs: []*Reflog{
					NewReflog("5d5d5d", "HEAD@{0}", "test message", "5d5d5d HEAD@{0}: test action: test message"),
					NewReflog("1x1x1x", "HEAD@{1}", "test message", "5d5d5d HEAD@{1}: test action: test message"),
					NewReflog("2x2x2x", "HEAD@{2}", "test message", "5d5d5d HEAD@{2}: test action: test message"),
				},
				id: "dummy",
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("reflog not found: %s", "dummy"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindReflogByHeadPoint(tt.args.reflogs, tt.args.id)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("FindReflogById() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("FindReflogById() error = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindReflogById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReflog_GetFzfInputForSelectActionType(t *testing.T) {
	type fields struct {
		Id          string
		HeadPoint   string
		Message     string
		RawReflog   string
		ActionTypes []ActionType
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
				Id:          "dummy",
				HeadPoint:   "HEAD@{0}",
				Message:     "commit message",
				RawReflog:   "dummy commit message",
				ActionTypes: []ActionType{ReflogActionTypes.ResetHard},
			},
			args: args{
				actionType: ReflogActionTypes.ResetHard,
			},
			want: "reset hard\tDescription : Hard reset to selected commit\tCommand     : git reset --hard HEAD@{0}\n",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Reflog{
				Id:          tt.fields.Id,
				HeadPoint:   tt.fields.HeadPoint,
				Message:     tt.fields.Message,
				RawReflog:   tt.fields.RawReflog,
				ActionTypes: tt.fields.ActionTypes,
			}
			if got := r.GetFzfInputForSelectActionType(tt.args.actionType); got != tt.want {
				t.Errorf("Reflog.GetFzfInputForSelectActionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseReflogs(t *testing.T) {
	type args struct {
		reflogs string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Reflog
		wantErr bool
	}{
		{
			name: "ブランチ/タグ装飾あり: message にも装飾を含める",
			args: args{
				reflogs: "" +
					"5d5d5d7 (HEAD -> main, tag: v1.0.0) HEAD@{0}: commit: initial commit\n" +
					"1a2b3c4 (origin/main, origin/HEAD) HEAD@{1}: commit (amend): fix typo\n" +
					"3f3f3f3 (feature/login) HEAD@{2}: checkout: moving from main to feature/login\n",
			},
			want: []*Reflog{
				NewReflog(
					"5d5d5d7",
					"HEAD@{0}",
					"(HEAD -> main, tag: v1.0.0) commit: initial commit",
					"5d5d5d7 (HEAD -> main, tag: v1.0.0) HEAD@{0}: commit: initial commit",
				),
				NewReflog(
					"1a2b3c4",
					"HEAD@{1}",
					"(origin/main, origin/HEAD) commit (amend): fix typo",
					"1a2b3c4 (origin/main, origin/HEAD) HEAD@{1}: commit (amend): fix typo",
				),
				NewReflog(
					"3f3f3f3",
					"HEAD@{2}",
					"(feature/login) checkout: moving from main to feature/login",
					"3f3f3f3 (feature/login) HEAD@{2}: checkout: moving from main to feature/login",
				),
			},
			wantErr: false,
		},
		{
			name: "ブランチ/タグ装飾なしでも解析できること",
			args: args{
				reflogs: "" +
					"aaaaaaaa HEAD@{0}: commit: chore: update deps\n" +
					"bbbbbbbb HEAD@{1}: commit: feat: add login\n",
			},
			want: []*Reflog{
				NewReflog(
					"aaaaaaaa",
					"HEAD@{0}",
					"commit: chore: update deps",
					"aaaaaaaa HEAD@{0}: commit: chore: update deps",
				),
				NewReflog(
					"bbbbbbbb",
					"HEAD@{1}",
					"commit: feat: add login",
					"bbbbbbbb HEAD@{1}: commit: feat: add login",
				),
			},
			wantErr: false,
		},
		{
			name: "空入力は空スライス",
			args: args{
				reflogs: "",
			},
			want:    []*Reflog{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseReflogs(tt.args.reflogs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReflogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseReflogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReflog_GetFullCommand(t *testing.T) {
	type fields struct {
		Id          string
		HeadPoint   string
		Message     string
		RawReflog   string
		ActionTypes []ActionType
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := Reflog{
				Id:          tt.fields.Id,
				HeadPoint:   tt.fields.HeadPoint,
				Message:     tt.fields.Message,
				RawReflog:   tt.fields.RawReflog,
				ActionTypes: tt.fields.ActionTypes,
			}
			if got := r.GetFullCommand(tt.args.actionType); got != tt.want {
				t.Errorf("Reflog.GetFullCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
