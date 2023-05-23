package network

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type WSocket struct {
	sync.RWMutex

	instance *websocket.Conn

	active   bool
	readBuf  chan []byte
	writeBuf chan []byte
}

func (this *WSocket) Constructor(socket *websocket.Conn) {
	this.instance = socket
	this.readBuf = make(chan []byte, 32)
	this.writeBuf = make(chan []byte, 32)
	fmt.Println("Websocket connect from: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))

	go this.readListen()
	go this.writeListen()
}
func (this *WSocket) Read() ([]byte, error) {
	var err error
	if this.active {
		if msg, ok := <-this.readBuf; ok {
			return msg, err
		} else {
			err = this.Close()
		}
	} else {
		err = errors.New("Websocket connection closed")
	}
	return nil, err
}
func (this *WSocket) Write(data []byte) {
	if this.active {
		if len(this.writeBuf) < cap(this.writeBuf) {
			this.writeBuf <- data
		}
	}
}
func (this *WSocket) Close() error {
	var err error
	if this.active {
		this.Lock()
		defer this.Unlock()
		this.active = false
		close(this.readBuf)
		close(this.writeBuf)
		fmt.Println("Websocket close: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))
		err = this.instance.Close()
	} else {
		err = errors.New("Websocket has been closed.")
	}
	return err
}
func (this *WSocket) readListen() {
	for this.active {
		_, msg, err := this.instance.ReadMessage()
		if err == nil {
			this.readBuf <- msg
		} else {
			err = this.Close()
		}
	}
}
func (this *WSocket) writeListen() {
	for this.active {
		if message, ok := <-this.writeBuf; ok {
			err := this.instance.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				this.Close()
			}
		} else {
			this.Close()
		}
	}
}
