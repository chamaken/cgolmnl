package cgolmnl

/*
#include <libmnl/libmnl.h>
*/
import "C"

import (
	"errors"
	"unsafe"
	"syscall"
)

// make copy

func (nlh *Nlmsghdr) MarshalBinary() ([]byte, error) {
	dst := make([]byte, nlh.Len)
	copy(dst, C.GoBytes(unsafe.Pointer(nlh), C.int(nlh.Len)))
	return dst, nil
}

func (nlh *Nlmsghdr) UnmarshalBinary(data []byte) error {
	if len(data) < int(MNL_NLMSG_HDRLEN) {
		return syscall.EINVAL // errors.New("too short length")
	}
	h := (*Nlmsghdr)(unsafe.Pointer(&data[0]))
	if int(h.Len) > len(data) {
		return syscall.EINVAL // errors.New("invalid length field")
	}

	dst := make([]byte, len(data))
	copy(dst, data)
	nlh = (*Nlmsghdr)(unsafe.Pointer(&data[0]))

	return nil
}

func (attr *Nlattr) MarshalBinary() ([]byte, error) {
	dst := make([]byte, attr.Len)
	copy(dst, C.GoBytes(unsafe.Pointer(attr), C.int(attr.Len)))
	return dst, nil
}

func (attr *Nlattr) UnmarshalBinary(data []byte) error {
	if len(data) < SizeofNlattr {
		return syscall.EINVAL
	}
	h := (*Nlattr)(unsafe.Pointer(&data[0]))
	if int(h.Len) > len(data) {
		return syscall.EINVAL
	}

	dst := make([]byte, len(data))
	copy(dst, data)
	attr = (*Nlattr)(unsafe.Pointer(&data[0]))

	return nil
}


// constructor specifying buffer size
func NewNlmsghdr(size int) (*Nlmsghdr, error) {
	if size < int(MNL_NLMSG_HDRLEN) {
		return nil, errors.New("too short size")
	}
	b := make([]byte, size)
	return (*Nlmsghdr)(unsafe.Pointer(&b[0])), nil
}

func (nlh *Nlmsghdr) PutHeader() {
	C.mnl_nlmsg_put_header(unsafe.Pointer(nlh))
}
