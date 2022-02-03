package main

import (
	"fmt"
	"io/fs"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files",
	Long:  "List files from a Google Storage Bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bucketSetupFunc(cmd, args)
		if err != nil {
			return err
		}

		if len(args) != 1 {
			return cmd.Usage()
		}

		files, err := fs.ReadDir(GCSFS, args[0])
		if err != nil {
			return err
		}

		for _, file := range files {
			fmt.Print(file.Name())
			if file.IsDir() {
				fmt.Print("/")
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
