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
		if len(args) != 1 {
			return cmd.Usage()
		}

		fsys := cmd.Context().Value(contextFSKey).(fs.FS)
		files, err := fs.ReadDir(fsys, args[0])
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
