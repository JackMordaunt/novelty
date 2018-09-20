package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackmordaunt/novelty"
	"github.com/jackmordaunt/novelty/protocol"
	"github.com/jackmordaunt/novelty/websocket"
	"github.com/pkg/errors"
)

// UseCases adapts websocket messages to novelty usecases.
type UseCases struct {
	Engine    *novelty.Engine
	Resources *Resources
}

// ListenAndServe listens for tcp connections and servers the usecases.
func (ws *UseCases) ListenAndServe(addr string) error {
	if ws.Engine == nil {
		return fmt.Errorf("novelty engine cannot be nil")
	}
	if ws.Resources == nil {
		ws.Resources = NewResources()
	}
	wsRouter := websocket.NewRouter(nil)
	wsRouter.Handle("open-show", ws.openShow)
	router := mux.NewRouter()
	router.Handle("/ws", wsRouter)
	router.Handle("/stream/{file-name}", ws.Resources)
	return http.ListenAndServe(addr, router)
}

func (ws *UseCases) openShow(s websocket.Sender, p websocket.Payload) {
	fmt.Printf("openShow\n")
	type Cmd struct {
		Name string `json:"name"`
		URI  string `json:"uri"`
	}
	var cmd Cmd
	if err := p.Bind(&cmd); err != nil {
		log.Fatalf("%v", err)
	}
	show := novelty.Show(cmd)
	show.Name = format(show.Name)
	r, err := ws.Engine.Open(show)
	if err != nil {
		log.Fatalf("%v", err)
	}
	closed := make(chan struct{})
	ws.Resources.Store(show.Name, resource{
		Resource: r,
		closed: func() {
			close(closed)
		},
	})
	type Response struct {
		StreamURL string `json:"stream_url"`
	}
	response := Response{
		StreamURL: fmt.Sprintf("stream/%s", show.Name),
	}
	payload, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("marshalling json: %v", err)
	}
	s.Send(websocket.Message{
		Name: "show-opened",
		Data: payload,
	})
	go func() {
		updates := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-updates.C:
				var status novelty.Status
				r.Status(&status)
			payload, err := json.Marshal(status)
			if err != nil {
				panic(errors.Wrap(err, "marshalling status update"))
			}
			// FIXME: When to stop pushing updates?
			// - return a bool from Send `if ok := s.Send(..); !ok {...}`
			s.Send(websocket.Message{
				Name: "show-status-update",
				Data: payload,
			})
			case <-closed:
				break
			}
		}
	}()
}

// resource decorates the novelty.Resource with a close handler.
// This allows us to respond to close events locally, in order to stop sending
// status updates.
type resource struct {
	novelty.Resource
	closed func()
}

func (r resource) Close() error {
	if r.closed != nil {
		r.closed()
	}
	return r.Resource.Close()
}

func format(name string) string {
	name = strings.Replace(name, " ", "-", -1)
	name = strings.Replace(name, "/", "", -1)
	name = strings.Replace(name, `\`, "", -1)
	return strings.ToLower(name)
}
