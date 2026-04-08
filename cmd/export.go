package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourname/cmd-launch-pad/internal/repository"
	"gopkg.in/yaml.v3"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "コマンド定義をYAML形式でエクスポート",
	Long:  `登録済みコマンドをYAML形式で標準出力に出力します。`,
	Example: `  clp export > my-commands.yaml
  clp export --output my-commands.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")

		repo, err := repository.NewCommandRepository()
		if err != nil {
			return err
		}

		commands, err := repo.ListCommands()
		if err != nil {
			return err
		}
		categories, err := repo.ListCategories()
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"categories": categories,
			"commands":   commands,
		}

		out, err := yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("YAMLシリアライズに失敗しました: %w", err)
		}

		if output != "" {
			if err := os.WriteFile(output, out, 0644); err != nil {
				return fmt.Errorf("ファイル書き込みに失敗しました: %w", err)
			}
			fmt.Printf("✅ エクスポート完了: %s\n", output)
		} else {
			fmt.Print(string(out))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("output", "o", "", "出力ファイルパス (省略時は標準出力)")
}
