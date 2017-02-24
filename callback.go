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

type MnlCb func(*Nlmsg, interface{}) (int, syscall.Errno)

// callback type for CbRun2
//
// Different from original, control message callback is not array of
// MnlCb type. This callback receives Netlink message type as second
// parameter. Your callback will dispatch by it.
type MnlCtlCb func(*Nlmsg, uint16, interface{}) (int, syscall.Errno)

type msgCbData struct {
	ctlcb MnlCb
	cb MnlCb
	data interface{}
}

// callback wrapper called from original cb_run(), cb_run2().
//export GoCb
func GoCb(nlh *C.struct_nlmsghdr, argp unsafe.Pointer) C.int {
	args := *(*msgCbData)(argp)
	if args.cb == nil {
		return MNL_CB_OK
	}
	h := nlmsgPointer(nlh)
	ret, err := args.cb(h, args.data) // returns (int, syscall.Errno)
	if err != 0 {
		// C.errno = int(err) // cannot refer to errno directly; see documentation
		C.SetErrno(C.int(err))
	}
	return C.int(ret)
}

// control message callback wrapper called from original cb_run2()
//export GoCtlCb
func GoCtlCb(nlh *C.struct_nlmsghdr, argp unsafe.Pointer) C.int {
	args := *(*msgCbData)(argp)
	if args.ctlcb == nil {
		C.SetErrno(C.int(syscall.EOPNOTSUPP))
		return MNL_CB_ERROR
	}
	h := nlmsgPointer(nlh)
	ret, err := args.ctlcb(h, args.data)
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
	cb_ctl MnlCb, ctltypes []uint16) (int, error) {
	if len(ctltypes) >= C.NLMSG_MIN_TYPE {
		return MNL_CB_ERROR, syscall.EINVAL
	}
	args := uintptr(unsafe.Pointer(&msgCbData{cb_ctl, cb_data, data}))
	ret, err := C.cb_run2_wrapper(unsafe.Pointer(&buf[0]), C.size_t(len(buf)),
		C.uint32_t(seq), C.uint32_t(portid), C.uintptr_t(args),
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
	args := uintptr(unsafe.Pointer(&msgCbData{nil, cb_data, data}))
	ret, err := C.cb_run_wrapper(unsafe.Pointer(&buf[0]), C.size_t(len(buf)),
		C.uint32_t(seq), C.uint32_t(portid), C.uintptr_t(args))
	return int(ret), err
}
