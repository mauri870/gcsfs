package main

import (
	"os"

	"github.com/mauri870/gcsfs"
	"github.com/spf13/cobra"
)

var GCSFS *gcsfs.FS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "gcsfs",
	Short:        "io/fs.FS interface to GCS",
	Long:         "Interacts with files inside a Google Storage Bucket using Golang's io/fs.FS",
	SilenceUsage: true,
}

func bucketSetupFunc(cmd *cobra.Command, args []string) error {
	bucket, err := cmd.Flags().GetString("bucket")
	if err != nil {
		return err
	}

	GCSFS, err = gcsfs.New(bucket)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().StringP("bucket", "b", "", "Bucket name to use")
	rootCmd.MarkFlagRequired("bucket")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
