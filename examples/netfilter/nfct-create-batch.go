package main

/*
#include <stdlib.h>
#include <arpa/inet.h>
#include <linux/netlink.h>

#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_conntrack.h>
#include <linux/netfilter/nf_conntrack_common.h>
#include <linux/netfilter/nf_conntrack_tcp.h>
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

func put_msg(nlh *mnl.Nlmsg, i uint16, seq uint32) {
	nlh.PutHeader()
	nlh.Type = (C.NFNL_SUBSYS_CTNETLINK << 8) | C.IPCTNL_MSG_CT_NEW
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_CREATE | C.NLM_F_EXCL | C.NLM_F_ACK
	nlh.Seq = seq

	nfh := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfh.Nfgen_family = C.AF_INET
	nfh.Version = C.NFNETLINK_V0
	nfh.Res_id = 0

	nest1 := nlh.NestStart(C.CTA_TUPLE_ORIG)
	nest2 := nlh.NestStart(C.CTA_TUPLE_IP)
	nlh.PutU32(C.CTA_IP_V4_SRC, inet.InetAddr("1.1.1.1"))
	nlh.PutU32(C.CTA_IP_V4_DST, inet.InetAddr("2.2.2.2"))
	nlh.NestEnd(nest2)

	nest2 = nlh.NestStart(C.CTA_TUPLE_PROTO)
	nlh.PutU8(C.CTA_PROTO_NUM, C.IPPROTO_TCP)
	nlh.PutU16(C.CTA_PROTO_SRC_PORT, inet.Htons(i))
	nlh.PutU16(C.CTA_PROTO_DST_PORT, inet.Htons(1025))
	nlh.NestEnd(nest2)
	nlh.NestEnd(nest1)

	nest1 = nlh.NestStart(C.CTA_TUPLE_REPLY)
	nest2 = nlh.NestStart(C.CTA_TUPLE_IP)
	nlh.PutU32(C.CTA_IP_V4_SRC, inet.InetAddr("2.2.2.2"))
	nlh.PutU32(C.CTA_IP_V4_DST, inet.InetAddr("1.1.1.1"))
	nlh.NestEnd(nest2)

	nest2 = nlh.NestStart(C.CTA_TUPLE_PROTO)
	nlh.PutU8(C.CTA_PROTO_NUM, C.IPPROTO_TCP)
	nlh.PutU16(C.CTA_PROTO_SRC_PORT, inet.Htons(1025))
	nlh.PutU16(C.CTA_PROTO_DST_PORT, inet.Htons(i))
	nlh.NestEnd(nest2)
	nlh.NestEnd(nest1)

	nest1 = nlh.NestStart(C.CTA_PROTOINFO)
	nest2 = nlh.NestStart(C.CTA_PROTOINFO_TCP)
	nlh.PutU8(C.CTA_PROTOINFO_TCP_STATE, C.TCP_CONNTRACK_SYN_SENT)
	nlh.NestEnd(nest2)
	nlh.NestEnd(nest1)

	nlh.PutU32(C.CTA_STATUS, inet.Htonl(C.IPS_CONFIRMED))
	nlh.PutU32(C.CTA_TIMEOUT, inet.Htonl(1000))
}

func send_batch(nl *mnl.Socket, b *mnl.NlmsgBatch, portid uint32) {
	var err error
	var epfd int
	var event syscall.EpollEvent
	var events [1]syscall.EpollEvent
	rcv_buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)

	if _, err := nl.Sendto(b.HeadBytes()); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}

	if epfd, err = syscall.EpollCreate1(0); err != nil {
		fmt.Fprintf(os.Stderr, "EpollCreate1: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer syscall.Close(epfd)

	event.Events = syscall.EPOLLIN
	event.Fd = int32(nl.Fd())
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, int(event.Fd), &event); err != nil {
		fmt.Fprintf(os.Stderr, "EpollCtl: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}

	ctl_cbs := []mnl.MnlCb{
		nil, nil,
		func(nlh *mnl.Nlmsg, data interface{}) (int, syscall.Errno) {
			err := (*mnl.Nlmsgerr)(nlh.Payload())
			if err.Error != 0 {
				errno := -err.Error
				fmt.Printf("mssage with seq %d has failed: %s\n", nlh.Seq, syscall.Errno(errno))
			}

			return mnl.MNL_CB_OK, 0
		}}
	ret := mnl.MNL_CB_OK
	for ret > 0 {
		var n int
		nevents, err := syscall.EpollWait(epfd, events[:], 0)
		if nevents == 0 {
			break
		}
		for n = 0; n < nevents; n++ {
			if events[n].Fd == event.Fd {
				break
			}
		}
		if n >= nevents {
			continue
		}

		nrecv, err := nl.Recvfrom(rcv_buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		_, err = mnl.CbRun2(rcv_buf[:nrecv], 0, portid, nil, nil, ctl_cbs)
	}
}

func main() {
	var err error
	var nl *mnl.Socket
	var b *mnl.NlmsgBatch

	if nl, err = mnl.NewSocket(C.NETLINK_NETFILTER); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer nl.Close()

	if err = nl.Bind(0, mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	portid := nl.Portid()

	snd_buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE*2, mnl.MNL_SOCKET_BUFFER_SIZE*2)
	if b, err = mnl.NewNlmsgBatch(snd_buf, mnl.Size_t(mnl.MNL_SOCKET_BUFFER_SIZE)); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_nlmsg_batch_start: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer b.Stop()

	seq := uint32(time.Now().Unix())
	for i := 1024; i < 65535; i++ {
		put_msg(b.CurrentNlmsg(), uint16(i), seq+uint32(i)-1024)
		if b.Next() {
			continue
		}

		send_batch(nl, b, portid)
		b.Reset()
	}

	if !b.IsEmpty() {
		send_batch(nl, b, portid)
	}
}
