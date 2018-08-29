package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Sender can send messages.
type Sender interface {
	Send(msg Message)
}

// Handler handles incoming websocket messages.
type Handler func(client Sender, payload Payload)

// Router provides dynamic routing of websocket messages to handlers.
type Router struct {
	upgrader *websocket.Upgrader
	rules    map[string]Handler
}

// NewRouter returns an initialised Router.
// A default websocket upgrader will be used if none is provided.
func NewRouter(upgrader *websocket.Upgrader) *Router {
	if upgrader == nil {
		upgrader = defaultUpgrader()
	}
	return &Router{
		upgrader: upgrader,
		rules:    make(map[string]Handler),
	}
}

// Handle registers a handler for a given message name.
func (mux *Router) Handle(name string, handler Handler) {
	if _, ok := mux.rules[name]; ok {
		panic(fmt.Sprintf("attempted to register a second handler for %q: handler already registered", name))
	}
	mux.rules[name] = handler
}

// FindHandler returns a handler for the given message name.
func (mux *Router) FindHandler(name string) (Handler, bool) {
	handler, ok := mux.rules[name]
	return handler, ok
}

// ErrorHandler handles websocket errors.
// The client will disconnect after the first error encountered.
// Errors handled here are websocket infrastructure errors, not application
// errors.
// NOTE: This would be a good candidate for telementry to analyse connection
// failures.
func (mux *Router) ErrorHandler(err error) {
	log.Printf("[websocket] error: %v", err)
}

func (mux *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	socket, err := mux.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not establish websocket connection", http.StatusInternalServerError)
		return
	}
	client := NewClient(socket, mux.FindHandler, mux.ErrorHandler)
	client.Run()
}

func defaultUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
}
