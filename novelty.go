package novelty

import (
	"io"

	"github.com/jackmordaunt/novelty/protocol"
	"github.com/pkg/errors"
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
	Status(*protocol.Status)
}

// Status contains information about the status of a Resource.
type Status struct{}

// Show contains the meta data for a TV Show episode or movie.
type Show struct {
	Name string
	// URI must be a valid form of one of the supported protocols.
	// That includes http, ftp, and magnet links.
	// TODO: Make this more strongly typed.
	URI string
}

// Engine provides the high level API.
type Engine struct{}

// Open the show's resource. This begins the download.
func (e Engine) Open(s Show) (Resource, error) {
	// FIXME: Let's just assume that URI is a magnet link for now.
	cfg := protocol.NewClientConfig()
	cfg.TorrentPath = s.URI
	client, err := protocol.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "opening torrent")
	}
	return client, nil
}
