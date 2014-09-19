package main

/*
#include <stdlib.h>
#include <sys/socket.h>
#include <linux/netlink.h>

#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_conntrack.h>
*/
import "C"

import (
	mnl "cgolmnl"
	"cgolmnl/inet"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func parse_counters_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.CTA_COUNTERS_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_COUNTERS_PACKETS:
		fallthrough
	case C.CTA_COUNTERS_BYTES:
		if err := attr.Validate(mnl.MNL_TYPE_U64); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func print_counters(nest *mnl.Nlattr) {
	tb := make(map[uint16]*mnl.Nlattr, C.CTA_COUNTERS_MAX+1)

	nest.ParseNested(parse_counters_cb, tb)
	if tb[C.CTA_COUNTERS_PACKETS] != nil {
		fmt.Printf("packets=%d ", inet.Be64toh(tb[C.CTA_COUNTERS_PACKETS].U64()))
	}
	if tb[C.CTA_COUNTERS_BYTES] != nil {
		fmt.Printf("bytes=%d ", inet.Be64toh(tb[C.CTA_COUNTERS_BYTES].U64()))
	}
}

func parse_ip_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.CTA_IP_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_IP_V4_SRC:
		fallthrough
	case C.CTA_IP_V4_DST:
		if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTA_IP_V6_SRC:
		fallthrough
	case C.CTA_IP_V6_DST:
		if err := attr.Validate2(mnl.MNL_TYPE_BINARY, net.IPv6len); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2: %s\n", err)
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
	if tb[C.CTA_IP_V6_SRC] != nil {
		fmt.Printf("src=%s ", net.IP(tb[C.CTA_IP_V6_SRC].PayloadBytes()))
	}
	if tb[C.CTA_IP_V6_DST] != nil {
		fmt.Printf("dst=%s ", net.IP(tb[C.CTA_IP_V6_DST].PayloadBytes()))
	}
}

func parse_proto_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.CTA_PROTO_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_PROTO_NUM:
		fallthrough
	case C.CTA_PROTO_ICMP_TYPE:
		fallthrough
	case C.CTA_PROTO_ICMP_CODE:
		if err := attr.Validate(mnl.MNL_TYPE_U8); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTA_PROTO_SRC_PORT:
		fallthrough
	case C.CTA_PROTO_DST_PORT:
		fallthrough
	case C.CTA_PROTO_ICMP_ID:
		if err := attr.Validate(mnl.MNL_TYPE_U16); err != nil {
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
		fmt.Printf("sport=%d ", inet.Ntohs(tb[C.CTA_PROTO_SRC_PORT].U16()))
	}
	if tb[C.CTA_PROTO_DST_PORT] != nil {
		fmt.Printf("dport=%d ", inet.Ntohs(tb[C.CTA_PROTO_SRC_PORT].U16()))
	}
	if tb[C.CTA_PROTO_ICMP_ID] != nil {
		fmt.Printf("id=%d ", inet.Ntohs(tb[C.CTA_PROTO_ICMP_ID].U16()))
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

	if err := attr.TypeValid(C.CTA_TUPLE_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_TUPLE_IP:
		if err := attr.Validate(mnl.MNL_TYPE_NESTED); err != nil {
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

	if err := attr.TypeValid(C.CTA_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_TUPLE_ORIG:
		fallthrough
	case C.CTA_COUNTERS_ORIG:
		fallthrough
	case C.CTA_COUNTERS_REPLY:
		if err := attr.Validate(mnl.MNL_TYPE_NESTED); err != nil {
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

	nlh.Parse(SizeofNfgenmsg, data_attr_cb, tb)
	if tb[C.CTA_TUPLE_ORIG] != nil {
		print_tuple(tb[C.CTA_TUPLE_ORIG])
	}

	if tb[C.CTA_MARK] != nil {
		fmt.Printf("mark=%d ", inet.Ntohl(tb[C.CTA_MARK].U32()))
	}

	if tb[C.CTA_SECMARK] != nil {
		fmt.Printf("secmark=%d ", inet.Ntohl(tb[C.CTA_SECMARK].U32()))
	}

	if tb[C.CTA_COUNTERS_ORIG] != nil {
		fmt.Printf("original ")
		print_counters(tb[C.CTA_COUNTERS_ORIG])
	}

	if tb[C.CTA_COUNTERS_REPLY] != nil {
		fmt.Printf("reply ")
		print_counters(tb[C.CTA_COUNTERS_REPLY])
	}

	fmt.Println()
	return mnl.MNL_CB_OK, 0
}

func main() {
	nl, err := mnl.NewSocket(C.NETLINK_NETFILTER)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer nl.Close()

	if err := nl.Bind(0, mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	nlh, _ := mnl.NlmsgPutHeaderBytes(buf)
	nlh.Type = (C.NFNL_SUBSYS_CTNETLINK << 8) | C.IPCTNL_MSG_CT_GET
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_DUMP
	seq := uint32(time.Now().Unix())
	nlh.Seq = seq

	nfh := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfh.Nfgen_family = C.AF_INET
	nfh.Version = C.NFNETLINK_V0
	nfh.Res_id = 0

	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	portid := nl.Portid()

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
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
