package main

/*
#include <stdlib.h>
#include <sys/socket.h>
#include <linux/rtnetlink.h>
*/
import "C"

import (
	"fmt"
	mnl "github.com/chamaken/cgolmnl"
	inet "github.com/chamaken/cgolmnl/inet"
	"os"
	"syscall"
	"time"
)

func data_attr_cb2(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)

	if err := attr.TypeValid(C.RTAX_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
		return mnl.MNL_CB_ERROR, err.(syscall.Errno)
	}

	tb[attr.GetType()] = attr
	return mnl.MNL_CB_OK, 0
}

func attributes_show_ipv4(tb map[uint16]*mnl.Nlattr) {
	if tb[C.RTA_TABLE] != nil {
		fmt.Printf("table=%d ", tb[C.RTA_TABLE].U32())
	}
	if tb[C.RTA_DST] != nil {
		addr := inet.InetNtoa(tb[C.RTA_DST].Payload())
		fmt.Printf("dst=%s ", addr)
	}
	if tb[C.RTA_SRC] != nil {
		addr := inet.InetNtoa(tb[C.RTA_DST].Payload())
		fmt.Printf("src=%s ", addr)
	}
	if tb[C.RTA_OIF] != nil {
		fmt.Printf("oif=%d ", tb[C.RTA_OIF].U32())
	}
	if tb[C.RTA_FLOW] != nil {
		fmt.Printf("flow=%d ", tb[C.RTA_FLOW].U32())
	}
	if tb[C.RTA_PREFSRC] != nil {
		addr := inet.InetNtoa(tb[C.RTA_PREFSRC].Payload())
		fmt.Printf("prefsrc=%s ", addr)
	}
	if tb[C.RTA_GATEWAY] != nil {
		addr := inet.InetNtoa(tb[C.RTA_GATEWAY].Payload())
		fmt.Printf("gw=%s ", addr)
	}
	if tb[C.RTA_PRIORITY] != nil {
		fmt.Printf("prio=%d ", tb[C.RTA_PRIORITY].U32())
	}
	if tb[C.RTA_METRICS] != nil {
		tbx := make([]*mnl.Nlattr, C.RTAX_MAX+1)
		tb[C.RTA_METRICS].ParseNested(data_attr_cb2, tbx)

		for i := 0; i < C.RTAX_MAX; i++ {
			if tbx[i] != nil {
				fmt.Printf("metrics[%d]=%u ", i, tbx[i].U32())
			}
		}
	}
	fmt.Println()
}

func attributes_show_ipv6(tb map[uint16]*mnl.Nlattr) {
	if tb[C.RTA_TABLE] != nil {
		fmt.Printf("table=%d ", tb[C.RTA_TABLE].U32())
	}
	if tb[C.RTA_DST] != nil {
		fmt.Printf("dst=%s ", inet.Inet6Ntoa(tb[C.RTA_DST].Payload()))
	}
	if tb[C.RTA_SRC] != nil {
		fmt.Printf("src=%s ", inet.Inet6Ntoa(tb[C.RTA_SRC].Payload()))
	}
	if tb[C.RTA_OIF] != nil {
		fmt.Printf("oif=%u ", tb[C.RTA_OIF].U32())
	}
	if tb[C.RTA_FLOW] != nil {
		fmt.Printf("flow=%u ", tb[C.RTA_FLOW].U32())
	}
	if tb[C.RTA_PREFSRC] != nil {
		fmt.Printf("prefsrc=%s ", inet.Inet6Ntoa(tb[C.RTA_PREFSRC].Payload()))
	}
	if tb[C.RTA_GATEWAY] != nil {
		fmt.Printf("gw=%s ", inet.Inet6Ntoa(tb[C.RTA_GATEWAY].Payload()))
	}
	if tb[C.RTA_PRIORITY] != nil {
		fmt.Printf("prio=%u ", tb[C.RTA_PRIORITY].U32())
	}
	if tb[C.RTA_METRICS] != nil {
		tbx := make([]*mnl.Nlattr, C.RTAX_MAX+1)
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

	if err := attr.TypeValid(C.RTA_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.RTA_TABLE:
		fallthrough
	case C.RTA_DST:
		fallthrough
	case C.RTA_SRC:
		fallthrough
	case C.RTA_OIF:
		fallthrough
	case C.RTA_FLOW:
		fallthrough
	case C.RTA_PREFSRC:
		fallthrough
	case C.RTA_GATEWAY:
		fallthrough
	case C.RTA_PRIORITY:
		if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.RTA_METRICS:
		if err := attr.Validate(mnl.MNL_TYPE_NESTED); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_ipv6_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.RTA_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.RTA_TABLE:
		fallthrough
	case C.RTA_OIF:
		fallthrough
	case C.RTA_FLOW:
	case C.RTA_PRIORITY:
		if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.RTA_DST:
	case C.RTA_SRC:
	case C.RTA_PREFSRC:
	case C.RTA_GATEWAY:
		if err := attr.Validate2(mnl.MNL_TYPE_BINARY, SizeofIn6Addr); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2")
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.RTA_METRICS:
		if err := attr.Validate(mnl.MNL_TYPE_NESTED); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate")
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_cb(nlh *mnl.Nlmsg, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr, C.RTA_MAX+1)
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
		attributes_show_ipv4(tb)
	case C.AF_INET6:
		nlh.Parse(SizeofRtmsg, data_ipv6_attr_cb, tb)
		attributes_show_ipv6(tb)
	}

	return mnl.MNL_CB_OK, 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <inet|inet6>\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	nlh, err := mnl.NlmsgPutHeaderBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
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

	nl, err := mnl.NewSocket(C.NETLINK_ROUTE)
	if err != nil {
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
		nrcv, err := nl.Recvfrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		ret, err = mnl.CbRun(buf[:nrcv], seq, portid, data_cb, nil)
	}
	if ret < mnl.MNL_CB_STOP {
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
