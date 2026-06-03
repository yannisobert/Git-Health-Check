package main

import (
	"fmt"
	"os"

	"github.com/yannisobert/git-health-check/server"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ghhealth",
		Short: "Audit the health of public GitHub repositories",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server.New().Run()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ghhealth dev")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
