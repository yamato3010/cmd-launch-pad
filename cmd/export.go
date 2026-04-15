package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yamato3010/cmd-launch-pad/internal/i18n"
	"github.com/yamato3010/cmd-launch-pad/internal/repository"
	"gopkg.in/yaml.v3"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: i18n.T("export.short"),
	Long:  i18n.T("export.long"),
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
			return fmt.Errorf(i18n.T("export.err.serialize"), err)
		}

		if output != "" {
			if err := os.WriteFile(output, out, 0644); err != nil {
				return fmt.Errorf(i18n.T("export.err.write"), err)
			}
			fmt.Printf(i18n.T("export.success")+"\n", output)
		} else {
			fmt.Print(string(out))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("output", "o", "", i18n.T("export.flag.output"))
}
