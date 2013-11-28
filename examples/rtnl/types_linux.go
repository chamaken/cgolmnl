// +build ignore
package main

/*
#include <linux/rtnetlink.h>
#include <netinet/in.h>
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

// rtnl-route-dump
const SizeofInAddr	= C.sizeof_struct_in_addr
type InAddr		C.struct_in_addr
const SizeofIn6Addr	= C.sizeof_struct_in6_addr
type In6Addr		C.struct_in6_addr
const SizeofRtmsg	= C.sizeof_struct_rtmsg
type Rtmsg		C.struct_rtmsg
