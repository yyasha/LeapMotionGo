package leap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"golang.org/x/net/websocket"
)

const (
	origin  = "http://localhost"
	port    = "6437"
	version = "/v6.json"
	debug   = false
)

// A Conn holds a connection to a leap motion device.
type Conn struct {
	conn *websocket.Conn
	dec  *json.Decoder
}

// Connect connects to a leap device running on the given host,
// and checks the version running is correct.
func Connect(host string) (*Conn, error) {
	u := url.URL{Scheme: "ws", Host: host + ":" + port, Path: version}
	ws, err := websocket.Dial(u.String(), "", origin)
	if err != nil {
		return nil, fmt.Errorf("open connection to leap: %v", err)
	}

	var r io.Reader = ws
	// Add logging to stdout if debug is true.
	if debug {
		pr, pw := io.Pipe()
		go func() {
			io.Copy(io.MultiWriter(os.Stdout, pw), ws)
		}()
		r = pr
	}

	conn := &Conn{ws, json.NewDecoder(r)}

	msgs := []interface{}{
		struct {
			Gestures bool `json:"enableGestures"`
		}{true},
		struct {
			Background bool `json:"focused"`
		}{true},
	}

	enc := json.NewEncoder(conn.conn)
	for i, msg := range msgs {
		err = enc.Encode(msg)
		if err != nil {
			ws.Close()
			return nil, fmt.Errorf("send config msg %d: %v", i, err)
		}
	}

	return conn, nil
}

// Frame returns the next Frame sent by the device.
func (c *Conn) Frame() (*Frame, error) {
	f := &Frame{}
	return f, c.dec.Decode(f)
}

// Decode uses the given value to decode the Frame sent from the device.
// This can be useful to minimize the amount of data decoded.
func (c *Conn) Decode(v interface{}) error {
	return c.dec.Decode(v)
}

// Close closes the connection to the device.
func (c *Conn) Close() {
	c.conn.Close()
}
