package network

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type UDPSocket struct {
	sync.RWMutex

	localAddress  *net.UDPAddr
	remoteAddress *net.UDPAddr
	instance      *net.UDPConn

	active bool

	readBuffer  chan []byte
	writeBuffer chan []byte
}

func (this *UDPSocket) Constructor(socket *net.UDPConn, laddr *net.UDPAddr, raddr *net.UDPAddr) (err error) {
	this.instance = socket
	this.localAddress = laddr
	this.remoteAddress = raddr

	if this.instance == nil && this.localAddress == nil && this.remoteAddress == nil {
		err = errors.New("UDP socket values are nil.")
		return
	}
	if this.instance == nil {
		this.instance, err = net.ListenUDP("udp", this.localAddress)
		if err != nil {
			fmt.Println(err)
			return
		}
		if this.localAddress == nil {
			this.localAddress, err = net.ResolveUDPAddr("udp", this.instance.LocalAddr().String())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	if this.localAddress != nil && this.remoteAddress != nil {
		fmt.Println("UDP local bind: " + this.localAddress.String() + ", remote to:" + this.remoteAddress.String())
	} else if this.localAddress != nil && this.remoteAddress == nil {
		fmt.Println("UDP listen at: " + this.localAddress.String())
	} else if this.localAddress == nil && this.remoteAddress != nil {
		fmt.Println("UDP connect to: " + this.remoteAddress.String())
	}

	this.readBuffer = make(chan []byte, 16)
	this.writeBuffer = make(chan []byte, 16)

	this.active = true

	go this.readListen()
	go this.writeListen()

	return
}
func (this *UDPSocket) Read() (data []byte, err error) {
	if this != nil {
		if this.active {
			if data, ok := <-this.readBuffer; ok {
				if len(data) > 0 {
					return data, nil
				}
			} else {
				err = errors.New("UDP read chan closed")
			}
		} else {
			err = errors.New("UDP connect closed")
		}
	}
	return
}
func (this *UDPSocket) Write(data []byte) {
	if this != nil {
		if this.active {
			if len(this.writeBuffer) < cap(this.writeBuffer) {
				this.writeBuffer <- data
			}
		}
	}
}
func (this *UDPSocket) Close() (err error) {
	if this != nil {
		this.Lock()
		defer this.Unlock()
		if this.active {
			this.active = false
			close(this.readBuffer)
			close(this.writeBuffer)

			fmt.Println("UDP close: " + this.localAddress.String())
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
func (this *UDPSocket) readListen() {
	if this != nil {
		for this.active {
			var buffer = make([]byte, 2048)
			var length, remote, err = this.instance.ReadFromUDP(buffer)
			if err == nil {
				if this.remoteAddress != nil {
					if this.remoteAddress.String() != remote.String() {
						continue
					}
				} else {
					this.remoteAddress = remote
				}
				this.readBuffer <- buffer[:length]
			}
		}
	}
}
func (this *UDPSocket) writeListen() {
	if this != nil {
		for this.active {
			if msg, ok := <-this.writeBuffer; ok {
				go func() {
					if this.remoteAddress != nil {
						var _, err = this.instance.WriteToUDP(msg, this.remoteAddress)
						if err != nil {
							fmt.Println(err)
							this.Close()
						}
					}
				}()
			} else {
				this.Close()
			}
		}
	}
}
