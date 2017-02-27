package cgolmnl

import (
	"reflect"
	"syscall"
	"unsafe"
	// "fmt"
	// "os"
)

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <stdlib.h>
#include <libmnl/libmnl.h>
#include "cb.h"
#include "set_errno.h"
*/
import "C"

// uint16_t mnl_attr_get_type(const struct nlattr *attr)
func attrGetType(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_type((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// uint16_t mnl_attr_get_len(const struct nlattr *attr)
func attrGetLen(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_len((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// uint16_t mnl_attr_get_payload_len(const struct nlattr *attr)
func attrGetPayloadLen(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_payload_len((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// void *mnl_attr_get_payload(const struct nlattr *attr)
func attrGetPayload(attr *Nlattr) unsafe.Pointer {
	return C.mnl_attr_get_payload((*C.struct_nlattr)(unsafe.Pointer(attr)))
}
func attrGetPayloadBytes(attr *Nlattr) []byte {
	return SharedBytes(attrGetPayload(attr), int(attrGetPayloadLen(attr)))
}

// bool mnl_attr_ok(const struct nlattr *attr, int len)
func attrOk(attr *Nlattr, size int) bool {
	return bool(C.mnl_attr_ok((*C.struct_nlattr)(unsafe.Pointer(attr)), C.int(size)))
}

// struct nlattr *mnl_attr_next(const struct nlattr *attr)
func attrNext(attr *Nlattr) *Nlattr {
	return (*Nlattr)(unsafe.Pointer(C.mnl_attr_next((*C.struct_nlattr)(unsafe.Pointer(attr)))))
}

// int mnl_attr_type_valid(const struct nlattr *attr, uint16_t max)
func attrTypeValid(attr *Nlattr, max uint16) error {
	_, err := C.mnl_attr_type_valid((*C.struct_nlattr)(unsafe.Pointer(attr)), C.uint16_t(max))
	return err
}

// int mnl_attr_validate(const struct nlattr *attr, enum mnl_attr_data_type type)
func attrValidate(attr *Nlattr, data_type AttrDataType) error {
	_, err := C.mnl_attr_validate((*C.struct_nlattr)(unsafe.Pointer(attr)), C.enum_mnl_attr_data_type(data_type))
	return err
}

// int
// mnl_attr_validate2(const struct nlattr *attr, enum mnl_attr_data_type type,
//		      size_t exp_len)
func attrValidate2(attr *Nlattr, data_type AttrDataType, exp_len Size_t) error {
	_, err := C.mnl_attr_validate2((*C.struct_nlattr)(unsafe.Pointer(attr)), C.enum_mnl_attr_data_type(data_type), C.size_t(exp_len))
	return err
}

// attribute callback wrapper called from original
type MnlAttrCb func(*Nlattr, interface{}) (int, syscall.Errno)

type attrCbData struct {
	cb MnlAttrCb
	data interface{}
}

//export GoAttrCb
func GoAttrCb(nla *C.struct_nlattr, argp unsafe.Pointer) C.int {
	args := *(*attrCbData)(argp)
	if args.cb == nil {
		return MNL_CB_OK
	}
	ret, err := args.cb((*Nlattr)(unsafe.Pointer(nla)), args.data) // returns (int, syscall.Errno)
	if err != 0 {
		// C.errno = int(err) // ``cannot refer to errno directly; see documentation''
		C.SetErrno(C.int(err))
	}
	return C.int(ret)
}

// int
// mnl_attr_parse(const struct nlmsghdr *nlh, unsigned int offset,
//	          mnl_attr_cb_t cb, void *data)
func attrParse(nlh *Nlmsg, offset Size_t, cb MnlAttrCb, data interface{}) (int, error) {
	args := uintptr(unsafe.Pointer(&attrCbData{cb, data}))
	ret, err := C.attr_parse_wrapper(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.size_t(offset), C.uintptr_t(args))
	return int(ret), err
}

// int
// mnl_attr_parse_nested(const struct nlattr *nested, mnl_attr_cb_t cb,
// 			 void *data)
func attrParseNested(nested *Nlattr, cb MnlAttrCb, data interface{}) (int, error) {
	args := uintptr(unsafe.Pointer(&attrCbData{cb, data}))
	ret, err := C.attr_parse_nested_wrapper((*C.struct_nlattr)(unsafe.Pointer(nested)), C.uintptr_t(args))
	return int(ret), err
}

// parse attributes in payload of Netlink message
//
// This function takes a pointer to the area that contains the attributes,
// commonly known as the payload of the Netlink message. Thus, you have to
// pass a pointer to the Netlink message payload, instead of the entire
// message.
//
// This function allows you to iterate over the sequence of attributes that are
// located at some payload offset. You can then put the attributes in one array
// as usual, or you can use any other data structure (such as lists or trees).
//
// This function propagates the return value of the callback, which can be
// MNL_CB_ERROR, MNL_CB_OK or MNL_CB_STOP.
func AttrParsePayload(payload []byte, cb MnlAttrCb, data interface{}) (int, error) {
	args := uintptr(unsafe.Pointer(&attrCbData{cb, data}))
	ret, err := C.attr_parse_payload_wrapper(unsafe.Pointer(&payload[0]), C.size_t(len(payload)), C.uintptr_t(args))
	return int(ret), err
}

// uint8_t mnl_attr_get_u8(const struct nlattr *attr)
func attrGetU8(attr *Nlattr) uint8 {
	return uint8(C.mnl_attr_get_u8((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// uint16_t mnl_attr_get_u16(const struct nlattr *attr)
func attrGetU16(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_u16((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// uint32_t mnl_attr_get_u32(const struct nlattr *attr)
func attrGetU32(attr *Nlattr) uint32 {
	return uint32(C.mnl_attr_get_u32((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// uint64_t mnl_attr_get_u64(const struct nlattr *attr)
func attrGetU64(attr *Nlattr) uint64 {
	return uint64(C.mnl_attr_get_u64((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// const char *mnl_attr_get_str(const struct nlattr *attr)
func attrGetStr(attr *Nlattr) string {
	return C.GoString(C.mnl_attr_get_str((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

// void
// mnl_attr_put(struct nlmsghdr *nlh, uint16_t type, size_t len, const void *data)
func attrPut(nlh *Nlmsg, attr_type uint16, size Size_t, data unsafe.Pointer) {
	C.mnl_attr_put(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type), C.size_t(size), data)
}
func attrPutPtr(nlh *Nlmsg, attr_type uint16, data interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		panic("pointer required for data")
	}
	t := reflect.Indirect(v).Type()
	attrPut(nlh, attr_type, Size_t(t.Size()), unsafe.Pointer(v.Pointer()))
}
func attrPutBytes(nlh *Nlmsg, attr_type uint16, data []byte) {
	C.mnl_attr_put(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type), C.size_t(len(data)), unsafe.Pointer(&data[0]))
}

// void mnl_attr_put_u8(struct nlmsghdr *nlh, uint16_t type, uint8_t data)
func attrPutU8(nlh *Nlmsg, attr_type uint16, data uint8) {
	C.mnl_attr_put_u8(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type),
		C.uint8_t(data))
}

// void mnl_attr_put_u16(struct nlmsghdr *nlh, uint16_t type, uint16_t data)
func attrPutU16(nlh *Nlmsg, attr_type uint16, data uint16) {
	C.mnl_attr_put_u16(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type), C.uint16_t(data))
}

// void mnl_attr_put_u32(struct nlmsghdr *nlh, uint16_t type, uint32_t data)
func attrPutU32(nlh *Nlmsg, attr_type uint16, data uint32) {
	C.mnl_attr_put_u32(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type), C.uint32_t(data))
}

// void mnl_attr_put_u64(struct nlmsghdr *nlh, uint16_t type, uint64_t data)
func attrPutU64(nlh *Nlmsg, attr_type uint16, data uint64) {
	C.mnl_attr_put_u64(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type), C.uint64_t(data))
}

// void mnl_attr_put_str(struct nlmsghdr *nlh, uint16_t type, const char *data)
func attrPutStr(nlh *Nlmsg, attr_type uint16, data string) {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	C.mnl_attr_put_str(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type), cs)
}

// void mnl_attr_put_strz(struct nlmsghdr *nlh, uint16_t type, const char *data)
func attrPutStrz(nlh *Nlmsg, attr_type uint16, data string) {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	C.mnl_attr_put_strz(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.uint16_t(attr_type), cs)
}

// struct nlattr *mnl_attr_nest_start(struct nlmsghdr *nlh, uint16_t type)
func attrNestStart(nlh *Nlmsg, attr_type uint16) *Nlattr {
	return (*Nlattr)(unsafe.Pointer(
		C.mnl_attr_nest_start(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.uint16_t(attr_type))))
}

// bool
// mnl_attr_put_check(struct nlmsghdr *nlh, size_t buflen,
//		      uint16_t type, size_t len, const void *data)
func attrPutCheck(nlh *Nlmsg, attr_type uint16, size Size_t, data unsafe.Pointer) bool {
	return bool(
		C.mnl_attr_put_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type),
			C.size_t(size), data))
}
func attrPutCheckPtr(nlh *Nlmsg, attr_type uint16, data interface{}) bool {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		panic("pointer required for data")
	}
	t := reflect.Indirect(v).Type()
	return attrPutCheck(nlh, attr_type, Size_t(t.Size()), unsafe.Pointer(v.Pointer()))
}
func attrPutCheckBytes(nlh *Nlmsg, attr_type uint16, data []byte) bool {
	return bool(
		C.mnl_attr_put_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type),
			C.size_t(len(data)), unsafe.Pointer(&data[0])))
}

// bool
// mnl_attr_put_u8_check(struct nlmsghdr *nlh, size_t buflen,
// 			 uint16_t type, uint8_t data)
func attrPutU8Check(nlh *Nlmsg, attr_type uint16, data uint8) bool {
	return bool(
		C.mnl_attr_put_u8_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type), C.uint8_t(data)))
}

// bool
// mnl_attr_put_u16_check(struct nlmsghdr *nlh, size_t buflen,
//			  uint16_t type, uint16_t data)
func attrPutU16Check(nlh *Nlmsg, attr_type uint16, data uint16) bool {
	return bool(
		C.mnl_attr_put_u16_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type), C.uint16_t(data)))
}

// bool
// mnl_attr_put_u32_check(struct nlmsghdr *nlh, size_t buflen,
//			  uint16_t type, uint32_t data)
func attrPutU32Check(nlh *Nlmsg, attr_type uint16, data uint32) bool {
	return bool(
		C.mnl_attr_put_u32_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type), C.uint32_t(data)))
}

// bool
// mnl_attr_put_u64_check(struct nlmsghdr *nlh, size_t buflen,
//			  uint16_t type, uint64_t data)
func attrPutU64Check(nlh *Nlmsg, attr_type uint16, data uint64) bool {
	return bool(
		C.mnl_attr_put_u64_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type), C.uint64_t(data)))
}

// bool
// mnl_attr_put_str_check(struct nlmsghdr *nlh, size_t buflen,
//			  uint16_t type, const char *data)
func attrPutStrCheck(nlh *Nlmsg, attr_type uint16, data string) bool {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	return bool(
		C.mnl_attr_put_str_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type), cs))
}

// bool
// mnl_attr_put_strz_check(struct nlmsghdr *nlh, size_t buflen,
//			   uint16_t type, const char *data)
func attrPutStrzCheck(nlh *Nlmsg, attr_type uint16, data string) bool {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	return bool(
		C.mnl_attr_put_strz_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type), cs))
}

// struct nlattr *
// mnl_attr_nest_start_check(struct nlmsghdr *nlh, size_t buflen, uint16_t type)
func attrNestStartCheck(nlh *Nlmsg, attr_type uint16) *Nlattr {
	return (*Nlattr)(unsafe.Pointer(
		C.mnl_attr_nest_start_check(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.size_t(len(nlh.buf)), C.uint16_t(attr_type))))
}

// void
// mnl_attr_nest_end(struct nlmsghdr *nlh, struct nlattr *start)
func attrNestEnd(nlh *Nlmsg, start *Nlattr) {
	C.mnl_attr_nest_end(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		(*C.struct_nlattr)(unsafe.Pointer(start)))
}

// void
// mnl_attr_nest_cancel(struct nlmsghdr *nlh, struct nlattr *start)
func attrNestCancel(nlh *Nlmsg, start *Nlattr) {
	C.mnl_attr_nest_cancel(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		(*C.struct_nlattr)(unsafe.Pointer(start)))
}
