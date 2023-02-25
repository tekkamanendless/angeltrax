package angeltrax

type GetKeyResponse struct {
	ErrorCode int `json:"errorcode"`
	Data      struct {
		Key string `json:"key"`
	} `json:"data"`
}
