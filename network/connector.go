package network

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liuqi0826/common/events"
)

func ConnectTCP(address string, token string) (connector *Connector, err error) {
	var raddr *net.TCPAddr
	raddr, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return
	}
	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return
	}

	if token != "" {
		var code = []byte(token)
		_, err = conn.Write(code)
		if err != nil {
			return
		}
	}

	var tcpsocket = &TCPSocket{}
	tcpsocket.Constructor(conn)
	connector = &Connector{}
	connector.Constructor(tcpsocket)

	return
}
func ConnectUDP(address string, token string) (connector *Connector, err error) {
	var remote *net.UDPAddr
	remote, err = net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}
	var conn *net.UDPConn
	conn, err = net.DialUDP("udp", nil, remote)
	if err != nil {
		return
	}

	if token != "" {
		var code = []byte(token)
		_, err = conn.Write(code)
		if err != nil {
			return
		}
	} else {
		var code = []byte(" ")
		_, err = conn.Write(code)
		if err != nil {
			return
		}
	}

	var buffer = make([]byte, 1024)
	var lenght int
	lenght, err = conn.Read(buffer)
	if err != nil {
		return
	}
	conn.Close()

	var targetAddress = string(buffer[:lenght])
	var addr *net.UDPAddr
	addr, err = net.ResolveUDPAddr("udp", targetAddress)
	if err != nil {
		return
	}
	var socket = &UDPSocket{}
	socket.Constructor(nil, nil, addr)
	connector = &Connector{}
	connector.Constructor(socket)

	return
}
func ConnectWS(address string, token string) (connector *Connector, err error) {
	var conn *websocket.Conn
	if token != "" {
		address = address + "&" + token
	}
	conn, _, err = websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return
	}

	var wsocket = &WSocket{}
	wsocket.Constructor(conn)
	connector = &Connector{}
	connector.Constructor(wsocket)

	return
}
func ConnectChan(address string, token string) (connector *Connector, err error) {
	if listener, has := chanListenerList[address]; has {
		var cs, ss *ChanSocket
		cs = &ChanSocket{}
		cs.Constructor()
		ss = &ChanSocket{}
		ss.Constructor()

		cs.Bind(ss)
		ss.Bind(cs)

		listener.pipe <- ss

		if token != "" {
			var code = []byte(token)
			cs.Write(code)
		}

		connector = &Connector{}
		connector.Constructor(cs)
	} else {
		err = errors.New("Chan listener is not found.")
	}

	return
}

type Connector struct {
	sync.RWMutex
	events.EventDispatcher

	active     bool
	createTime int64
	deadline   int64
	markTime   int64

	buffer chan []byte

	socket ISocket
}

func (this *Connector) Constructor(socket ISocket) {
	this.EventDispatcher.Constructor(this)

	this.socket = socket
	if this.socket != nil {
		this.createTime = time.Now().Unix()
		this.buffer = make(chan []byte)
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
		if this.active {
			return <-this.buffer
		}
	}
	return
}
func (this *Connector) Close() (err error) {
	if this != nil {
		this.Lock()
		defer this.Unlock()

		err = this.socket.Close()
		this.socket = nil

		this.active = false

		this.DispatchEvent(&events.Event{
			Type: events.CLOSE,
		})
	} else {
		err = errors.New("Connector is closed!")
	}
	return
}
func (this *Connector) SetHeartBeat(interval time.Duration) (err error) {
	if this != nil {
		go func() {
			for {
				if this.active && this.socket != nil {
					var p = []byte{0, 0}
					this.Send(p)
					time.Sleep(interval)
				} else {
					break
				}
			}
		}()
	} else {
		err = errors.New("Connector is closed!")
	}
	return
}
func (this *Connector) listen() {
	if this != nil {
		for this.active {
			var buffer, err = this.socket.Read()
			if err != nil {
				this.DispatchEvent(&events.Event{
					Type: events.IO_ERROR,
					Data: err.Error(),
				})
			} else {
				this.buffer <- buffer
			}
		}
	}
}
func (this *Connector) pong() {
	if this != nil {
		if this.markTime == 0 {
			go func() {
				for {
					if this.active && this.socket != nil {
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
}
