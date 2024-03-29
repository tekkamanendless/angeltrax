package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

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
		groupCmd := &cobra.Command{
			Use:   "task",
			Short: "Task-related commands",
		}
		rootCmd.AddCommand(groupCmd)

		{
			var deviceID string
			var deviceName string
			cmd := &cobra.Command{
				Use:   "monitor",
				Short: "Monitor the tasks",
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
						logrus.Debugf("Device: %s (%s)", device.DeviceID, device.CarLicense)
						if deviceName != "" && device.CarLicense != deviceName {
							continue
						}
						if deviceID != "" && device.DeviceID != deviceID {
							continue
						}

						output, err := client.MonitorAutoDownload(ctx, angeltrax.MonitorAutoDownloadInput{DeviceID: device.DeviceID})
						if err != nil {
							logrus.Errorf("Error: [%T] %v", err, err)
							os.Exit(1)
						}
						logrus.Debugf("Total: %d", output.Total)
						for _, task := range output.Rows {
							fmt.Printf("Task %d: %s (%s): %s | %s %s - %s\n", task.TaskID, task.DeviceID, task.CarLicense, task.TaskName, task.Date, task.StartTime, task.EndTime)
						}
					}
				},
			}
			cmd.Flags().StringVar(&deviceID, "device-id", "", "The device ID (optional)")
			cmd.Flags().StringVar(&deviceName, "device-name", "", "The device name (optional)")
			groupCmd.AddCommand(cmd)
		}

		{
			var deviceID string
			var deviceName string
			var effectiveDays int
			var startDate string
			var endDate string
			var startTime string
			var endTime string
			var taskName string
			cmd := &cobra.Command{
				Use:   "create",
				Short: "Create a new task",
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

					input := angeltrax.CreateAutoDownloadTaskInput{
						TaskName:      taskName,
						StartExecute:  startDate,
						EndExecute:    endDate,
						StartTime:     startTime,
						EndTime:       endTime,
						EffectiveDays: effectiveDays,
						TaskChannels:  []int{},
						TaskType:      angeltrax.TaskTypeVideo,
						Period:        angeltrax.TaskPeriodOnce,
					}
					for _, device := range getCenterDevicesResponse.Data {
						logrus.Debugf("Device: %s (%s)", device.DeviceID, device.CarLicense)
						if deviceName != "" && device.CarLicense != deviceName {
							continue
						}
						if deviceID != "" && device.DeviceID != deviceID {
							continue
						}
						input.DeviceID = device.DeviceID
						for i := 0; i < device.ChannelCount; i++ {
							input.TaskChannels = append(input.TaskChannels, i+1)
						}
						break
					}
					if input.DeviceID == "" {
						logrus.Errorf("Could not find device.")
						os.Exit(1)
					}

					output, err := client.CreateAutoDownloadTask(ctx, input)
					if err != nil {
						logrus.Errorf("Error: [%T] %v", err, err)
						os.Exit(1)
					}
					fmt.Printf("Success: %t\n", output.Result)
				},
			}
			cmd.Flags().StringVar(&deviceID, "device-id", "", "The device ID (you may omit this if you use --device-name)")
			cmd.Flags().StringVar(&deviceName, "device-name", "", "The device name (you may omit this if you use --device-id)")
			cmd.Flags().IntVar(&effectiveDays, "effective-days", 7, "The effective days")
			cmd.Flags().StringVar(&startDate, "start-date", "", "The start date (yyyy-mm-dd)")
			cmd.Flags().StringVar(&endDate, "end-date", "", "The end date (yyyy-mm-dd)")
			cmd.Flags().StringVar(&startTime, "start-time", "", "The start time (hh:mm:ss)")
			cmd.Flags().StringVar(&endTime, "end-time", "", "The end time (hh:mm:ss)")
			cmd.Flags().StringVar(&taskName, "task-name", "", "The task name")
			// TODO: Add a flag for cameras (right now we just do them all).
			groupCmd.AddCommand(cmd)
		}

		{
			cmd := &cobra.Command{
				Use:   "monitor-task ${id}",
				Short: "Get the information about the given task",
				Args:  cobra.ExactArgs(1),
				Run: func(cmd *cobra.Command, args []string) {
					loginOrFail()

					taskID := args[0]

					_, err := client.RegisterLogin(ctx)
					if err != nil {
						logrus.Errorf("Error: [%T] %v", err, err)
						os.Exit(1)
					}

					output, err := client.MonitorAutoDownloadTask(ctx, taskID)
					if err != nil {
						logrus.Errorf("Error: [%T] %v", err, err)
						os.Exit(1)
					}
					fmt.Printf("Task #%d (%s) - %s (%s)\n", output.TaskID, output.TaskName, output.DeviceID, output.CarLicense)
					fmt.Printf("   %s - %s, %s - %s\n", output.StartExecute, output.EndExecute, output.StartTime, output.EndTime)
				},
			}
			groupCmd.AddCommand(cmd)
		}

		{
			var deviceID string
			var deviceName string
			var status int
			var startDate string
			var endDate string
			cmd := &cobra.Command{
				Use:   "global-report",
				Short: "Global report",
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
						logrus.Debugf("Device: %s (%s)", device.DeviceID, device.CarLicense)
						if deviceName != "" && device.CarLicense != deviceName {
							continue
						}
						if deviceID != "" && device.DeviceID != deviceID {
							continue
						}

						input := angeltrax.GlobalReportAutoDownloadInput{
							DeviceID:  device.DeviceID,
							Status:    angeltrax.TaskStatus(status),
							StartDate: startDate,
							EndDate:   endDate,
						}
						output, err := client.GlobalReportAutoDownload(ctx, input)
						if err != nil {
							logrus.Errorf("Error: [%T] %v", err, err)
							os.Exit(1)
						}
						for _, task := range output.Rows {
							fmt.Printf("Task #%d (%s) - %s (%s)\n", task.TaskID, task.TaskName, task.DeviceID, task.CarLicense)
							fmt.Printf("   %s - %s\n", task.StartTime, task.EndTime)
							fmt.Printf("   %s (by %s)\n", task.Status, task.Username)
						}
					}
				},
			}
			cmd.Flags().StringVar(&deviceID, "device-id", "", "The device ID (you may omit this if you use --device-name)")
			cmd.Flags().StringVar(&deviceName, "device-name", "", "The device name (you may omit this if you use --device-id)")
			cmd.Flags().IntVar(&status, "status", 0, "The status")
			cmd.Flags().StringVar(&startDate, "start-date", time.Now().Add(7*24*time.Hour).Format("2006-01-02"), "The start date (yyyy-mm-dd)")
			cmd.Flags().StringVar(&endDate, "end-date", time.Now().Format("2006-01-02"), "The end date (yyyy-mm-dd)")
			groupCmd.AddCommand(cmd)
		}

		{
			var deviceID string
			var date string
			cmd := &cobra.Command{
				Use:   "global-report-task ${id}",
				Short: "Get the information about the given task",
				Args:  cobra.ExactArgs(1),
				Run: func(cmd *cobra.Command, args []string) {
					loginOrFail()

					taskID := args[0]

					_, err := client.RegisterLogin(ctx)
					if err != nil {
						logrus.Errorf("Error: [%T] %v", err, err)
						os.Exit(1)
					}

					input := angeltrax.GlobalReportAutoDownloadTaskInput{
						DeviceID: deviceID,
						TaskID:   taskID,
						Date:     date,
					}
					output, err := client.GlobalReportAutoDownloadTask(ctx, input)
					if err != nil {
						logrus.Errorf("Error: [%T] %v", err, err)
						os.Exit(1)
					}
					for _, task := range output.Rows {
						fmt.Printf("%s, %s of %s (%s%%)\n", task.FileSource, task.CurrentSize, task.TotalSize, task.Percent)
						fmt.Printf("   %s, %s - %s\n", task.Date, task.StartTime, task.EndTime)
						fmt.Printf("   status: %s\n", task.Status)
						fmt.Printf("   speed: %s\n", task.Speed)
					}
				},
			}
			cmd.Flags().StringVar(&deviceID, "device-id", "", "The device ID (optional)")
			cmd.Flags().StringVar(&date, "date", "", "The date (yyyy-mm-dd) (optional)")
			groupCmd.AddCommand(cmd)
		}
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
