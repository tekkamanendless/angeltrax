package angeltrax

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultRowCount = 100

type TaskStatus int

const (
	TaskStatusPaused           TaskStatus = -6
	TaskStatusConnectionLimit  TaskStatus = -5
	TaskStatusAnalyzing        TaskStatus = -4
	TaskStatusNotFinished      TaskStatus = -3
	TaskStatusInsufficientDisk TaskStatus = -2
	TaskStatusWaiting          TaskStatus = -1
	TaskStatusAnalyzed         TaskStatus = 0
	TaskStatusDownloading      TaskStatus = 1
	TaskStatusNoFiles          TaskStatus = 2
	TaskStatusFinished         TaskStatus = 3
	TaskStatusDownloadFailed   TaskStatus = 4 // This might be simply "failed" and the other one is for "download failed"
	TaskStatusDelete           TaskStatus = 5
	TaskStatusDownloadFailed2  TaskStatus = 6
	TaskStatusTimeout          TaskStatus = 8
)

func (t TaskStatus) String() string {
	var value string
	switch t {
	case TaskStatusPaused:
		value = "TaskStatusPaused"
	case TaskStatusConnectionLimit:
		value = "TaskStatusConnectionLimit"
	case TaskStatusAnalyzing:
		value = "TaskStatusAnalyzing"
	case TaskStatusNotFinished:
		value = "TaskStatusNotFinished"
	case TaskStatusInsufficientDisk:
		value = "TaskStatusInsufficientDisk"
	case TaskStatusWaiting:
		value = "TaskStatusWaiting"
	case TaskStatusAnalyzed:
		value = "TaskStatusAnalyzed"
	case TaskStatusDownloading:
		value = "TaskStatusDownloading"
	case TaskStatusNoFiles:
		value = "TaskStatusNoFiles"
	case TaskStatusFinished:
		value = "TaskStatusFinished"
	case TaskStatusDownloadFailed:
		value = "TaskStatusDownloadFailed"
	case TaskStatusDelete:
		value = "TaskStatusDelete"
	case TaskStatusDownloadFailed2:
		value = "TaskStatusDownloadFailed2"
	case TaskStatusTimeout:
		value = "TaskStatusTimeout"
	}
	return fmt.Sprintf("%s(%d)", value, t)
}

type TaskPeriod int

const (
	TaskPeriodManual     TaskPeriod = -1
	TaskPeriodOnce       TaskPeriod = 0
	TaskPeriodEveryDay   TaskPeriod = 1
	TaskPeriodEveryWeek  TaskPeriod = 2
	TaskPeriodEveryMonth TaskPeriod = 3
)

type TaskType int

const (
	TaskTypeBlackBox      TaskType = 0
	TaskTypeVideo         TaskType = 1
	TaskTypeBlackBoxVideo TaskType = 2 // This was "default" in a switch statement.
)

type RegisterLoginResponse struct {
	Code   int  `json:"Code"`
	Result bool `json:"Result"`
}

type MonitorAutoDownloadInput struct {
	DeviceID string `form:"id"`
}

type MonitorAutoDownloadResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		TaskID      int        `json:"TaskID"`
		Status      TaskStatus `json:"Status"`
		DeviceID    string     `json:"Device"`
		CarLicense  string     `json:"Carlicense"`
		TaskName    string     `json:"TaskName"`
		Period      TaskPeriod `json:"Period"`
		TaskType    TaskType   `json:"TaskType"`
		Date        string     `json:"Date"`       // yyyy-mm-dd
		StartTime   string     `json:"StartTime"`  // hh:mm:ss
		EndTime     string     `json:"EndTime"`    // hh:mm:ss
		ChannelList string     `json:"Channel"`    // CSV of channel numbers, starting from "1".
		CreateTime  string     `json:"CreateTime"` // yyyy-mm-dd hh:mm:ss

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
	Period       TaskPeriod    `json:"Period"`
	TaskType     TaskType      `json:"TaskType"`
	TaskPeriod   []interface{} `json:"TaskPeriod"`  // TODO: WHAT IS THIS FORMAT?
	TaskChannel  []int         `json:"TaskChannel"` // List of one-indexed channels.
	TaskEvent    []interface{} `json:"TaskEvent"`   // TODO: WHAT IS THIS FORMAT?
	TaskIO       []interface{} `json:"TaskIO"`      // TODO: WHAT IS THIS FORMAT?
	Relation     string        `json:"Relation"`
	CarLicense   string        `json:"Carlicense"`
	NetMode      string        `json:"NetMode"` // 1: lan, 2: wifi, 3: wifiandlan, 4: 3G, 7: all
	Effective    int           `json:"Effective"`
	Stream       int           `json:"Stream"`    // 0: sub, 1: main
	VideoType    int           `json:"VideoType"` // 0: all, 1: normal, 2: alarm
	StoreType    int           `json:"Storetype"` // 0: main, 1: sub, 2: both
}

