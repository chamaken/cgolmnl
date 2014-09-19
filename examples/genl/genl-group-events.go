package main

/*
#include <stdlib.h>
#include <linux/genetlink.h>
*/
import "C"

import (
	mnl "cgolmnl"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

var group int

func data_cb(nlh *mnl.Nlmsghdr, data interface{}) (int, syscall.Errno) {
	fmt.Printf("received event type=%d from genetlink group %d\n", nlh.Type, group)
	return mnl.MNL_CB_OK, 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("%s [group]\n", os.Args[0])
		os.Exit(C.EXIT_FAILURE)
	}
	group, _ = strconv.Atoi(os.Args[1])

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
	if err := nl.SetsockoptCint(C.NETLINK_ADD_MEMBERSHIP, group); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_setsockopt")
		os.Exit(C.EXIT_FAILURE)
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	ret := mnl.MNL_CB_OK
	for ret > mnl.MNL_CB_STOP {
		nrcv, err := nl.Recvfrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_socket_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		ret, err = mnl.CbRun(buf[:nrcv], 0, 0, data_cb, nil)
	}

	if ret < mnl.MNL_CB_STOP {
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
