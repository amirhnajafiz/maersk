package main

import (
	"fmt"
	"log"
	"time"

	"github.com/amirhnajafiz/maersk"
)

func main() {
	order := &maersk.ShippingOrder{
		Chunks:  10,
		Workers: 10,
		URL:     "somewhere.com",
		Out:     "data.zip",
		Timeout: 5 * time.Second,
		Mode:    maersk.DEBUG,
	}

	center := maersk.Build(order)
	defer func() {
		if e := center.Cancel(); e != nil {
			log.Printf("center kill switch failed: %v\n", e)
		}
	}()

	go func() {
		if err := center.Ship(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second * 3)

	report := center.Reports()
	fmt.Printf("created: %s, downloads: %d\n", report.Created, report.DownloadedChunks)

	select {}
}
