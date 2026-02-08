package udprobe

import (
	"net"

	"golang.org/x/sys/unix" // The successor to syscall
)

// LocalUDPAddr returns the UDPAddr and net for the provided UDPConn.
//
// For UDPConn instances, net is generaly 'udp'.
func LocalUDPAddr(conn *net.UDPConn) (*net.UDPAddr, string, error) {
	addr := conn.LocalAddr()
	network := addr.Network()
	udpAddr, err := net.ResolveUDPAddr(network, addr.String())
	if err != nil {
		return udpAddr, network, err
	}
	return udpAddr, network, nil
}

// SetTos will set the IP_TOS value for the unix socket for the provided conn.
func SetTos(conn *net.UDPConn, tos byte) {
	file, err := conn.File()
	defer FileCloseHandler(file)
	HandleError(err)
	err = unix.SetsockoptByte(int(file.Fd()), unix.IPPROTO_IP,
		unix.IP_TOS, tos)
	HandleError(err)
}

// GetTos will get the IP_TOS value for the unix socket for the provided conn.
func GetTos(conn *net.UDPConn) byte {
	file, err := conn.File()
	defer FileCloseHandler(file)
	HandleError(err)
	value, err := unix.GetsockoptInt(int(file.Fd()), unix.IPPROTO_IP,
		unix.IP_TOS)
	HandleError(err)
	// Convert it to a byte and return
	return byte(value)
}

// EnableTimestamps enables kernel receive timestamping of packets on the
// provided conn.
//
// The timestamp values can later be extracted in the oob data from
// Receive.
func EnableTimestamps(conn *net.UDPConn) {
	file, err := conn.File()
	defer FileCloseHandler(file)
	HandleError(err)
	err = unix.SetsockoptInt(int(file.Fd()), unix.SOL_SOCKET,
		unix.SO_TIMESTAMPNS, 1)
	HandleError(err)
}
