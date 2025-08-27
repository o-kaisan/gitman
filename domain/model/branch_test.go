package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindBranchByBranchName(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := FindBranchByBranchName(tt.args.branches, tt.args.branchName)
			if (err != nil) != tt.wantErr || err != nil && err.Error() != tt.wantErrMessage.Error() {
				t.Errorf("FindBranchByBranchName() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("FindBranchByBranchName() error message = %v, wantErrMessage %v", err.Error(), tt.wantErrMessage.Error())
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindBranchByBranchName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBranches(t *testing.T) {
	type args struct {
		log string
	}
	tests := []struct {
		name string
		args args
		want []*Branch
	}{
		{
			name: "branchを取得すること",
			args: args{
				log: "* branch1 xxxxxx [ahead 4] [from now] dummy\n" +
					"branch2 yyyyyy [from now] dummy\n" +
					"feature/foo e57408f [gone] Fix deploy\n",
			},
			want: []*Branch{
				NewBranch(true, "branch1", "xxxxxx", "[ahead 4] [from now] dummy", "branch1 xxxxxx [ahead 4] [from now] dummy"),
				NewBranch(false, "branch2", "yyyyyy", "[from now] dummy", "branch2 yyyyyy [from now] dummy"),
				NewBranch(false, "feature/foo", "e57408f", "[gone] Fix deploy", "feature/foo e57408f [gone] Fix deploy"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseBranches(tt.args.log)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBranches() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestBranch_GetFullCommand(t *testing.T) {
	t.Parallel()
	type fields struct {
		Current             bool
		Name                string
		LastCommitId        string
		LastCommitMessage   string
		RawGirBranchMessage string
		ActionTypes         []ActionType
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
				Current:             true,
				Name:                "dummy",
				LastCommitId:        "dummy",
				LastCommitMessage:   "dummy",
				RawGirBranchMessage: "dummy",
				ActionTypes:         []ActionType{BranchActionTypes.Diff},
			},
			args: args{
				actionType: BranchActionTypes.Diff,
			},
			want: "git diff dummy",
		},
		{
			name: "アクションを選択するためにfzfに渡す文字列を生成できること(複数オプション)",
			fields: fields{
				Current:             true,
				Name:                "dummy",
				LastCommitId:        "dummy",
				LastCommitMessage:   "dummy",
				RawGirBranchMessage: "dummy",
				ActionTypes:         []ActionType{BranchActionTypes.RebaseInteractive},
			},
			args: args{
				actionType: BranchActionTypes.RebaseInteractive,
			},
			want: "git rebase -i dummy",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := Branch{
				Current:             tt.fields.Current,
				Name:                tt.fields.Name,
				LastCommitId:        tt.fields.LastCommitId,
				LastCommitMessage:   tt.fields.LastCommitMessage,
				RawGirBranchMessage: tt.fields.RawGirBranchMessage,
				ActionTypes:         tt.fields.ActionTypes,
			}
			if got := b.GetFullCommand(tt.args.actionType); got != tt.want {
				t.Errorf("Branch.GetFullCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranch_GetOptionsWithBranchInfo(t *testing.T) {
	t.Parallel()
	type fields struct {
		Current             bool
		Name                string
		LastCommitId        string
		LastCommitMessage   string
		RawGirBranchMessage string
		ActionTypes         []ActionType
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
				Current:             true,
				Name:                "dummy",
				LastCommitId:        "dummy",
				LastCommitMessage:   "dummy",
				RawGirBranchMessage: "dummy",
				ActionTypes:         []ActionType{BranchActionTypes.Diff},
			},
			args: args{
				actionType: BranchActionTypes.Diff,
			},
			want: []string{"diff", "dummy"},
		},
		{
			name: "アクションを選択するためにfzfに渡す文字列を生成できること(複数オプション)",
			fields: fields{
				Current:             true,
				Name:                "dummy",
				LastCommitId:        "dummy",
				LastCommitMessage:   "dummy",
				RawGirBranchMessage: "dummy",
				ActionTypes:         []ActionType{BranchActionTypes.RebaseInteractive},
			},
			args: args{
				actionType: BranchActionTypes.RebaseInteractive,
			},
			want: []string{"rebase", "-i", "dummy"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := Branch{
				Current:             tt.fields.Current,
				Name:                tt.fields.Name,
				LastCommitId:        tt.fields.LastCommitId,
				LastCommitMessage:   tt.fields.LastCommitMessage,
				RawGirBranchMessage: tt.fields.RawGirBranchMessage,
				ActionTypes:         tt.fields.ActionTypes,
			}
			if got := b.GetOptionsWithBranchInfo(tt.args.actionType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Branch.GetOptionsWithBranchInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
