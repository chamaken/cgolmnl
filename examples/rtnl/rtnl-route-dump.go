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

func data_attr_cb2(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)

	if ret, _ := attr.TypeValid(C.RTAX_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}
	
	if ret, err := attr.Validate(mnl.MNL_TYPE_U32); ret < 0 {
		fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
		return mnl.MNL_CB_ERROR, err.(syscall.Errno)
	}
	
	tb[attr.GetType()] = attr
	return mnl.MNL_CB_OK, 0
}

func attribute_show_ipv4(tb map[uint16]*mnl.Nlattr) {
	if tb[C.RTA_TABLE] != nil {
		fmt.Printf("table=%d ", tb[C.RTA_TABLE].U32())
	}
	if tb[C.RTA_DST] != nil {
		addr, _ := C.inet_ntoa(*(*C.struct_in_addr)(tb[C.RTA_DST].Payload()))
		fmt.Printf("dst=%s ", C.GoString(addr))
	}
	if tb[C.RTA_SRC] != nil {
		addr, _ := C.inet_ntoa(*(*C.struct_in_addr)(tb[C.RTA_DST].Payload()))
		fmt.Printf("src=%s ", C.GoString(addr))
	}
	if tb[C.RTA_OIF] != nil {
		fmt.Printf("oif=%d ", tb[C.RTA_OIF].U32())
	}
	if tb[C.RTA_FLOW] != nil {
		fmt.Printf("flow=%d ", tb[C.RTA_FLOW].U32())
	}
	if tb[C.RTA_PREFSRC] != nil {
		addr, _ := C.inet_ntoa(*(*C.struct_in_addr)(tb[C.RTA_PREFSRC].Payload()))
		fmt.Printf("prefsrc=%s ", C.GoString(addr))
	}
	if tb[C.RTA_GATEWAY] != nil {
		addr, _ := C.inet_ntoa(*(*C.struct_in_addr)(tb[C.RTA_GATEWAY].Payload()))
		fmt.Printf("gw=%s ", C.GoString(addr))
	}
	if tb[C.RTA_PRIORITY] != nil {
		fmt.Printf("prio=%d ", tb[C.RTA_PRIORITY].U32())
	}
	if tb[C.RTA_METRICS] != nil {
		tbx := make([]*mnl.Nlattr, C.RTAX_MAX + 1)
		tb[C.RTA_METRICS].ParseNested(data_attr_cb2, tbx)

		for i := 0; i < C.RTAX_MAX; i++ {
			if tbx[i] != nil {
				fmt.Printf("metrics[%d]=%u ", i, tbx[i].U32())
			}
		}
	}
	fmt.Println()
}

func inet6_ntoa(in6 C.struct_in6_addr) string {
	buf := make([]byte, C.INET6_ADDRSTRLEN)
	return C.GoString(C.inet_ntop(C.AF_INET6, unsafe.Pointer(&in6.__in6_u),
		          (*C.char)(unsafe.Pointer(&buf[0])), C.socklen_t(len(buf))))
}

func attributes_show_ipv6(tb []*mnl.Nlattr) {
	if tb[C.RTA_TABLE] != nil {
		fmt.Printf("table=%d ", tb[C.RTA_TABLE].U32())
	}
	if tb[C.RTA_DST] != nil {
		fmt.Printf("dst=%s ", inet6_ntoa(*(*C.struct_in6_addr)(tb[C.RTA_DST].Payload())))
	}
	if tb[C.RTA_SRC] != nil {
		fmt.Printf("src=%s ", inet6_ntoa(*(*C.struct_in6_addr)(tb[C.RTA_SRC].Payload())))
	}
	if tb[C.RTA_OIF] != nil {
		fmt.Printf("oif=%u ", tb[C.RTA_OIF].U32())
	}
	if tb[C.RTA_FLOW] != nil {
		fmt.Printf("flow=%u ", tb[C.RTA_FLOW].U32())
	}
	if tb[C.RTA_PREFSRC] != nil {
		fmt.Printf("prefsrc=%s ", inet6_ntoa(*(*C.struct_in6_addr)(tb[C.RTA_PREFSRC].Payload())))
	}
	if tb[C.RTA_GATEWAY] != nil {
		fmt.Printf("gw=%s ", inet6_ntoa(*(*C.struct_in6_addr)(tb[C.RTA_GATEWAY].Payload())))
	}
	if tb[C.RTA_METRICS] != nil {
		tbx := make([]*mnl.Nlattr, C.RTAX_MAX + 1)
		tb[C.RTA_METRICS].ParseNested(data_attr_cb2, tbx)

		for i := 0; i < C.RTA_MAX; i++ {
			if tbx[i] != nil {
				fmt.Printf("metrics[%d]=%d ", i, tbx[i].U32())
			}
		}
	}
	fmt.Println()
}

