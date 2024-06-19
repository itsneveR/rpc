package client

type ClientCodec interface {
	WriteReq(*Request, any) (int, error)
	ReadRes(*Response) error

	cloes() error
}
