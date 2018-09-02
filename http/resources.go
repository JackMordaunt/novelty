package http

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackmordaunt/novelty"
)

// Resources handles serving resources.
type Resources struct {
	active *sync.Map
}

// NewResources initialises a resource handler.
func NewResources() *Resources {
	return &Resources{
		active: &sync.Map{},
	}
}

func (h Resources) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, ok := mux.Vars(r)["file-name"]
	if !ok {
		fmt.Printf("expected show name in url\n")
		return
	}
	resource, ok := h.Load(name)
	if !ok {
		fmt.Printf("show name %q does not exist\n", name)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	http.ServeContent(w, r, "", time.Now(), resource)
}

// Store the resource at the key.
func (h Resources) Store(key string, r novelty.Resource) {
	h.active.Store(key, r)
}

// Load the resource from the key.
func (h Resources) Load(key string) (novelty.Resource, bool) {
	r, ok := h.active.Load(key)
	if !ok {
		return nil, false
	}
	re, ok := r.(novelty.Resource)
	if !ok {
		return nil, false
	}
	return re, true
}
