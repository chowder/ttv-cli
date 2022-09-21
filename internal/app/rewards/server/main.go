package server

import (
	_ "embed"
	"log"
	"net/http"
)

//go:embed index.html
var indexHtml []byte

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(indexHtml)
}

func Run(addr string, streamer string) {
	if len(streamer) == 0 {
		log.Fatalln("Streamer name not provided")
	}

	hub := newHub()
	go hub.run(streamer)

	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("new connection received")
		serveWebsocket(hub, w, r)
	})

	log.Println("Serving on:", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
