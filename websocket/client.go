package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/gorilla/websocket"
)

// Client is the client side of a websocket interaction.
type Client struct {
	socket  *websocket.Conn
	updates chan torrent.TorrentStats
}

// Close the client connection.
func (c *Client) Close() {
	c.socket.Close()
	close(c.updates)
}

func (c *Client) loop() {
	for update := range c.updates {
		if err := c.socket.WriteJSON(update); err != nil {
			log.Printf("[websocket client] writing: %v", err)
			break
		}
	}
	if err := c.socket.Close(); err != nil {
		log.Printf("[websocket client] closing: %v", err)
	}
}

// Torrent is a handle that provides access to a torrent's information via
// ticked updates.
type Torrent struct {
	Handle  *torrent.Torrent
	clients map[*Client]bool
	join    chan *Client
	leave   chan *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (t *Torrent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[Torrent] ServeHTTP: %v", err)
		return
	}
	client := &Client{
		socket:  socket,
		updates: make(chan torrent.TorrentStats),
	}
	t.join <- client
	defer func() { t.leave <- client }()
	client.loop()
}

func (t *Torrent) loop() {
	if t.clients == nil {
		t.clients = make(map[*Client]bool)
	}
	if t.join == nil {
		t.join = make(chan *Client)
	}
	if t.leave == nil {
		t.leave = make(chan *Client)
	}
	update := time.NewTicker(time.Second / 60)
	for {
		select {
		case client := <-t.join:
			t.clients[client] = true
		case client := <-t.leave:
			delete(t.clients, client)
			client.Close()
		case <-update.C:
			for client := range t.clients {
				client.updates <- t.Handle.Stats()
			}
		}
	}
}
