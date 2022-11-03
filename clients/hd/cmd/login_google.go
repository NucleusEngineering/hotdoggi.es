package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
		authenticateGoogle()
	},
}

func authenticateGoogle() {
	config := configGoogle()
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	codeChan := make(chan string, 1)
	go codeViaTerminalGoogle(authURL, codeChan)
	go codeViaCallbackGoogle(authURL, codeChan)
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
	viper.Set("token.identity", token.Extra("id_token").(string))

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

func codeViaTerminalGoogle(url string, codeChan chan string) {
	fmt.Printf("Go to the following link: \n%v\n\n", url)
	fmt.Printf("Paste back the code you received in this terminal: \n")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
	codeChan <- authCode
}

func codeViaCallbackGoogle(url string, codeChan chan string) {
	server := &http.Server{Addr: fmt.Sprintf(":%s", callbackPort)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			codeChan <- r.URL.Query().Get("code")
			err := server.Shutdown(context.Background())
			if err != nil {
				fail(err)
			}
		}
	})

	err := openBrowser(url)
	if err != nil {
		fail(err)
	}
	server.ListenAndServe()
}

func init() {
	loginCmd.AddCommand(loginGoogleCmd)
}