func data_ipv4_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.RTA_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.RTA_TABLE:	fallthrough
	case C.RTA_DST:		fallthrough
	case C.RTA_SRC:		fallthrough
	case C.RTA_OIF:		fallthrough
	case C.RTA_FLOW:	fallthrough
	case C.RTA_PREFSRC:	fallthrough
	case C.RTA_GATEWAY:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U32); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.RTA_METRICS:
		if ret, err := attr.Validate(mnl.MNL_TYPE_NESTED); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_ipv6_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.RTA_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.RTA_TABLE:	fallthrough
	case C.RTA_OIF:		fallthrough
	case C.RTA_FLOW:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U32); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.RTA_DST:
	case C.RTA_SRC:
	case C.RTA_PREFSRC:
	case C.RTA_GATEWAY:
		if ret, err := attr.Validate2(mnl.MNL_TYPE_BINARY, SizeofIn6Addr); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2")
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.RTA_METRICS:
		if ret, err := attr.Validate(mnl.MNL_TYPE_NESTED); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate")
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr, C.RTA_MAX + 1)
	rm := (*Rtmsg)(nlh.Payload())

	fmt.Printf("family=%d ", rm.Family)
	fmt.Printf("dst_len=%d ", rm.Dst_len)
	fmt.Printf("src_len=%d ", rm.Src_len)
	fmt.Printf("tos=%d ", rm.Tos)
	fmt.Printf("table=%d ", rm.Table)
	fmt.Printf("type=%d ", rm.Type)
	fmt.Printf("scope=%d ", rm.Scope)
	fmt.Printf("proto=%d ", rm.Protocol)
	fmt.Printf("flags=%d ", rm.Flags)
	switch rm.Family {
	case C.AF_INET:
		nlh.Parse(SizeofRtmsg, data_ipv4_attr_cb, tb)
		attribute_show_ipv4(tb)
	case C.AF_INET6:
		nlh.Parse(SizeofRtmsg, data_ipv6_attr_cb, tb)
	}

	return mnl.MNL_CB_OK, 0
}

func main() {
	var nl *mnl.SocketDescriptor
	var err error
	rcv_buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <inet|inet6>\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}

	nlh, err := mnl.PutNewNlmsghdr(mnl.MNL_SOCKET_BUFFER_SIZE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewNlmsghdr: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	nlh.Type = C.RTM_GETROUTE
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_DUMP
	seq := uint32(time.Now().Unix())
	nlh.Seq = seq
	rtm := (*Rtgenmsg)(nlh.PutExtraHeader(SizeofRtgenmsg))
	if os.Args[1] == "inet" {
		rtm.Family = C.AF_INET
	} else if os.Args[1] == "inet6" {
		rtm.Family = C.AF_INET6
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

	if _, err = nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mln_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	ret := mnl.MNL_CB_OK
	for ret > mnl.MNL_CB_STOP {
		rsize, err := nl.Recvfrom(rcv_buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		if ret, err = mnl.CbRun(rcv_buf[:rsize], seq, portid, data_cb, nil); ret < 0 {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(C.EXIT_FAILURE)
		}
	}
}
