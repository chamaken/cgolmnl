// +build ignore

package cgolmnl

// cgo CFLAGS: does not work?

/*
#cgo CFLAGS: -I./include
#include <stdbool.h>
#include <stdio.h>
#include <stdint.h>
#include <unistd.h>
#include <sys/socket.h>
#include <linux/netlink.h>
*/
import "C"

type (
	Size_t    C.size_t
	Pid_t     C.pid_t
	Ssize_t   C.ssize_t
	Socklen_t C.socklen_t
)

type Nlmsghdr C.struct_nlmsghdr
const SizeofNlmsg = C.sizeof_struct_nlmsghdr

type Nlmsgerr C.struct_nlmsgerr
const SizeofNlmsgerr = C.sizeof_struct_nlmsgerr

type NlPktinfo C.struct_nl_pktinfo
const SizeofNlPktinfo = C.sizeof_struct_nl_pktinfo

type nlattr C.struct_nlattr
const SizeofNlattr = C.sizeof_struct_nlattr
