package main

/*
#include <stdlib.h>
#include <sys/socket.h>
#include <linux/if.h>
#include <linux/rtnetlink.h>
*/
import "C"

import (
	"fmt"
	"os"
	"strings"
	"time"
	mnl "cgolmnl"
)

func main() {
	var change, flags uint32

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s [ifname] [up|down]\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}

	if strings.ToLower(os.Args[2]) == "up" {
		change |= C.IFF_UP
		flags |= C.IFF_UP
	} else if strings.ToLower(os.Args[2]) == "down" {
		change |= C.IFF_UP
		flags &= ^uint32(C.IFF_UP)
	} else {
		fmt.Fprintf(os.Stderr, "%s is not `up' nor `down'\n", os.Args[2])
		os.Exit(C.EXIT_FAILURE)
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	nlh, err := mnl.NlmsgPutHeaderBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh.Type = C.RTM_NEWLINK
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_ACK
	seq := uint32(time.Now().Unix())
	nlh.Seq = seq
	ifm := (*Ifinfomsg)(nlh.PutExtraHeader(SizeofIfinfomsg))
	ifm.Change = change
	ifm.Flags = flags

	nlh.PutStr(C.IFLA_IFNAME, os.Args[1])

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

	nlh.Fprint(os.Stdout, SizeofIfinfomsg)

	if _, err = nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
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
