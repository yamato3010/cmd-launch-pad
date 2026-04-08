package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourname/cmd-launch-pad/internal/models"
	"github.com/yourname/cmd-launch-pad/internal/repository"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "コマンドをCLIから追加する",
	Long:  `CLIからコマンドランチャーにコマンドを追加します。`,
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
			return fmt.Errorf("--name と --command は必須です")
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
			return fmt.Errorf("コマンドの追加に失敗しました: %w", err)
		}

		fmt.Printf("✅ コマンドを追加しました: %s (%s)\n", name, command)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("name", "", "コマンド名 (必須)")
	addCmd.Flags().String("command", "", "実行コマンド (必須)")
	addCmd.Flags().String("category", "custom", "カテゴリID")
	addCmd.Flags().String("desc", "", "コマンドの説明")
	addCmd.Flags().String("icon", "⚡", "アイコン (絵文字)")
	addCmd.Flags().String("args", "", "引数 (スペース区切り)")
}
