package main

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Displays files and folders as a tree",
	Long:  "Shows a hierarchical tree of files and folders in the bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		fsys, ok := cmd.Context().Value(contextFSKey).(fs.FS)
		if !ok {
			return fmt.Errorf("failed to get fs from context")
		}

		rootDir := "."
		if len(args) > 0 {
			rootDir = args[0]
		}

		curDir := rootDir
		tree := treeprint.New()
		cur := tree

		err := fs.WalkDir(fsys, rootDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if path == "." {
				return nil
			}

			if d.IsDir() {
				if !strings.Contains(path, curDir) {
					cur = tree
				}
				curDir = path

				// add dir branch
				cur = cur.AddBranch(d.Name())
				return nil
			}

			// add file node
			cur.AddNode(d.Name())
			return nil
		})

		if err != nil {
			return err
		}

		fmt.Println(tree.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}
