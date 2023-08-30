package main

import (
	"errors"
	"flag"
	"log"
	"time"

	maersk "github.com/amirhnajafiz/maersk/internal"
)

const (
	EmptyStr = ""
)

var (
	ErrEmptyOutput = errors.New("output field cannot be empty")
	ErrURL         = errors.New("url field cannot be empty")
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

	// check flags
	if *OutPutFlag == EmptyStr {
		log.Fatal(ErrEmptyOutput)
	}

	if *URLFlag == EmptyStr {
		log.Fatal(ErrURL)
	}

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

	log.Println("request is build")

	// start download
	if err := center.Ship(); err != nil {
		log.Fatal(err)
	}

	// print last report
	report := center.Reports()
	log.Printf(
		"created: %s, downloads: %d out of %d with %d failures\n",
		report.Created,
		report.DownloadedChunks,
		report.NumberOfChunks,
		report.NumberOfErrors,
	)
}
