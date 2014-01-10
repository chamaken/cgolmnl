package cgolmnl

import (
	"reflect"
	"syscall"
	"unsafe"
	// "fmt"
	// "os"
)

/*
#cgo LDFLAGS: -lmnl
#include <stdlib.h>
#include <libmnl/libmnl.h>
#include "cb.h"
#include "set_errno.h"
*/
import "C"

/**
 * mnl_attr_get_type - get type of netlink attribute
 *
 * uint16_t mnl_attr_get_type(const struct nlattr *attr)
 */
func AttrGetType(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_type((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_get_len - get length of netlink attribute
 *
 * uint16_t mnl_attr_get_len(const struct nlattr *attr)
 */
func AttrGetLen(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_len((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_get_payload_len - get the attribute payload-value length
 *
 * uint16_t mnl_attr_get_payload_len(const struct nlattr *attr)
 */
func AttrGetPayloadLen(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_payload_len((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_get_payload - get pointer to the attribute payload
 *
 * void *mnl_attr_get_payload(const struct nlattr *attr)
 */
func AttrGetPayload(attr *Nlattr) unsafe.Pointer {
	return C.mnl_attr_get_payload((*C.struct_nlattr)(unsafe.Pointer(attr)))
}

func AttrGetPayloadBytes(attr *Nlattr) []byte {
	return SharedBytes(AttrGetPayload(attr), int(AttrGetPayloadLen(attr)))
}

/**
 * mnl_attr_ok - check if there is room for an attribute in a buffer
 *
 * bool mnl_attr_ok(const struct nlattr *attr, int len)
 */
func AttrOk(attr *Nlattr, size int) bool {
	return bool(C.mnl_attr_ok((*C.struct_nlattr)(unsafe.Pointer(attr)), C.int(size)))
}

/**
 * mnl_attr_next - get the next attribute in the payload of a netlink message
 *
 * struct nlattr *mnl_attr_next(const struct nlattr *attr)
 */
func AttrNext(attr *Nlattr) *Nlattr {
	return (*Nlattr)(unsafe.Pointer(C.mnl_attr_next((*C.struct_nlattr)(unsafe.Pointer(attr)))))
}

/**
 * mnl_attr_type_valid - check if the attribute type is valid
 *
 * converting error returned from C. to int because
 * this function may be called from callback in mnl_attr_parse(), mnl_attr_parse_nested(),
 * mnl_attr_parse_payload()
 *
 * int mnl_attr_type_valid(const struct nlattr *attr, uint16_t max)
 */
func AttrTypeValid(attr *Nlattr, max uint16) (int, error) {
	ret, err := C.mnl_attr_type_valid((*C.struct_nlattr)(unsafe.Pointer(attr)), C.uint16_t(max))
	return int(ret), err
}

/**
 * mnl_attr_validate - validate netlink attribute (simplified version)
 *
 * int mnl_attr_validate(const struct nlattr *attr, enum mnl_attr_data_type type)
 */
func AttrValidate(attr *Nlattr, data_type AttrDataType) (int, error) {
	ret, err := C.mnl_attr_validate((*C.struct_nlattr)(unsafe.Pointer(attr)), C.enum_mnl_attr_data_type(data_type))
	return int(ret), err
}

/**
 * mnl_attr_validate2 - validate netlink attribute (extended version)
 *
 * int
 * mnl_attr_validate2(const struct nlattr *attr, enum mnl_attr_data_type type,
 *		      size_t exp_len)
 */
func AttrValidate2(attr *Nlattr, data_type AttrDataType, exp_len Size_t) (int, error) {
	ret, err := C.mnl_attr_validate2((*C.struct_nlattr)(unsafe.Pointer(attr)), C.enum_mnl_attr_data_type(data_type), C.size_t(exp_len))
	return int(ret), err
}

/** https://groups.google.com/forum/#!topic/golang-nuts/PRcvOJqItow
 *
 * Unfortunately you can't pass a Go func to C code and have the C code
 * call it.  The best you can do is pass a Go func to C code and have the C
 * code turn around and pass the Go func back to a Go function that then
 * calls the func. 
 */

type MnlAttrCb func(*Nlattr, interface{}) (int, syscall.Errno)

//export GoAttrCb
func GoAttrCb(nla *C.struct_nlattr, argp unsafe.Pointer) C.int {
	args := *(*[2]unsafe.Pointer)(argp)
	cb := *(*MnlAttrCb)(args[0])
	if cb == nil {
		return MNL_CB_OK
	}
	data := *(*interface{})(args[1])
	ret, err := cb((*Nlattr)(unsafe.Pointer(nla)), data) // returns (int, syscall.Errno)
	if err != 0 {
		// C.errno = int(err) // ``cannot refer to errno directly; see documentation''
		C.SetErrno(C.int(err))
	}
	return C.int(ret)
}

/**
 * mnl_attr_parse - parse attributes
 *
 * int
 * mnl_attr_parse(const struct nlmsghdr *nlh, unsigned int offset,
 *	          mnl_attr_cb_t cb, void *data)
 */
func AttrParse(nlh *Nlmsghdr, offset Size_t, cb MnlAttrCb, data interface{}) (int, error) {
	args := [2]unsafe.Pointer{unsafe.Pointer(&cb), unsafe.Pointer(&data)}
	ret, err := C.attr_parse_wrapper((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(offset), unsafe.Pointer(&args))
	return int(ret), err
}

/**
 * mnl_attr_parse_nested - parse attributes inside a nest
 *
 * int
 * mnl_attr_parse_nested(const struct nlattr *nested, mnl_attr_cb_t cb,
 * 			 void *data)
 */
func AttrParseNested(nested *Nlattr, cb MnlAttrCb, data interface{}) (int, error) {
	args := [2]unsafe.Pointer{unsafe.Pointer(&cb), unsafe.Pointer(&data)}
	ret, err := C.attr_parse_nested_wrapper((*C.struct_nlattr)(unsafe.Pointer(nested)), unsafe.Pointer(&args))
	return int(ret), err
}
		
/**
 * mnl_attr_parse_payload - parse attributes in payload of Netlink message
 *
 * int
 * mnl_attr_parse_payload(const void *payload, size_t payload_len,
 *			  mnl_attr_cb_t cb, void *data)
 */
func AttrParsePayload(payload []byte, cb MnlAttrCb, data interface{}) (int, error) {
	args := [2]unsafe.Pointer{unsafe.Pointer(&cb), unsafe.Pointer(&data)}
	ret, err := C.attr_parse_payload_wrapper(unsafe.Pointer(&payload[0]), C.size_t(len(payload)), unsafe.Pointer(&args))
	return int(ret), err
}

/**
 * mnl_attr_get_u8 - returns 8-bit unsigned integer attribute payload
 *
 * uint8_t mnl_attr_get_u8(const struct nlattr *attr)
 */
func AttrGetU8(attr *Nlattr) uint8 {
	return uint8(C.mnl_attr_get_u8((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_get_u16 - returns 16-bit unsigned integer attribute payload
 *
 * uint16_t mnl_attr_get_u16(const struct nlattr *attr)
 */
func AttrGetU16(attr *Nlattr) uint16 {
	return uint16(C.mnl_attr_get_u16((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_get_u32 - returns 32-bit unsigned integer attribute payload
 *
 * uint32_t mnl_attr_get_u32(const struct nlattr *attr)
 */
func AttrGetU32(attr *Nlattr) uint32 {
	return uint32(C.mnl_attr_get_u32((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_get_u64 - returns 64-bit unsigned integer attribute.
 *
 * uint64_t mnl_attr_get_u64(const struct nlattr *attr)
 */
func AttrGetU64(attr *Nlattr) uint64 {
	return uint64(C.mnl_attr_get_u64((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_get_str - returns pointer to string attribute.
 *
 * const char *mnl_attr_get_str(const struct nlattr *attr)
 */
func AttrGetStr(attr *Nlattr) string {
	return C.GoString(C.mnl_attr_get_str((*C.struct_nlattr)(unsafe.Pointer(attr))))
}

/**
 * mnl_attr_put - add an attribute to netlink message
 *
 * void
 * mnl_attr_put(struct nlmsghdr *nlh, uint16_t type, size_t len, const void *data)
 */
func AttrPut(nlh *Nlmsghdr, attr_type uint16, size Size_t, data unsafe.Pointer) {
	C.mnl_attr_put((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), C.uint16_t(attr_type), C.size_t(size), data)
}
func AttrPutData(nlh *Nlmsghdr, attr_type uint16, data interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		panic("pointer required for data")
	}
	t := reflect.Indirect(v).Type()
	AttrPut(nlh, attr_type, Size_t(t.Size()), unsafe.Pointer(v.Pointer()))
}
func AttrPutBytes(nlh *Nlmsghdr, attr_type uint16, data []byte) {
	C.mnl_attr_put((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.uint16_t(attr_type), C.size_t(len(data)), unsafe.Pointer(&data[0]))
}

/**
 * mnl_attr_put_u8 - add 8-bit unsigned integer attribute to netlink message
 *
 * void mnl_attr_put_u8(struct nlmsghdr *nlh, uint16_t type, uint8_t data)
 */
func AttrPutU8(nlh *Nlmsghdr, attr_type uint16, data uint8) {
	C.mnl_attr_put_u8((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.uint16_t(attr_type), C.uint8_t(data))
}

/**
 * mnl_attr_put_u16 - add 16-bit unsigned integer attribute to netlink message
 *
 * void mnl_attr_put_u16(struct nlmsghdr *nlh, uint16_t type, uint16_t data)
 */
func AttrPutU16(nlh *Nlmsghdr, attr_type uint16, data uint16) {
	C.mnl_attr_put_u16((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.uint16_t(attr_type), C.uint16_t(data))
}

/**
 * mnl_attr_put_u32 - add 32-bit unsigned integer attribute to netlink message
 *
 * void mnl_attr_put_u32(struct nlmsghdr *nlh, uint16_t type, uint32_t data)
 */
func AttrPutU32(nlh *Nlmsghdr, attr_type uint16, data uint32) {
	C.mnl_attr_put_u32((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.uint16_t(attr_type), C.uint32_t(data))
}

/**
 * mnl_attr_put_u64 - add 64-bit unsigned integer attribute to netlink message
 *
 * void mnl_attr_put_u64(struct nlmsghdr *nlh, uint16_t type, uint64_t data)
 */
func AttrPutU64(nlh *Nlmsghdr, attr_type uint16, data uint64) {
	C.mnl_attr_put_u64((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.uint16_t(attr_type), C.uint64_t(data))
}

/**
 * mnl_attr_put_str - add string attribute to netlink message
 *
 * void mnl_attr_put_str(struct nlmsghdr *nlh, uint16_t type, const char *data)
 */
func AttrPutStr(nlh *Nlmsghdr, attr_type uint16, data string) {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	C.mnl_attr_put_str((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.uint16_t(attr_type), cs)
}

/**
 * mnl_attr_put_strz - add string attribute to netlink message
 *
 * void mnl_attr_put_strz(struct nlmsghdr *nlh, uint16_t type, const char *data)
 */
func AttrPutStrz(nlh *Nlmsghdr, attr_type uint16, data string) {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	C.mnl_attr_put_strz((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.uint16_t(attr_type), cs)
}

/**
 * mnl_attr_nest_start - start an attribute nest
 *
 * struct nlattr *mnl_attr_nest_start(struct nlmsghdr *nlh, uint16_t type)
 */
func AttrNestStart(nlh *Nlmsghdr, attr_type uint16) *Nlattr {
	return (*Nlattr)(unsafe.Pointer(C.mnl_attr_nest_start((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), C.uint16_t(attr_type))))
}

/**
 * mnl_attr_put_check - add an attribute to netlink message
 *
 * bool
 * mnl_attr_put_check(struct nlmsghdr *nlh, size_t buflen,
 *		      uint16_t type, size_t len, const void *data)
 */
func AttrPutCheck(nlh *Nlmsghdr, buflen Size_t, attr_type uint16, data []byte) bool {
	return bool(C.mnl_attr_put_check((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(buflen), C.uint16_t(attr_type), C.size_t(len(data)), unsafe.Pointer(&data[0])))
}

/**
 * mnl_attr_put_u8_check - add 8-bit unsigned int attribute to netlink message
 *
 * bool
 * mnl_attr_put_u8_check(struct nlmsghdr *nlh, size_t buflen,
 * 			 uint16_t type, uint8_t data)
 */
func AttrPutU8Check(nlh *Nlmsghdr, buflen Size_t, attr_type uint16, data uint8) bool {
	return bool(C.mnl_attr_put_u8_check((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(buflen), C.uint16_t(attr_type), C.uint8_t(data)))
}

/**
 * mnl_attr_put_u16_check - add 16-bit unsigned int attribute to netlink message
 *
 * bool
 * mnl_attr_put_u16_check(struct nlmsghdr *nlh, size_t buflen,
 *			  uint16_t type, uint16_t data)
 */
func AttrPutU16Check(nlh *Nlmsghdr, buflen Size_t, attr_type uint16, data uint16) bool {
	return bool(C.mnl_attr_put_u16_check((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(buflen), C.uint16_t(attr_type), C.uint16_t(data)))
}

/**
 * mnl_attr_put_u32_check - add 32-bit unsigned int attribute to netlink message
 *
 * bool
 * mnl_attr_put_u32_check(struct nlmsghdr *nlh, size_t buflen,
 *			  uint16_t type, uint32_t data)
 */
func AttrPutU32Check(nlh *Nlmsghdr, buflen Size_t, attr_type uint16, data uint32) bool {
	return bool(C.mnl_attr_put_u32_check((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(buflen), C.uint16_t(attr_type), C.uint32_t(data)))
}

/**
 * mnl_attr_put_u64_check - add 64-bit unsigned int attribute to netlink message
 *
 * bool
 * mnl_attr_put_u64_check(struct nlmsghdr *nlh, size_t buflen,
 *			  uint16_t type, uint64_t data)
 */
func AttrPutU64Check(nlh *Nlmsghdr, buflen Size_t, attr_type uint16, data uint64) bool {
	return bool(C.mnl_attr_put_u64_check((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(buflen), C.uint16_t(attr_type), C.uint64_t(data)))
}

/**
 * mnl_attr_put_str_check - add string attribute to netlink message
 *
 * bool
 * mnl_attr_put_str_check(struct nlmsghdr *nlh, size_t buflen,
 *			  uint16_t type, const char *data)
 */
func AttrPutStrCheck(nlh *Nlmsghdr, buflen Size_t, attr_type uint16, data string) bool {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	return bool(C.mnl_attr_put_str_check((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(buflen), C.uint16_t(attr_type), cs))
}

/**
 * mnl_attr_put_strz_check - add string attribute to netlink message
 *
 * bool
 * mnl_attr_put_strz_check(struct nlmsghdr *nlh, size_t buflen,
 *			   uint16_t type, const char *data)
 */
func AttrPutStrzCheck(nlh *Nlmsghdr, buflen Size_t, attr_type uint16, data string) bool {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))
	return bool(C.mnl_attr_put_strz_check((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)),
		C.size_t(buflen), C.uint16_t(attr_type), cs))
}

/**
 * mnl_attr_nest_end - end an attribute nest
 *
 * void
 * mnl_attr_nest_end(struct nlmsghdr *nlh, struct nlattr *start)
 */
func AttrNestEnd(nlh *Nlmsghdr, start *Nlattr) {
	C.mnl_attr_nest_end((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), (*C.struct_nlattr)(unsafe.Pointer(start)))
}

/**
 * mnl_attr_nest_cancel - cancel an attribute nest
 *
 * void
 * mnl_attr_nest_cancel(struct nlmsghdr *nlh, struct nlattr *start)
 */
func AttrNestCancel(nlh *Nlmsghdr, start *Nlattr) {
	C.mnl_attr_nest_cancel((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), (*C.struct_nlattr)(unsafe.Pointer(start)))
}
