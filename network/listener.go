package network

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func CreateTCPListener() (listener *TCPListener) {
	listener = &TCPListener{}
	return
}
func CreateUDPListener() (listener *UDPListener) {
	listener = &UDPListener{}
	return
}
func CreateWSListener() (listener *WSListener) {
	listener = &WSListener{}
	return
}
func CreateChanListener() (listener *ChanListener) {
	listener = &ChanListener{}
	return
}

// ==================== TCP Listener ====================
type TCPListener struct {
	sync.Mutex

	active bool

	address  *net.TCPAddr
	instance *net.TCPListener

	pipe chan ISocket
}

func (this *TCPListener) Listen(address string) error {
	var err error
	this.pipe = make(chan ISocket, 128)
	this.address, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	this.instance, err = net.ListenTCP("tcp", this.address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	this.active = true
	go func() {
		for this.active {
			var conn, err = this.instance.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go func() {
				var tcp = &TCPSocket{}
				tcp.Constructor(conn)
				this.pipe <- tcp
			}()
		}
		err = this.instance.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	return err
}
func (this *TCPListener) Accept() *Connector {
	var c = <-this.pipe
	var connector = &Connector{}
	connector.Constructor(c)
	return connector
}
func (this *TCPListener) StopListen() {
	this.Lock()
	defer this.Unlock()
	this.active = false
}

// ==================== UDP Listener ====================
type UDPListener struct {
	sync.Mutex

	active bool

	address  *net.UDPAddr
	instance *net.UDPConn
	list     map[string]*UDPSocket

	pipe chan ISocket
}

func (this *UDPListener) Listen(address string) error {
	var err error
	this.pipe = make(chan ISocket, 128)
	this.list = make(map[string]*UDPSocket)
	this.address, err = net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	this.instance, err = net.ListenUDP("udp", this.address)
	if err != nil {
		fmt.Println(err)
		return err
	}
	this.active = true
	go func() {
		for this.active {
			var buffer = make([]byte, 4096)
			var lenght, addr, err = this.instance.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println(err)
				continue
			}
			go func() {
				if socket, has := this.list[addr.IP.String()]; has {
					socket.readBuf <- buffer[:lenght]
				} else {
					var socket = &UDPSocket{}
					socket.Constructor(addr, this.instance)
					this.list[addr.IP.String()] = socket
					this.pipe <- socket
				}
			}()
		}
		err = this.instance.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	return err
}
func (this *UDPListener) Accept() *Connector {
	var c = <-this.pipe
	var connector = &Connector{}
	connector.Constructor(c)
	return connector
}
func (this *UDPListener) StopListen() {
	this.Lock()
	defer this.Unlock()
	this.active = false
}

// ==================== WebSocket Listener ====================
type WSListener struct {
	sync.Mutex

	active bool

	instance *websocket.Upgrader

	pipe chan ISocket
}

func (this *WSListener) Listen(address string, sslCert string, sslKey string) error {
	var err error
	this.pipe = make(chan ISocket, 128)
	this.instance = &websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	this.active = true
	go func() {
		http.HandleFunc(address, func(w http.ResponseWriter, r *http.Request) {
			var ws, err = this.instance.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			var socket = &WSocket{}
			socket.Constructor(ws)
			this.pipe <- socket
		})
		if sslCert != "" && sslKey != "" {
			http.ListenAndServeTLS(address, sslCert, sslKey, nil)
		} else {
			http.ListenAndServe(address, nil)
		}
	}()
	return err
}
func (this *WSListener) Accept() *Connector {
	var c = <-this.pipe
	var connector = &Connector{}
	connector.Constructor(c)
	return connector
}
func (this *WSListener) StopListen() {
	this.Lock()
	defer this.Unlock()
	this.active = false
}

// ==================== ChanSocket Listener ====================
type ChanListener struct {
	sync.Mutex

	active   bool
	instance chan string

	pipe chan ISocket
}

func (this *ChanListener) Listen(address string) error {
	var err error
	this.pipe = make(chan ISocket, 128)
	this.instance = make(chan string, 256)
	this.active = true
	go func() {
		for this.active {
			for _, c := range <-this.instance {
				fmt.Println(c)
			}
		}
	}()
	return err
}
func (this *ChanListener) Accept() *Connector {
	var c = <-this.pipe
	var connector = &Connector{}
	connector.Constructor(c)
	return connector
}
func (this *ChanListener) StopListen() {
	this.Lock()
	defer this.Unlock()
	this.active = false
}
