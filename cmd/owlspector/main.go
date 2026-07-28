package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/yannisobert/git-health-check/server"
)

func main() {
	_ = godotenv.Load()

	rootCmd := &cobra.Command{
		Use:   "owlspector",
		Short: "Audit the health of public GitHub repositories",
	}

	rootCmd.AddCommand(newCheckCmd())

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
			fmt.Println("owlspector dev")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
