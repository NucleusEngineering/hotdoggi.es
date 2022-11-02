package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	callbackPort = "8934"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: fmt.Sprintf("login to %s.", appName),
	Long:  fmt.Sprintf("login to %s.", appName),
}

func getTokenWithConfig() (*oauth2.Token, *oauth2.Config, error) {
	token := new(oauth2.Token)
	token.AccessToken = viper.GetString("token.access")
	token.RefreshToken = viper.GetString("token.refresh")
	token.TokenType = viper.GetString("token.type")
	token.Expiry = viper.GetTime("token.expiry")

	issuer := viper.GetString("token.issuer")

	if issuer == "" {
		return nil, nil, fmt.Errorf("no credentials found. Try running '%s login'", appNameShort)
	}

	switch issuer {
	case "google":
		return token, configGoogle(), nil
	case "github":
		return nil, nil, fmt.Errorf("unimplemented token issuer: github")
	default:
		return nil, nil, fmt.Errorf("invalid token issuer: %s", issuer)
	}
}

func codeViaTerminal(url string, codeChan chan string) {
	fmt.Printf("Go to the following link: \n%v\n\n", url)
	fmt.Printf("Paste back the code you received in this terminal: \n")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
	codeChan <- authCode
}

func codeViaCallback(url string, codeChan chan string) {
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

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
