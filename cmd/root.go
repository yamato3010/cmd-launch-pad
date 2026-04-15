package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yamato3010/cmd-launch-pad/internal/config"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "clp",
	Short: i18n.T("root.short"),
	Long:  i18n.T("root.long"),
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

func init() {
	// 設定ファイルから言語を読み込み、i18n を初期化する
	// エラーが発生した場合はデフォルト(英語)のまま続行する
	if cfg, err := config.LoadAppConfig(); err == nil {
		lang := i18n.DetectLang(cfg.Language)
		i18n.SetLang(lang)
	} else {
		lang := i18n.DetectLang("")
		i18n.SetLang(lang)
	}
}
