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
	size := 0
	channel := make(chan chunk, c.Workers)
	c.chunks = make([][]byte, c.Workers)

	client := &http.Client{}

	req, err := http.NewRequest("HEAD", c.URL, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	//defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)

	log.Println("Headers : ", resp.Header["Content-Length"])

	if header, ok := resp.Header["Content-Length"]; ok {
		fileSize, err := strconv.Atoi(header[0])

		if err != nil {
			log.Fatal("File size could not be determined : ", err)
		}

		size = fileSize / c.Workers

	} else {
		log.Fatal("File size was not provided!")
	}

	for i := 0; i < 5; i++ {
		w := worker{
			channel: channel,
			url:     c.URL,
		}
		go w.work(i, size, i == c.Workers-1)
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
		log.Fatal(err)
	}

	return nil
}
