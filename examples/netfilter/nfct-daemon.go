package main

/*
#include <unistd.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <linux/netlink.h>

#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_conntrack.h>
*/
import "C"

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
	"unsafe"
	mnl "cgolmnl"
	"cgolmnl/inet"
)

type Nstats struct {
	Addr		net.IP
	Pkts, Bytes	uint64
}

var nstats_map = make(map[string]*Nstats)

func parse_counters_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.CTA_COUNTERS_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_COUNTERS_PACKETS: fallthrough
	case C.CTA_COUNTERS_BYTES:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U64); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func parse_counters(nest *mnl.Nlattr, ns *Nstats) {
	tb := make(map[uint16]*mnl.Nlattr, C.CTA_COUNTERS_MAX + 1)

	nest.ParseNested(parse_counters_cb, tb)
	if tb[C.CTA_COUNTERS_PACKETS] != nil {
		ns.Pkts += inet.Be64toh(tb[C.CTA_COUNTERS_PACKETS].U64())
	}
	if tb[C.CTA_COUNTERS_BYTES] != nil {
		ns.Bytes += inet.Be64toh(tb[C.CTA_COUNTERS_BYTES].U64())
	}
}

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
	case C.CTA_IP_V6_SRC: fallthrough
	case C.CTA_IP_V6_DST:
		if ret, err := attr.Validate2(mnl.MNL_TYPE_BINARY, net.IPv6len); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func parse_ip(nest *mnl.Nlattr, ns *Nstats) {
	tb := make(map[uint16]*mnl.Nlattr, C.CTA_IP_MAX + 1)

	nest.ParseNested(parse_ip_cb, tb)
	if tb[C.CTA_IP_V4_SRC] != nil {
		ns.Addr = net.IP(tb[C.CTA_IP_V4_SRC].PayloadBytes())
	}
	if tb[C.CTA_IP_V6_SRC] != nil {
		ns.Addr = net.IP(tb[C.CTA_IP_V6_SRC].PayloadBytes())
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
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func parse_tuple(nest *mnl.Nlattr, ns *Nstats) {
	tb := make(map[uint16]*mnl.Nlattr, C.CTA_TUPLE_MAX + 1)

	nest.ParseNested(parse_tuple_cb, tb)
	if tb[C.CTA_TUPLE_IP] != nil {
		parse_ip(tb[C.CTA_TUPLE_IP], ns)
	}
}

func data_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.CTA_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTA_TUPLE_ORIG: fallthrough
	case C.CTA_COUNTERS_ORIG: fallthrough
	case C.CTA_COUNTERS_REPLY:
		if ret, err := attr.Validate(mnl.MNL_TYPE_NESTED); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr, C.CTA_MAX + 1)
	ns := &Nstats{}

	nlh.Parse(SizeofNfgenmsg, data_attr_cb, tb)
	if tb[C.CTA_TUPLE_ORIG] != nil {
		parse_tuple(tb[C.CTA_TUPLE_ORIG], ns)
	}

	if tb[C.CTA_COUNTERS_ORIG] != nil {
		parse_counters(tb[C.CTA_COUNTERS_ORIG], ns)
	}

	if tb[C.CTA_COUNTERS_REPLY] != nil {
		parse_counters(tb[C.CTA_COUNTERS_REPLY], ns)
	}

	cur := nstats_map[ns.Addr.String()]
	if cur == nil {
		cur = ns
		nstats_map[ns.Addr.String()] = ns
	}
	cur.Pkts += ns.Pkts
	cur.Bytes += ns.Bytes

	return mnl.MNL_CB_OK, 0
}

