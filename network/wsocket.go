package network

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

type WSocket struct {
	sync.RWMutex

	instance *websocket.Conn

	active bool

	readBuffer  chan []byte
	writeBuffer chan []byte
}

func (this *WSocket) Constructor(socket *websocket.Conn) {
	this.instance = socket
	if this.instance != nil {
		fmt.Println("Websocket connect from: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))
		this.readBuffer = make(chan []byte, 16)
		this.writeBuffer = make(chan []byte, 16)

		go this.readListen()
		go this.writeListen()
	}
}
func (this *WSocket) Read() (data []byte, err error) {
	if this != nil {
		if this.active {
			if data, ok := <-this.readBuffer; ok {
				return data, err
			} else {
				err = errors.New("Websocket read chan closed")
			}
		} else {
			err = errors.New("Websocket connect closed")
		}
	}
	return
}
func (this *WSocket) Write(data []byte) {
	if this != nil {
		if this.active {
			if len(this.writeBuffer) < cap(this.writeBuffer) {
				this.writeBuffer <- data
			}
		}
	}
}
func (this *WSocket) Close() (err error) {
	if this != nil {
		this.Lock()
		defer this.Unlock()
		if this.active {
			this.active = false
			close(this.readBuffer)
			close(this.writeBuffer)

			fmt.Println("Websocket close: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))
			err = this.instance.Close()
			if err != nil {
				return
			}
		} else {
			err = errors.New("Websocket has been closed.")
		}
	}
	return
}
func (this *WSocket) readListen() {
	if this != nil {
		for this.active {
			_, data, err := this.instance.ReadMessage()
			if err == nil {
				this.readBuffer <- data
			} else {
				err = this.Close()
			}
		}
	}
}
func (this *WSocket) writeListen() {
	if this != nil {
		for this.active {
			if message, ok := <-this.writeBuffer; ok {
				go func() {
					var err = this.instance.WriteMessage(websocket.BinaryMessage, message)
					if err != nil {
						this.Close()
					}
				}()
			} else {
				this.Close()
			}
		}
	}
}
