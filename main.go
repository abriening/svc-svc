package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"time"
)

func main() {
	port := flag.Int("port", 0, "Port to listen on (default \"0\", lets the OS assign a free port)")
	var upstreamUrl string
	flag.StringVar(&upstreamUrl, "upstream", "", "a fqdn")
	flag.Parse()

	ping := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(204)
	}

	upstream := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		resp, err := http.Get(upstreamUrl)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		var data interface{}
		data_err := json.Unmarshal(body, &data)
		if data_err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			return
		}

		result := map[string]interface{}{name: data}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
	}

	random := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200)
		length := rand.Intn(10)
		a := make([]int, length)
		for i := 0; i < length; i++ {
			a[i] = rand.Intn(10000)
		}
		if err := json.NewEncoder(w).Encode(a); err != nil {
			panic(err)
		}
	}

	logError := func(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					log.Print("ERROR: ", r, string(debug.Stack()))
					w.WriteHeader(500)
				}
			}()
			f(w, r)
		}
	}

	logRequest := func(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			defer func() {
				log.Print("[" + r.Method + "] " + r.URL.String() + " " + time.Since(t).String())
			}()
			f(w, r)
		}
	}

	if upstreamUrl != "" {
		http.HandleFunc("/upstream", logError(logRequest(upstream)))
	}
	http.HandleFunc("/ping", logError(logRequest(ping)))
	http.HandleFunc("/random", logError(logRequest(random)))
	log.Print(fmt.Sprintf("starting on localhost:%d", *port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
