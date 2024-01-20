package main

import (
	"io"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "Concatenate files",
	Long:  "Concatenate files from a Google Storage Bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		fsys := cmd.Context().Value(contextFSKey).(fs.FS)

		for _, filename := range args {
			f, err := fsys.Open(filename)
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
