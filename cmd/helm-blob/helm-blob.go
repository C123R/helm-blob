package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blob",
	Short: "manage helm repositories on Cloud Storage Blobs",
	Long: `
Blob plugin supports operations like pushing or deletion
of charts from remote Helm Repository hosted on Blob Storage.

This could be use to initilize new Helm Repository on Blob Storage.

Currently, Blob plugin supports GCS, S3 and Azure Blob storage.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
