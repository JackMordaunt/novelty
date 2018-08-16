package main

import (
	"github.com/zserge/webview"
)

func main() {
	// url := os.Args[1]
	wv := webview.New(webview.Settings{
		Title:     "novelty",
		URL:       "http://localhost:8080",
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     true,
	})
	wv.Run()
}
