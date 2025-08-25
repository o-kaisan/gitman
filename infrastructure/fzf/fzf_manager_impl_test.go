package fzf

import (
	"gitman/domain/model"
	"reflect"
	"testing"
)

func Test_parseSelectedActionType(t *testing.T) {
	t.Parallel()
	type args struct {
		selected string
	}
	tests := []struct {
		name    string
		args    args
		want    model.ActionType
		wantErr bool
	}{
		{
			name: "model.Commit.GetFzfInputForSelectActionTypeを取得すること",
			args: args{
				selected: "diff",
			},
			want:    model.CommitActionTypes.Diff,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseSelectedActionType(tt.args.selected)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSelectedActionType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSelectedActionType() = %v, want %v", got, tt.want)
			}
		})
	}
}
