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

type MonitorAutoDownloadTaskResponse struct {
	TaskID       int           `json:"TaskID"`
	TaskName     string        `json:"TaskName"`
	DeviceID     string        `json:"Device"`
	StartExecute string        `json:"StartExecute"` // yyyy-mm-dd
	EndExecute   string        `json:"EndExecute"`   // yyyy-mm-dd
	StartTime    string        `json:"StartTime"`    // hh:mm:ss
	EndTime      string        `json:"EndTime"`      // hh:mm:ss
	Period       int           `json:"Period"`
	TaskType     int           `json:"TaskType"`
	TaskPeriod   []interface{} `json:"TaskPeriod"`  // TODO: WHAT IS THIS FORMAT?
	TaskChannel  []int         `json:"TaskChannel"` // List of one-indexed channels.
	TaskEvent    []interface{} `json:"TaskEvent"`   // TODO: WHAT IS THIS FORMAT?
	TaskIO       []interface{} `json:"TaskIO"`      // TODO: WHAT IS THIS FORMAT?
	Relation     string        `json:"Relation"`
	CarLicense   string        `json:"Carlicense"`
	NetMode      string        `json:"NetMode"`
	Effective    int           `json:"Effective"`
	Stream       int           `json:"Stream"`
	VideoType    int           `json:"VideoType"`
	StoreType    int           `json:"StoreType"`
}

type GlobalReportAutoDownloadInput struct {
	DeviceID  string `form:"Device"`
	StartDate string `form:"StartTime"` // yyyy-mm-dd
	EndDate   string `form:"EndTime"`   // yyyy-mm-dd
}

type GlobalReportAutoDownloadResponse struct {
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

type GlobalReportAutoDownloadTaskResponse struct {
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
	TaskName  string `form:"TaskName"`
	DeviceID  string `form:"nodeName"`
	StartTime string `form:"StartTime"` // hh:mm:ss
	EndTime   string `form:"EndTime"`   // hh:mm:ss
	//TaskType int `form:"TaskType"` // 1
	StartExecute string `form:"StartExecute"` // yyyy-mm-dd
	EndExecute   string `form:"EndExecute"`   // yyyy-mm-dd
	Period       int    `form:"Period"`       // 0
	TaskChannels []int  `form:"TaskChannels"` // TODO: This will be a comma-separated string in the form.
	//TaskPeriod string `form:"TaskPeriod"`
	//TaskIO string `form:"TaskIO"`
	//TaskEvent string `form:"TaskEvent"` // []
	//NetMode int `form:"NetMode"` // 7
	EffectiveDays int `form:"EffectiveDays"`
	Stream        int `form:"Stream"`    // 1
	Storetype     int `form:"StoreType"` // 2
	VideoType     int `form:"VideoType"` // 0
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

func (c *Client) MonitorAutoDownloadTask(ctx context.Context, taskID string) (*MonitorAutoDownloadTaskResponse, error) {
	c.init()

	values := url.Values{}

	inputValues := url.Values{}
	inputValues.Set("action", "getTask")
	inputValues.Set("id", taskID)
	inputValuesString := inputValues.Encode()

	var output MonitorAutoDownloadTaskResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/AutoDownload/Monitor/Default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) GlobalReportAutoDownload(ctx context.Context, input GlobalReportAutoDownloadInput) (*GlobalReportAutoDownloadResponse, error) {
	c.init()

	values := url.Values{}

	inputValues := url.Values{}
	inputValues.Set("action", "queryTask")
	inputValues.Set("NodeType", "1")
	inputValues.Set("Device", input.DeviceID)
	inputValues.Set("StartTime", input.StartDate)
	inputValues.Set("EndTime", input.EndDate)
	inputValues.Set("Status", "0")
	inputValues.Set("Type", "1")
	inputValues.Set("page", "1")
	inputValues.Set("rows", "10")
	inputValuesString := inputValues.Encode()

	var output GlobalReportAutoDownloadResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/AutoDownload/GlobalReport/Default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) GlobalReportAutoDownloadTask(ctx context.Context, deviceID string, taskID string) (*GlobalReportAutoDownloadTaskResponse, error) {
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

	var output GlobalReportAutoDownloadTaskResponse
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
