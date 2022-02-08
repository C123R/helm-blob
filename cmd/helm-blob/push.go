package cmd

import (
	"os"
	"path/filepath"

	"github.com/PTC-Global/helm-blob/pkg/repo"
	"github.com/spf13/cobra"
)

var force bool

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push <chart-name> <repo-name>",
	Short: "Uploads a chart into a repository",
	Args:  cobra.ExactValidArgs(2),
	Long: `
Push command uploads a chart into a remote repository and merge
with remote index file.
`,
	Example: `helm blob push sample-chart.tgz sample-repository`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var charts []string
		chartpath, repositoryName := args[0], args[1]
		r, err := repo.NewRepoByName(repositoryName)
		if err != nil {
			return err
		}
		f, err := os.Stat(chartpath)
		switch {
		case err != nil:
			return err
		case f.IsDir():
			charts, err = filepath.Glob(filepath.Join(chartpath, "*.tgz"))
			if err != nil {
				return err
			}
		default:
			charts = append(charts, chartpath)
		}

		for _, chart := range charts {
			if err = r.Push(chart, force); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVarP(&force, "force", "f", false, "force upload/replace existing chart")
}
