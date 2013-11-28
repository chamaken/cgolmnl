package cgolmnl

import "unsafe"
// import "fmt"

/*
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
#include "cb.h"
*/
import "C"

type SocketDescriptor [0]byte // unsafe.Pointer

/**
 * mnl_socket_get_fd - obtain file descriptor from netlink socket
 *
 * int mnl_socket_get_fd(const struct mnl_socket *nl)
 */
func SocketGetFd(nl *SocketDescriptor) int {
	return int(C.mnl_socket_get_fd((*[0]byte)(nl)))
}

/**
 * mnl_socket_get_portid - obtain Netlink PortID from netlink socket
 *
 * uint32_t mnl_socket_get_portid(const struct mnl_socket *nl)
 */
func SocketGetPortid(nl *SocketDescriptor) uint32 {
	return uint32(C.mnl_socket_get_portid((*[0]byte)(nl)))
}

/**
 * mnl_socket_open - open a netlink socket
 * 
 * struct mnl_socket *mnl_socket_open(int bus)
 */
func SocketOpen(bus int) (*SocketDescriptor, error) {
	// return C.mnl_socket_open(C.int(bus))
	ret, err := C.mnl_socket_open(C.int(bus))
	return (*SocketDescriptor)(ret), err
}

/**
 * mnl_socket_bind - bind netlink socket
 *
 * int mnl_socket_bind(struct mnl_socket *nl, unsigned int groups, pid_t pid)
 */
func SocketBind(nl *SocketDescriptor, groups uint, pid Pid_t) error {
	_, err := C.mnl_socket_bind((*[0]byte)(nl), C.uint(groups), C.pid_t(pid))
	return err
}

/**
 * mnl_socket_sendto - send a netlink message of a certain size
 *
 * ssize_t
 * mnl_socket_sendto(const struct mnl_socket *nl, const void *buf, size_t len)
 */
func SocketSendto(nl *SocketDescriptor, buf []byte) (Ssize_t, error) {
	ret, err := C.mnl_socket_sendto((*[0]byte)(nl), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return Ssize_t(ret), err
}

func SocketSendNlmsg(nl *SocketDescriptor, nlh *Nlmsghdr) (Ssize_t, error) {
	b := C.GoBytes(unsafe.Pointer(nlh), C.int(nlh.Len))
	return SocketSendto(nl, b)
}
/**
 * mnl_socket_recvfrom - receive a netlink message
 *
 * ssize_t
 * mnl_socket_recvfrom(const struct mnl_socket *nl, void *buf, size_t bufsiz)
 */
func SocketRecvfrom(nl *SocketDescriptor, buf []byte) (Ssize_t, error) {
	ret, err := C.mnl_socket_recvfrom((*[0]byte)(nl), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return Ssize_t(ret), err
}

/**
 * mnl_socket_close - close a given netlink socket
 *
 * int mnl_socket_close(struct mnl_socket *nl)
 */
func SocketClose(nl *SocketDescriptor) error {
	_, err := C.mnl_socket_close((*[0]byte)(nl))
	return err
}

/**
 * mnl_socket_setsockopt - set Netlink socket option
 *
 * int mnl_socket_setsockopt(const struct mnl_socket *nl, int type,
 *			     void *buf, socklen_t len)
 */
func SocketSetsockopt(nl *SocketDescriptor, optype int, buf []byte) error {
	_, err := C.mnl_socket_setsockopt((*[0]byte)(nl), C.int(optype), unsafe.Pointer(&buf[0]), C.socklen_t(len(buf)))
	return err
}

/**
 * mnl_socket_getsockopt - get a Netlink socket option
 *
 * int mnl_socket_getsockopt(const struct mnl_socket *nl, int type,
 * 			     void *buf, socklen_t *len)
 */
func SocketGetsockopt(nl *SocketDescriptor, optype int, size Socklen_t) ([]byte, error) {
	c_size := C.socklen_t(size)
	buf := make([]byte, int(size))
	_, err := C.mnl_socket_getsockopt((*[0]byte)(nl), C.int(optype), unsafe.Pointer(&buf[0]), &c_size)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
