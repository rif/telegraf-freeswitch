package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

func handler(w http.ResponseWriter, route string) {
}

func main() {
	flag.Parse()
	l := log.New(os.Stderr, "", 0)
	fetcher, err := utils.NewFetcher(*host, *port, *pass)
	if err != nil {
		l.Print("error connecting to fs: ", err)
	}
	defer fetcher.Close()
	if !*serve {
		sessions, sofiaProfiles, err := fetcher.GetData()
		if err != nil {
			l.Print(err.Error())
		}
		fmt.Print(formatOutputInflux(sessions, sofiaProfiles))
		os.Exit(0)
	}

	http.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		sessions, sofiaProfiles, err := fetcher.GetData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		status, _ := fetcher.FormatOutput(utils.JSONFormat)
		if _, err := w.Write(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/profiles/", func(w http.ResponseWriter, r *http.Request) {
		sessions, sofiaProfiles, err := fetcher.GetData()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, profiles := fetcher.FormatOutput(utils.JSONFormat)
		if _, err := w.Write(profiler); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	listen := fmt.Sprintf("%s:%d", *listen_address, *listen_port)
	fmt.Printf("Listening on %s...", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
