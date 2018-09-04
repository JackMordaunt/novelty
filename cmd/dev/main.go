package main

import (
	"flag"
	"log"
	"sync"

	"github.com/jackmordaunt/novelty"
	"github.com/jackmordaunt/novelty/http"
	"github.com/zserge/webview"
)

func main() {
	var (
		headless bool
		uiURL    string
	)
	flag.BoolVar(&headless, "headless", false, "Run without the UI.")
	flag.StringVar(&uiURL, "ui", "http://127.0.0.1:8080", "URL that serves the web UI.")
	flag.Parse()
	work := &sync.WaitGroup{}
	work.Add(1)
	go func() {
		defer work.Done()
		uc := &http.UseCases{
			Engine: &novelty.Engine{},
		}
		err := uc.ListenAndServe(":9090")
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
