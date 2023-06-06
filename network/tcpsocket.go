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
	sync.RWMutex

	instance *net.TCPConn

	active     bool
	headRead   bool
	packageLen uint32
	buffer     []byte

	readBuffer  chan []byte
	writeBuffer chan []byte
}

func (this *TCPSocket) Constructor(socket *net.TCPConn) {
	this.instance = socket
	if this.instance != nil {
		fmt.Println("TCP connect from: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))
		this.readBuffer = make(chan []byte, 16)
		this.writeBuffer = make(chan []byte, 16)

		this.active = true

		go this.readListen()
		go this.writeListen()
	}
}
func (this *TCPSocket) Read() (data []byte, err error) {
	if this != nil {
		if this.active {
			if data, ok := <-this.readBuffer; ok {
				return data, nil
			} else {
				err = errors.New("TCP read chan closed")
			}
		} else {
			err = errors.New("TCP connect closed")
		}
	}
	return
}
func (this *TCPSocket) Write(data []byte) {
	if this != nil {
		if this.active {
			if len(this.writeBuffer) < cap(this.writeBuffer) {
				var length = uint32(len(data))
				var stream = bytes.NewBuffer([]byte{})
				binary.Write(stream, binary.BigEndian, length)
				binary.Write(stream, binary.BigEndian, data)
				this.writeBuffer <- stream.Bytes()
			}
		}
	}
}
func (this *TCPSocket) Close() (err error) {
	if this != nil {
		this.Lock()
		defer this.Unlock()
		if this.active {
			this.active = false
			close(this.readBuffer)
			close(this.writeBuffer)

			fmt.Println("TCP close: " + fmt.Sprintf("%s", this.instance.RemoteAddr()))
			err = this.instance.Close()
			if err != nil {
				return
			}
		} else {
			err = errors.New("TCP has been closed!")
		}
	}
	return
}

func (this *TCPSocket) readListen() {
	if this != nil {
		for this.active {
			var buffer = make([]byte, 2048)
			var length, err = this.instance.Read(buffer)
			if err == nil {
				this.buffer = append(this.buffer, buffer[:length]...)
				for {
					if this.headRead {
						if len(this.buffer) >= int(this.packageLen) {
							this.headRead = false
							var data = bytes.NewBuffer(this.buffer[:this.packageLen])
							this.buffer = this.buffer[this.packageLen:len(this.buffer)]
							this.readBuffer <- data.Bytes()
						} else {
							break
						}
					} else {
						if len(this.buffer) >= 4 {
							var lenBuffer = bytes.NewBuffer(this.buffer[0:4])
							var e = binary.Read(lenBuffer, binary.BigEndian, &this.packageLen)
							if e != nil {
								fmt.Println(e)
							}
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
}
func (this *TCPSocket) writeListen() {
	if this != nil {
		for this.active {
			if msg, ok := <-this.writeBuffer; ok {
				go func() {
					var _, err = this.instance.Write(msg)
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
