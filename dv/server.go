package dv

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

var upgrader = websocket.Upgrader{} // use default options

// echo is a simple socket used in gorilla's example
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	println("opening socket", r.RemoteAddr)
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	println("closing socket", r.RemoteAddr)
}

func stream(w http.ResponseWriter, r *http.Request) {
	wconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer wconn.Close()
	log.Printf("opening socket: %q\n", r.RemoteAddr)
	recv := registerReceiver(wconn)
	for {
		mt, message, err := wconn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %d %s", mt, message)
	}
	log.Printf("closing socket %q\n", r.RemoteAddr)
	recv.Close()
}

func setFileHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("cache-control", "no-cache, must-revalidate")
		h.ServeHTTP(rw, req)
	})
}

func Server(c *cli.Context) error {
	r := mux.NewRouter()

	r.HandleFunc("/echo", echo)
	r.HandleFunc("/stream", stream)
	r.HandleFunc("/data", data)

	// This will serve files under http://localhost:8000/static/<filename>
	fr := r.PathPrefix("/").Subrouter()
	fr.Use(setFileHeader)
	fr.PathPrefix("").Handler(http.FileServer(http.Dir("web")))

	addr := c.String("addr")

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Started server on " + addr)

	log.Fatal(srv.ListenAndServe())
	return nil
}
