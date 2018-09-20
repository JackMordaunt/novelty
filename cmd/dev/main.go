package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/jackmordaunt/novelty"
	"github.com/jackmordaunt/novelty/http"
	"github.com/jackmordaunt/novelty/protocol"
	"github.com/zserge/webview"
)

func main() {
	var (
		headless bool
		develop  bool
	)
	flag.BoolVar(&headless, "headless", false, "Run without the UI.")
	flag.BoolVar(&develop, "develop", false, "Run in development mode.")
	flag.Parse()
	work := &sync.WaitGroup{}
	work.Add(1)
	go func() {
		defer work.Done()
		engine := &novelty.Engine{}
		// Register torrent protocol handler.
		engine.Register(func(s novelty.Show) (novelty.Resource, bool) {
			cfg := protocol.NewClientConfig()
			cfg.TorrentPath = s.URI
			client, err := protocol.NewClient(cfg)
			if err != nil {
				return nil, false
			}
			return client, true
		})
		uc := &http.UseCases{
			Engine: engine,
		}
		if err := uc.ListenAndServe(":9090"); err != nil {
			log.Fatalf("[server] error: %v", err)
		}
	}()
	if !headless {
		work.Add(1)
		go func() {
			defer work.Done()
			fmt.Printf("new webview\n")
			wv := webview.New(webview.Settings{
				Title:     "novelty",
				URL:       "http://127.0.0.1:8080",
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
