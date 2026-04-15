package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/models"
	"github.com/yamato3010/cmd-launch-pad/internal/repository"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: i18n.T("add.short"),
	Long:  i18n.T("add.long"),
	Example: `  clp add --name "Neovim" --command "nvim" --category "editor" --desc "テキストエディタ"
  clp add --name "lazygit" --command "lazygit" --category "git" --icon "🌿"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		command, _ := cmd.Flags().GetString("command")
		category, _ := cmd.Flags().GetString("category")
		desc, _ := cmd.Flags().GetString("desc")
		icon, _ := cmd.Flags().GetString("icon")
		argsStr, _ := cmd.Flags().GetString("args")

		if name == "" || command == "" {
			return fmt.Errorf("%s", i18n.T("add.err.required"))
		}

		var cmdArgs []string
		if argsStr != "" {
			cmdArgs = strings.Fields(argsStr)
		}

		repo, err := repository.NewCommandRepository()
		if err != nil {
			return err
		}

		newCmd := &models.Command{
			Name:        name,
			Command:     command,
			Args:        cmdArgs,
			Description: desc,
			CategoryID:  category,
			Icon:        icon,
		}

		if err := repo.AddCommand(newCmd); err != nil {
			return fmt.Errorf("%s: %w", i18n.T("add.err.failed"), err)
		}

		fmt.Printf(i18n.T("add.success")+"\n", name, command)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("name", "", i18n.T("add.flag.name"))
	addCmd.Flags().String("command", "", i18n.T("add.flag.command"))
	addCmd.Flags().String("category", "custom", i18n.T("add.flag.category"))
	addCmd.Flags().String("desc", "", i18n.T("add.flag.desc"))
	addCmd.Flags().String("icon", "⚡", i18n.T("add.flag.icon"))
	addCmd.Flags().String("args", "", i18n.T("add.flag.args"))
}
