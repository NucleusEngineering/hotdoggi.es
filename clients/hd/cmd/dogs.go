package cmd

import (
	"github.com/spf13/cobra"
)

var dogsCmd = &cobra.Command{
	Use:   "dogs",
	Short: "interact with dogs.",
	Long:  `interact with dogs.`,
}

func init() {
	rootCmd.AddCommand(dogsCmd)
}
