package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "Concatenate files",
	Long:  "Concatenate files from a Google Storage Bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bucketSetupFunc(cmd, args)
		if err != nil {
			return err
		}

		for _, filename := range args {
			f, err := GCSFS.Open(filename)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(os.Stdout, f)
			if err != nil {
				return err
			}

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(catCmd)
}
