package main

/*
#include <stdlib.h>
#include <sys/socket.h>
#include <linux/netlink.h>

#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_log.h>
*/
import "C"

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	mnl "cgolmnl"
	"cgolmnl/inet"
)

func parse_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if ret, _ := attr.TypeValid(C.NFULA_MAX); ret < 0 {
		return mnl.MNL_CB_OK, 0
	}
	switch int(attr_type) {
	case C.NFULA_MARK: fallthrough
	case C.NFULA_IFINDEX_INDEV: fallthrough
	case C.NFULA_IFINDEX_OUTDEV: fallthrough
	case C.NFULA_IFINDEX_PHYSINDEV: fallthrough
	case C.NFULA_IFINDEX_PHYSOUTDEV:
		if ret, err := attr.Validate(mnl.MNL_TYPE_U32); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.NFULA_TIMESTAMP:
		if ret, err := attr.Validate2(mnl.MNL_TYPE_UNSPEC, SizeofNfulnlMsgPacketTimestamp); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.NFULA_HWADDR:
		if ret, err := attr.Validate2(mnl.MNL_TYPE_UNSPEC, SizeofNfulnlMsgPacketHw); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.NFULA_PREFIX:
		if ret, err := attr.Validate(mnl.MNL_TYPE_NUL_STRING); ret < 0 {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.NFULA_PAYLOAD:
		// do something
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func log_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	var ph *NfulnlMsgPacketHdr
	var prefix string
	var mark uint32
	tb := make(map[uint16]*mnl.Nlattr, C.NFULA_MAX + 1)

	nlh.Parse(SizeofNfgenmsg, parse_attr_cb, tb)
	if tb[C.NFULA_PACKET_HDR] != nil {
		ph = (*NfulnlMsgPacketHdr)(tb[C.NFULA_PACKET_HDR].Payload())
	}
	if tb[C.NFULA_PREFIX] != nil {
		prefix = tb[C.NFULA_PREFIX].Str()
	}
	if tb[C.NFULA_MARK] != nil {
		mark = inet.Ntohl(tb[C.NFULA_MARK].U32())
	}

	fmt.Printf("log received (prefix=\"%s\" hw=0x%04x hook=%d mark=%d)\n",
		prefix, inet.Ntohs(ph.Protocol), ph.Hook, mark)

	return mnl.MNL_CB_OK, 0
}

func nflog_build_cfg_pf_request(buf []byte, command uint8) *mnl.Nlmsghdr {
	nlh, err := mnl.NlmsgPutHeaderBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh.Type = (C.NFNL_SUBSYS_ULOG << 8) | C.NFULNL_MSG_CONFIG
	nlh.Flags = C.NLM_F_REQUEST

	nfg := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfg.Nfgen_family = C.AF_INET
	nfg.Version = C.NFNETLINK_V0

	cmd := &NfulnlMsgConfigCmd{Command: command}
	nlh.PutData(C.NFULA_CFG_CMD, cmd)

	return nlh
}

func nflog_build_cfg_request(buf []byte, command uint8, qnum int) *mnl.Nlmsghdr {
	nlh, err := mnl.NlmsgPutHeaderBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh.Type = (C.NFNL_SUBSYS_ULOG << 8) | C.NFULNL_MSG_CONFIG
	nlh.Flags = C.NLM_F_REQUEST

	nfg := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfg.Nfgen_family = C.AF_INET
	nfg.Version = C.NFNETLINK_V0
	nfg.Res_id = inet.Htons(uint16(qnum))

	cmd := &NfulnlMsgConfigCmd{Command: command}
	nlh.PutData(C.NFULA_CFG_CMD, cmd)

	return nlh
}

func nflog_build_cfg_params(buf []byte, copy_mode uint8, copy_range, qnum int) *mnl.Nlmsghdr {
	nlh, err := mnl.NlmsgPutHeaderBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh.Type = (C.NFNL_SUBSYS_ULOG << 8) | C.NFULNL_MSG_CONFIG
	nlh.Flags = C.NLM_F_REQUEST

	nfg := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfg.Nfgen_family = C.AF_UNSPEC
	nfg.Version = C.NFNETLINK_V0
	nfg.Res_id = inet.Htons(uint16(qnum))

	params := &NfulnlMsgConfigMode{	Range: inet.Htonl(uint32(copy_range)), Mode: copy_mode }
	nlh.PutData(C.NFULA_CFG_MODE, params)

	return nlh
}

func mnl_socket_poll(nl *mnl.Socket) int {
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EpollCreate1: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer syscall.Close(epfd)

	var event syscall.EpollEvent
	events := make([]syscall.EpollEvent, 1)
	event.Events = syscall.EPOLLIN | syscall.EPOLLERR
	event.Fd = int32(nl.Fd())
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, int(event.Fd), &event); err != nil {
		fmt.Fprintf(os.Stderr, "EpollCtl: %s", err)
		os.Exit(C.EXIT_FAILURE)
	}

	for true {
		nevents, err := syscall.EpollWait(epfd, events, -1)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			fmt.Fprintf(os.Stderr, "EpollWait: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		for i := 0; i < nevents; i++ {
			if events[i].Fd != event.Fd {
				continue
			}
			if events[i].Events & syscall.EPOLLIN == syscall.EPOLLIN {
				return 0
			}
			if events[i].Events & syscall.EPOLLERR == syscall.EPOLLERR {
				return -1
			}
		}
	}
	return -1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [queue_num]\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}
	qnum, _ := strconv.Atoi(os.Args[1])

	nl, err := mnl.SocketOpen(C.NETLINK_NETFILTER)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer nl.Close()

	req := &mnl.NlMmapReq{Block_size: 4096 * 16, Block_nr: 16, Frame_size: 2048, Frame_nr: 16 * 16 * 2}
	if ret, err := nl.SetRingopt(req, mnl.MNL_RING_RX); ret == -1 {
		fmt.Fprintf(os.Stderr, "mnl_socket_set_ringopt: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	} else {
		fmt.Fprintf(os.Stderr, "mnl_socket_set_ringopt: %d\n", ret)
	}

	if ret, err := nl.MapRing(); ret == - 1 {
		fmt.Fprintf(os.Stderr, "mnl_socket_map_ring: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	if err = nl.Bind(0, mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	portid := nl.Portid()

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)

	nlh := nflog_build_cfg_pf_request(buf, C.NFULNL_CFG_CMD_PF_UNBIND)
	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh = nflog_build_cfg_pf_request(buf, C.NFULNL_CFG_CMD_PF_BIND)
	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh = nflog_build_cfg_request(buf, C.NFULNL_CFG_CMD_BIND, qnum)
	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh = nflog_build_cfg_params(buf, C.NFULNL_COPY_PACKET, 0xFFFF, qnum)

	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	var ptr []byte
	var bsize int
	ret := mnl.MNL_CB_OK
	for ret >= mnl.MNL_CB_STOP {
		hdr := nl.Frame(mnl.MNL_RING_RX)
		if hdr.Status == C.NL_MMAP_STATUS_VALID {
			if hdr.Len == 0 {
				goto release
			}
			bsize = int(hdr.Len)
			ptr = mnl.RingMsghdr(hdr)
		} else if hdr.Status == C.NL_MMAP_STATUS_COPY {
			nrecv, err := nl.Recvfrom(buf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
				os.Exit(C.EXIT_FAILURE)
			}
			bsize = int(nrecv)
			ptr = buf
		} else {
			if (mnl_socket_poll(nl) == -1) {
				fmt.Fprintf(os.Stderr, "mnl_socket_poll")
				os.Exit(C.EXIT_FAILURE)
			}
			continue
		}
		ret, err = mnl.CbRun(ptr[:bsize], 0, portid, log_cb, nil)
	release:
		hdr.Status = C.NL_MMAP_STATUS_UNUSED
		nl.AdvanceRing(mnl.MNL_RING_RX)
	}

	if ret < mnl.MNL_CB_STOP {
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
