package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rif/go-eventsocket/eventsocket"
)

const (
	JSONFormat   = "JSON"
	InfluxFormat = "INFLUX"
)

type Fetcher struct {
	conn          *eventsocket.Connection
	sessions      *Sessions
	sofiaProfiles []*SofiaProfile
	cacheTime     time.Time
}

func NewFetcher(host string, port int, pass string) (*Fetcher, error) {
	conn := &eventsocket.Connection{}
	conn, err := eventsocket.Dial(fmt.Sprintf("%s:%d", host, port), pass)
	if err != nil {
		return nil, err
	}

	return &Fetcher{conn: conn}, nil
}

type status struct {
	Active       int `json:"active"`
	Peak         int `json:"peak"`
	Peak5min     int `json:"peak_5min"`
	Total        int `json:"total"`
	RateCurrent  int `json:"rate_current"`
	RateMax      int `json:"rate_max"`
	RatePeak     int `json:"rate_peak"`
	RatePeak5min int `json:"rate_peak_5min"`
}

type profile struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Running string `json:"running"`
}

func (f *Fetcher) Close() {
	f.conn.Close()
}

func (f *Fetcher) GetData() error {
	if !f.cacheTime.IsZero() && time.Since(f.cacheTime) < 5*time.Second {
		return nil
	}
	var c *Command
	ev, err := f.conn.Send(`api json {"command" : "status", "data" : ""}`)
	if err != nil {
		ev, err = f.conn.Send(`api status`)
		if err != nil {
			return fmt.Errorf("error sending status command: %v", err)
		}
		c, err = LoadStatusText(ev.Body)
		if err != nil || c.Status != "success" {
			return fmt.Errorf("error parsing status command: %v %+v", err, c)
		}
	} else {
		c, err = LoadStatusJSON(ev.Body)
		if err != nil || c.Status != "success" {
			return fmt.Errorf("error parsing status command: %v %+v", err, c)
		}
	}
	sessions := &c.Response.Sessions
	ev, err = f.conn.Send("api sofia xmlstatus")
	if err != nil {
		return fmt.Errorf("error sending xmlstatus: %v", err)
	}
	sofiaProfiles, err := ParseSofiaStatus(ev.Body)
	if err != nil {
		return fmt.Errorf("error parsing xmlstatus: %v", err)
	}
	f.sessions = sessions
	f.sofiaProfiles = sofiaProfiles
	f.cacheTime = time.Now()
	return nil
}

func (f *Fetcher) FormatOutput(format string) (string, string) {
	if f.sessions == nil || f.sofiaProfiles == nil {
		return "", ""
	}
	if format == JSONFormat {
		s := status{
			Active:       f.sessions.Count.Active,
			Peak:         f.sessions.Count.Peak,
			Peak5min:     f.sessions.Count.Peak5min,
			Total:        f.sessions.Count.Total,
			RateCurrent:  f.sessions.Rate.Current,
			RateMax:      f.sessions.Rate.Max,
			RatePeak:     f.sessions.Rate.Peak,
			RatePeak5min: f.sessions.Rate.Peak5min,
		}
		status, _ := json.MarshalIndent(s, "", " ")
		pfs := make([]profile, len(f.sofiaProfiles))
		for i, sofiaProfile := range f.sofiaProfiles {
			pfs[i] = profile{
				Name:    sofiaProfile.Name,
				IP:      sofiaProfile.Address,
				Running: sofiaProfile.Running,
			}
		}
		profiles, _ := json.MarshalIndent(pfs, "", " ")
		return string(status), string(profiles)
	}
	status := fmt.Sprintf("freeswitch_sessions active=%d,peak=%d,peak_5min=%d,total=%d,rate_current=%d,rate_max=%d,rate_peak=%d,rate_peak_5min=%d\n",
		f.sessions.Count.Active,
		f.sessions.Count.Peak,
		f.sessions.Count.Peak5min,
		f.sessions.Count.Total,
		f.sessions.Rate.Current,
		f.sessions.Rate.Max,
		f.sessions.Rate.Peak,
		f.sessions.Rate.Peak5min,
	)
	profiles := ""
	for _, sofiaProfile := range f.sofiaProfiles {
		profiles += fmt.Sprintf("freeswitch_profile_sessions,profile=%s,ip=%s running=%s\n",
			sofiaProfile.Name,
			sofiaProfile.Address,
			sofiaProfile.Running)
	}
	return status, profiles
}
