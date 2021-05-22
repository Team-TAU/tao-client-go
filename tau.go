package go_tau

import "sync"
import "github.com/gorilla/websocket"

// Client represents a client connected to TAU
type Client struct {
	conn      *websocket.Conn
	hostname  string
	port      string
	writeLock *sync.Mutex

	// callback functions

}
