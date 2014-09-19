package cgolmnl

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
#include <linux/netlink.h>
#include "cb.h"
#include "set_errno.h"
*/
import "C"

import (
	"syscall"
	"unsafe"
	// "fmt"
	// "os"
)

type MnlCb func(*Nlmsghdr, interface{}) (int, syscall.Errno)

// callback type for CbRun2
//
// Different from original, control message callback is not array of
// MnlCb type. This callback receives Netlink message type as second
// parameter. Your callback will dispatch by it.
type MnlCtlCb func(*Nlmsghdr, uint16, interface{}) (int, syscall.Errno)

// argp (param data in original function) has three elements
//   0: function pointer to MnlCtlCb (see below, mnl_cb_run2)
//   1: function pointer to MnlCb
//   2: (real) data
//
// packing at CbRun2, CbRun
// unpacking at GoCb, GoCtlCb

// callback wrapper called from original cb_run(), cb_run2().
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

// control message callback wrapper called from original cb_run2()
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

// callback runqueue for netlink messages
//
// You can set the cb_ctl to nil if you want to use the default control
// callback handlers. Or set it with ctltypes which contains control message
// types to handle.
//
// Your callback may return three possible values:
// 	- MNL_CB_ERROR (<=-1): an error has occurred. Stop callback runqueue.
// 	- MNL_CB_STOP (=0): stop callback runqueue.
// 	- MNL_CB_OK (>=1): no problem has occurred.
//
// This function propagates the callback return value. On error, it returns
// -1 and errno is explicitly set. If the portID is not the expected, errno
// is set to ESRCH. If the sequence number is not the expected, errno is set
// to EPROTO. If the dump was interrupted, errno is set to EINTR and you should
// request a new fresh dump again.
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

// callback runqueue for netlink messages (simplified version)
//
// This function is like CbRun2() but it does not allow you to set
// the control callback handlers.
//
// Your callback may return three possible values:
// 	- MNL_CB_ERROR (<=-1): an error has occurred. Stop callback runqueue.
// 	- MNL_CB_STOP (=0): stop callback runqueue.
// 	- MNL_CB_OK (>=1): no problems has occurred.
//
// This function propagates the callback return value.
func CbRun(buf []byte, seq, portid uint32, cb_data MnlCb, data interface{}) (int, error) {
	args := [3]unsafe.Pointer{unsafe.Pointer(nil), unsafe.Pointer(&cb_data), unsafe.Pointer(&data)}
	ret, err := C.cb_run_wrapper(unsafe.Pointer(&buf[0]), C.size_t(len(buf)),
		C.uint32_t(seq), C.uint32_t(portid), unsafe.Pointer(&args))
	return int(ret), err
}
