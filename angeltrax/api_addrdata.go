package angeltrax

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type GetCenterGroupsResponse struct {
	ErrorCode int           `json:"errorcode"`
	Data      []CenterGroup `json:"data"`
}

type CenterGroup struct {
	GroupFatherID int    `json:"groupfatherid"`
	GroupID       int    `json:"groupid"`
	GroupName     string `json:"groupname"`
	Remark        string `json:"string"`
}

type GetCenterDevicesResponse struct {
	ErrorCode int            `json:"errorcode"`
	Data      []CenterDevice `json:"data"`
}

type CenterDevice struct {
	CarLicense     string `json:"carlicence"`
	ChannelCount   int    `json:"channelcount"`
	CameraNames    string `json:"cname"` // A comma-separated list of channel names.  Unnamed channels are empty strings.
	DeviceID       string `json:"deviceid"`
	DevicePassword string `json:"devicepassword"`
	DeviceType     string `json:"devicetype"` // 4 is what we have on all the trucks.
	DeviceUsername string `json:"deviceusername"`
	En             int    `json:"en"`          // TODO: WHAT IS THIS?  I've seen 31, 15, and -1.
	GroupID        int    `json:"groupid"`     // This refers to a CenterGroup.GroupID.
	LinkType       string `json:"linktype"`    // TODO: WHAT IS THIS?
	PrevChannel    int    `json:"prevchannel"` // TODO: WHAT IS THIS?  The channel to preview?  Zero-indexed?
	RegisterIP     string `json:"registerip"`
	RegisterPort   int    `json:"registerport"`
	Remark         string `json:"remark"`
	TransmitIP     string `json:"transmitip"`
	TransmitPort   int    `json:"transmitport"`
}

func (c *Client) GetCenterGroups(ctx context.Context) (*GetCenterGroupsResponse, error) {
	c.init()

	values := url.Values{}
	values.Set("key", c.Key)
	values.Set("random", fmt.Sprintf("%d", time.Now().Unix()))

	var output GetCenterGroupsResponse
	err := c.RawServiceRequest(ctx, "addrdata", http.MethodGet, "/center/group", values, nil, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) GetCenterDevices(ctx context.Context) (*GetCenterDevicesResponse, error) {
	c.init()

	values := url.Values{}
	values.Set("key", c.Key)
	values.Set("random", fmt.Sprintf("%d", time.Now().Unix()))

	var output GetCenterDevicesResponse
	err := c.RawServiceRequest(ctx, "addrdata", http.MethodGet, "/center/device", values, nil, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
