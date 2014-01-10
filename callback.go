package cgolmnl

/*
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
#include <linux/netlink.h>
#include "cb.h"
#include "set_errno.h"
*/
import "C"

import (
	"unsafe"
	"syscall"
	// "fmt"
	// "os"
)

type MnlCb func(*Nlmsghdr, interface{}) (int, syscall.Errno)
type MnlCtlCb func(*Nlmsghdr, uint16, interface{}) (int, syscall.Errno)

// argp (param data in original function) has three elements
//   0: function pointer to MnlCtlCb (see below, mnl_cb_run2)
//   1: function pointer to MnlCb
//   2: (real) data
//
// packing at CbRun2, CbRun
// unpacking below GoCb, GoCtlCb


//export GoCb
func GoCb(nlh *C.struct_nlmsghdr, argp unsafe.Pointer) C.int {
	args := *(*[3]unsafe.Pointer)(argp)
	cb := *(*MnlCb)(args[1])
	if cb == nil {
		return MNL_CB_OK
	}
	data := *(*interface{})(args[2])
	ret, err := cb((*Nlmsghdr)(unsafe.Pointer(nlh)), data) // returns (int, syscall.Errno)
	if err != 0 {
		// C.errno = int(err) // cannot refer to errno directly; see documentation
		C.SetErrno(C.int(err))
	}
	return C.int(ret)
}
	
//export GoCtlCb
func GoCtlCb(nlh *C.struct_nlmsghdr, argp unsafe.Pointer) C.int {
	args := *(*[3]unsafe.Pointer)(argp)
	cb := *(*MnlCtlCb)(args[0])
	if cb == nil {
		C.SetErrno(C.int(syscall.EOPNOTSUPP))
		return MNL_CB_ERROR
	}
	data := *(*interface{})(args[2])
	h := (*Nlmsghdr)(unsafe.Pointer(nlh))
	ret, err := cb(h, h.Type, data)
	if err != 0 {
		C.SetErrno(C.int(err))
	}
	return C.int(ret)
}

/**
 * mnl_cb_run2 - callback runqueue for netlink messages
 *
 * It seems callback Go function must be exported by //export
 * This means we can not dynamically create mnl_cb_t[]
 * To alleviate I introduct new C type mnl_ctl_cb_t
 *
 *   int (*mnl_ctl_cb_t)(const struct nlmsghdr *nlh, uint16_t type, void *data)
 *
 * int
 * mnl_cb_run2(const void *buf, size_t numbytes, unsigned int seq,
 *	       unsigned int portid, mnl_cb_t cb_data, void *data,
 *	       mnl_cb_t *cb_ctl_array, unsigned int cb_ctl_array_len)
 */
func CbRun2(buf []byte, seq, portid uint32, cb_data MnlCb, data interface{},
	cb_ctl MnlCtlCb, ctltypes []uint16) (int, error) {
	if len(ctltypes) >= C.NLMSG_MIN_TYPE {
		return MNL_CB_ERROR, syscall.EINVAL
	}

	args := [3]unsafe.Pointer{unsafe.Pointer(&cb_ctl), unsafe.Pointer(&cb_data), unsafe.Pointer(&data)}
	ret, err := C.cb_run2_wrapper(unsafe.Pointer(&buf[0]), C.size_t(len(buf)),
				      C.uint32_t(seq), C.uint32_t(portid), unsafe.Pointer(&args),
				      (*C.uint16_t)(&ctltypes[0]), C.size_t(len(ctltypes)))
	return int(ret), err
}

/**
 * mnl_cb_run - callback runqueue for netlink messages (simplified version)
 *
 * int
 * mnl_cb_run(const void *buf, size_t numbytes, uint32_t seq,
 *	      uint32_t portid, mnl_cb_t cb_data, void *data)
 */
func CbRun(buf []byte, seq, portid uint32, cb_data MnlCb, data interface{}) (int, error) {
	args := [3]unsafe.Pointer{unsafe.Pointer(nil), unsafe.Pointer(&cb_data), unsafe.Pointer(&data)}
	ret, err := C.cb_run_wrapper(unsafe.Pointer(&buf[0]), C.size_t(len(buf)),
				     C.uint32_t(seq), C.uint32_t(portid), unsafe.Pointer(&args))
	return int(ret), err
}
