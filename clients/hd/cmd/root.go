//  Copyright 2022 Google

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

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
		viper.Set("google.apikey", "")
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
