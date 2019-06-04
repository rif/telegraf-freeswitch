package utils

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strconv"
)

var (
	countActiveParser = regexp.MustCompile(`(\d+) session\(s\) - peak (\d+), last 5min (\d+)`)
	countTotalParser  = regexp.MustCompile(`(\d+) session\(s\) since startup`)
	countMaxParser    = regexp.MustCompile(`(\d+) session\(s\) max`)
	rateParser        = regexp.MustCompile(`(\d+) session\(s\) per Sec out of max (\d+), peak (\d+), last 5min (\d+)`)
)

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

func LoadStatusText(data string) (*Command, error) {
	activeSlice := countActiveParser.FindStringSubmatch(data)
	if len(activeSlice) != 4 {
		log.Print(activeSlice)
		return nil, errors.New("could not parse active session count: " + data)
	}
	totalSlice := countTotalParser.FindStringSubmatch(data)
	if len(totalSlice) != 2 {
		return nil, errors.New("could not parse total session count: " + data)
	}
	maxSlice := countMaxParser.FindStringSubmatch(data)
	if len(maxSlice) != 2 {
		return nil, errors.New("could not parse rate session count: " + data)
	}
	rateSlice := rateParser.FindStringSubmatch(data)
	if len(rateSlice) != 5 {
		return nil, errors.New("could not parse rate session count: " + data)
	}
	c := &Command{
		Response: StatusResponse{
			Sessions: Sessions{
				Count: Count{
					Total:    parseInt(totalSlice[1]),
					Active:   parseInt(activeSlice[1]),
					Peak:     parseInt(activeSlice[2]),
					Peak5min: parseInt(activeSlice[3]),
					Limit:    parseInt(maxSlice[1]),
				},
				Rate: Rate{
					Current:  parseInt(rateSlice[1]),
					Max:      parseInt(rateSlice[2]),
					Peak:     parseInt(rateSlice[3]),
					Peak5min: parseInt(rateSlice[4]),
				},
			},
		},
		Status: "success",
	}
	return c, nil
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
