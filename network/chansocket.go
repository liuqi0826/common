package network

import (
	"errors"
	"fmt"
	"sync"
)

type ChanSocket struct {
	sync.RWMutex

	buffer chan []byte

	other *ChanSocket

	active bool
}

func (this *ChanSocket) Constructor() {
	this.buffer = make(chan []byte, 16)
	this.active = true
}
func (this *ChanSocket) Bind(other *ChanSocket) (err error) {
	this.other = other
	fmt.Println("Chan connect from: " + fmt.Sprintf("%p", this.other))
	return
}
func (this *ChanSocket) Read() ([]byte, error) {
	var err error
	if this.active {
		if msg, ok := <-this.buffer; ok {
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
		if this.other != nil && this.other.active {
			this.other.buffer <- data
		}
	}
}
func (this *ChanSocket) Close() (err error) {
	this.Lock()
	defer this.Unlock()

	if this.active {
		this.active = false
		close(this.buffer)

		fmt.Println("Chan close: " + fmt.Sprintf("%p", this.buffer))
	} else {
		err = errors.New("Chan Closed!")
	}
	return
}