func handle(nl *mnl.Socket) int {
	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)

	nrcv, err := nl.Recvfrom(buf)
	if err != nil {
		// It only happens if NETLINK_NO_ENOBUFS is not set, it means
		// we are leaking statistics.
		if err == syscall.ENOBUFS {
			fmt.Fprintf(os.Stderr, "The daemon has hit ENOBUFS, you can " +
				"increase the size of your receiver " +
				"buffer to mitigate this or enable " +
				"reliable delivery.\n")
			// http://stackoverflow.com/questions/7933460/how-do-you-write-multiline-strings-in-go
			// `line1
			// line2
			// line3`
		} else {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
		}
		return -1
	}

	if ret, err := mnl.CbRun(buf[:nrcv], 0, 0, data_cb, nil); ret <= mnl.MNL_CB_ERROR {
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s", err)
		return -1
	} else if ret <= mnl.MNL_CB_STOP {
		return 0
	}

	return 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <poll-secs>\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}
	secs, _ := strconv.Atoi(os.Args[1])

	fmt.Printf("Polling every %d seconds from kernel...\n", secs)

	// Set high priority for this process, less chances to overrun
	// the netlink receiver buffer since the scheduler gives this process
	// more chances to run.
	C.nice(C.int(-20))

	// Open netlink socket to operate with netfilter
	nl, err := mnl.SocketOpen(C.NETLINK_NETFILTER)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	// Subscribe to destroy events to avoid leaking counters. The same
	// socket is used to periodically atomically dump and reset counters.
	if err := nl.Bind(C.NF_NETLINK_CONNTRACK_DESTROY, mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	// Set netlink receiver buffer to 16 MBytes, to avoid packet drops */
	buffersize := (1 << 22)
	C.setsockopt(C.int(nl.Fd()), C.SOL_SOCKET, C.SO_RCVBUFFORCE,
		unsafe.Pointer(&buffersize), SizeofSocklen_t)

	// The two tweaks below enable reliable event delivery, packets may
	// be dropped if the netlink receiver buffer overruns. This happens ...
	//
	// a) if the kernel spams this user-space process until the receiver
	//    is filled.
	//
	// or:
	//
	// b) if the user-space process does not pull messages from the
	//    receiver buffer so often.
	nl.SetsockoptCint(C.NETLINK_BROADCAST_ERROR, 1)
	nl.SetsockoptCint(C.NETLINK_NO_ENOBUFS, 1)

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	nlh, _ := mnl.NlmsgPutHeaderBytes(buf)
	// Counters are atomically zeroed in each dump
	nlh.Type = (C.NFNL_SUBSYS_CTNETLINK << 8) | C.IPCTNL_MSG_CT_GET_CTRZERO
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_DUMP

	nfh := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfh.Nfgen_family = C.AF_INET
	nfh.Version = C.NFNETLINK_V0
	nfh.Res_id = 0

	// Filter by mark: We only want to dump entries whose mark is zero
	nlh.PutU32(C.CTA_MARK, inet.Htonl(0))
	nlh.PutU32(C.CTA_MARK_MASK, inet.Htonl(0xffffffff))

	// prepare for epoll
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EpollCreate1: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer syscall.Close(epfd)

	var event syscall.EpollEvent
	events := make([]syscall.EpollEvent, 64) // XXX: magic number
	event.Events = syscall.EPOLLIN
	event.Fd = int32(nl.Fd())
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, int(event.Fd), &event); err != nil {
		fmt.Fprintf(os.Stderr, "EpollCtl: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}

	// Every N seconds ...
	ticker := time.NewTicker(time.Second * time.Duration(secs))
	go func() {
		for _ = range ticker.C {
			if _, err := nl.SendNlmsg(nlh);  err != nil {
				fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
				os.Exit(C.EXIT_FAILURE)
			}
		}
	}()

	for true {
		nevents, err := syscall.EpollWait(epfd, events, -1)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			fmt.Fprintf(os.Stderr, "EpollWait: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}

		// Handled event and periodic atomic-dump-and-reset messages
		for i := 0; i < nevents; i++ {
			if events[i].Fd != event.Fd {
				continue
			}
			if ret := handle(nl); ret < 0 {
				fmt.Fprintf(os.Stderr, "handle failed: %d\n", ret)
				os.Exit(C.EXIT_FAILURE)
			}

			// print the content of the list
			for k, v := range nstats_map {
				fmt.Printf("src=%s ", k)
				fmt.Printf("counters %d %d\n", v.Pkts, v.Bytes)
			}
		}
	}		
}
