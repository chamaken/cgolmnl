package main

/*
#include <stdlib.h>
#include <arpa/inet.h>
#include <linux/netlink.h>

#include <linux/netfilter/nfnetlink.h>
#include <linux/netfilter/nfnetlink_queue.h>
#include <linux/netfilter.h>
*/
import "C"

import (
	"fmt"
	mnl "github.com/chamaken/cgolmnl"
	inet "github.com/chamaken/cgolmnl/inet"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"os"
	"strconv"
	"syscall"
)

func parse_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.NFQA_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch int(attr_type) {
	case C.NFQA_MARK:
		fallthrough
	case C.NFQA_IFINDEX_INDEV:
		fallthrough
	case C.NFQA_IFINDEX_OUTDEV:
		fallthrough
	case C.NFQA_IFINDEX_PHYSINDEV:
		fallthrough
	case C.NFQA_IFINDEX_PHYSOUTDEV:
		if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.NFQA_TIMESTAMP:
		if err := attr.Validate2(mnl.MNL_TYPE_UNSPEC, SizeofNfqnlMsgPacketTimestamp); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.NFQA_HWADDR:
		if err := attr.Validate2(mnl.MNL_TYPE_UNSPEC, SizeofNfqnlMsgPacketHw); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate2: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.NFQA_PAYLOAD:
		// do something
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func queue_cb(nlh *mnl.Nlmsg, data interface{}) (int, syscall.Errno) {
	var ph *NfqnlMsgPacketHdr
	var id uint32
	var etherType uint16

	tb := make(map[uint16]*mnl.Nlattr, C.NFQA_MAX+1)

	nlh.Parse(SizeofNfgenmsg, parse_attr_cb, tb)
	if tb[C.NFQA_PACKET_HDR] != nil {
		ph = (*NfqnlMsgPacketHdr)(tb[C.NFQA_PACKET_HDR].Payload())
		id = inet.Ntohl(ph.Packet_id)

		fmt.Printf("packet received (id=%d, hw=0x%04x hook=%d)\n", id, inet.Ntohs(ph.Hw_protocol), ph.Hook)
		etherType = inet.Ntohs(ph.Hw_protocol)
	}

	if tb[C.NFQA_PAYLOAD] != nil {
		pbuf := tb[C.NFQA_PAYLOAD].PayloadBytes()
		switch etherType {
		case syscall.ETH_P_IP:
			pkt := layers.IPv4{}
			if err := pkt.DecodeFromBytes(pbuf, gopacket.NilDecodeFeedback); err != nil {
				fmt.Fprintf(os.Stderr, "IPv4.DecodeFromBytes: %s\n", err)
				return mnl.MNL_CB_ERROR, syscall.EINVAL
			}
			fmt.Printf("payload src: %v, dst: %v, protocol: %v\n",
				pkt.SrcIP, pkt.DstIP, pkt.Protocol)
		case syscall.ETH_P_IPV6:
			pkt := layers.IPv6{}
			if err := pkt.DecodeFromBytes(pbuf, gopacket.NilDecodeFeedback); err != nil {
				fmt.Fprintf(os.Stderr, "IPv6.DecodeFromBytes: %s\n", err)
				return mnl.MNL_CB_ERROR, syscall.EINVAL
			}
			fmt.Printf("payload src: %v, dst: %v, next header: %v\n",
				pkt.SrcIP, pkt.DstIP, pkt.NextHeader)
		}
	}
	return mnl.MNL_CB_OK + int(id), 0
}

func nfq_build_cfg_pf_request(buf []byte, command uint8) *mnl.Nlmsg {
	nlh, err := mnl.NewNlmsgBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh.Type = (C.NFNL_SUBSYS_QUEUE << 8) | C.NFQNL_MSG_CONFIG
	nlh.Flags = C.NLM_F_REQUEST

	nfg := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfg.Nfgen_family = C.AF_UNSPEC
	nfg.Version = C.NFNETLINK_V0

	cmd := &NfqnlMsgConfigCmd{Command: command, Pf: inet.Htons(C.AF_INET)}
	nlh.PutPtr(C.NFQA_CFG_CMD, cmd)

	return nlh
}

func nfq_build_cfg_request(buf []byte, command uint8, queue_num int) *mnl.Nlmsg {
	nlh, err := mnl.NewNlmsgBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh.Type = (C.NFNL_SUBSYS_QUEUE << 8) | C.NFQNL_MSG_CONFIG
	nlh.Flags = C.NLM_F_REQUEST

	nfg := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfg.Nfgen_family = C.AF_UNSPEC
	nfg.Version = C.NFNETLINK_V0
	nfg.Res_id = inet.Htons(uint16(queue_num))

	cmd := &NfqnlMsgConfigCmd{Command: command, Pf: inet.Htons(C.AF_INET)}
	nlh.PutPtr(C.NFQA_CFG_CMD, cmd)

	return nlh
}

func nfq_build_cfg_params(buf []byte, copy_mode uint8, copy_range, queue_num int) *mnl.Nlmsg {
	nlh, err := mnl.NewNlmsgBytes(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nlmsg_put_header: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	nlh.Type = (C.NFNL_SUBSYS_QUEUE << 8) | C.NFQNL_MSG_CONFIG
	nlh.Flags = C.NLM_F_REQUEST

	nfg := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfg.Nfgen_family = C.AF_UNSPEC
	nfg.Version = C.NFNETLINK_V0
	nfg.Res_id = inet.Htons(uint16(queue_num))

	params := &NfqnlMsgConfigParams{Range: inet.Htonl(uint32(copy_range)), Mode: copy_mode}
	nlh.PutPtr(C.NFQA_CFG_PARAMS, params)

	return nlh
}

func nfq_build_verdict(buf []byte, id, queue_num, verd int) *mnl.Nlmsg {
	nlh, _ := mnl.NewNlmsgBytes(buf)
	nlh.Type = (C.NFNL_SUBSYS_QUEUE << 8) | C.NFQNL_MSG_VERDICT
	nlh.Flags = C.NLM_F_REQUEST
	nfg := (*Nfgenmsg)(nlh.PutExtraHeader(SizeofNfgenmsg))
	nfg.Nfgen_family = C.AF_UNSPEC
	nfg.Version = C.NFNETLINK_V0
	nfg.Res_id = inet.Htons(uint16(queue_num))

	vh := &NfqnlMsgVerdictHdr{Verdict: inet.Htonl(uint32(verd)), Id: inet.Htonl(uint32(id))}
	nlh.PutPtr(C.NFQA_VERDICT_HDR, vh)

	return nlh
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [queue_num]\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}
	queue_num, _ := strconv.Atoi(os.Args[1])

	nl, err := mnl.NewSocket(C.NETLINK_NETFILTER)
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

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)

	nlh := nfq_build_cfg_pf_request(buf, C.NFQNL_CFG_CMD_PF_UNBIND)
	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh = nfq_build_cfg_pf_request(buf, C.NFQNL_CFG_CMD_PF_BIND)
	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh = nfq_build_cfg_request(buf, C.NFQNL_CFG_CMD_BIND, queue_num)
	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	nlh = nfq_build_cfg_params(buf, C.NFQNL_COPY_PACKET, 0xFFFF, queue_num)
	if _, err := nl.SendNlmsg(nlh); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	ret := mnl.MNL_CB_OK
	for ret >= mnl.MNL_CB_STOP {
		nrcv, err := nl.Recvfrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		ret, err = mnl.CbRun(buf[:nrcv], 0, portid, queue_cb, nil)
		if ret < mnl.MNL_CB_STOP {
			fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		id := ret - mnl.MNL_CB_OK
		nlh = nfq_build_verdict(buf, id, queue_num, C.NF_ACCEPT)
		if _, err := nl.SendNlmsg(nlh); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_sendto: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
	}
}
