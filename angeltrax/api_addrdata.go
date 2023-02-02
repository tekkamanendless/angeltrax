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
