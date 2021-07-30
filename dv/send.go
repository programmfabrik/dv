package dv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

func Send(c *cli.Context) error {
	// read data from stdin
	mimeType, recycleReader, err := recycleReader(os.Stdin)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Sending %q...", mimeType)
	resp, err := http.Post(c.String("url")+"/data", mimeType, recycleReader)
	if err != nil {
		log.Print("post:", err.Error())
	} else if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		log.Print("post:", resp.Status, string(msg))
	}
	fmt.Fprintf(os.Stderr, "Data sent\n")
	return nil
}

// recycleReader returns the MIME type of input and a new reader
// containing the whole data from input.
func recycleReader(input io.Reader) (mimeType string, recycled io.Reader, err error) {
	// header will store the bytes mimetype uses for detection.
	header := bytes.NewBuffer(nil)

	// After DetectReader, the data read from input is copied into header.
	mtype, err := mimetype.DetectReader(io.TeeReader(input, header))

	// Concatenate back the header to the rest of the file.
	// recycled now contains the complete, original data.
	recycled = io.MultiReader(header, input)

	return mtype.String(), recycled, err
}

func data(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dead := []*Receiver{}

	wg := sync.WaitGroup{}
	for _, recv := range recvs {
		wg.Add(1)
		go func(recv *Receiver) {
			defer wg.Done()

			// header is sent as json
			bs, err := json.Marshal(r.Header)
			if err != nil {
				// recv is dead
				log.Print("send marshal header:", err.Error())
				return
			}

			err = recv.conn.WriteMessage(websocket.TextMessage, bs)
			if err != nil {
				// recv is dead
				log.Print("send header:", err.Error())
				dead = append(dead, recv)
				return
			}

			err = recv.conn.WriteMessage(websocket.BinaryMessage, body)
			if err != nil {
				// recv is dead
				log.Print("send body:", err.Error())
				dead = append(dead, recv)
				return
			}
			log.Print("sent to:", recv.id, "bytes:", len(bs))
		}(recv)
	}

	for _, recv := range dead {
		recv.Close()
	}
}
