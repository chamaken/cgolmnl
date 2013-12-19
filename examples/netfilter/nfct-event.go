package main

/*
#include <unistd.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <linux/if.h>
#include <linux/if_link.h>
#include <linux/rtnetlink.h>

#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_conntrack.h>
*/
import "C"

import (
	"fmt"
	"net"
	"os"
	"syscall"
	mnl "cgolmnl"
	. "cgolmnl/inet"
)

func parse_ip_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.CTA_IP_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_IP_V4_SRC: fallthrough
	case C.CTA_IP_V4_DST:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U32); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func print_ip(nest *mnl.Nlattr) {
	tb := make(map[uint16]*mnl.Nlattr)

	nest.ParseNested(parse_ip_cb, tb)
	if tb[C.CTA_IP_V4_SRC] != nil {
		fmt.Printf("src=%s ", net.IP(tb[C.CTA_IP_V4_SRC].PayloadBytes()))
	}
	if tb[C.CTA_IP_V4_DST] != nil {
		fmt.Printf("dst=%s ", net.IP(tb[C.CTA_IP_V4_DST].PayloadBytes()))
	}
}

func parse_proto_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.CTA_PROTO_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_PROTO_NUM:		fallthrough
	case C.CTA_PROTO_ICMP_TYPE:	fallthrough
	case C.CTA_PROTO_ICMP_CODE:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U8); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTA_PROTO_SRC_PORT:	fallthrough
	case C.CTA_PROTO_DST_PORT:	fallthrough
	case C.CTA_PROTO_ICMP_ID:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U16); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func print_proto(nest *mnl.Nlattr) {
	tb := make(map[uint16]*mnl.Nlattr)

	nest.ParseNested(parse_proto_cb, tb)
	if tb[C.CTA_PROTO_NUM] != nil {
		fmt.Printf("proto=%d ", tb[C.CTA_PROTO_NUM].U8())
	}
	if tb[C.CTA_PROTO_SRC_PORT] != nil {
		fmt.Printf("sport=%d ", Ntohs(tb[C.CTA_PROTO_SRC_PORT].U16()))
	}
	if tb[C.CTA_PROTO_DST_PORT] != nil {
		fmt.Printf("dport=%d ", Ntohs(tb[C.CTA_PROTO_SRC_PORT].U16()))
	}
	if tb[C.CTA_PROTO_ICMP_ID] != nil {
		fmt.Printf("id=%d ", Ntohs(tb[C.CTA_PROTO_ICMP_ID].U16()))
	}
	if tb[C.CTA_PROTO_ICMP_TYPE] != nil {
		fmt.Printf("type=%d ", tb[C.CTA_PROTO_ICMP_TYPE].U8())
	}
	if tb[C.CTA_PROTO_ICMP_CODE] != nil {
		fmt.Printf("code=%d ", tb[C.CTA_PROTO_ICMP_CODE].U8())
	}
}

func parse_tuple_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.CTA_TUPLE_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_TUPLE_IP:
		if ret, err := attr.Validate(mnl.MNL_TYPE_NESTED); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTA_TUPLE_PROTO:
		if ret, err := attr.Validate(mnl.MNL_TYPE_NESTED); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func print_tuple(nest *mnl.Nlattr) {
	tb := make(map[uint16]*mnl.Nlattr)

	nest.ParseNested(parse_tuple_cb, tb)
	if tb[C.CTA_TUPLE_IP] != nil {
		print_ip(tb[C.CTA_TUPLE_IP])
	}
	if tb[C.CTA_TUPLE_PROTO] != nil {
		print_proto(tb[C.CTA_TUPLE_PROTO])
	}
}

func data_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.CTA_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_TUPLE_ORIG:
		if ret, err := attr.Validate(mnl.MNL_TYPE_NESTED); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTA_TIMEOUT:	fallthrough
	case C.CTA_MARK:	fallthrough
	case C.CTA_SECMARK:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U32); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr)
	// nfg := (*Nfgenmsg)(nlh.Payload())

	switch nlh.Type & 0xFF {
	case C.IPCTNL_MSG_CT_NEW:
		if nlh.Flags & (C.NLM_F_CREATE|C.NLM_F_EXCL) != 0 {
			fmt.Printf("%9s ", "[NEW] ")
		} else {
			fmt.Printf("%9s ", "[UPDATE] ")
		}
	case C.IPCTNL_MSG_CT_DELETE:
		fmt.Printf("%9s ", "[DESTROY] ")
	}

	nlh.Parse(SizeofNfgenmsg, data_attr_cb, tb)
	if tb[C.CTA_TUPLE_ORIG] != nil {
		print_tuple(tb[C.CTA_TUPLE_ORIG])
	}
	if tb[C.CTA_MARK] != nil {
		fmt.Printf("mark=%d ", Ntohl(tb[C.CTA_MARK].U32()))
	}
	if tb[C.CTA_SECMARK] != nil {
		fmt.Printf("secmark=%d ", Ntohl(tb[C.CTA_SECMARK].U32()))
	}
	fmt.Println()
	return mnl.MNL_CB_OK, 0
}

func main() {
	nl, err := mnl.SocketOpen(C.NETLINK_NETFILTER)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer nl.Close()

	if err := nl.Bind(C.NF_NETLINK_CONNTRACK_NEW |
		C.NF_NETLINK_CONNTRACK_UPDATE |
		C.NF_NETLINK_CONNTRACK_DESTROY,
		mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	ret := mnl.MNL_CB_OK
	for ret >= mnl.MNL_CB_STOP {
		nrcv, err := nl.Recvfrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		ret, err = mnl.CbRun(buf[:nrcv], 0, 0, data_cb, nil)
	}

	if ret < mnl.MNL_CB_STOP {
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
