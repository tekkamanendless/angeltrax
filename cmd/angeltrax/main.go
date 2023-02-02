package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tekkamanendless/angeltrax/angeltrax"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
}

func main() {
	var defaultConfigFilename string
	{
		userConfigDirectory, _ := os.UserConfigDir()
		if userConfigDirectory != "" {
			userConfigDirectory = userConfigDirectory + string(os.PathSeparator) + "angeltrax"
			_ = os.Mkdir(userConfigDirectory, 0755)

			defaultConfigFilename = userConfigDirectory + string(os.PathSeparator) + "config.json"
		}
	}

	var configFilename string
	var debug bool

	ctx := context.Background()
	var client angeltrax.Client

	rootCmd := cobra.Command{
		Use: "angeltrax",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
			}

			if configFilename != "" {
				_, err := os.Stat(configFilename)
				if errors.Is(err, os.ErrNotExist) {
					// Okay; no such file.
				} else if err != nil {
					logrus.Errorf("Could not stat config file %q: %v", configFilename, err)
					os.Exit(1)
				} else {
					var config Config
					contents, err := os.ReadFile(configFilename)
					if err != nil {
						logrus.Errorf("Could not read config file %q: %v", configFilename, err)
						os.Exit(1)
					}
					err = json.Unmarshal(contents, &config)
					if err != nil {
						logrus.Errorf("Could not parse config file %q: %v", configFilename, err)
						os.Exit(1)
					}
					client.Server = config.Server
					client.Username = config.Username
					client.Password = config.Password
				}
			}
		},
	}
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable this to show more verbose logging.")
	rootCmd.PersistentFlags().StringVar(&configFilename, "config-file", defaultConfigFilename, "Enable this to show output as JSON.")

	{
		var server string
		var username string
		var password string
		cmd := &cobra.Command{
			Use:   "login",
			Short: "Login",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if client.Server == "" && server == "" {
					logrus.Errorf("Missing server.")
					os.Exit(1)
				}
				if client.Username == "" && username == "" {
					logrus.Errorf("Missing username.")
					os.Exit(1)
				}
				if client.Password == "" && password == "" {
					logrus.Errorf("Missing password.")
					os.Exit(1)
				}

				if configFilename != "" {
					var config Config
					config.Server = server
					config.Username = username
					config.Password = password

					contents, err := json.Marshal(config)
					if err != nil {
						logrus.Errorf("Could not marshal config contents: %v", err)
						os.Exit(1)
					}
					err = os.WriteFile(configFilename, contents, 0644)
					if err != nil {
						logrus.Errorf("Could not read config file %q: %v", configFilename, err)
						os.Exit(1)
					}
				}

				err := client.Login(ctx, server, username, password)
				if err != nil {
					logrus.Errorf("Error: [%T] %v", err, err)
					os.Exit(1)
				}
			},
		}
		cmd.Flags().StringVar(&server, "server", "", "The server")
		cmd.Flags().StringVar(&username, "username", "", "The username")
		cmd.Flags().StringVar(&password, "password", "", "The password")
		rootCmd.AddCommand(cmd)
	}

	err := rootCmd.Execute()
	if err != nil {
		logrus.Errorf("Error: [%T] %v", err, err)
		os.Exit(1)
	}
}
