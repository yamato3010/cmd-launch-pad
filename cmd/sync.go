package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yamato3010/cmd-launch-pad/internal/config"
	gitpkg "github.com/yamato3010/cmd-launch-pad/internal/git"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: i18n.T("sync.short"),
	Long:  i18n.T("sync.long"),
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: i18n.T("sync.push.short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		mgr, err := gitpkg.NewGitManager(cfgDir)
		if err != nil {
			return fmt.Errorf(i18n.T("sync.err.not_init"), err)
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
		fmt.Println(i18n.T("sync.push.success"))
		return nil
	},
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: i18n.T("sync.pull.short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		mgr, err := gitpkg.NewGitManager(cfgDir)
		if err != nil {
			return fmt.Errorf(i18n.T("sync.err.not_init_bare"), err)
		}
		appCfg, err := config.LoadAppConfig()
		if err != nil {
			return err
		}
		if err := mgr.Pull("origin", appCfg.Git.Branch, nil); err != nil {
			return err
		}
		fmt.Println(i18n.T("sync.pull.success"))
		return nil
	},
}

var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: i18n.T("sync.status.short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		mgr, err := gitpkg.NewGitManager(cfgDir)
		if err != nil {
			return fmt.Errorf(i18n.T("sync.err.not_init_bare"), err)
		}
		status, err := mgr.Status()
		if err != nil {
			return err
		}
		if status == "" {
			fmt.Println(i18n.T("sync.status.clean"))
		} else {
			fmt.Println(status)
		}
		return nil
	},
}

var syncInitCmd = &cobra.Command{
	Use:   "init",
	Short: i18n.T("sync.init.short"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgDir, err := config.EnsureConfigDir()
		if err != nil {
			return err
		}
		if _, err := gitpkg.Init(cfgDir); err != nil {
			return err
		}
		fmt.Printf(i18n.T("sync.init.success")+"\n", cfgDir)
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
