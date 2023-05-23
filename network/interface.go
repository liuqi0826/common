package network

type ISocket interface {
	Read() ([]byte, error)
	Write([]byte)
	Close() error
}
