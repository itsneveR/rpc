package server

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"sync"
)

const (
	defaultPacketSize = 1024
	headerLen         = 4
)

type server struct {
	services services
	conn     chan net.Conn
	err      chan error
	close    chan struct{}
	h        *handler
}

type handler struct {
	err  chan error
	data chan []byte
}

func newService() *services {
	return &services{
		mu:  sync.Mutex{},
		svc: make(map[string]*object),
	}
}

func newServer(network string, addr string) *server {
	s := &server{
		services: *newService(),
		conn:     make(chan net.Conn),
		err:      make(chan error),
		close:    make(chan struct{}),
		h:        new(handler),
	}

	go s.SetNetwork(network, addr)

	return s
}

func (s *server) SetNetwork(network string, addr string) {

	ln, err := net.Listen(network, addr)

	if err != nil {
		s.err <- err
	}

	go func() {
		for {
			conn, err := ln.Accept()
			s.err <- err
			s.conn <- conn
		}
	}()

}

func (s *server) Register(name string, object any) error {
	return s.services.RegisterService(name, object)
}

func (s *server) handler() {

	select {
	case conn := <-s.conn:
		buf := make([]byte, headerLen)
		byteBuf := bytes.NewBuffer(buf)
		headData := byteBuf.Bytes()
		//defer byteBuf.Reset()

		_, err := io.ReadFull(conn, headData)

		dataLen := binary.BigEndian.Uint32(headData)

		data := make([]byte, dataLen)
		dataBuf := bytes.NewBuffer(data)
		dataByte := dataBuf.Bytes()
		//defer dataBuf.Reset()
		_, err = io.ReadFull(conn, dataByte[headerLen:])

		s.h.err <- err

		s.h.data <- data
	}

}
