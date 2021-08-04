package dv

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type message struct {
	Header http.Header
	Body   json.RawMessage
}

type receiverID string

type receiver struct {
	id   receiverID
	conn *websocket.Conn
	mtx  sync.Mutex
}

func newReceiver(conn *websocket.Conn) (recv *receiver) {
	recv = &receiver{
		id:   receiverID(uuid.New().String()),
		conn: conn,
	}
	return recv
}

func (recv *receiver) writeMessage(messageType int, data []byte) error {
	recv.mtx.Lock()
	defer recv.mtx.Unlock()
	return recv.conn.WriteMessage(messageType, data)
}

func (recv *receiver) close() {
	delete(recvs, recv.id)
	log.Printf("Receiver %q closed\n", recv.id)
}

func (m message) String() string {
	return string(m.Body)
}

// receiver regisry
var recvs = map[receiverID]*receiver{}

// newReceiver returns a channel to receive streamed data
// from. Each websocket
func registerReceiver(wconn *websocket.Conn) *receiver {
	recv := newReceiver(wconn)
	recvs[recv.id] = recv
	log.Printf("Receiver %q opened\n", recv.id)
	return recv
}
