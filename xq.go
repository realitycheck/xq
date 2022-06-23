package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	addr        = "0.0.0.0:9091"
	excuseURL   = "http://www.programmerexcuses.com/"
	excuseRegex = "<center .*><a .*>(.*)</a></center>"
	excuseRe    *regexp.Regexp
	serve       bool
)

type Excuse struct {
	Text string `json:"text"`
}

func main() {
	flag.StringVar(&addr, "addr", addr, "Server address")
	flag.StringVar(&excuseURL, "url", excuseURL, "Url to excuses site")
	flag.StringVar(&excuseRegex, "regex", excuseRegex, "Regex for excuses text")
	flag.BoolVar(&serve, "s", serve, "Run server")
	flag.Parse()

	excuseRe = regexp.MustCompile(excuseRegex)

	if serve {
		runServer()
	} else {
		excuse, err := fetchExcuse()
		if err != nil {
			log.Fatalf("xq: %v", err)
		}
		os.Stdout.WriteString(excuse.Text + "\n")
	}
}

func runServer() {
	http.HandleFunc("/xq", func(w http.ResponseWriter, r *http.Request) {
		excuse, err := fetchExcuse()
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(excuse)
		if err != nil {
			log.Fatalf("xq: %v", err)
		}
	})

	log.Println("Listening...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func fetchExcuse() (*Excuse, error) {
	resp, err := http.Get(excuseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Excuse{
		Text: string(excuseRe.FindSubmatch(body)[1]),
	}, nil
}
