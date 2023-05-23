package network

import (
	"errors"
	"sync"
	"time"

	"common/events"
)

type Connector struct {
	sync.Mutex
	events.EventDispatcher

	active     bool
	createTime int64
	deadline   int64
	markTime   int64

	pipe chan []byte

	socket ISocket
}

func (this *Connector) Constructor(socket ISocket) {
	this.EventDispatcher.Constructor(this)
	this.socket = socket
	if this.socket != nil {
		this.createTime = time.Now().Unix()
		this.pipe = make(chan []byte, 128)
		this.active = true

		go this.listen()
	}
}
func (this *Connector) Send(data []byte) {
	if this != nil {
		if this.active {
			go this.socket.Write(data)
		}
	}
}
func (this *Connector) Read() (data []byte) {
	if this != nil {
		buff, success := <-this.pipe
		if success {
			return buff
		}
	}
	return
}
func (this *Connector) Close() error {
	var err error
	if this != nil {
		this.Lock()
		defer this.Unlock()

		err = this.socket.Close()
		this.socket = nil

		this.active = false

		var evt = &events.Event{
			Type: events.CLOSE,
		}
		this.DispatchEvent(evt)
	} else {
		err = errors.New("Network is closed!")
	}
	return err
}
func (this *Connector) SetHeartBeat(interval time.Duration) error {
	var err error
	if this != nil {
		go func() {
			for {
				if this != nil && this.active && this.socket != nil {
					var p = []byte{0, 0}
					this.Send(p)
					time.Sleep(interval)
				} else {
					break
				}
			}
		}()
	} else {
		err = errors.New("Network is closed!")
	}
	return err
}
func (this *Connector) listen() {
	if this != nil {
		for this.active {
			var buffer, err = this.socket.Read()
			if err != nil {

			} else {
				if buffer != nil {
					var evt = &events.Event{
						Type: events.IO_ERROR,
						Data: "",
					}
					this.DispatchEvent(evt)
					break
				} else {
					this.pipe <- buffer
				}
			}
		}
	}
}
func (this *Connector) pong() {
	if this.markTime == 0 {
		go func() {
			for {
				if this != nil && this.active && this.socket != nil {
					if time.Now().Unix()-this.markTime > this.deadline {
						this.Close()
					}
					time.Sleep(time.Second * 10)
				} else {
					break
				}
			}
		}()
	}
	this.markTime = time.Now().Unix()
}
