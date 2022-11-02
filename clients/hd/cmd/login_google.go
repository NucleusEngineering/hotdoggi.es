package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var loginGoogleCmd = &cobra.Command{
	Use:   "google",
	Short: fmt.Sprintf("use Google to login to %s.", appName),
	Long:  fmt.Sprintf("use Google to login to %s.", appName),
	Run: func(cmd *cobra.Command, args []string) {
		accessTokenGoogle()
	},
}

func accessTokenGoogle() {
	config := &oauth2.Config{
		ClientID:     viper.GetString("google.clientID"),
		ClientSecret: viper.GetString("google.clientSecret"),
		RedirectURL:  fmt.Sprintf("http://localhost:%s", callbackPort),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	codeChan := make(chan string, 1)
	go codeViaTerminal(authURL, codeChan)
	go codeViaCallback(authURL, codeChan)

	authCode := <-codeChan

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	viper.Set("token.issuer", "google")
	viper.Set("token.access", token.AccessToken)
	viper.Set("token.expiry", token.Expiry)
	viper.Set("token.refresh", token.RefreshToken)
	viper.Set("token.type", token.TokenType)
	err = viper.WriteConfig()
	if err != nil {
		fail(err)
	}
}

func init() {
	loginCmd.AddCommand(loginGoogleCmd)
}
