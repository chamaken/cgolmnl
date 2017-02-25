package cgolmnl

import (
	"C"
	"reflect"
	"syscall"
	"unsafe"
)

// make copy

func (nlh *Nlmsg) MarshalBinary() ([]byte, error) {
	dst := make([]byte, nlh.Len)
	copy(dst, C.GoBytes(unsafe.Pointer(nlh.Nlmsghdr), C.int(nlh.Len)))
	return dst, nil
}

// confide receiver len
func (nlh *Nlmsg) UnmarshalBinary(data []byte) error {
	if len(data) < int(MNL_NLMSG_HDRLEN) {
		return syscall.EINVAL // errors.New("too short data length")
	}
	if int(nlh.Len) < len(data) {
		return syscall.EINVAL // errors.New("too short receiver size")
	}

	var dst []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&dst))
	h.Cap = int(nlh.Len)
	h.Len = int(nlh.Len)
	h.Data = uintptr(unsafe.Pointer(nlh.Nlmsghdr))
	copy(dst, data)

	return nil
}

func (attr *Nlattr) MarshalBinary() ([]byte, error) {
	dst := make([]byte, attr.Len)
	copy(dst, C.GoBytes(unsafe.Pointer(attr), C.int(attr.Len)))
	return dst, nil
}

// confide receiver len
func (attr *Nlattr) UnmarshalBinary(data []byte) error {
	if len(data) < SizeofNlattr {
		return syscall.EINVAL
	}
	if int(attr.Len) < len(data) {
		return syscall.EINVAL
	}

	var dst []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&dst))
	h.Cap = int(attr.Len)
	h.Len = int(attr.Len)
	h.Data = uintptr(unsafe.Pointer(attr))
	copy(dst, data)

	return nil
}

// C.GoBytes() returns copy
func SharedBytes(p unsafe.Pointer, plen int) []byte {
	var b []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Cap = plen
	h.Len = plen
	h.Data = uintptr(p)
	return b
}
