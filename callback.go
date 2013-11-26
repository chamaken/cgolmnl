package cgolmnl

/*
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
#include "cb.h"
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
//   0: function pointer to MnlCtlCb
//   1: function pointer to MnlCb
//   2: (real) data
//
// packing at CbRun.?()
// unpacking below GoCb, GoCtlCb


//export GoCb
func GoCb(nlh *C.struct_nlmsghdr, argp unsafe.Pointer) C.int {
	args := *(*[3]unsafe.Pointer)(argp)
	cb := *(*MnlCb)(args[1])
	data := *(*interface{})(args[2])
	ret, _ := cb((*Nlmsghdr)(unsafe.Pointer(nlh)), data) // returns (int, syscall.Errno)
	// if err != nil {
	//	C.errno = int(err) // cannot refer to errno directly; see documentation
	// }
	return C.int(ret)
}
	
//export GoCtlCb
func GoCtlCb(nlh *C.struct_nlmsghdr, msgtype C.uint16_t, argp unsafe.Pointer) C.int {
	args := *(*[3]unsafe.Pointer)(argp)
	cb := *(*MnlCtlCb)(args[0])
	data := *(*interface{})(args[2])
	ret, _ := cb((*Nlmsghdr)(unsafe.Pointer(nlh)), uint16(msgtype), data)
	return C.int(ret)
}
	
/**
 * mnl_cb_run2 - callback runqueue for netlink messages
 *
 * int
 * mnl_cb_run2(const void *buf, size_t numbytes, unsigned int seq,
 *	       unsigned int portid, mnl_cb_t cb_data, void *data,
 *	       mnl_cb_t *cb_ctl_array, unsigned int cb_ctl_array_len)
 */
func CbRun3(buf []byte, seq, portid uint32,
	cb_data MnlCb, data interface{}, cb_ctl MnlCtlCb) int {
	args := [3]unsafe.Pointer{unsafe.Pointer(&cb_ctl), unsafe.Pointer(&cb_data), unsafe.Pointer(&data)}
	return int(C.cb_run3_wrapper(unsafe.Pointer(&buf[0]), C.size_t(len(buf)),
		C.uint32_t(seq), C.uint32_t(portid), unsafe.Pointer(&args)))
}

/**
 * mnl_cb_run - callback runqueue for netlink messages (simplified version)
 *
 * int
 * mnl_cb_run(const void *buf, size_t numbytes, uint32_t seq,
 *	      uint32_t portid, mnl_cb_t cb_data, void *data)
 */
func CbRun(buf []byte, seq, portid uint32, cb_data MnlCb, data interface{}) int {
	args := [3]unsafe.Pointer{unsafe.Pointer(nil), unsafe.Pointer(&cb_data), unsafe.Pointer(&data)}
	return int(C.cb_run_wrapper(unsafe.Pointer(&buf[0]), C.size_t(len(buf)),
		C.uint32_t(seq), C.uint32_t(portid), unsafe.Pointer(&args)))
}
