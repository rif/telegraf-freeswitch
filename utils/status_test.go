package utils

import (
	"reflect"
	"testing"
)

var (
	statusDataJSON = `{
	"command": "status",
	"data": "",
	"status": "success",
	"response": {
		"systemStatus": "ready",
		"uptime": {
			"years": 0,
			"days": 2,
			"hours": 1,
			"minutes": 11,
			"seconds": 41,
			"milliseconds": 464,
			"microseconds": 27
		},
		"version": "1.6.19 -36-7a77e0b 64bit",
		"sessions": {
			"count": {
				"total": 7437938,
				"active": 2886,
				"peak": 4969,
				"peak5Min": 2934,
				"limit": 10000
			},
			"rate": {
				"current": 31,
				"max": 300,
				"peak": 283,
				"peak5Min": 55
			}
		},
		"idleCPU": {
			"used": 0,
			"allowed": 72.266667
		},
		"stackSizeKB": {
			"current": 240,
			"max": 8192
		}
	}
}
`
	statusDataText = `UP 0 years, 207 days, 9 hours, 5 minutes, 40 seconds, 521 milliseconds, 268 microseconds
FreeSWITCH (Version 1.2.23  64bit) is ready
18169728 session(s) since startup
140 session(s) - peak 342, last 5min 142
11 session(s) per Sec out of max 80, peak 42, last 5min 14
1000 session(s) max
min idle cpu 0.00/62.00
Current Stack Size/Max 240K/8192K`
)

func TestLoadStatusJSON(t *testing.T) {
	c, err := LoadStatusJSON(statusDataJSON)
	if err != nil {
		t.Error("error parsing respons: ", err)
	}
	if c.Status != "success" {
		t.Error("bad command status: ", c.Status)
	}
	sr := StatusResponse{
		Sessions: Sessions{
			Count: Count{
				Total:    7437938,
				Active:   2886,
				Peak:     4969,
				Peak5min: 2934,
				Limit:    10000,
			},
			Rate: Rate{
				Current:  31,
				Max:      300,
				Peak:     283,
				Peak5min: 55,
			},
		},
		StackSize: StackSize{
			Current: 240,
			Max:     8192,
		},
	}
	if !reflect.DeepEqual(c.Response, sr) {
		t.Errorf("bad response: %+v", c.Response)
	}
}

func TestLoadStatusText(t *testing.T) {
	c, err := LoadStatusText(statusDataText)
	if err != nil {
		t.Error("error parsing respons: ", err)
		t.FailNow()
	}
	sr := StatusResponse{
		Sessions: Sessions{
			Count: Count{
				Total:    18169728,
				Active:   140,
				Peak:     342,
				Peak5min: 142,
				Limit:    1000,
			},
			Rate: Rate{
				Current:  11,
				Max:      80,
				Peak:     42,
				Peak5min: 14,
			},
		},
	}
	if !reflect.DeepEqual(c.Response, sr) {
		t.Errorf("bad response: %+v", c.Response)
	}
}
