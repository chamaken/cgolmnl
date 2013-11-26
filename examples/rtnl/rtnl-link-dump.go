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

	if ret, _ := attr.TypeValid(C.IFLA_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch (attr_type) {
	case C.IFLA_ADDRESS:
		if ret, err := attr.Validate(mnl.MNL_TYPE_BINARY); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, 0
		}
	case C.IFLA_MTU:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U32); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, 0
		}
	case C.IFLA_IFNAME:
		if ret, err := attr.Validate(mnl.MNL_TYPE_STRING); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_valudate: %s\n", err)
			return mnl.MNL_CB_ERROR, 0
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr, C.IFLA_MAX + 1)

	ifm := (*Ifinfomsg)(nlh.Payload())
	fmt.Printf("index=%d type=%d flags=%d family=%d ", ifm.Index, ifm.Type, ifm.Flags, ifm.Family)

	if ifm.Flags & C.IFF_RUNNING == C.IFF_RUNNING {
		fmt.Printf("[RUNNING] ")
	} else {
		fmt.Printf("[NOT RUNNING] ")
	}

	nlh.Parse(SizeofIfinfomsg, data_attr_cb, tb)
	if tb[C.IFLA_MTU] != nil {
		fmt.Printf("mtu=%d ", tb[C.IFLA_MTU].U32())
	}
	if tb[C.IFLA_IFNAME] != nil {
		fmt.Printf("name=%s ", tb[C.IFLA_IFNAME].Str())
	}

	if tb[C.IFLA_ADDRESS] != nil {
		hwaddr := *(*[8]byte)(tb[C.IFLA_ADDRESS].Payload())
		// hwaddr := tb[C.IFLA_ADDRESS].PayloadBytes()
		fmt.Printf("hwaddr=")
		var i uint16
		for i = 0; i < tb[C.IFLA_ADDRESS].PayloadLen(); i++ {
			fmt.Printf("%.2x", hwaddr[i] & 0xff)
			if i + 1 != tb[C.IFLA_ADDRESS].PayloadLen() {
				fmt.Printf(":")
			}
		}
	}
	fmt.Printf("\n")
	return mnl.MNL_CB_OK, 0
}

func main() {
	var nl *mnl.SocketDescriptor
	var err error
	rcv_buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)

	seq := uint32(time.Now().Unix())
	nlh := mnl.NewNlmsghdr(mnl.MNL_SOCKET_BUFFER_SIZE)
	nlh.PutHeader()
	nlh.Type = C.RTM_GETLINK
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_DUMP
	nlh.Seq = seq
	rt := (*Rtgenmsg)(nlh.PutExtraHeader(SizeofRtgenmsg))
	rt.Family = C.AF_PACKET

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
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
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
