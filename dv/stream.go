package dv

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	Header http.Header
	Body   json.RawMessage
}

type receiverID string

type Receiver struct {
	id   receiverID
	conn *websocket.Conn
}

func newReceiver(conn *websocket.Conn) (recv *Receiver) {
	recv = &Receiver{
		id:   receiverID(uuid.New().String()),
		conn: conn,
	}
	return recv
}

func (recv *Receiver) Close() {
	delete(recvs, recv.id)
	log.Printf("Receiver %q closed\n", recv.id)
}

func (m Message) String() string {
	return string(m.Body)
}

// receiver regisry
var recvs = map[receiverID]*Receiver{}

// newReceiver returns a channel to receive streamed data
// from. Each websocket
func registerReceiver(wconn *websocket.Conn) *Receiver {
	recv := newReceiver(wconn)
	recvs[recv.id] = recv
	log.Printf("Receiver %q opened\n", recv.id)
	return recv
}
