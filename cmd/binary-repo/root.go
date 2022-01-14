package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
}

var rootCmd = &cobra.Command{
	Use:   "binaryrepo",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains`,
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
