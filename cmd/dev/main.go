package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackmordaunt/novelty"
	"github.com/pkg/errors"
)

func main() {
	// n := novelty.NewClient()
	srv := &MediaServer{}
	err := http.ListenAndServe(":9090", srv)
	if err != nil {
		panic(err)
	}
	// url := os.Args[1]
	// wv := webview.New(webview.Settings{
	// 	Title:     "novelty",
	// 	URL:       "http://127.0.0.1:8080",
	// 	Width:     800,
	// 	Height:    600,
	// 	Resizable: true,
	// 	Debug:     true,
	// })
	// wv.Run()
}

// ID is an integer representation of an ID.
type ID int64

// MediaServer transforms a torrent into a file for the client to consume.
type MediaServer struct {
	nextID           ID
	active           map[ID]*novelty.Client
	statsConnections map[ID][]*websocket.Conn
	upgrader         *websocket.Upgrader

	isRunning bool
	loopMutex *sync.Mutex
}

func (s *MediaServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.upgrader == nil {
		s.upgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
	}
	if s.loopMutex == nil {
		s.loopMutex = &sync.Mutex{}
	}
	if s.active == nil {
		s.active = make(map[ID]*novelty.Client)
		s.statsConnections = make(map[ID][]*websocket.Conn)
	}
	defer r.Body.Close()
	if strings.Contains(r.URL.String(), "/open") {
		if r.Method != "POST" {
			w.WriteHeader(400)
			w.Write([]byte("expected POST request"))
			return
		}
		if err := s.Open(w, r); err != nil {
			s.error(w, r, err)
		}
		return
	} else if strings.Contains(r.URL.String(), "/stats") {
		if err := s.Stats(w, r); err != nil {
			s.error(w, r, err)
		}
		return
	} else if strings.Contains(r.URL.String(), "/stream") {
		if err := s.Stream(w, r); err != nil {
			s.error(w, r, err)
		}
		return
	}
	fmt.Printf("um, what\n")
}

// Open the provided torrent.
// Starts downloading.
// Returns two 3 endpoints
// 1. `<torrent-id>/stream` which provides the video file for download.
// 2. `<torrent-id>/stats` which provides a websocket connection that pushes updates about
// 	the torrent download.
// 3. `<torrent-id>/close` which stops the torrent download.
func (s *MediaServer) Open(w http.ResponseWriter, r *http.Request) error {
	defer func() {
		s.loopMutex.Lock()
		go s.loop()
		s.loopMutex.Unlock()
	}()
	config := novelty.NewClientConfig()
	torrentPath, err := loadTorrentPath(r.Body)
	if err != nil {
		return err
	}
	config.TorrentPath = torrentPath
	client, err := novelty.NewClient(config)
	if err != nil {
		return errors.Wrap(err, "initialising torrent client")
	}
	s.nextID++
	s.active[s.nextID] = &client
	type response struct {
		Stream string `json:"stream"`
		Stats  string `json:"stats"`
		Close  string `json:"close"`
	}
	resp := response{
		Stream: fmt.Sprintf("/%d/stream", s.nextID),
		Stats:  fmt.Sprintf("/%d/stats", s.nextID),
		Close:  fmt.Sprintf("/%d/close", s.nextID),
	}
	return writeJSON(w, resp)
}

// Stream from a given torrent client.
// `/<id>/stream`
func (s *MediaServer) Stream(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return errors.Wrap(err, "parsing ID")
	}
	if _, ok := s.active[id]; !ok {
		return errors.Wrapf(err, "no active client for id (%d)", id)
	}
	client := s.active[id]
	client.GetFile(w, r)
	return nil
}

// Stats opens a connection that periodically pushes statistics for a given client.
func (s *MediaServer) Stats(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return err
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.Wrap(err, "upgrading connection")
	}
	s.statsConnections[id] = append(s.statsConnections[id], conn)
	return nil
}

func (s *MediaServer) loop() {
	if s.isRunning {
		return
	}
	s.isRunning = true
	defer func() {
		s.isRunning = false
	}()
	update := time.NewTicker(time.Millisecond * 16)
	for {
		select {
		case <-update.C:
			for id, conns := range s.statsConnections {
				client := s.active[id]
				for ii, conn := range conns {
					err := conn.WriteJSON(client.Torrent.Stats())
					if err != nil {
						client.Close()
						conns = append(conns[:ii], conns[ii+1:]...)
						log.Printf("[error] [stats] [%d] %v\n", id, err)
					}
				}
			}
		}
	}
}

func (s MediaServer) error(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	w.WriteHeader(500)
	w.Write([]byte(errors.Wrap(err, "error occured during request").Error()))
}

func loadTorrentPath(r io.Reader) (string, error) {
	data := map[string]string{}
	if err := loadJSON(r, &data); err != nil {
		return "", errors.Wrap(err, "loading json from request body")
	}
	torrentPath, ok := data["torrent_path"]
	if !ok {
		return "", fmt.Errorf("invalid json payload; missing 'torrent_path'")
	}
	return torrentPath, nil
}

func loadJSON(r io.Reader, v interface{}) error {
	by, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(by, v)
}

func writeJSON(w io.Writer, v interface{}) error {
	by, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, bytes.NewBuffer(by)); err != nil {
		return err
	}
	return nil
}

// parseID will return the last integer fragment in the URL.
// Very primitive.
func parseID(r *http.Request) (ID, error) {
	var ret = -1
	for _, fragment := range strings.Split(r.URL.String(), "/") {
		if number, err := strconv.Atoi(fragment); err == nil {
			ret = number
		}
	}
	if ret == -1 {
		return ID(ret), fmt.Errorf("no ID fragment found in URL")
	}
	return ID(ret), nil
}
