package cmd

import (
	"github.com/C123R/helm-blob/pkg/repo"
	"github.com/spf13/cobra"
)

var version string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <chart-name> <repo-name>",
	Short: "Deletes a chart from repository",
	Args:  cobra.MinimumNArgs(2),
	Long: `
Delete command deletes a chart from a remote repository.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chartName, repositoryName := args[0], args[1]
		r, err := repo.NewRepoByName(repositoryName)
		if err != nil {
			return err
		}
		if err = r.Delete(chartName, version); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&version, "version", "v", "", "version of the chart to delete")
}
