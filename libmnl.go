// libmnl is a minimalistic user-space library oriented to Netlink developers.
// There are a lot of common tasks in parsing, validating, constructing of
// both the Netlink header and TLVs that are repetitive and easy to get wrong.
// This library aims to provide simple helpers that allows you to avoid
// re-inventing the wheel in common Netlink tasks.
// 
//     "Simplify, simplify" -- Henry David Thoureau. Walden (1854)
// 
// The acronym libmnl stands for LIBrary Minimalistic NetLink.
// 
// libmnl homepage is:
//      http://www.netfilter.org/projects/libmnl/
// 
// Main Features
// - Small: the shared library requires around 30KB for an x86-based computer.
// - Simple: this library avoids complex abstractions that tend to hide Netlink
//   details. It avoids elaborated object-oriented infrastructure and complex
//   callback-based workflow.
// - Easy to use: the library simplifies the work for Netlink-wise developers.
//   It provides functions to make socket handling, message building,
//   validating, parsing and sequence tracking, easier.
// - Easy to re-use: you can use this library to build your own abstraction
//   layer upon this library, if you want to provide another library that
//   hides Netlink details to your users.
// - Decoupling: the interdependency of the main bricks that compose this
//   library is reduced, i.e. the library provides many helpers, but the
//   programmer is not forced to use them.
// 
// Licensing terms
//   This library is released under the LGPLv2.1 or any later (at your option).
// 
// Dependencies
//   You have to install the Linux kernel headers that you want to use to develop
//   your application. Moreover, this library requires that you have some basics
//   on Netlink.
// 
// Git Tree
//   The current development version of libmnl can be accessed at:
//   http://git.netfilter.org/cgi-bin/gitweb.cgi?p=libmnl.git;a=summary
// 
// Using libmnl
//   You can access several example files under examples/ in the libmnl source
//   code tree.
// 
package cgolmnl

/*
// assume C.int to uint32 to follow nlmsghdr.nlm_len
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

const MNL_SOCKET_AUTOPID	= C.MNL_SOCKET_AUTOPID
var MNL_SOCKET_BUFFER_SIZE int
const MNL_ALIGNTO		= uint32(C.MNL_ALIGNTO)

func MnlAlign(mnl_len uint32) uint32 {
	return (mnl_len + MNL_ALIGNTO - 1) & ^(MNL_ALIGNTO - 1)
}

var MNL_NLMSG_HDRLEN	= MnlAlign(SizeofNlmsghdr)

func init() {
	pagesize := os.Getpagesize()
	if pagesize < 8192 {
		MNL_SOCKET_BUFFER_SIZE = pagesize
	} else {
		MNL_SOCKET_BUFFER_SIZE = 8192
	}
}

var MNL_ATTR_HDRLEN	= MnlAlign(SizeofNlattr)

type AttrDataType C.enum_mnl_attr_data_type
const (
	MNL_TYPE_UNSPEC		AttrDataType = C.MNL_TYPE_UNSPEC
	MNL_TYPE_U8		AttrDataType = C.MNL_TYPE_U8
	MNL_TYPE_U16		AttrDataType = C.MNL_TYPE_U16
	MNL_TYPE_U32		AttrDataType = C.MNL_TYPE_U32
	MNL_TYPE_U64		AttrDataType = C.MNL_TYPE_U64
	MNL_TYPE_STRING		AttrDataType = C.MNL_TYPE_STRING
	MNL_TYPE_FLAG		AttrDataType = C.MNL_TYPE_FLAG
	MNL_TYPE_MSECS		AttrDataType = C.MNL_TYPE_MSECS
	MNL_TYPE_NESTED		AttrDataType = C.MNL_TYPE_NESTED
	MNL_TYPE_NESTED_COMPAT	AttrDataType = C.MNL_TYPE_NESTED_COMPAT
	MNL_TYPE_NUL_STRING	AttrDataType = C.MNL_TYPE_NUL_STRING
	MNL_TYPE_BINARY		AttrDataType = C.MNL_TYPE_BINARY
	MNL_TYPE_MAX		AttrDataType = C.MNL_TYPE_MAX
)

const MNL_CB_ERROR		= C.MNL_CB_ERROR
const MNL_CB_STOP		= C.MNL_CB_STOP
const MNL_CB_OK			= C.MNL_CB_OK
