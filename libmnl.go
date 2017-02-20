package cgolmnl

/*
// assume C.int as uint32 in nlmsghdr.nlm_len context
// to declare constants, macros

#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <linux/netlink.h>
#include <libmnl/libmnl.h>
*/
import "C"

import (
	"os"
)

const MNL_SOCKET_AUTOPID = C.MNL_SOCKET_AUTOPID

var MNL_SOCKET_BUFFER_SIZE int

const MNL_ALIGNTO = uint32(C.MNL_ALIGNTO)

func MnlAlign(mnl_len uint32) uint32 {
	return (mnl_len + MNL_ALIGNTO - 1) & ^(MNL_ALIGNTO - 1)
}

var MNL_NLMSG_HDRLEN = MnlAlign(SizeofNlmsg)

func init() {
	pagesize := os.Getpagesize()
	if pagesize < 8192 {
		MNL_SOCKET_BUFFER_SIZE = pagesize
	} else {
		MNL_SOCKET_BUFFER_SIZE = 8192
	}
}

var MNL_ATTR_HDRLEN = MnlAlign(SizeofNlattr)

type AttrDataType C.enum_mnl_attr_data_type

const (
	MNL_TYPE_UNSPEC        AttrDataType = C.MNL_TYPE_UNSPEC
	MNL_TYPE_U8            AttrDataType = C.MNL_TYPE_U8
	MNL_TYPE_U16           AttrDataType = C.MNL_TYPE_U16
	MNL_TYPE_U32           AttrDataType = C.MNL_TYPE_U32
	MNL_TYPE_U64           AttrDataType = C.MNL_TYPE_U64
	MNL_TYPE_STRING        AttrDataType = C.MNL_TYPE_STRING
	MNL_TYPE_FLAG          AttrDataType = C.MNL_TYPE_FLAG
	MNL_TYPE_MSECS         AttrDataType = C.MNL_TYPE_MSECS
	MNL_TYPE_NESTED        AttrDataType = C.MNL_TYPE_NESTED
	MNL_TYPE_NESTED_COMPAT AttrDataType = C.MNL_TYPE_NESTED_COMPAT
	MNL_TYPE_NUL_STRING    AttrDataType = C.MNL_TYPE_NUL_STRING
	MNL_TYPE_BINARY        AttrDataType = C.MNL_TYPE_BINARY
	MNL_TYPE_MAX           AttrDataType = C.MNL_TYPE_MAX
)

const MNL_CB_ERROR = C.MNL_CB_ERROR
const MNL_CB_STOP = C.MNL_CB_STOP
const MNL_CB_OK = C.MNL_CB_OK
