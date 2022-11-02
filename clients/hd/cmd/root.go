package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName      = "hotdoggi.es"
	appNameShort = "hd"
)

var rootCmd = &cobra.Command{
	Use:   appNameShort,
	Short: fmt.Sprintf("interact with %s.", appName),
	Long:  fmt.Sprintf("interact with %s.", appName),
}

func Execute() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fail(err)
	}

	configPath := filepath.Join(dirname, ".config", appName)
	configType := "yaml"
	viper.AddConfigPath(configPath)
	viper.SetConfigName(appName)
	viper.SetConfigType(configType)

	err = os.Mkdir(configPath, fs.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		fail(err)
	}

	_, err = os.Stat(filepath.Join(configPath, fmt.Sprintf("%s.%s", appName, configType)))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fail(err)
	}
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// Create default config
		viper.Set("google.clientID", "")
		viper.Set("google.clientSecret", "")
		viper.Set("", "")
		viper.Set("oauthStateString", "")
		err = viper.WriteConfig()
		if err != nil {
			fail(err)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		fail(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Println(err)
	os.Exit(1)
}
