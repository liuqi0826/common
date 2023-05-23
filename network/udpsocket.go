package network

import (
	"net"
	"sync"
)

type UDPSocket struct {
	sync.Mutex

	addr   *net.UDPAddr
	socket *net.UDPConn

	active     bool
	headRead   bool
	packageLen uint32
	buffer     []byte
	readBuf    chan []byte
	writeBuf   chan []byte
}

func (this *UDPSocket) Constructor(addr *net.UDPAddr, socket *net.UDPConn) {
	this.addr = addr
	this.socket = socket
	if this.addr != nil && this.socket != nil {
		this.readBuf = make(chan []byte, 32)
		this.writeBuf = make(chan []byte, 32)

		go this.readListen()
		go this.writeListen()

		this.active = true
	}
}
func (this *UDPSocket) Read() ([]byte, error) {
	return nil, nil
}
func (this *UDPSocket) Write(data []byte) {
}
func (this *UDPSocket) Close() error {
	return nil
}
func (this *UDPSocket) readListen() {
}
func (this *UDPSocket) writeListen() {
}
