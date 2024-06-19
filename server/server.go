package server

import (
	"net"

	"log"

	"github.com/itsneveR/rpc/util"
)

type Server struct {
	services  serviceRegistry
	codecType util.ServerCodec
}

func NewServer(proto string, address string) *Server {

	connEstablish()

}

func connEstablish(proto string, address string) *ServerCodec {
	switch proto {
	case "tcp":
		tcpServer(address)
	case "ws":
	case "http":
	case "":
	default:
	}
}

func tcpServer(address string) {
	conn, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	go loop()
}

func (s *Server) serverListener(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if util.IsTemporaryError(err) {
			log.Println("TCP/RPC accept error", err)
			continue
		} else if err != nil {
			return err
		}

		log.Println("Accepted RPC connection from", conn.RemoteAddr())

		go s.handler(conn, nil)
	}

}

func (s *Server) hanlder(conn net.Conn, codec util.ServerCodec) any {

	if codec == nil {
		s.codecType = util.Codec.New(conn)
	}

}

// *default codec implementation*//
type jsonCodec struct {
}
