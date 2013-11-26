// +build ignore
package main

/*
#include <linux/rtnetlink.h>
*/
import "C"

// rtnl-link-dump
const SizeofIfinfomsg	= C.sizeof_struct_ifinfomsg
type Ifinfomsg		C.struct_ifinfomsg
const SizeofRtgenmsg	= C.sizeof_struct_rtgenmsg
type Rtgenmsg		C.struct_rtgenmsg

// rtnl-addr-dump
const SizeofIfaddrmsg	= C.sizeof_struct_ifaddrmsg
type Ifaddrmsg		C.struct_ifaddrmsg
