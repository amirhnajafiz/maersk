package maersk

import (
	"errors"
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
	created time.Time
	chunks  [][]byte
	failed  int
}

func (c *Cargo) Ship() error {
	var (
		// size is the number of jobs
		size int
		// the size of our file
		fileSize int
		// channel is for sending chunks
		channel = make(chan chunk, c.Workers)
		// jobs channel
		jobs = make(chan job, c.Chunks)
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

	// getting the header of content length
	if header, ok := resp.Header["Content-Length"]; ok {
		fileSize, err = strconv.Atoi(header[0])
		if err != nil {
			return err
		}

		size = fileSize / c.Chunks
	} else {
		return errors.New("new error")
	}

	// generating workers
	for i := 0; i < c.Workers; i++ {
		// creating one worker
		w := worker{
			channel: channel,
			jobs:    jobs,
			url:     c.URL,
		}
		// starting worker
		go func(j int) {
			err := w.work()
			if err != nil {
				log.Println(err)
			}
		}(i)
	}

	// create jobs based on the number of chunks
	for i := 0; i < c.Chunks; i++ {
		j := job{
			index: i,
			size:  size,
			last:  i == c.Chunks-1,
		}

		jobs <- j
	}

	// counting until we get each of the chunks
	counter := 0
	for part := range channel {
		counter++

		c.chunks[part.index] = part.data
		if counter == c.Chunks {
			break
		}
	}

	// storing the chunks into the output file
	file := make([]byte, fileSize)

	// appending into fies array
	for _, part := range c.chunks {
		file = append(file, part...)
	}

	// Set permissions accordingly, 0700 may not be the best choice
	if err = ioutil.WriteFile(c.Out, file, 0700); err != nil {
		return err
	}

	return nil
}
