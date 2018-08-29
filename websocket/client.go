package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// FindHandler maps the message name to a handler for that message.
type FindHandler func(string) (Handler, bool)

// ErrorHandler receives any errors that occur during message reading or writing.
type ErrorHandler func(error)

// Message contains the name and payload for a websocket message.
type Message struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

// Client manages the websocket connection between the server and a user.
type Client struct {
	socket       *websocket.Conn
	send         chan Message
	done         chan struct{}
	findHandler  FindHandler
	errorHandler ErrorHandler
	once         sync.Once
}

// NewClient returns an initialised Client.
func NewClient(socket *websocket.Conn, finder FindHandler, onErr ErrorHandler) *Client {
	return &Client{
		socket:       socket,
		send:         make(chan Message),
		done:         make(chan struct{}),
		findHandler:  finder,
		errorHandler: onErr,
		once:         sync.Once{},
	}
}

// Run starts the pumps and blocks until the connection errors out.
func (c *Client) Run() {
	defer c.socket.Close()
	go c.readPump()
	go c.writePump()
	<-c.done
}

// Send a message to the client.
func (c *Client) Send(msg Message) {
	c.send <- msg
}

// readPump reads messages from the client socket as they come in.
func (c *Client) readPump() {
	var msg Message
	for {
		if err := c.socket.ReadJSON(&msg); err != nil {
			c.error(errors.Wrap(err, "reading message"))
			break
		}
		if handler, ok := c.findHandler(msg.Name); ok {
			handler(c, Payload{Raw: msg.Data})
		}
	}
}

// writePump writes messages to the client socket as they come in.
func (c *Client) writePump() {
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			c.error(errors.Wrap(err, "writing message"))
			break
		}
	}
}

func (c *Client) error(err error) {
	if c.errorHandler != nil {
		c.errorHandler(err)
	}
	c.once.Do(func() {
		close(c.done)
	})
}

// Payload provides a convenient way to bind json to a value.
type Payload struct {
	Raw []byte
}

// Bind json data to value.
func (p Payload) Bind(v interface{}) error {
	return json.Unmarshal(p.Raw, v)
}
