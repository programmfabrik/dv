package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

type Sender struct {
	opts SenderOpts
}

type SenderOpts struct {
	URL string
}

var DefaultOpts = SenderOpts{
	URL: "http://localhost:10000",
}

func NewSender(opts SenderOpts) (s *Sender) {
	s = &Sender{opts: opts}
	s.SendData("New sender initialized")
	return s
}

func (s *Sender) SendData(data interface{}) (err error) {

	rv := reflect.ValueOf(data)
	println("rt", rv.String(), s.opts.URL)

	var (
		dr *bytes.Buffer
		ct string
	)

	switch rv.Kind() {
	case reflect.String:
		dr = bytes.NewBuffer([]byte(rv.String()))
		ct = "text/plain"
	default:
		bs, err := json.Marshal(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal data: %q", err.Error())
			return err
		}
		ct = "application/json"
		dr = bytes.NewBuffer(bs)
	}

	fmt.Fprintf(os.Stderr, "Sending %q... %d bytes to %s", ct, dr.Len(), s.opts.URL)
	resp, err := http.Post(s.opts.URL+"/data", ct, dr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "post: %s\n", err.Error())
		return nil
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