type GlobalReportAutoDownloadInput struct {
	DeviceID  string     `form:"Device"`
	Status    TaskStatus `form:"Status"`
	StartDate string     `form:"StartTime"` // yyyy-mm-dd
	EndDate   string     `form:"EndTime"`   // yyyy-mm-dd
}

type GlobalReportAutoDownloadResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		TaskID      int        `json:"TaskID"`
		Status      TaskStatus `json:"Status"`
		DeviceID    string     `json:"Device"`
		CarLicense  string     `json:"Carlicense"`
		TaskName    string     `json:"TaskName"`
		Period      TaskPeriod `json:"Period"`
		TaskType    TaskType   `json:"TaskType"`
		Date        string     `json:"Date"`       // yyyy-mm-dd
		StartTime   string     `json:"StartTime"`  // hh:mm:ss
		EndTime     string     `json:"EndTime"`    // hh:mm:ss
		ChannelList string     `json:"Channel"`    // CSV of channel numbers, starting from "1".
		CreateTime  string     `json:"CreateTime"` // yyyy-mm-dd hh:mm:ss

		FinishTime string `json:"FinishTime"` // yyyy-mm-dd hh:mm:ss
		Username   string `json:"UserName"`
	} `json:"rows"`
}

type GlobalReportAutoDownloadTaskInput struct {
	Date     string `form:"Date"`
	DeviceID string `form:"Device"`
	TaskID   string `form:"TaskID"`
}

type GlobalReportAutoDownloadTaskResponse struct {
	Total int `json:"total"`
	Rows  []struct {
		DeviceID    string     `json:"Device"`
		Status      TaskStatus `json:"Status"`
		Percent     string     `json:"Percent"`   // This appears to be a float encoded as a string.
		Speed       string     `json:"Speed"`     // This appears to be a float encoded as a string.
		Date        string     `json:"Date"`      // yyyy-mm-dd
		StartTime   string     `json:"StartTime"` // hh:mm:ss
		EndTime     string     `json:"EndTime"`   // hh:mm:ss
		TotalSize   string     `json:"TotalSize"` // This appears to be a float encoded as a string.
		CurrentSize string     `json:"CurSize"`   // This appears to be a float encoded as a string.
		Channel     int        `json:"Channel"`   // The one-index of the channel.
		Error       string     `json:"Error"`
		TaskID      int        `json:"TaskID"`
		FileSource  string     `json:"FileSource"`
		PreAlarm    int        `json:"PreAlarm"`
		NextAlarm   int        `json:"NextAlarm"`
	} `json:"rows"`
}

type CreateAutoDownloadTaskInput struct {
	TaskName     string     `form:"TaskName"`
	DeviceID     string     `form:"nodeName"`
	StartTime    string     `form:"StartTime"` // hh:mm:ss
	EndTime      string     `form:"EndTime"`   // hh:mm:ss
	TaskType     TaskType   `form:"TaskType"`
	StartExecute string     `form:"StartExecute"` // yyyy-mm-dd
	EndExecute   string     `form:"EndExecute"`   // yyyy-mm-dd
	Period       TaskPeriod `form:"Period"`       // 0
	TaskChannels []int      `form:"TaskChannels"` // TODO: This will be a comma-separated string in the form.
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

func (c *Client) MonitorAutoDownload(ctx context.Context, input MonitorAutoDownloadInput) (*MonitorAutoDownloadResponse, error) {
	c.init()

	values := url.Values{}

	inputValues := url.Values{}
	inputValues.Set("action", "refreshTask")
	inputValues.Set("id", input.DeviceID)
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
	inputValues.Set("Status", fmt.Sprintf("%d", input.Status))
	inputValues.Set("Type", "1")
	inputValues.Set("page", "1")
	inputValues.Set("rows", fmt.Sprintf("%d", DefaultRowCount))
	inputValuesString := inputValues.Encode()

	var output GlobalReportAutoDownloadResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/AutoDownload/GlobalReport/Default.ashx", values, inputValuesString, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) GlobalReportAutoDownloadTask(ctx context.Context, input GlobalReportAutoDownloadTaskInput) (*GlobalReportAutoDownloadTaskResponse, error) {
	c.init()

	values := url.Values{}

	inputValues := url.Values{}
	inputValues.Set("action", "queryVideo")
	inputValues.Set("Device", input.DeviceID)
	inputValues.Set("Date", input.Date)
	inputValues.Set("TaskID", input.TaskID)
	inputValues.Set("page", "1")
	inputValues.Set("rows", fmt.Sprintf("%d", DefaultRowCount))
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
	inputValues.Set("TaskType", fmt.Sprintf("%d", input.TaskType))
	inputValues.Set("StartExecute", input.StartExecute)
	inputValues.Set("EndExecute", input.EndExecute)
	inputValues.Set("Period", fmt.Sprintf("%d", input.Period))
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
