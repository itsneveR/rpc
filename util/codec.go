package util

import "net"

type Codec interface {
	New(conn net.Conn) ServerCodec
	Encode(v any) error
	Decode(v any) error
}

type ServerCodec interface {
}
