package maersk

import (
	"fmt"
	"io"
	"net/http"
)

// Manages to download each chunk from given url.
type worker struct {
	// host url for downloading file.
	url string
	// channel for sending the downloaded chunk.
	channel chan chunk
	// input channel for getting the jobs.
	jobs chan job
}

// work function starts the worker to listen
// on jobs channel.
func (w *worker) work() error {
	for j := range w.jobs {
		if err := w.process(j.index, j.size, j.last); err != nil {
			return err
		}
	}

	return nil
}

// process
// downloads a chunk of file from given url.
func (w *worker) process(index, size int, isLast bool) error {
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
	req, err := http.NewRequest("GET", w.url, nil)
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
	w.channel <- chunk{index: index, data: body}

	return nil
}
