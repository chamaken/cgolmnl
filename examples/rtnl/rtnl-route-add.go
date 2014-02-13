package main

/*
#include <stdlib.h>
#include <sys/socket.h>
#include <linux/rtnetlink.h>
*/
import "C"

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	mnl "cgolmnl"
	"cgolmnl/inet"
)

func main() {
	if len(os.Args) <= 3 {
		fmt.Printf("Usage: %s iface destination cidr [gateway]\n", os.Args[0]);
		fmt.Printf("Example: %s eth0 10.0.1.12 32 10.0.1.11\n", os.Args[0]);
		fmt.Printf("	 %s eth0 ffff::10.0.1.12 128 fdff::1\n", os.Args[0]);
		os.Exit(C.EXIT_FAILURE)
	}

	iface, err := inet.IfNametoindex(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "if_nametoindex: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	family := C.AF_INET
	var dst, gw net.IP

	dst = net.ParseIP(os.Args[2])
	if dst == nil {
		fmt.Fprintf(os.Stderr, "ParseIP - invalid dst IP: %s\n", os.Args[2])
		os.Exit(C.EXIT_FAILURE)
	}
	if dst.To4() == nil {
		family = C.AF_INET6
	}

	prefix, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Atoi - invalid mask num: %s\n", os.Args[3])
		os.Exit(C.EXIT_FAILURE)
	}

	if len(os.Args) == 5 {
		gw = net.ParseIP(os.Args[4])
		if gw == nil || (gw.To4() == nil && family == C.AF_INET) {
			fmt.Fprintf(os.Stderr, "ParseIP - invalid gw IP: %s\n", os.Args[5])
			os.Exit(C.EXIT_FAILURE)
		}
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	nlh, err := mnl.NlmsgPutHeaderBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	nlh.Type = C.RTM_NEWROUTE
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_CREATE | C.NLM_F_ACK
	seq := uint32(time.Now().Unix())
	nlh.Seq = seq

	rtm := (*Rtmsg)(nlh.PutExtraHeader(SizeofRtmsg))
	rtm.Family = uint8(family)
	rtm.Dst_len = uint8(prefix)
	rtm.Src_len = 0
	rtm.Tos = 0
	rtm.Protocol = C.RTPROT_STATIC
	rtm.Table = C.RT_TABLE_MAIN
	rtm.Type = C.RTN_UNICAST
	if len(os.Args) == 4 {
		rtm.Scope = C.RT_SCOPE_LINK
	} else {
		rtm.Scope = C.RT_SCOPE_UNIVERSE
	}
	rtm.Flags = 0

	var binaddr []byte
	if family == C.AF_INET {
		binaddr = ([]byte)(dst.To4())
	} else {
		binaddr = ([]byte)(dst.To16())
	}
	nlh.PutBytes(C.RTA_DST, binaddr)
	nlh.PutU32(C.RTA_OIF, uint32(iface))
	if len(os.Args) == 5 {
		if family == C.AF_INET {
			binaddr = ([]byte)(gw.To4())
		} else {
			binaddr = ([]byte)(gw.To16())
		}
		nlh.PutBytes(C.RTA_GATEWAY, binaddr)
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
		os.Exit(C.EXIT_FAILURE)
	}

	nrcv, err := nl.Recvfrom(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	ret, err := mnl.CbRun(buf[:nrcv], seq, portid, nil, nil)
	if ret < mnl.MNL_CB_STOP {
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
