package util

import (
	"encoding/json"
	"net"
)

/*
JSON -> String: Marshal
String -> JSON: Unmarshal
JSON -> Stream: Encode
Stream -> JSON: Decode

stream = continuous streams of data, sequence of data elements that are being processed continuously
*/

func NewCodec(conn net.Conn) ServerCodec {
	dec := json.NewDecoder(conn)
	dec.UseNumber()
	enc := json.NewEncoder(conn)

}
