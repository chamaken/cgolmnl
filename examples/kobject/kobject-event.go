package main

/*
#include <stdlib.h>
#include <linux/netlink.h>
*/
import "C"

import (
	"fmt"
	mnl "github.com/chamaken/cgolmnl"
	"os"
)

func main() {
	nl, err := mnl.NewSocket(C.NETLINK_KOBJECT_UEVENT)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_open: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	if err := nl.Bind((1 << 0), mnl.MNL_SOCKET_AUTOPID); err != nil {
		fmt.Fprintf(os.Stderr, "mnl_socket_bind: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}

	buf := make([]byte, mnl.MNL_SOCKET_BUFFER_SIZE)
	ret := mnl.MNL_CB_OK
	for ret > mnl.MNL_CB_STOP {
		nrcv, err := nl.Recvfrom(buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mnl_cb_recvfrom: %s\n", err)
			os.Exit(C.EXIT_FAILURE)
		}
		for i := 0; i < int(nrcv); i++ {
			fmt.Printf("%c", buf[i])
		}
		fmt.Println()
	}

	if ret < mnl.MNL_CB_STOP {
		fmt.Fprintf(os.Stderr, "mnl_cb_run: %s\n", err)
		os.Exit(C.EXIT_FAILURE)
	}
}
