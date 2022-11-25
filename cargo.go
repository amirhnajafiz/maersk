package maersk

import "time"

// Cargo
// is the main module of maersk. In cargo the operation
// of downloading chunks with workers is handled.
type Cargo struct {
	Out     string
	URL     string
	Mode    string
	Workers int
	Chunks  int
	Timeout time.Duration
	chunks  [][]byte
}
