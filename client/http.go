package client

import "net/http"

type HTTPAuth func(h http.Header) error
