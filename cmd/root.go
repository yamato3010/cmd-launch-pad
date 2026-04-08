package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yamato3010/cmd-launch-pad/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "clp",
	Short: "cmd-launch-pad - TUIコマンドランチャー",
	Long: `cmd-launch-pad (clp) はターミナルユーザー向けのTUIコマンドランチャーです。
nvim、lazygit、lazydockerなどのコマンドをGUIの「Launchpad」のように
視覚的に管理・起動できます。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
}

// Execute はルートコマンドを実行する
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
