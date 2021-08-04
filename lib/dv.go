package dv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

type Logger struct {
	opts LoggerOpts
}

type LoggerOpts struct {
	URL string
}

var DefaultOpts = LoggerOpts{
	URL: "http://localhost:10000",
}

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(DefaultOpts)
}

func NewLogger(opts LoggerOpts) (s *Logger) {
	s = &Logger{opts: opts}
	s.Log("New logger initialized")
	return s
}

func Log(data interface{}) (err error) {
	return defaultLogger.Log(data)
}

func (s *Logger) Log(data interface{}) (err error) {

	rv := reflect.ValueOf(data)

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

	fmt.Fprintf(os.Stderr, "dv: sending %q %d bytes to %s...", ct, dr.Len(), s.opts.URL)
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
	fmt.Fprint(os.Stderr, " sent\n")
	return nil
}
