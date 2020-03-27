package cmd

import (
	"github.com/C123R/helm-blob/pkg/repo"
	"github.com/spf13/cobra"
)

var force bool

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push <chart-name> <repo-name>",
	Short: "uploads a chart into a repository",
	Args:  cobra.ExactValidArgs(2),
	Long: `
Push command uploads a chart into a remote repository.
`,
	Example: `helm blob push sample-chart.tgz sample-repository`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chart, repositoryName := args[0], args[1]
		r, err := repo.NewRepoByName(repositoryName)
		if err != nil {
			return err
		}
		if err = r.Push(chart, force); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVarP(&force, "force", "f", false, "force upload/replace existing chart")
}
