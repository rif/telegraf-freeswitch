package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rif/telegraf-freeswitch/utils"
)

var (
	host          = flag.String("host", "localhost", "freeswitch host address")
	port          = flag.Int("port", 8021, "freeswitch port")
	pass          = flag.String("pass", "ClueCon", "freeswitch password")
	serve         = flag.Bool("serve", false, "run as a server")
	execd         = flag.Bool("execd", false, "run as an execd server")
	listenAddress = flag.String("listen_address", "127.0.0.1", "listen on address")
	listenPort    = flag.Int("listen_port", 9191, "listen on port")
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
		if *execd {
			reader := bufio.NewReader(os.Stdin)
			for {
				text, err := reader.ReadString('\n')
				if err != nil {
					l.Print("error reading from stdin: ", err)
					continue
				}
				if strings.TrimSpace(text) != "" {
					break
				}
				if err := fetcher.GetData(); err != nil {
					l.Print(err.Error())
				}
				fmt.Print(fetcher.FormatOutput(utils.InfluxFormat))
			}
			os.Exit(0)
		}
		if err := fetcher.GetData(); err != nil {
			l.Print(err.Error())
		}
		fmt.Print(fetcher.FormatOutput(utils.InfluxFormat))
		os.Exit(0)
	}

	http.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		if err := fetcher.GetData(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		status, _ := fetcher.FormatOutput(utils.JSONFormat)
		if _, err := w.Write([]byte(status)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/profiles/", func(w http.ResponseWriter, r *http.Request) {
		if err := fetcher.GetData(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, profiles := fetcher.FormatOutput(utils.JSONFormat)
		if _, err := w.Write([]byte(profiles)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	listen := fmt.Sprintf("%s:%d", *listenAddress, *listenPort)
	fmt.Printf("Listening on %s...", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
