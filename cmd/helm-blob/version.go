package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current helm-blob version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		byte, err := ioutil.ReadFile("VERSION")
		if err != nil {
			return fmt.Errorf("Error getting version of helm-blob")
		}
		fmt.Printf("Helm Blob Plugin Version: %s\n", string(byte))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
