package maersk

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

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

func (c *Cargo) Ship() error {
	var (
		// size is the number of jobs
		size int
		// size of the file
		fileSize int
		// channel is for sending chunks
		channel = make(chan chunk, c.Workers)
	)

	// generating the list of chunks
	c.chunks = make([][]byte, c.Chunks)

	// creating a http request to get the file information
	client := &http.Client{}

	req, err := http.NewRequest("HEAD", c.URL, nil)
	if err != nil {
		return err
	}

	// making http request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if header, ok := resp.Header["Content-Length"]; ok {
		fileSize, err = strconv.Atoi(header[0])
		if err != nil {
			return err
		}

		size = fileSize / c.Workers
	} else {
		log.Fatal("File size was not provided!")
	}

	for i := 0; i < c.Workers; i++ {
		w := worker{
			channel: channel,
			url:     c.URL,
		}

		go func(j int) {
			err := w.work(j, size, j == c.Workers-1)
			if err != nil {
				log.Println(err)
			}
		}(i)
	}

	counter := 0

	for part := range channel {
		counter++

		c.chunks[part.index] = part.data
		if counter == c.Workers {
			break
		}
	}

	file := make([]byte, 1024)

	for _, part := range c.chunks {
		file = append(file, part...)
	}

	// Set permissions accordingly, 0700 may not
	// be the best choice
	err = ioutil.WriteFile("./data.zip", file, 0700)

	if err != nil {
		return err
	}

	return nil
}
