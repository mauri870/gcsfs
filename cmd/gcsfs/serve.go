package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve bucket files using an http server",
	Long:  "Creates an HTTP server to serve the contents of the bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bucketSetupFunc(cmd, args)
		if err != nil {
			return err
		}

		http.Handle("/", http.FileServer(http.FS(GCSFS)))

		addr := fmt.Sprintf(":%d", port)
		fmt.Println("Server listening to " + addr)
		return http.ListenAndServe(addr, nil)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Http server port to listen on")
}
