package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rif/go-eventsocket/eventsocket"
	"github.com/rif/telegraf-freeswitch/utils"
)

var (
	host = flag.String("host", "localhost", "freeswitch host address")
	port = flag.Int("port", 8021, "freeswitch port")
	pass = flag.String("passs", "ClueCon", "freeswitch password")
)

func main() {
	flag.Parse()
	l := log.New(os.Stderr, "", 0)
	conn, err := eventsocket.Dial(fmt.Sprintf("%s:%d", *host, *port), *pass)
	if err != nil {
		l.Print("error connecting to fs: ", err)
	}
	var sessions utils.Sessions
	if ev, err := conn.Send(`api json {"command" : "status", "data" : ""}`); err != nil {
		l.Print("error sending status command: ", err)
	} else {
		c, err := utils.LoadStatusJSON(ev.Body)
		if err != nil || c.Status != "success" {
			l.Printf("error parsing status command: %v %+v", err, c)
		} else {
			sessions = c.Response.Sessions
		}
	}
	var sofiaProfiles []*utils.SofiaProfile
	if ev, err := conn.Send("api sofia xmlstatus"); err != nil {
		l.Print("error sending xmlstatus: ", err)
	} else {
		sofiaProfiles, err = utils.ParseSofiaStatus(ev.Body)
		if err != nil {
			l.Print("error parsing xmlstatus: ", err)
		}
	}
	fmt.Printf("freeswitch_sessions active=%d,peak=%d,peak_5min=%d,total=%d,rate_current=%d,rate_max=%d,rate_peak=%d,rate_peak_5min=%d\n",
		sessions.Count.Active,
		sessions.Count.Peak,
		sessions.Count.Peak5min,
		sessions.Count.Total,
		sessions.Rate.Current,
		sessions.Rate.Max,
		sessions.Rate.Peak,
		sessions.Rate.Peak5min,
	)
	for _, sofiaProfile := range sofiaProfiles {
		fmt.Printf("freeswitch_profile_sessions,profile=%s,ip=%s running=%s\n",
			sofiaProfile.Name,
			sofiaProfile.Address,
			sofiaProfile.Running)
	}
}
