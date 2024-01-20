package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve bucket files using an http server",
	Long:  "Creates an HTTP server to serve the contents of the bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		fsys := cmd.Context().Value(contextFSKey).(fs.FS)

		http.Handle("/", http.FileServer(http.FS(fsys)))

		addr := fmt.Sprintf(":%d", port)
		fmt.Println("Server listening to " + addr)
		return http.ListenAndServe(addr, nil)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Http server port to listen on")
}
