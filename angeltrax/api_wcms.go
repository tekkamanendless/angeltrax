package angeltrax

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RegisterLoginResponse struct {
	Code   int  `json:"Code"`
	Result bool `json:"Result"`
}

type MonitorAutoDownloadResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		TaskID      int    `json:"TaskID"`
		Status      int    `json:"Status"`
		DeviceID    string `json:"Device"`
		CarLicense  string `json:"Carlicense"`
		TaskName    string `json:"TaskName"`
		Period      int    `json:"Period"`
		TaskType    int    `json:"TaskType"`
		Date        string `json:"Date"`       // yyyy-mm-dd
		StartTime   string `json:"StartTime"`  // hh:mm:ss
		EndTime     string `json:"EndTime"`    // hh:mm:ss
		ChannelList string `json:"Channel"`    // CSV of channel numbers, starting from "1".
		CreateTime  string `json:"CreateTime"` // yyyy-mm-dd hh:mm:ss

		NetMode string `json:"NetMode"`
	} `json:"rows"`
}

type QueryAutoDownloadResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		TaskID      int    `json:"TaskID"`
		Status      int    `json:"Status"`
		DeviceID    string `json:"Device"`
		CarLicense  string `json:"Carlicense"`
		TaskName    string `json:"TaskName"`
		Period      int    `json:"Period"`
		TaskType    int    `json:"TaskType"`
		Date        string `json:"Date"`       // yyyy-mm-dd
		StartTime   string `json:"StartTime"`  // hh:mm:ss
		EndTime     string `json:"EndTime"`    // hh:mm:ss
		ChannelList string `json:"Channel"`    // CSV of channel numbers, starting from "1".
		CreateTime  string `json:"CreateTime"` // yyyy-mm-dd hh:mm:ss

		FinishTime string `json:"FinishTime"` // yyyy-mm-dd hh:mm:ss
		Username   string `json:"UserName"`
	} `json:"rows"`
}

type QueryAutoDownloadTaskResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		DeviceID    string `json:"Device"`
		Status      int    `json:"Status"`
		Percent     string `json:"Percent"`     // This appears to be a float encoded as a string.
		Speed       string `json:"Speed"`       // This appears to be a float encoded as a string.
		Date        string `json:"Date"`        // yyyy-mm-dd
		StartTime   string `json:"StartTime"`   // hh:mm:ss
		EndTime     string `json:"EndTime"`     // hh:mm:ss
		TotalSize   string `json:"TotalTime"`   // This appears to be a float encoded as a string.
		CurrentSize string `json:"CurrentTime"` // This appears to be a float encoded as a string.
		Channel     int    `json:"Channel"`     // The one-index of the channel.
		Error       string `json:"Error"`
		TaskID      int    `json:"TaskID"`
		FileSource  string `json:"FileSource"`
		PreAlarm    int    `json:"PreAlarm"`
		NextAlarm   int    `json:"NextAlarm"`
	} `json:"rows"`
}

type CreateAutoDownloadTaskInput struct {
	TaskName  string
	DeviceID  string
	StartTime string // hh:mm:ss
	EndTime   string // hh:mm:ss
	//TaskType int // 1
	StartExecute string // yyyy-mm-dd
	EndExecute   string // yyyy-mm-dd
	Period       int    // 0
	TaskChannels []int
	//TaskPeriod string
	//TaskIO string
	//TaskEvent string // []
	//NetMode int // 7
	EffectiveDays int
	Stream        int // 1
	Storetype     int // 2
	VideoType     int // 0
}

type CreateAutoDownloadTaskResponse struct {
	Result bool `json:"result"`
}

