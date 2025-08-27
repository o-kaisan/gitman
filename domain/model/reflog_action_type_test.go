package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflogActionTypeMap_All(t *testing.T) {
	type fields struct {
		ResetHard ActionType
		Unknown   ActionType
	}
	tests := []struct {
		name   string
		fields fields
		want   []ActionType
	}{
		{
			name: "全てのReflogアクションを取得すること",
			fields: fields{
				ResetHard: ReflogActionTypes.ResetHard,
			},
			want: []ActionType{
				ReflogActionTypes.ResetHard,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ReflogActionTypeMap{
				ResetHard: tt.fields.ResetHard,
				Unknown:   tt.fields.Unknown,
			}
			if got := r.All(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReflogActionTypeMap.All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSelectedReflogActionType(t *testing.T) {
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
			name: "fzfの選択結果を元に、対応するReflogアクションを取得すること",
			args: args{
				selectedLine: "reset hard\tDescription : hogehoge\tCommand     : fugafuga\n",
			},
			want:           ReflogActionTypes.ResetHard,
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "不明な文字列が指定された場合、Unknownを返却すること",
			args: args{
				selectedLine: "dummy\tDescription : hogehoge\tCommand     : fugafuga\n",
			},
			want:           ReflogActionTypes.Unknown,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("unknown action: %s", "dummy"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseSelectedReflogActionType(tt.args.selectedLine)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("ParseSelectedReflogActionType() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("ParseSelectedReflogActionType() error = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSelectedReflogActionType() = %v, want %v", got, tt.want)
			}
		})
	}
}
