package internal

import "time"

// ShippingOrder
// stores the information needed for creating cargo.
type ShippingOrder struct {
	Out     string
	URL     string
	Mode    string
	Workers int
	Chunks  int
	Timeout time.Duration
}

// Build
// generates a new cargo center with given orders.
func Build(order *ShippingOrder) *Cargo {
	return &Cargo{
		Out:     order.Out,
		URL:     order.URL,
		Mode:    order.Mode,
		Workers: order.Workers,
		Chunks:  order.Chunks,
		Timeout: order.Timeout,
	}
}
