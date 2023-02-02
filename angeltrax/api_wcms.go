package angeltrax

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type RegisterLoginResponse struct {
	Code   int  `json:"Code"`
	Result bool `json:"Result"`
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
	input := inputValues.Encode()

	var output GetCenterGroupsResponse
	err := c.RawServiceRequest(ctx, "wcms", http.MethodPost, "/Plugin/RegisterLogin/default.ashx", values, input, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
