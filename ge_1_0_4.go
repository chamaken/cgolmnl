// +build ge_1_0_4

package cgolmnl

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
*/
import "C"

// struct mnl_socket *mnl_socket_fdopen(int fd)
func socketFdopen(fd int) (*Socket, error) {
	ret, err := C.mnl_socket_fdopen(C.int(fd))
	return (*Socket)(ret), err
}

// associates a mnl_socket object with pre-existing socket.
//
// On error, it returns NULL and errno is appropriately set. Otherwise, it
// returns a valid pointer to the mnl_socket structure. It also sets the portID
// if the socket fd is already bound and it is AF_NETLINK.
//
// Note that mnl_socket_get_portid() returns 0 if this function is used with
// non-netlink socket.
func NewSocketFd(fd int) (*Socket, error) {
	return socketFdopen(fd)
}

// struct mnl_socket *mnl_socket_open2(int bus, int flags)
func socketOpen2(bus int, flags int) (*Socket, error) {
	ret, err := C.mnl_socket_open2(C.int(bus), C.int(flags))
	return (*Socket)(ret), err
}

// open a netlink socket with appropriate flags
// - bus the netlink socket bus ID (see NETLINK_* constants)
// - flags the netlink socket flags (see SOCK_* constants in socket(2))
//
// This is similar to mnl_socket_open(), but allows to set flags like
// SOCK_CLOEXEC at socket creation time (useful for multi-threaded programs
// performing exec calls).
//
// On error, it returns NULL and errno is appropriately set. Otherwise, it
// returns a valid pointer to the mnl_socket structure.
func NewSocket2(bus int, flags int) (*Socket, error) {
	return socketOpen2(bus, flags)
}
