package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const endpoint = "https://api.hotdoggies.stamer.demo.altostrat.com/v1/dogs/"

var dogsListenCmd = &cobra.Command{
	Use:   "list",
	Short: "list all your dogs",
	Long:  `list all your dogs`,
	Run: func(cmd *cobra.Command, args []string) {
		listDogs()
	},
}

func listDogs() {
	fail(fmt.Errorf("not yet implemented"))
}

func init() {
	dogsCmd.AddCommand(dogsListenCmd)
}
