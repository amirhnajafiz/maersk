package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Manages to download each chunk from given url.
type ship struct {
	// host url for downloading file.
	url string
	// channel for sending the downloaded chunk.
	channel chan chunk
	// input channel for getting the jobs.
	jobs chan job
	// failed jobs channel
	failed chan job
	// kill switch channel
	killSwitch chan int
	// timeout time of processing
	timeout time.Duration
}

// berth function starts the ship worker to listen
// on jobs channel.
func (s *ship) berth() error {
	// listen on jobs channel
	for j := range s.jobs {
		// creating context
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		ch := make(chan int, 1)

		// starting the process
		go func() {
			if err := s.sail(j.index, j.size, j.last); err != nil {
				s.failed <- j
			}

			ch <- 0
		}()

		// wait for timeout
		select {
		case <-ch:
			cancel()
			continue
		case <-ctx.Done():
			cancel()
			s.failed <- j
		case <-s.killSwitch:
			cancel()
			break
		}
	}

	return nil
}

// sail
// downloads a chunk of file from given url.
func (s *ship) sail(index, size int, isLast bool) error {
	var (
		// creating http client
		client = &http.Client{}
		// start bit
		start = index * size
		// header value
		dataRange = fmt.Sprintf("bytes=%d-%d", start, start+size-1)
	)

	// check if the worker is working on last chunk
	if isLast {
		dataRange = fmt.Sprintf("bytes=%d-", start)
	}

	// make http request
	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		return fmt.Errorf("make http request failed: %v", err)
	}

	// adding the chunk header
	req.Header.Add("Range", dataRange)

	// executing http request
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %v", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	// read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body failed: %v", err)
	}

	// publish the chunk in channel
	s.channel <- chunk{index: index, data: body}

	return nil
}
