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
	"strconv"
	"time"
	mnl "cgolmnl"
	. "cgolmnl/inet"
)

func main() {
	if len(os.Args) <= 3 {
		fmt.Printf("Usage: %s iface destination cidr [gateway]\n", os.Args[0]);
		fmt.Printf("Example: %s eth0 10.0.1.12 32 10.0.1.11\n", os.Args[0]);
		fmt.Printf("	 %s eth0 ffff::10.0.1.12 128 fdff::1\n", os.Args[0]);
		os.Exit(C.EXIT_FAILURE)
	}

	iface, err := IfNametoindex(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "if_nametoindex: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}

	var dst, gw IPAddr
	dst.UnmarshalText(([]byte)(os.Args[2]))
	prefix, _ := strconv.Atoi(os.Args[3])
	if len(os.Args) == 5 {
		gw.UnmarshalText(([]byte)(os.Args[4]))
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
	if len(dst) == 4 {
		rtm.Family = C.AF_INET
	} else if len(dst) == 6 {
		rtm.Family = C.AF_INET6
	}
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

	mb, _ := dst.MarshalBinary()
	nlh.PutBytes(C.RTA_DST, mb)
	nlh.PutU32(C.RTA_OIF, uint32(iface))
	if len(os.Args) == 5 {
		mb, _ = gw.MarshalBinary()
		nlh.PutBytes(C.RTA_GATEWAY, mb)
	}

	nl, err := mnl.SocketOpen(C.NETLINK_ROUTE)
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
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
