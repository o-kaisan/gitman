package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindReflogById(t *testing.T) {
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
				id: "2x2x2x",
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
			got, err := FindReflogById(tt.args.reflogs, tt.args.id)
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
				Message:     "commit message",       // 使わない
				RawReflog:   "dummy commit message", // 使わない
				ActionTypes: []ActionType{ReflogActionTypes.ResetHard},
			},
			args: args{
				actionType: ReflogActionTypes.ResetHard,
			},
			want: "reset hard\tDescription : Hard reset to selected commit\tCommand     : git reset --hard dummy\n",
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
