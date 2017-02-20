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

func data_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.IFA_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.IFA_ADDRESS:
		if err := attr.Validate(mnl.MNL_TYPE_BINARY); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_cb(nlh *mnl.Nlmsg, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr, C.IFLA_MAX+1)
	ifa := (*Ifaddrmsg)(nlh.Payload())

	fmt.Printf("index=%d family=%d ", ifa.Index, ifa.Family)

	nlh.Parse(SizeofIfaddrmsg, data_attr_cb, tb)
	fmt.Printf("addr=")
	if tb[C.IFA_ADDRESS] != nil {
		fmt.Printf("%s ", inet.InetNtop(int(ifa.Family), tb[C.IFA_ADDRESS].Payload()))
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
	nlh.Type = C.RTM_GETADDR
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_DUMP
	seq := uint32(time.Now().Unix())
	nlh.Seq = seq
	rt := (*Rtgenmsg)(nlh.PutExtraHeader(SizeofRtgenmsg))
	if os.Args[1] == "inet" {
		rt.Family = C.AF_INET
	} else if os.Args[1] == "inet6" {
		rt.Family = C.AF_INET6
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
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
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
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
