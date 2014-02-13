package main

/*
#include <stdlib.h>
#include <linux/genetlink.h>
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"time"
	mnl "cgolmnl"
)

func parse_mc_grps_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.CTRL_ATTR_MCAST_GRP_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTRL_ATTR_MCAST_GRP_ID:
		if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTRL_ATTR_MCAST_GRP_NAME:
		if err := attr.Validate(mnl.MNL_TYPE_STRING); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func parse_genl_mc_grps(nested *mnl.Nlattr) {
	for pos := range nested.Nesteds() {
		tb := make(map[uint16]*mnl.Nlattr, C.CTRL_ATTR_MCAST_GRP_MAX + 1)

		pos.ParseNested(parse_mc_grps_cb, tb)
		if tb[C.CTRL_ATTR_MCAST_GRP_ID] != nil {
			fmt.Printf("id-0x%x ", tb[C.CTRL_ATTR_MCAST_GRP_ID].U32())
		}
		if tb[C.CTRL_ATTR_MCAST_GRP_NAME] != nil {
			fmt.Printf("name: %s ", tb[C.CTRL_ATTR_MCAST_GRP_NAME].Str())
		}
		fmt.Println()
	}
}

func parse_family_ops_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.CTRL_ATTR_OP_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTRL_ATTR_OP_ID:
		if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTRL_ATTR_OP_MAX:
		/* just break */
	default:
		return mnl.MNL_CB_OK, 0
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func parse_genl_family_ops(nested *mnl.Nlattr) {
	for pos := range nested.Nesteds() {
		tb := make(map[uint16]*mnl.Nlattr, C.CTRL_ATTR_OP_MAX + 1)

		pos.ParseNested(parse_family_ops_cb, tb)
		if tb[C.CTRL_ATTR_OP_ID] != nil {
			fmt.Printf("id-0x%x ", tb[C.CTRL_ATTR_OP_ID].U32())
		}
		if tb[C.CTRL_ATTR_OP_MAX] != nil {
			fmt.Printf("flags ");
		}
		fmt.Println()
	}
}

func data_attr_cb(attr *mnl.Nlattr, data interface{}) (int, syscall.Errno) {
	tb := data.(map[uint16]*mnl.Nlattr)
	attr_type := attr.GetType()

	if err := attr.TypeValid(C.CTRL_ATTR_MAX); err != nil {
		return mnl.MNL_CB_OK, 0
	}

	switch attr_type {
	case C.CTRL_ATTR_FAMILY_NAME:
		if err := attr.Validate(mnl.MNL_TYPE_STRING); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTRL_ATTR_FAMILY_ID:
		if err := attr.Validate(mnl.MNL_TYPE_U16); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTRL_ATTR_VERSION:	fallthrough
	case C.CTRL_ATTR_HDRSIZE:	fallthrough
	case C.CTRL_ATTR_MAXATTR:
		if err := attr.Validate(mnl.MNL_TYPE_U32); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	case C.CTRL_ATTR_OPS:		fallthrough
	case C.CTRL_ATTR_MCAST_GROUPS:
		if err := attr.Validate(mnl.MNL_TYPE_NESTED); err != nil {
			fmt.Fprintf(os.Stderr, "mnl_attr_validate: %s\n", err)
			return mnl.MNL_CB_ERROR, err.(syscall.Errno)
		}
	}
	tb[attr_type] = attr
	return mnl.MNL_CB_OK, 0
}

func data_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	tb := make(map[uint16]*mnl.Nlattr, C.CTRL_ATTR_MAX + 1)
	// genl := (*Genlmsghdr)(nlh.Payload())

	nlh.Parse(SizeofGenlmsghdr, data_attr_cb, tb)
	if tb[C.CTRL_ATTR_FAMILY_NAME] != nil {
		fmt.Printf("name=%s\t", tb[C.CTRL_ATTR_FAMILY_NAME].Str())
	}
	if tb[C.CTRL_ATTR_FAMILY_ID] != nil {
		fmt.Printf("id=%d\t", tb[C.CTRL_ATTR_FAMILY_ID].U16())
	}
	if tb[C.CTRL_ATTR_HDRSIZE] != nil {
		fmt.Printf("hrsize=%d\t", tb[C.CTRL_ATTR_HDRSIZE].U32())
	}
	if tb[C.CTRL_ATTR_MAXATTR] != nil {
		fmt.Printf("maxattr=%d\t", tb[C.CTRL_ATTR_MAXATTR].U32())
	}
	fmt.Println()
	if tb[C.CTRL_ATTR_OPS] != nil {
		fmt.Println("ops:")
		parse_genl_family_ops(tb[C.CTRL_ATTR_OPS])
	}
	if tb[C.CTRL_ATTR_MCAST_GROUPS] != nil {
		fmt.Println("grps:")
		parse_genl_mc_grps(tb[C.CTRL_ATTR_MCAST_GROUPS])
	}
	fmt.Println()
	return mnl.MNL_CB_OK, 0
}

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("%s [family name]\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	nlh, _ := mnl.NlmsgPutHeaderBytes(buf)
	nlh.Type = C.GENL_ID_CTRL
	nlh.Flags = C.NLM_F_REQUEST | C.NLM_F_ACK
	seq := uint32(time.Now().Unix())
	nlh.Seq = seq

	genl := (*Genlmsghdr)(nlh.PutExtraHeader(SizeofGenlmsghdr))
	genl.Cmd = C.CTRL_CMD_GETFAMILY
	genl.Version = 1

	nlh.PutU32(C.CTRL_ATTR_FAMILY_ID, C.GENL_ID_CTRL)
	if len(os.Args) >= 2 {
		nlh.PutStrz(C.CTRL_ATTR_FAMILY_NAME, os.Args[1])
	} else {
		nlh.Flags |= C.NLM_F_DUMP
	}

	nl, err := mnl.NewSocket(C.NETLINK_GENERIC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	defer nl.Close()

	if err := nl.Bind(0, mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
	portid := nl.Portid()

	if _, err := nl.SendNlmsg(nlh); err != nil {
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
