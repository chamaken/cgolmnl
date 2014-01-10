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
	Size_t		C.size_t
	Pid_t		C.pid_t
	Ssize_t		C.ssize_t
	Socklen_t	C.socklen_t
)

const SizeofNlmsghdr	= C.sizeof_struct_nlmsghdr
type Nlmsghdr		  C.struct_nlmsghdr

const SizeofNlmsgerr	= C.sizeof_struct_nlmsgerr
type Nlmsgerr		  C.struct_nlmsgerr

const SizeofNlPktinfo	= C.sizeof_struct_nl_pktinfo
type NlPktinfo		  C.struct_nl_pktinfo

const SizeofNlMmapReq	= C.sizeof_struct_nl_mmap_req
type NlMmapReq		  C.struct_nl_mmap_req

const SizeofNlMmapHdr	= C.sizeof_struct_nl_mmap_hdr
type NlMmapHdr		  C.struct_nl_mmap_hdr

const SizeofNlattr	= C.sizeof_struct_nlattr
type Nlattr		  C.struct_nlattr
