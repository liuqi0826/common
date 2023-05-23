package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
)

type TCPSocket struct {
	sync.Mutex

	instance *net.TCPConn

	active     bool
	headRead   bool
	packageLen uint32
	buffer     []byte
	readBuf    chan []byte
	writeBuf   chan []byte
}

func (this *TCPSocket) Constructor(socket *net.TCPConn) {
	this.instance = socket
	if this.instance != nil {
		this.readBuf = make(chan []byte, 32)
		this.writeBuf = make(chan []byte, 32)
		fmt.Println("TCP instance from: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))

		go this.readListen()
		go this.writeListen()

		this.active = true
	}
}
func (this *TCPSocket) Read() ([]byte, error) {
	var err error
	if this.active {
		if msg, ok := <-this.readBuf; ok {
			return msg, nil
		} else {
			err = errors.New("TCP read chan closed")
		}
	} else {
		err = errors.New("TCP connect closed")
	}
	return nil, err
}
func (this *TCPSocket) Write(data []byte) {
	if this.active {
		if len(this.writeBuf) < cap(this.writeBuf) {
			var length = uint32(len(data))
			var stream = bytes.NewBuffer([]byte{})
			binary.Write(stream, binary.BigEndian, length)
			binary.Write(stream, binary.BigEndian, data)
			this.writeBuf <- stream.Bytes()
		}
	}
}
func (this *TCPSocket) Close() error {
	var err error
	this.Lock()
	defer this.Unlock()
	if this.active {
		this.active = false
		close(this.readBuf)
		close(this.writeBuf)
		fmt.Println("TCP close: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))
		err = this.instance.Close()
	} else {
		err = errors.New("TCP closed!")
	}
	return err
}

func (this *TCPSocket) readListen() {
	for this != nil && this.active {
		var buffer = make([]byte, 4096)
		length, err := this.instance.Read(buffer)
		if err == nil {
			this.buffer = append(this.buffer, buffer[:length]...)
			for {
				if this.headRead {
					if len(this.buffer) >= int(this.packageLen) {
						this.headRead = false
						var data = bytes.NewBuffer(this.buffer[:this.packageLen])
						this.buffer = this.buffer[this.packageLen:len(this.buffer)]
						this.readBuf <- data.Bytes()
					} else {
						break
					}
				} else {
					if len(this.buffer) >= 4 {
						lenBuffer := bytes.NewBuffer(this.buffer[0:4])
						err = binary.Read(lenBuffer, binary.BigEndian, &this.packageLen)
						fmt.Println(err)
						this.headRead = true
						this.buffer = this.buffer[4:len(this.buffer)]
					} else {
						break
					}
				}
			}
		} else {
			this.Close()
		}
	}
}
func (this *TCPSocket) writeListen() {
	for this != nil && this.active {
		if msg, ok := <-this.writeBuf; ok {
			go func() {
				var _, err = this.instance.Write(msg)
				if err != nil {

				}
			}()
		} else {
			this.Close()
		}
	}
}
