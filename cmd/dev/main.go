package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jackmordaunt/novelty"
	"github.com/jackmordaunt/novelty/websocket"
	"github.com/zserge/webview"
)

func main() {
	var (
		headless bool
		uiURL    string
	)
	flag.BoolVar(&headless, "headless", false, "Run without the UI.")
	flag.StringVar(&uiURL, "ui", "http:/127.0.0.1:8080", "URL that serves the web UI.")
	flag.Parse()
	work := &sync.WaitGroup{}
	work.Add(1)
	go func() {
		defer work.Done()

		// What api do you want?
		// n := novelty.NewEngine()
		// show := r.Body
		// resource := n.Open(show)
		//
		n := &novelty.Engine{}
		r := websocket.NewRouter(nil)
		r.Handle("open-show", func(s websocket.Sender, p websocket.Payload) {
			type Cmd struct {
				Name string `json:"name"`
				URI  string `json:"uri"`
			}
			var cmd Cmd
			if err := p.Bind(&cmd); err != nil {
				log.Fatalf("%v", err)
			}
			show := novelty.Show(cmd)
			resource, err := n.Open(show)
			if err != nil {
				log.Fatalf("%v", err)
			}
			// Store the resource in some object that exposes it via http.
			http.Handle(fmt.Sprintf("/%s", show.Name), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Need filename of resource.
				w.Header().Set("Content-Disposition", "attachment; filename=\""+show.Name+"\"")
				http.ServeContent(w, r, "", time.Now(), resource)
			}))
			type Response struct {
				StreamURL string `json:"stream_url"`
			}
			response := Response{
				StreamURL: fmt.Sprintf("/%s", show.Name),
			}
			payload, err := json.Marshal(response)
			if err != nil {
				log.Fatalf("marshalling json: %v", err)
			}
			s.Send(websocket.Message{
				Name: "show-opened",
				Data: payload,
			})

		})
		// r.Handle(torrentStatusSubscribe, handleTorrentStatusSubscribe)
		// r.Handle("torrent status - subscribe", &TorrentStatusSubscribeHandler{})
		http.Handle("/", r)
		err := http.ListenAndServe(":9090", r)
		if err != nil {
			log.Fatalf("[server] error: %v", err)
		}
	}()
	if !headless {
		work.Add(1)
		go func() {
			defer work.Done()
			wv := webview.New(webview.Settings{
				Title:     "novelty",
				URL:       uiURL,
				Width:     800,
				Height:    600,
				Resizable: true,
				Debug:     true,
			})
			wv.Run()
		}()
	}
	work.Wait()
}
