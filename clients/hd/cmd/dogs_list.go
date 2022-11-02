package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"

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
	token, config, err := getTokenWithConfig()
	if err != nil {
		fail(err)
	}

	client := config.Client(context.Background(), token)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fail(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		fail(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fail(fmt.Errorf("received HTTP %d", resp.StatusCode))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fail(err)
	}
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)
}

func init() {
	dogsCmd.AddCommand(dogsListenCmd)
}
