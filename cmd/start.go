package cmd

import (
	"appscode/fr-sarker/golang-chi-crud-api/apiHandler"
	"fmt"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var (
	// Port stores port number for starting a connection
	Port     int
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "start cmd starts the server on a port",
		Long: `It starts the server on a given port number, 
				Port number will be given in the cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Port)
			apiHandler.RunServer(Port)
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().IntVarP(&Port, "port", "p", 3000, "Port number for starting server")

}
