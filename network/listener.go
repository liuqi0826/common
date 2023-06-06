package network

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
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
	sync.RWMutex

	address  *net.TCPAddr
	instance *net.TCPListener

	verify func(*net.TCPConn) bool

	list chan *net.TCPConn
	pipe chan ISocket

	active bool
}

func (this *TCPListener) Listen(address string) (err error) {
	this.Lock()
	defer this.Unlock()

	this.list = make(chan *net.TCPConn, 1024)
	this.pipe = make(chan ISocket, 1024)
	this.address, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	this.instance, err = net.ListenTCP("tcp", this.address)
	if err != nil {
		fmt.Println(err)
		return
	}
	this.active = true

	go this.listen()
	go this.handle()

	return
}
func (this *TCPListener) Accept() (connector *Connector) {
	var conn, ok = <-this.pipe
	if ok {
		connector = &Connector{}
		connector.Constructor(conn)
	}
	return
}
func (this *TCPListener) StopListen() {
	this.Lock()
	defer this.Unlock()

	var err error
	if this.active {
		this.active = false
		err = this.instance.Close()
		if err != nil {
			fmt.Println(err)
		}
		close(this.list)
		close(this.pipe)
	}
}
func (this *TCPListener) SetVerify(value func(*net.TCPConn) bool) {
	this.verify = value
}
func (this *TCPListener) listen() {
	for this.active {
		var conn, err = this.instance.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if this.verify != nil {
			if this.verify(conn) {
				this.list <- conn
			}
		} else {
			this.list <- conn
		}
	}
}
func (this *TCPListener) handle() {
	for this.active {
		if conn, success := <-this.list; success {
			var tcp = &TCPSocket{}
			tcp.Constructor(conn)
			this.pipe <- tcp
		}
	}
}

// ==================== UDP Listener ====================
type UDPListener struct {
	sync.RWMutex

	address  *net.UDPAddr
	instance *net.UDPConn
	port     uint16

	verify func(*net.UDPAddr, []byte) bool

	list chan *net.UDPAddr
	pipe chan ISocket

	active bool
}

func (this *UDPListener) Listen(address string) (err error) {
	this.Lock()
	defer this.Unlock()

	this.list = make(chan *net.UDPAddr, 1024)
	this.pipe = make(chan ISocket, 1024)
	this.address, err = net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	this.instance, err = net.ListenUDP("udp", this.address)
	if err != nil {
		fmt.Println(err)
		return
	}
	this.active = true

	go this.listen()
	go this.handle()

	return
}
func (this *UDPListener) Accept() (connector *Connector) {
	var conn, ok = <-this.pipe
	if ok {
		connector = &Connector{}
		connector.Constructor(conn)
	}
	return
}
func (this *UDPListener) StopListen() {
	this.Lock()
	defer this.Unlock()

	var err error
	if this.active {
		this.active = false
		err = this.instance.Close()
		if err != nil {
			fmt.Println(err)
		}
		close(this.list)
		close(this.pipe)
	}
}
func (this *UDPListener) SetVerify(value func(*net.UDPAddr, []byte) bool) {
	this.verify = value
}
func (this *UDPListener) listen() {
	for this.active {
		var buffer = make([]byte, 2048)
		var lenght, addr, err = this.instance.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		if this.verify != nil {
			if this.verify(addr, buffer[:lenght]) {
				this.list <- addr
			}
		} else {
			this.list <- addr
		}
	}
}
func (this *UDPListener) handle() {
	for this.active {
		if remote, success := <-this.list; success {
			var err error
			var localAddressString string
			var localAddress *net.UDPAddr
			var listenConnect *net.UDPConn
			for {
				if this.port == 0 {
					//UDP监听，从5000开始
					this.port = 4999
				}
				this.port++

				localAddressString = this.address.IP.String() + ":" + strconv.Itoa(int(this.port))
				localAddress, err = net.ResolveUDPAddr("udp", this.address.IP.String()+":"+strconv.Itoa(int(this.port)))
				if err != nil {
					fmt.Println(err)
					break
				}
				listenConnect, err = net.ListenUDP("udp", localAddress)
				if err != nil {
					fmt.Println(err)
				} else {
					break
				}
			}
			if listenConnect != nil {
				var socket = &UDPSocket{}
				socket.Constructor(listenConnect, localAddress, nil)
				this.pipe <- socket

				this.instance.WriteToUDP([]byte(localAddressString), remote)
			}
		}
	}
}

// ==================== WebSocket Listener ====================
type WSListener struct {
	sync.RWMutex

	address  string
	path     string
	instance *websocket.Upgrader

	verify func(*http.Request) bool

	list chan *websocket.Conn
	pipe chan ISocket

	active bool
}

func (this *WSListener) Listen(address string, path string, tlsc *TLSConfig) (err error) {
	this.Lock()
	defer this.Unlock()

	this.address = address
	this.path = path
	this.list = make(chan *websocket.Conn, 1024)
	this.pipe = make(chan ISocket, 1024)
	this.instance = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	this.active = true

	if this.path == "" {
		this.path = "/"
	}
	http.HandleFunc(this.path, this.listen)

	go func() {
		if tlsc != nil && tlsc.CertFile != "" && tlsc.KeyFile != "" {
			http.ListenAndServeTLS(this.address, tlsc.CertFile, tlsc.KeyFile, nil)
		} else {
			http.ListenAndServe(this.address, nil)
		}
	}()
	go this.handle()

	return
}
func (this *WSListener) Accept() (connector *Connector) {
	var conn, ok = <-this.pipe
	if ok {
		connector = &Connector{}
		connector.Constructor(conn)
	}
	return
}
func (this *WSListener) StopListen() {
	this.Lock()
	defer this.Unlock()

	if this.active {
		this.active = false
		close(this.list)
		close(this.pipe)
	}
}
func (this *WSListener) SetVerify(value func(*http.Request) bool) {
	this.verify = value
}
func (this *WSListener) listen(w http.ResponseWriter, r *http.Request) {
	var ws, err = this.instance.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if this.verify != nil {
		if this.verify(r) {
			this.list <- ws
		}
	} else {
		this.list <- ws
	}
}
func (this *WSListener) handle() {
	for this.active {
		if conn, success := <-this.list; success {
			var socket = &WSocket{}
			socket.Constructor(conn)
			this.pipe <- socket
		}
	}
}

// ==================== ChanSocket Listener ====================
var chanListenerList map[string]*ChanListener

func init() {
	chanListenerList = make(map[string]*ChanListener)
}

type ChanListener struct {
	sync.RWMutex

	address string

	pipe chan ISocket

	active bool
}

func (this *ChanListener) Listen(address string) (err error) {
	if _, has := chanListenerList[address]; has {
		err = errors.New("Chan listner is occupied.")
	} else {
		this.address = address
		this.pipe = make(chan ISocket, 1024)
		this.active = true

		this.Lock()
		defer this.Unlock()

		chanListenerList[address] = this
	}
	return
}
func (this *ChanListener) Accept() (connector *Connector) {
	var conn, ok = <-this.pipe
	if ok {
		connector = &Connector{}
		connector.Constructor(conn)
	}
	return
}
func (this *ChanListener) StopListen() {
	this.Lock()
	defer this.Unlock()

	if this.active {
		this.active = false
		close(this.pipe)
	}
}
