package main

import (
	"context"
	"os"

	gcs "cloud.google.com/go/storage"
	"github.com/mauri870/gcsfs"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

type contextKey string

var (
	GCSFS        *gcsfs.FS
	contextFSKey contextKey = "fs"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "gcsfs",
	Short:        "io/fs.FS interface to GCS",
	Long:         "Interacts with files inside a Google Storage Bucket using Golang's io/fs.FS",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		bucket, err := cmd.Flags().GetString("bucket")
		if err != nil {
			return err
		}

		var opts []option.ClientOption

		if cmd.Flags().Changed("without-authentication") {
			opts = append(opts, option.WithoutAuthentication())
		}

		gcsClient, err := gcs.NewClient(cmd.Context(), opts...)
		if err != nil {
			return err
		}

		fsys := gcsfs.NewWithClient(gcsClient, bucket)

		ctx := context.WithValue(cmd.Context(), contextFSKey, fsys)
		cmd.SetContext(ctx)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("bucket", "b", "", "Bucket name to use")
	rootCmd.MarkFlagRequired("bucket")

	rootCmd.PersistentFlags().BoolP("without-authentication", "w", false, "Disables authentication. Useful to access public buckets")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
