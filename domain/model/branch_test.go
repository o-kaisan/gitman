package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindBranchById(t *testing.T) {
	type args struct {
		branches   []*Branch
		branchName string
	}
	tests := []struct {
		name           string
		args           args
		want           *Branch
		wantErr        bool
		wantErrMessage error
	}{
		{
			name: "branchNameを指定して、branchを取得すること",
			args: args{
				branches: []*Branch{
					NewBranch(true, "branch1", "branch1", "branch1", "branch1 branch1 branch1"),
					NewBranch(false, "branch2", "branch2", "branch2", "branch2 branch2 branch2"),
					NewBranch(false, "branch3", "branch3", "branch3", "branch3 branch3 branch3"),
				},
				branchName: "branch2",
			},
			want:           NewBranch(false, "branch2", "branch2", "branch2", "branch2 branch2 branch2"),
			wantErr:        false,
			wantErrMessage: nil,
		},
		{
			name: "指定したbranchNameに該当するbranchが存在しない場合にエラーを返すこと",
			args: args{
				branches: []*Branch{
					NewBranch(true, "branch1", "branch1", "branch1", "branch1 branch1 branch1"),
					NewBranch(false, "branch2", "branch2", "branch2", "branch2 branch2 branch2"),
					NewBranch(false, "branch3", "branch3", "branch3", "branch3 branch3 branch3"),
				},
				branchName: "dummy",
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("branch %s not found", "dummy"),
		},
		{
			name: "空のbranchに該当するbranchが存在しない場合にエラーを返すこと",
			args: args{
				branches:   []*Branch{},
				branchName: "branch2",
			},
			want:           nil,
			wantErr:        true,
			wantErrMessage: fmt.Errorf("branch %s not found", "branch2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindBranchByBranchName(tt.args.branches, tt.args.branchName)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("FindBranchById() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("FindBranchById() error message = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindBranchById() = %v, want %v", got, tt.want)
			}
		})
	}
}
