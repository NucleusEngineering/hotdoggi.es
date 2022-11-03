package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	callbackPort = "8934"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: fmt.Sprintf("login to %s.", appName),
	Long:  fmt.Sprintf("login to %s.", appName),
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
