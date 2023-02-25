package angeltrax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	Server        string
	ServerPort    int
	Username      string
	Password      string
	Key           string
	serviceMap    map[string]ClientService
	hostCookieMap map[string][]string
	httpClient    http.Client
}

func (c *Client) init() {
	if c.ServerPort == 0 {
		c.ServerPort = 7264
	}
	if c.serviceMap == nil {
		c.serviceMap = map[string]ClientService{}
	}
	if c.hostCookieMap == nil {
		c.hostCookieMap = map[string][]string{}
	}
}

func (c *Client) GetServers(ctx context.Context, server string) (*GetServersResponse, error) {
	c.init()

	values := url.Values{}
	values.Set("did", "bbb")

	var output GetServersResponse
	err := c.RawRequest(ctx, http.MethodGet, "http://"+server+":"+fmt.Sprintf("%d", c.ServerPort)+"/serversforclient/BalanceServer.ashx", values, nil, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *Client) RawServiceRequest(ctx context.Context, server, method, path string, values url.Values, requestData, responseData interface{}) error {
	c.init()

	info, ok := c.serviceMap[server]
	if !ok {
		return fmt.Errorf("no server info for: %s", server)
	}
	logrus.Debugf("Server %q: %+v", server, info)

	var base string
	if info.UseSecure > 0 {
		base = "https://"
		if info.SecureAddress == "0.0.0.0" {
			base += c.Server
		} else {
			base += info.SecureAddress
		}
		base += ":" + fmt.Sprintf("%d", info.SecurePort)
	} else {
		base = "http://"
		if info.Address == "0.0.0.0" {
			base += c.Server
		} else {
			base += info.Address
		}
		base += ":" + fmt.Sprintf("%d", info.Port)
	}

	return c.RawRequest(ctx, method, base+"/"+strings.TrimPrefix(path, "/"), values, requestData, responseData)
}
func (c *Client) RawRequest(ctx context.Context, method, path string, values url.Values, requestData, responseData interface{}) error {
	c.init()
	if len(values) > 0 {
		path = path + "?" + values.Encode()
	}
	logrus.Debugf("Making request: %s %s", method, path)

	var contentType string

	var requestBody []byte
	if requestData != nil {
		logrus.WithContext(ctx).Debugf("Request: requestData: [%T]", requestData)
		if v, ok := requestData.(string); ok {
			requestBody = []byte(v)
			contentType = "application/x-www-form-urlencoded"
		} else {
			contents, err := json.Marshal(requestData)
			if err != nil {
				return err
			}
			requestBody = contents
		}
	}
	var requestBodyReader io.Reader
	if requestBody != nil {
		logrus.WithContext(ctx).Debugf("Request: Body length: %d", len(requestBody))
		logrus.WithContext(ctx).Debugf("Request: Body: %s", requestBody)
		requestBodyReader = bytes.NewReader(requestBody)
	}

	request, err := http.NewRequest(http.MethodGet, path, requestBodyReader)
	if err != nil {
		return err
	}

	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	for _, cookie := range c.hostCookieMap[request.Host] {
		request.Header.Add("Cookie", cookie)
	}

	for key, values := range request.Header {
		logrus.WithContext(ctx).Debugf("> %s: %v", key, values)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if cookies := response.Header.Values("Set-Cookie"); len(cookies) > 0 {
		c.hostCookieMap[request.Host] = cookies
	}

	for key, values := range response.Header {
		logrus.WithContext(ctx).Debugf("< %s: %v", key, values)
	}

	logrus.Debugf("Response status: %d", response.StatusCode)
	if response.StatusCode > 299 {
		return fmt.Errorf("http status %d", response.StatusCode)
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	logrus.Debugf("Response body: %s", contents)

	if responseData != nil {
		err = json.Unmarshal(contents, responseData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Login(ctx context.Context, server, username, password string) error {
	c.init()

	getServersResponse, err := c.GetServers(ctx, server)
	if err != nil {
		return err
	}
	c.serviceMap = getServersResponse.ServiceMap

	values := url.Values{}
	values.Set("username", username)
	values.Set("password", password)

	var output GetKeyResponse
	err = c.RawServiceRequest(ctx, "webclient", http.MethodGet, "/api/v1/inner/key", values, nil, &output)
	if err != nil {
		return err
	}

	// For whatever reason, these guys return the key already URL-escaped.
	c.Key, err = url.PathUnescape(output.Data.Key)
	if err != nil {
		return err
	}

	return nil
}
