package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

		r := mux.NewRouter()
		r.PathPrefix("/").Handler(http.FileServer(http.FS(GCSFS)))

		addr := fmt.Sprintf(":%d", port)
		fmt.Println("Server listening to " + addr)
		return http.ListenAndServe(addr, r)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Http server port to listen on")
}
