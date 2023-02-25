package angeltrax

import (
	"encoding/json"
	"fmt"
)

type GetServersResponse struct {
	ServiceMap        map[string]ClientService `json:"-"`
	ClienPath         string                   `json:"clienpath"`      // This is probably a typo on their part.
	LicenseTimeout    string                   `json:"licensetimeout"` // This is probably a date.
	LicenseTimeoutTip int                      `json:"licensetimeouttip"`
	ServerDate        string                   `json:"serverdate"`
	Support           []string                 `json:"support"`
	Upgrade           int                      `json:"upgrade"`
	Version           string                   `json:"version"`
}

func (r *GetServersResponse) UnmarshalJSON(contents []byte) error {
	r.ServiceMap = map[string]ClientService{}

	var m map[string]json.RawMessage
	err := json.Unmarshal(contents, &m)
	if err != nil {
		return fmt.Errorf("could not unmarshal into intermediate structure: %v", err)
	}
	for key, value := range m {
		switch key {
		case "clienpath":
			err = json.Unmarshal(value, &r.ClienPath)
			if err != nil {
				return fmt.Errorf("could not unmarshal key %q: %v", key, err)
			}
		case "licensetimeout":
			err = json.Unmarshal(value, &r.LicenseTimeout)
			if err != nil {
				return fmt.Errorf("could not unmarshal key %q: %v", key, err)
			}
		case "licensetimeouttip":
			err = json.Unmarshal(value, &r.LicenseTimeoutTip)
			if err != nil {
				return fmt.Errorf("could not unmarshal key %q: %v", key, err)
			}
		case "serverdate":
			err = json.Unmarshal(value, &r.ServerDate)
			if err != nil {
				return fmt.Errorf("could not unmarshal key %q: %v", key, err)
			}
		case "support":
			err = json.Unmarshal(value, &r.Support)
			if err != nil {
				return fmt.Errorf("could not unmarshal key %q: %v", key, err)
			}
		case "upgrade":
			err = json.Unmarshal(value, &r.Upgrade)
			if err != nil {
				return fmt.Errorf("could not unmarshal key %q: %v", key, err)
			}
		case "version":
			err = json.Unmarshal(value, &r.Version)
			if err != nil {
				return fmt.Errorf("could not unmarshal key %q: %v", key, err)
			}
		default:
			var clientServer ClientService
			err = json.Unmarshal(value, &clientServer)
			if err != nil {
				return fmt.Errorf("could not unmarshal client server %q: %v", key, err) // TODO: CONSIDER JUST WARNING HERE INSTEAD
			}
			r.ServiceMap[key] = clientServer
		}
	}
	return nil
}

type ClientService struct {
	Address       string `json:"ip,omitempty"`
	Port          int    `json:"port,omitempty"`
	SecureAddress string `json:"ips,omitempty"`
	SecurePort    int    `json:"ports,omitempty"`
	Enable        int    `json:"enable"`
	UseSecure     int    `json:"usesecure,omitempty"`
}
