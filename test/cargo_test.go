package test

import (
	"os"
	"testing"
	"time"

	maersk "github.com/amirhnajafiz/maersk/internal"
)

// TestCargo
// testing cargo initialization and file downloading feature.
func TestCargo(t *testing.T) {
	url := "http://212.183.159.230/5MB.zip"
	out := "5MB.zip"

	defer func() {
		if err := os.RemoveAll(out); err != nil {
			t.Error(err)
		}
	}()

	order := maersk.ShippingOrder{
		Out:     out,
		URL:     url,
		Mode:    maersk.DEBUG,
		Workers: 5,
		Chunks:  5,
		Timeout: 5 * time.Second,
	}

	center := maersk.Build(&order)

	if err := center.Ship(); err != nil {
		t.Error(err)
	}

	t.Log("succeed")
}
