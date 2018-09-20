package novelty

import (
	"fmt"
	"io"
)

// Resource represents the media file belonging to a Show.
// It can be streamed, meaning partially downloaded, such that Seeks passed the
// downloaded data must prioritise the new position for download.
// Close should cancel the download if it hasn't finished.
type Resource interface {
	io.Reader
	io.Seeker
	io.Closer
	// Status populates the Status object with the current status.
	// A pointer is passed in to avoid object allocations.
	Status(*Status)
}

// Status contains data about the current status of a streaming show.
type Status struct {
	Progress         int64  `json:"progress"`
	Uploaded         int64  `json:"uploaded"`
	Throughput       string `json:"throughput"`
	Size             int64  `json:"size"`
	ReadyForPlayback bool   `json:"ready_for_playback"`
}

// Show contains the meta data for a TV Show episode or movie.
type Show struct {
	Name string
	// URI must be a valid form of one of the supported protocols.
	// That includes http, ftp, and magnet links.
	// TODO: Make this more strongly typed.
	URI string
}

// ProtocolHandler creates a Resource from the Show using an appriate protocol
// implementation, if registered.
type ProtocolHandler func(s Show) (r Resource, ok bool)

// Engine provides the high level API.
type Engine struct {
	protocols []ProtocolHandler
}

// Register a function that handles a particular protocol.
func (e *Engine) Register(ph ProtocolHandler) {
	e.protocols = append(e.protocols, ph)
}

// Open the show's resource using a registered protocol handler.
func (e *Engine) Open(s Show) (Resource, error) {
	for _, ph := range e.protocols {
		if resource, ok := ph(s); ok && resource != nil {
			return resource, nil
		}
	}
	return nil, fmt.Errorf("no valid handler registered for '%s'", s.URI)
}
