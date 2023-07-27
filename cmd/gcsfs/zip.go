package main

import (
	"archive/zip"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Creates a zip archive with contents of the filesystem",
	Long:  "Walks the directory structure provided by the bucket, adding files to the zip archive",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing zip filename")
		}

		err := bucketSetupFunc(cmd, args)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(args[0], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer file.Close()

		zw := zip.NewWriter(file)
		defer zw.Close()

		return zw.AddFS(GCSFS)
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
}
