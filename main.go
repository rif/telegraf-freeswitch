package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rif/go-eventsocket/eventsocket"
	"github.com/rif/telegraf-freeswitch/utils"
)

var (
	host           = flag.String("host", "localhost", "freeswitch host address")
	port           = flag.Int("port", 8021, "freeswitch port")
	pass           = flag.String("pass", "ClueCon", "freeswitch password")
	serve          = flag.Bool("serve", false, "run as a server")
	listen_address = flag.String("listen_address", "127.0.0.1", "listen on address")
	listen_port    = flag.Int("listen_port", 9191, "listen on port")
)

func getData(conn *eventsocket.Connection) (sessions *utils.Sessions, sofiaProfiles []*utils.SofiaProfile, err error) {
	ev, err := conn.Send(`api json {"command" : "status", "data" : ""}`)
	if err != nil {
		return nil, nil, fmt.Errorf("error sending status command: %v", err)
	}
	c, err := utils.LoadStatusJSON(ev.Body)
	if err != nil || c.Status != "success" {
		return nil, nil, fmt.Errorf("error parsing status command: %v %+v", err, c)
	}
	sessions = &c.Response.Sessions
	ev, err = conn.Send("api sofia xmlstatus")
	if err != nil {
		return nil, nil, fmt.Errorf("error sending xmlstatus: %v", err)
	}
	sofiaProfiles, err = utils.ParseSofiaStatus(ev.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing xmlstatus: %v", err)
	}
	return
}

func getOutput(sessions *utils.Sessions, sofiaProfiles []*utils.SofiaProfile) string {
	out := fmt.Sprintf("freeswitch_sessions active=%d,peak=%d,peak_5min=%d,total=%d,rate_current=%d,rate_max=%d,rate_peak=%d,rate_peak_5min=%d\n",
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
		out += fmt.Sprintf("freeswitch_profile_sessions,profile=%s,ip=%s running=%s\n",
			sofiaProfile.Name,
			sofiaProfile.Address,
			sofiaProfile.Running)
	}
	return out
}

func main() {
	flag.Parse()
	l := log.New(os.Stderr, "", 0)
	conn, err := eventsocket.Dial(fmt.Sprintf("%s:%d", *host, *port), *pass)
	if err != nil {
		l.Print("error connecting to fs: ", err)
	}
	defer conn.Close()
	if !*serve {
		sessions, sofiaProfiles, err := getData(conn)
		if err != nil {
			l.Print(err.Error())
		}
		fmt.Print(getOutput(sessions, sofiaProfiles))
		os.Exit(0)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sessions, sofiaProfiles, err := getData(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(getOutput(sessions, sofiaProfiles)))
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *listen_address, *listen_port), nil))
}
