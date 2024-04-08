package debugger

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Streamer struct {
	svr *http.Server
	ch  chan []byte
}

// Write implements io.Writer.
func (s *Streamer) Write(p []byte) (n int, err error) {
	// try to write or drop
	select {
	case s.ch <- p:
	default:
	}
	return len(p), nil
}

func NewStreamer() *Streamer {
	return &Streamer{
		ch: make(chan []byte, 3),
	}
}

func (s *Streamer) Start(port string) error {
	http.HandleFunc("/logs", s.handleConnection)

	s.svr = &http.Server{Addr: ":" + port}
	return s.svr.ListenAndServe()
}

var upgrader = websocket.Upgrader{} // use default options

func (s *Streamer) handleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for msg := range s.ch {
		err := c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Print("write:", err)
			return
		}
	}
}
