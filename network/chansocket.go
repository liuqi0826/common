package network

import (
	"errors"
	"sync"
)

type ChanSocket struct {
	sync.Mutex

	active   bool
	readBuf  chan []byte
	writeBuf chan []byte
}

func (this *ChanSocket) Constructor() {
	this.readBuf = make(chan []byte, 32)
	this.writeBuf = make(chan []byte, 32)
	this.active = true
}
func (this *ChanSocket) Read() ([]byte, error) {
	var err error
	if this.active {
		if msg, ok := <-this.readBuf; ok {
			return msg, nil
		} else {
			err = errors.New("Read chan closed")
		}
	} else {
		err = errors.New("Connect closed")
	}
	return nil, err
}
func (this *ChanSocket) Write(data []byte) {
	if this.active {
		this.writeBuf <- data
	}
}
func (this *ChanSocket) Close() error {
	var err error
	this.Lock()
	defer this.Unlock()
	if this.active {
		this.active = false
		close(this.readBuf)
		close(this.writeBuf)
	} else {
		err = errors.New("Closed!")
	}
	return err
}
