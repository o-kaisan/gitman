package model

import "testing"

func TestActionType_GetOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		actionType ActionType
		want       string
	}{
		{
			name: "アクションタイプにオプションがない場合は空文字を返すこと",
			actionType: ActionType{
				Command: "diff",
				Options: nil,
				Help:    "Show changes between commits",
			},
			want: "",
		},
		{
			name: "アクションタイプにオプションが1つ設定されている場合に文字列を返せること",
			actionType: ActionType{
				Command: "commit",
				Options: []string{"--amend"},
				Help:    "Show changes between commits",
			},
			want: "--amend",
		},
		{
			name: "アクションタイプにオプションが複数設定されている場合には文字列を返すこと",
			actionType: ActionType{
				Command: "log",
				Options: []string{"--oneline", "--graph", "-n", "10"},
				Help:    "Show changes between commits",
			},
			want: "--oneline --graph -n 10",
		},
		{
			name: "アクションタイプにオプションが空配列の場合には空字列を返すこと",
			actionType: ActionType{
				Command: "log",
				Options: []string{},
				Help:    "Show changes between commits",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if result := tt.actionType.GetOptions(); result != tt.want {
				t.Errorf("GetOptions() = %v, want %v", result, tt.want)
			}
		})
	}
}
