package main

import (
	"log"

	"github.com/amirhnajafiz/maersk"
)

func main() {
	center := maersk.Cargo{
		URL:     "someplace.com",
		Out:     "file.zip",
		Workers: 2,
		Chunks:  2,
	}

	log.Println(center.Ship())
}
