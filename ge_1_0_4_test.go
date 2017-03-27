// +build ge_1_0_4

package cgolmnl_test

import (
	. "github.com/chamaken/cgolmnl"
	. "github.com/chamaken/cgolmnl/testlib"
	. "github.com/onsi/ginkgo"

	"fmt"
	"os"
	"syscall"
)


var _ = Describe("Socket", func() {
	fmt.Fprintf(os.Stdout, "Hello, socket fd tester!\n") // to import os, sys for debugging
	var (
		nl *Socket
	)

	BeforeEach(func() {
		fd, _ := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, NETLINK_NETFILTER)
		nl, _ = NewSocketFd(fd)
	})

	AfterEach(func() {
		nl.Close()
	})

	socketContexts(nl)
})

var _ = Describe("Socket", func() {
	fmt.Fprintf(os.Stdout, "Hello, socket open2 tester!\n") // to import os, sys for debugging
	var (
		nl *Socket
	)

	BeforeEach(func() {
		nl, _ = NewSocket2(NETLINK_NETFILTER, 0)
	})

	AfterEach(func() {
		nl.Close()
	})

	socketContexts(nl)
})
