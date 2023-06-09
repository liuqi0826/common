package network

func CreateTCPListener(address string) (listener *TCPListener, err error) {
	listener = &TCPListener{}
	err = listener.Listen(address)
	return
}
func CreateUDPListener(address string) (listener *UDPListener, err error) {
	listener = &UDPListener{}
	err = listener.Listen(address)
	return
}
func CreateWSListener(address, path string, tlsc *TLSConfig) (listener *WSListener, err error) {
	listener = &WSListener{}
	err = listener.Listen(address, path, tlsc)
	return
}
func CreateChanListener(address string) (listener *ChanListener, err error) {
	listener = &ChanListener{}
	err = listener.Listen(address)
	return
}
