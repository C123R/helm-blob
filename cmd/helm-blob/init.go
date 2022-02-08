package cmd

import (
	"fmt"

	"github.com/PTC-Global/helm-blob/pkg/repo"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <blob-url>",
	Short: "Initialize new repository",
	Args:  cobra.ExactValidArgs(1),
	Long: `
Init command will initialize a new helm repository on provided blob url.

Note: This command will not create new blob storage, moreover
it will just add empty index.yaml file.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoUrl := args[0]
		r, err := repo.NewRepoByUrl(repoUrl)
		if err != nil {
			return err
		}
		if err = r.Init(); err != nil {
			return err
		}
		fmt.Printf(`Successfully initialized %s.

Now you can add this repo using 'helm repo add <repo-name> %s'`, repoUrl, repoUrl)
		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
