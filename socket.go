package cgolmnl

import "unsafe"

// import "fmt"

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
*/
import "C"

// Netlink socket helpers
type Socket C.struct_mnl_socket // [0]byte

// int mnl_socket_get_fd(const struct mnl_socket *nl)
func socketGetFd(nl *Socket) int {
	return int(C.mnl_socket_get_fd((*C.struct_mnl_socket)(nl)))
}

// uint32_t mnl_socket_get_portid(const struct mnl_socket *nl)
func socketGetPortid(nl *Socket) uint32 {
	return uint32(C.mnl_socket_get_portid((*C.struct_mnl_socket)(nl)))
}

// struct mnl_socket *mnl_socket_open(int bus)
func socketOpen(bus int) (*Socket, error) {
	// return C.mnl_socket_open(C.int(bus))
	ret, err := C.mnl_socket_open(C.int(bus))
	return (*Socket)(ret), err
}

// int mnl_socket_bind(struct mnl_socket *nl, unsigned int groups, pid_t pid)
func socketBind(nl *Socket, groups uint, pid Pid_t) error {
	_, err := C.mnl_socket_bind((*C.struct_mnl_socket)(nl), C.uint(groups), C.pid_t(pid))
	return err
}

// ssize_t
// mnl_socket_sendto(const struct mnl_socket *nl, const void *buf, size_t len)
func socketSendto(nl *Socket, buf []byte) (Ssize_t, error) {
	ret, err := C.mnl_socket_sendto((*C.struct_mnl_socket)(nl), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return Ssize_t(ret), err
}
func socketSendNlmsg(nl *Socket, nlh *Nlmsg) (Ssize_t, error) {
	ret, err := C.mnl_socket_sendto((*C.struct_mnl_socket)(nl), unsafe.Pointer(nlh.Nlmsghdr), C.size_t(nlh.Len))
	return Ssize_t(ret), err
}

// ssize_t
// mnl_socket_recvfrom(const struct mnl_socket *nl, void *buf, size_t bufsiz)
func socketRecvfrom(nl *Socket, buf []byte) (Ssize_t, error) {
	ret, err := C.mnl_socket_recvfrom((*C.struct_mnl_socket)(nl), unsafe.Pointer(&buf[0]), C.size_t(len(buf)))
	return Ssize_t(ret), err
}

// int mnl_socket_close(struct mnl_socket *nl)
func socketClose(nl *Socket) error {
	_, err := C.mnl_socket_close((*C.struct_mnl_socket)(nl))
	return err
}

// int mnl_socket_setsockopt(const struct mnl_socket *nl, int type,
//			     void *buf, socklen_t len)
func socketSetsockopt(nl *Socket, optype int, optval unsafe.Pointer, optlen Socklen_t) error {
	_, err := C.mnl_socket_setsockopt((*C.struct_mnl_socket)(nl), C.int(optype), optval, C.socklen_t(optlen))
	return err
}
func socketSetsockoptBytes(nl *Socket, optype int, buf []byte) error {
	_, err := C.mnl_socket_setsockopt((*C.struct_mnl_socket)(nl), C.int(optype), unsafe.Pointer(&buf[0]), C.socklen_t(len(buf)))
	return err
}
func socketSetsockoptByte(nl *Socket, opt int, value byte) error {
	v := C.uint8_t(value)
	_, err := C.mnl_socket_setsockopt((*C.struct_mnl_socket)(nl), C.int(opt), unsafe.Pointer(&v), 1)
	return err
}
func socketSetsockoptCint(nl *Socket, opt int, value int) error {
	v := C.int(value)
	_, err := C.mnl_socket_setsockopt((*C.struct_mnl_socket)(nl), C.int(opt), unsafe.Pointer(&v), C.sizeof_int)
	return err
}

// int mnl_socket_getsockopt(const struct mnl_socket *nl, int type,
// 			     void *buf, socklen_t *len)
func socketGetsockopt(nl *Socket, optype int, size Socklen_t) ([]byte, error) {
	c_size := C.socklen_t(size)
	buf := make([]byte, int(size))
	_, err := C.mnl_socket_getsockopt((*C.struct_mnl_socket)(nl), C.int(optype), unsafe.Pointer(&buf[0]), &c_size)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
