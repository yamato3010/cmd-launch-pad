package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yamato3010/cmd-launch-pad/internal/config"
	gitpkg "github.com/yamato3010/cmd-launch-pad/internal/git"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Gitによる設定同期",
	Long:  `Gitリポジトリを使って設定ファイルを同期します。`,
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "設定をリモートにプッシュ",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		mgr, err := gitpkg.NewGitManager(cfgDir)
		if err != nil {
			return fmt.Errorf("Gitリポジトリが初期化されていません。先に `clp sync init` を実行してください: %w", err)
		}
		appCfg, err := config.LoadAppConfig()
		if err != nil {
			return err
		}
		if err := mgr.AddAll(); err != nil {
			return err
		}
		if err := mgr.Commit("Update commands via CLI"); err != nil {
			return err
		}
		if err := mgr.Push("origin", appCfg.Git.Branch, nil); err != nil {
			return err
		}
		fmt.Println("✅ プッシュ完了")
		return nil
	},
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "リモートから設定をプル",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		mgr, err := gitpkg.NewGitManager(cfgDir)
		if err != nil {
			return fmt.Errorf("Gitリポジトリが初期化されていません: %w", err)
		}
		appCfg, err := config.LoadAppConfig()
		if err != nil {
			return err
		}
		if err := mgr.Pull("origin", appCfg.Git.Branch, nil); err != nil {
			return err
		}
		fmt.Println("✅ プル完了")
		return nil
	},
}

var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Gitステータスを表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		mgr, err := gitpkg.NewGitManager(cfgDir)
		if err != nil {
			return fmt.Errorf("Gitリポジトリが初期化されていません: %w", err)
		}
		status, err := mgr.Status()
		if err != nil {
			return err
		}
		if status == "" {
			fmt.Println("✅ 変更なし (クリーン)")
		} else {
			fmt.Println(status)
		}
		return nil
	},
}

var syncInitCmd = &cobra.Command{
	Use:   "init",
	Short: "設定ディレクトリをGitリポジトリとして初期化",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.EnsureConfigDir()
		if err != nil {
			return err
		}
		if _, err := gitpkg.Init(cfgDir); err != nil {
			return err
		}
		fmt.Printf("✅ Gitリポジトリを初期化しました: %s\n", cfgDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(syncPushCmd)
	syncCmd.AddCommand(syncPullCmd)
	syncCmd.AddCommand(syncStatusCmd)
	syncCmd.AddCommand(syncInitCmd)
}
