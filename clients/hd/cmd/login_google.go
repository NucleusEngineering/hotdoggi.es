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
	config := configGoogle()
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

	// TODO need id_token as described in https://developers.google.com/identity/openid-connect/openid-connect#exchangecode
	// Probably with an OIDC lib

	err = viper.WriteConfig()
	if err != nil {
		fail(err)
	}
}

func configGoogle() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("google.clientID"),
		ClientSecret: viper.GetString("google.clientSecret"),
		RedirectURL:  fmt.Sprintf("http://localhost:%s", callbackPort),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func init() {
	loginCmd.AddCommand(loginGoogleCmd)
}
