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

func SendReader(c *cli.Context, readFrom io.Reader) error {
	// read data from stdin
	mimeType, recycleReader, err := recycleReader(readFrom)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Sending %q...", mimeType)
	resp, err := http.Post(c.String("url")+"/data", mimeType, recycleReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "post: %s\n", err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "post: %s %s\n", resp.Status, string(msg))
		return nil
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
	println("reading data...")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	println("read", len(body), len(recvs), "receivers")

	dead := []*receiver{}

	wg := sync.WaitGroup{}
	for _, recv := range recvs {
		wg.Add(1)
		go func(recv *receiver) {
			defer wg.Done()

			// header is sent as json
			bs, err := json.Marshal(r.Header)
			if err != nil {
				// recv is dead
				log.Print("send marshal header:", err.Error())
				return
			}

			err = recv.writeMessage(websocket.TextMessage, bs)
			if err != nil {
				// recv is dead
				log.Print("send header:", err.Error())
				dead = append(dead, recv)
				return
			}

			err = recv.writeMessage(websocket.BinaryMessage, body)
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
		recv.close()
	}
}
