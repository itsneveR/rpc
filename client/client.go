package client

import (
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	TCP = iota
	WS
	HTTP
)

type handler struct {
}
type clientConn struct {
	proto   int
	handler *handler

	//tcp settings
	tcpDialer net.Dialer

	// WebSocket settings
	wsDialer           *websocket.Dialer
	wsMessageSizeLimit *int64 // wsMessageSizeLimit nil = default, 0 = no limit

	// HTTP settings
	httpClient  *http.Client
	httpHeaders http.Header
	httpAuth    HTTPAuth
}

type Client struct {
	Conn  clientConn
	idgen func() ID // for subscriptions

	codec ClientCodec

	reqMutex      sync.Mutex // protects following
	request       Request
	services      *serviceRegistry
	mutex         sync.Mutex // protects following
	seq           uint64
	pending       map[uint64]*Call
	reconnectFunc reconnectFunc

	close    chan struct{} // user has called Close
	shutdown chan struct{} // server has told us to stop
	closing  chan struct{} // closed when client is quitting
	didClose chan struct{}

	// This function, if non-nil, is called when the connection is lost.
	reconnectFunc reconnectFunc
}

func Dial() {

}

func initClient(codec ClientCodec)
func (c *Client) Start() {

}

func (c *Client) WriteReq(r *Request, data any) (int, error) {
	switch c.Conn.proto {
	case 0:
	case 1:
	case 2:

	}
}