func (c *Client) RegisterLogin(ctx context.Context) (*GetCenterGroupsResponse, error) {
	c.init()

	values := url.Values{}
	values.Set("Action", "Login")
	values.Set("Type", "post")
	values.Set("DataType", "Json")
	values.Set("Guid", fmt.Sprintf("%d", time.Now().UnixMilli()))

	inputValues := url.Values{}
	inputValues.Set("UserName", "")
	inputValues.Set("UserPassword", "")
	inputValues.Set("Token", url.QueryEscape(c.Key)) // Remember, the API wants the key to be escaped.
	inputValues.Set("Page", "alarmcenter")
	inputValues.Set("AuthCode", "")
	inputValues.Set("IsDES", "false")
	inputValuesString := inputValues.Encode()

	var output GetCenterGroupsResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/RegisterLogin/default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) MonitorAutoDownload(ctx context.Context, deviceID string) (*MonitorAutoDownloadResponse, error) {
	c.init()

	values := url.Values{}

	inputValues := url.Values{}
	inputValues.Set("action", "refreshTask")
	inputValues.Set("id", deviceID)
	inputValues.Set("nodetype", "1")
	inputValuesString := inputValues.Encode()

	var output MonitorAutoDownloadResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/AutoDownload/Monitor/Default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) QueryAutoDownload(ctx context.Context, deviceID string) (*QueryAutoDownloadResponse, error) {
	c.init()

	values := url.Values{}

	inputValues := url.Values{}
	inputValues.Set("action", "queryTask")
	inputValues.Set("NodeType", "1")
	inputValues.Set("Device", deviceID)
	inputValues.Set("StartTime", "2023-01-01") // TODO: PARAMETER
	inputValues.Set("EndTime", "2023-02-02")   // TODO: PARAMETER
	inputValues.Set("Status", "0")
	inputValues.Set("Type", "1")
	inputValues.Set("page", "1")
	inputValues.Set("rows", "10")
	inputValuesString := inputValues.Encode()

	var output QueryAutoDownloadResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/AutoDownload/GlobalReport/Default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) QueryAutoDownloadTask(ctx context.Context, deviceID string, taskID string) (*QueryAutoDownloadTaskResponse, error) {
	c.init()

	values := url.Values{}

	inputValues := url.Values{}
	inputValues.Set("action", "queryVideo")
	inputValues.Set("Device", deviceID)
	//inputValues.Set("Date", "2023-01-12") // TODO: PARAMETER
	inputValues.Set("TaskID", taskID)
	inputValues.Set("page", "1")
	inputValues.Set("rows", "10")
	inputValuesString := inputValues.Encode()

	var output QueryAutoDownloadTaskResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/AutoDownload/GlobalReport/Default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) CreateAutoDownloadTask(ctx context.Context, input CreateAutoDownloadTaskInput) (*CreateAutoDownloadTaskResponse, error) {
	c.init()

	values := url.Values{}

	var taskChannelStrings []string
	for _, channel := range input.TaskChannels {
		taskChannelStrings = append(taskChannelStrings, fmt.Sprintf("%d", channel))
	}
	inputValues := url.Values{}
	inputValues.Set("action", "saveTask")
	inputValues.Set("TaskName", input.TaskName)
	inputValues.Set("nodeType", "1")
	inputValues.Set("nodeName", input.DeviceID)
	inputValues.Set("StartTime", input.StartTime)
	inputValues.Set("EndTime", input.EndTime)
	inputValues.Set("TaskType", "1")
	inputValues.Set("StartExecute", input.StartExecute)
	inputValues.Set("EndExecute", input.EndExecute)
	inputValues.Set("Period", "0")
	inputValues.Set("TaskChannel", strings.Join(taskChannelStrings, ","))
	inputValues.Set("TaskPeriod", "")
	inputValues.Set("TaskIO", "")
	inputValues.Set("TaskEvent", "[]")
	inputValues.Set("NetMode", "7")
	inputValues.Set("Effective", fmt.Sprintf("%d", input.EffectiveDays))
	inputValues.Set("Stream", "1")
	inputValues.Set("Storetype", "2")
	inputValues.Set("VideoType", "0")
	inputValuesString := inputValues.Encode()

	var output CreateAutoDownloadTaskResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/AutoDownload/Task/Default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
