package cmd

import (
	"strings"

	"github.com/C123R/helm-blob/pkg/repo"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch <blob-url>",
	Short: "fetches file from blob",
	Args:  cobra.MinimumNArgs(1),
	Long: `
Fetch command fetches file from remote chart and prints it on standard output.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// workaround: "helm repo add" commands sends url
		// with index.yaml path appended in it
		repoUrl := strings.TrimRight(args[0], "index.yaml")
		r, err := repo.NewRepoByUrl(repoUrl)
		if err != nil {
			return err
		}
		if err = r.Fetch(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
