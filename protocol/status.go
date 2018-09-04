package protocol

// Status containst he status information of a given resource.
// This object is effectively a union type which covers all supported protocols.
type Status struct {
	Progress   int64  `json:"progress"`
	Uploaded   int64  `json:"uploaded"`
	Throughput string `json:"throughput"`
	Size       int64  `json:"size"`
}
