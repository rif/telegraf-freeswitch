package utils

import "encoding/json"

type Count struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Peak     int `json:"peak"`
	Peak5min int `json:"peak5Min"`
	Limit    int `json:"limit"`
}

type Rate struct {
	Current  int `json:"current"`
	Max      int `json:"max"`
	Peak     int `json:"peak"`
	Peak5min int `json:"peak5Min"`
}

type StackSize struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

type Sessions struct {
	Count Count `json:"count"`
	Rate  Rate  `json:"rate"`
}

type StatusResponse struct {
	Sessions  Sessions  `json:"sessions"`
	StackSize StackSize `json:"stackSizeKB"`
}

type Command struct {
	Status   string         `json:"status"`
	Response StatusResponse `json:"response"`
}

func LoadStatusJSON(data string) (*Command, error) {
	c := &Command{}
	err := json.Unmarshal([]byte(data), c)
	return c, err
}
