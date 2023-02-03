package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tekkamanendless/angeltrax/angeltrax"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Key      string `json:"key"`
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

	loginOrFail := func() {
		if client.Server == "" {
			logrus.Errorf("Missing server.")
			os.Exit(1)
		}
		if client.Username == "" {
			logrus.Errorf("Missing username.")
			os.Exit(1)
		}
		if client.Password == "" {
			logrus.Errorf("Missing password.")
			os.Exit(1)
		}

		err := client.Login(ctx, client.Server, client.Username, client.Password)
		if err != nil {
			logrus.Errorf("Could not log in: [%T] %v", err, err)
			os.Exit(1)
		}
	}

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
					client.Key = config.Key
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
				if server == "" {
					server = client.Server
				}
				if username == "" {
					username = client.Username
				}
				if password == "" {
					password = client.Password
				}
				if server == "" {
					logrus.Errorf("Missing server.")
					os.Exit(1)
				}
				if username == "" {
					logrus.Errorf("Missing username.")
					os.Exit(1)
				}
				if password == "" {
					logrus.Errorf("Missing password.")
					os.Exit(1)
				}

				err := client.Login(ctx, server, username, password)
				if err != nil {
					logrus.Errorf("Error: [%T] %v", err, err)
					os.Exit(1)
				}

				if configFilename != "" {
					var config Config
					config.Server = client.Server
					config.Username = client.Username
					config.Password = client.Password
					config.Key = client.Key

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
			},
		}
		cmd.Flags().StringVar(&server, "server", "", "The server")
		cmd.Flags().StringVar(&username, "username", "", "The username")
		cmd.Flags().StringVar(&password, "password", "", "The password")
		rootCmd.AddCommand(cmd)
	}

	{
		cmd := &cobra.Command{
			Use:   "center-groups",
			Short: "",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				loginOrFail()

				getCenterGroupsResponse, err := client.GetCenterGroups(ctx)
				if err != nil {
					logrus.Errorf("Error: [%T] %v", err, err)
					os.Exit(1)
				}

				getCenterDevicesResponse, err := client.GetCenterDevices(ctx)
				if err != nil {
					logrus.Errorf("Error: [%T] %v", err, err)
					os.Exit(1)
				}

				for _, group := range getCenterGroupsResponse.Data {
					var groupName string
					{
						currentGroup := &group
						for currentGroup != nil {
							if groupName == "" {
								groupName = currentGroup.GroupName
							} else {
								groupName = currentGroup.GroupName + "/" + groupName
							}

							parentID := currentGroup.GroupFatherID
							currentGroup = nil
							for _, g := range getCenterGroupsResponse.Data {
								if g.GroupID == parentID {
									currentGroup = &g
									break
								}
							}
						}
					}

					fmt.Printf("Group #%d: %s\n", group.GroupID, groupName)
					for _, device := range getCenterDevicesResponse.Data {
						if device.GroupID != group.GroupID {
							continue
						}
						fmt.Printf("   Device #%s: %s | Channels: %d\n", device.DeviceID, device.CarLicense, device.ChannelCount)
					}
				}
			},
		}
		rootCmd.AddCommand(cmd)
	}

	{
		cmd := &cobra.Command{
			Use:   "tasks",
			Short: "",
			Args:  cobra.ExactArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				loginOrFail()

				getCenterDevicesResponse, err := client.GetCenterDevices(ctx)
				if err != nil {
					logrus.Errorf("Error: [%T] %v", err, err)
					os.Exit(1)
				}

				_, err = client.RegisterLogin(ctx)
				if err != nil {
					logrus.Errorf("Error: [%T] %v", err, err)
					os.Exit(1)
				}

				for _, device := range getCenterDevicesResponse.Data {
					logrus.Infof("Device: %s (%s)", device.DeviceID, device.CarLicense)
					output, err := client.MonitorAutoDownload(ctx, device.DeviceID)
					if err != nil {
						logrus.Errorf("Error: [%T] %v", err, err)
						os.Exit(1)
					}
					logrus.Infof("Total: %d", output.Total)
				}
			},
		}
		rootCmd.AddCommand(cmd)
	}

	{
		var service string
		var method string
		cmd := &cobra.Command{
			Use:   "raw-request [--service ${service}] [--method ${method}] ${path}",
			Short: "Perform a raw request",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				loginOrFail()

				path := args[0]

				var output json.RawMessage
				err := client.RawServiceRequest(ctx, service, method, path, nil, nil, &output)
				if err != nil {
					logrus.Errorf("Error: [%T] %v", err, err)
					os.Exit(1)
				}
				fmt.Printf("%s\n", output)
			},
		}
		cmd.Flags().StringVar(&service, "service", "webclient", "The service")
		cmd.Flags().StringVar(&method, "method", http.MethodGet, "The method")
		rootCmd.AddCommand(cmd)
	}

	err := rootCmd.Execute()
	if err != nil {
		logrus.Errorf("Error: [%T] %v", err, err)
		os.Exit(1)
	}
}
