package cmd

import (
	"strings"

	"github.com/PTC-Global/helm-blob/pkg/repo"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch <blob-url>",
	Short: "Fetches file from blob",
	Args:  cobra.MinimumNArgs(1),
	Long: `
Fetch command fetches file from remote chart and prints it on standard output.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		urlParts := strings.Split(args[0], "/")
		file2download := urlParts[len(urlParts)-1]
		repoUrl := strings.TrimRight(args[0], file2download)

		r, err := repo.NewRepoByUrl(repoUrl)
		if err != nil {
			return err
		}
		if err = r.Fetch(file2download); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
