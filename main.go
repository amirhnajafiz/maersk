package main

import (
	"flag"
	"time"

	maersk "github.com/amirhnajafiz/maersk/internal"
)

const (
	EmptyStr = ""
)

func main() {
	var (
		OutPutFlag  = flag.String("output", EmptyStr, "output file name")
		URLFlag     = flag.String("url", EmptyStr, "url address of the file you want")
		WorkersFlag = flag.Int("workers", 1, "number of workers")
		ChunksFlag  = flag.Int("chunks", 10, "number of chunks")
		TimeoutFlag = flag.Int("timeout", 10, "timeout in seconds")
		ModeFlag    = flag.String("mode", maersk.INFO, "debug mode")
	)

	// parse flags
	flag.Parse()

	// create order
	order := &maersk.ShippingOrder{
		Chunks:  *ChunksFlag,
		Workers: *WorkersFlag,
		URL:     *URLFlag,
		Out:     *OutPutFlag,
		Timeout: time.Duration(*TimeoutFlag) * time.Second,
		Mode:    *ModeFlag,
	}

	// create center
	center := maersk.Build(order)
}
