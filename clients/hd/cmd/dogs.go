package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dogsCmd = &cobra.Command{
	Use:   "version",
	Short: "print version information.",
	Long:  `print version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", appNameShort, version)
	},
}

func init() {
	rootCmd.AddCommand(dogsCmd)
}
