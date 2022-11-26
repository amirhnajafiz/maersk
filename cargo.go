package maersk

import (
	"fmt"
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

// Ship
// starts the workers and gets the chunks to
// build the output file.
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
		// failed jobs channel
		failed = make(chan job)
		// kill switch channel
		killSwitch = make(chan int)
	)

	// building the crane
	crane := crane{
		jobs:   jobs,
		failed: failed,
		program: func(id int) error {
			c.failed++

			if c.Mode == DEBUG || c.Mode == INFO {
				log.Printf("chunk failed:\n\tindex: %d\n", id)
			}

			return nil
		},
	}

	// starting crane
	go func() {
		if err := crane.start(); err != nil {
			if c.Mode != OFF {
				log.Printf("no crane:\n\t%v\n", err)
			}
		}
	}()

	// generating the list of chunks
	c.chunks = make([][]byte, c.Chunks)

	// creating a http request to get the file information
	client := &http.Client{}

	req, err := http.NewRequest("HEAD", c.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create http request: %v", err)
	}

	// making http request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make http request: %v", err)
	}

	// getting the header of content length
	if header, ok := resp.Header["Content-Length"]; ok {
		fileSize, err = strconv.Atoi(header[0])
		if err != nil {
			return fmt.Errorf("failed to get content length: %v", err)
		}

		size = fileSize / c.Chunks
	} else {
		return fmt.Errorf("failed to get content length")
	}

	// generating workers
	for i := 0; i < c.Workers; i++ {
		// creating one worker
		w := worker{
			channel:    channel,
			jobs:       jobs,
			failed:     failed,
			killSwitch: killSwitch,
			timeout:    c.Timeout,
			url:        c.URL,
		}
		// starting worker
		go func(j int) {
			if e := w.work(); e != nil {
				if c.Mode != OFF {
					log.Println(fmt.Errorf("failed to start worker:\n\t%v\n", e))
				}
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
		return fmt.Errorf("failed to assemble the chunks: %v", err)
	}

	return nil
}

// Reports
// returns a status of cargo.
func (c *Cargo) Reports() *Report {
	return &Report{
		Created:          c.created,
		NumberOfChunks:   c.Chunks,
		DownloadedChunks: c.numberOfDownloads(),
		NumberOfErrors:   c.failed,
	}
}

// numberOfDownloads
// counts the number of downloaded chunks.
func (c *Cargo) numberOfDownloads() int {
	counter := 0

	for index := range c.chunks {
		if c.chunks[index] != nil {
			counter++
		}
	}

	return counter
}
