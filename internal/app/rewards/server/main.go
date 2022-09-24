package server

import (
	_ "embed"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"ttv-cli/internal/pkg/twitch/gql/query/channel"
)

//go:embed index.html
var indexHtml []byte

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("alive"))
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(indexHtml)
}

func Run(addr string) {
	hubsByStreamer := make(map[string]*Hub)
	mutex := sync.Mutex{}

	http.DefaultServeMux.HandleFunc("/", serveHome)

	http.DefaultServeMux.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		streamer := strings.TrimPrefix(r.URL.Path, "/ws/")
		log.Println("New connection from:", r.RemoteAddr, "for streamer:", streamer)

		hub, ok := hubsByStreamer[streamer]
		if !ok {
			hub, err := validateAndMakeNewHub(streamer)
			if err != nil {
				log.Println("Could not create hub for streamer:", streamer, ", error:", err)
				http.Error(w, "streamer not found", http.StatusNotFound)
				return
			}
			log.Println("Making new hub for streamer:", streamer)
			hubsByStreamer[streamer] = hub
			serveWebsocket(hub, w, r)
		} else {
			serveWebsocket(hub, w, r)
		}
	})

	log.Println("Serving on:", addr)
	err := http.ListenAndServe(addr, http.DefaultServeMux)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

// validateAndMakeNewHub first validates that the streamer is valid, and returns a hub if so
func validateAndMakeNewHub(streamer string) (*Hub, error) {
	c, err := channel.GetChannel(streamer)
	if err != nil {
		return nil, err
	}
	if len(c.Name) == 0 {
		return nil, errors.New("streamer does not exist")
	}

	hub := newHub()
	go hub.run(streamer)
	return hub, nil
}
