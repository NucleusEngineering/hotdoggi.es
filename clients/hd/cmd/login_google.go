package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	providerIDGoogle  = "google.com"
	mimeTypeAPIGoogle = "application/json"

	authURIEndpointGoogle = "https://identitytoolkit.googleapis.com/v1/accounts:createAuthUri"
	signInEndpointGoogle  = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithIdp"
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
	authURI, sessionID, err := getAuthUriGoogle()
	if err != nil {
		fail(err)
	}

	authResponseChan := make(chan string, 1)
	go codeViaCallbackGoogle(authURI, authResponseChan)
	authResponse := <-authResponseChan

	fmt.Println(authResponse)

	exchangeResponse, err := exchangeTokenGoogle(authResponse, sessionID)
	if err != nil {
		fail(err)
	}

	fmt.Println(exchangeResponse)

	// viper.Set("token.issuer", "google")
	// viper.Set("token.access", token.AccessToken)
	// viper.Set("token.expiry", token.Expiry)
	// viper.Set("token.refresh", token.RefreshToken)
	// viper.Set("token.type", token.TokenType)
	// viper.Set("token.identity", token.Extra("id_token").(string))

	err = viper.WriteConfig()
	if err != nil {
		fail(err)
	}
}

func exchangeTokenGoogle(authResponse string, sessionID string) (string, error) {
	apiKey := viper.GetString("google.apikey")

	payload := []byte(
		fmt.Sprintf(`{
				"requestUri":"%s",
				"postBody": "%s", 
				"sessionId":"%s",
				"returnRefreshToken": true, 
				"returnSecureToken": true, 
				"returnIdpCredential": true
			}`, redirectURI, authResponse, sessionID),
	)
	payloadReader := bytes.NewReader(payload)

	url := fmt.Sprintf("%s?key=%s", signInEndpointGoogle, apiKey)

	req, err := http.NewRequest(http.MethodPost, url, payloadReader)
	if err != nil {
		return "", err
	}

	req.Header.Set("content-type", mimeTypeAPIGoogle)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	body := string(bodyBytes)

	return body, nil
}

func getAuthUriGoogle() (string, string, error) {
	apiKey := viper.GetString("google.apikey")

	payload := []byte(
		fmt.Sprintf(`{
			"providerId":"%s",
			"continueUri": "%s"
		}`, providerIDGoogle, redirectURI),
	)
	payloadReader := bytes.NewReader(payload)

	url := fmt.Sprintf("%s?key=%s", authURIEndpointGoogle, apiKey)

	req, err := http.NewRequest(http.MethodPost, url, payloadReader)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("content-type", mimeTypeAPIGoogle)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return "", "", err
	}

	authURI := data["authUri"].(string)
	sessionID := data["sessionId"].(string)

	return authURI, sessionID, nil
}

func codeViaCallbackGoogle(url string, authResponseChan chan string) {
	server := &http.Server{Addr: fmt.Sprintf(":%s", callbackPort)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// TODO this needs to return client side code to extract and forward the Hash URI #state as POST on the same handler
			state := r.URL.Query().Get("state")
			if state != "" {
				authResponseChan <- state
				w.WriteHeader(200)
				w.Write([]byte("200 OK"))
				return
			}
			// Garbage until here more or less?
		}
		if r.Method == http.MethodPost {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte("400 BAD REQUEST"))
				return
			}
			body := string(bodyBytes)
			if body != "" {
				authResponseChan <- body
				w.WriteHeader(200)
				w.Write([]byte("200 OK"))
				return
			}
		}

		w.WriteHeader(400)
		w.Write([]byte("400 BAD REQUEST"))
	})

	go func() {
		err := openBrowser(url)
		if err != nil {
			fail(err)
		}
	}()

	server.ListenAndServe()
}

func init() {
	loginCmd.AddCommand(loginGoogleCmd)
}
