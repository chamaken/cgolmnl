package main

/*
#include <stdlib.h>
#include <arpa/inet.h>
#include <linux/if.h>
#include <linux/if_link.h>
#include <linux/rtnetlink.h>
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
	mnl "cgolmnl"
)

func data_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.IFA_MAX); ret < 0 {
		return mnl.MNL_CB_OK, syscall.Errno(0)
	}

	switch attr_type {
	case C.IFA_ADDRESS:
		if ret, err := attr.Validate(mnl.MNL_TYPE_BINARY); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, syscall.Errno(0)
}

func data_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr, C.IFLA_MAX + 1)
	ifa := (*Ifaddrmsg)(nlh.Payload())

	fmt.Printf("index=%d family=%d ", ifa.Index, ifa.Family)

	nlh.Parse(SizeofIfaddrmsg, data_attr_cb, tb)
	fmt.Printf("addr=")
	if tb[C.IFA_ADDRESS] != nil {
		addr := *(*[C.INET6_ADDRSTRLEN]byte)(tb[C.IFA_ADDRESS].Payload())
		out := make([]byte, C.INET6_ADDRSTRLEN)
		if s, err := C.inet_ntop(C.int(ifa.Family), unsafe.Pointer(&addr[0]), (*C.char)(unsafe.Pointer(&out[0])), C.socklen_t(len(out))); err != nil {
			fmt.Fprintf(os.Stderr, "C.inet_ntop: %s\n", err)
		} else {
			// fmt.Printf("%#v ", out)
			fmt.Printf("%s ", C.GoString(s))
		}
	}
	fmt.Printf("scope=")
	switch ifa.Scope {
	case 0:
		fmt.Printf("global ")
	case 200:
		fmt.Printf("site ")
	case 253:
		fmt.Printf("link ")
	case 254:
		fmt.Printf("host ")
	case 255:
		fmt.Printf("nowhere ")
	default:
		fmt.Printf("%d ", ifa.Scope)
	}

	fmt.Printf("\n")
	return mnl.MNL_CB_OK, syscall.Errno(0)
}

func main() {
	var nl *mnl.SocketDescriptor
	var err error
	var snd_buf []byte
	rcv_buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <inet|inet6>\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}

	seq := uint32(time.Now().Unix())
	nlh, err := mnl.NewNlmsghdr(mnl.MNL_SOCKET_BUFFER_SIZE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewNlmsghdr: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	nlh.PutHeader()
	nlh.Type = C.RTM_GETADDR
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_DUMP
	nlh.Seq = seq
	rt := (*Rtgenmsg)(nlh.PutExtraHeader(SizeofRtgenmsg))
	if os.Args[1] == "inet" {
		rt.Family = C.AF_INET
	} else if os.Args[1] == "inet6" {
		rt.Family = C.AF_INET6
	}

	if nl, err = mnl.SocketOpen(C.NETLINK_ROUTE); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer nl.Close()

	if err = nl.Bind(0, mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	portid := nl.Portid()

	if snd_buf, err = nlh.MarshalBinary(); err != nil {
		fmt.Fprintf(os.Stderr, "nlh.MarshalBinary: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	if err = nl.Sendto(snd_buf); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_soket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	ret := mnl.MNL_CB_OK
	for ret > mnl.MNL_CB_STOP {
		rsize, err := nl.Recvfrom(rcv_buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		if ret = mnl.CbRun(rcv_buf[:rsize], seq, portid, data_cb, nil); ret < 0 {
			fmt.Fprintf(os.Stderr, "error")
			os.Exit(C.EXIT_FAILURE)
		}
	}
}