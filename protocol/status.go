package protocol

import "github.com/anacrolix/torrent"

// Status containst he status information of a given resource.
// This object is effectively a union type which covers all supported protocols.
type Status struct {
	*torrent.TorrentStats
}
