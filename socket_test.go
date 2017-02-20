package cgolmnl_test

import (
	. "github.com/chamaken/cgolmnl"
	. "github.com/chamaken/cgolmnl/testlib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"syscall"
)

func socketContexts(nl *Socket) func() {
	return func() {
		Context("Socket Descriptor", func() {
			It("should be a socket", func() {
				buf := &syscall.Stat_t{}
				fd := nl.Fd()
				_ = syscall.Fstat(fd, buf)
				Expect(IsSock(buf.Mode)).To(BeTrue())
			})
			It("should be closed and become invalid fd", func() {
				fd := nl.Fd()
				Expect(nl.Close()).To(BeNil())
				Expect(IsValidFd(fd)).To(BeFalse())
				nl, _ = NewSocket(NETLINK_NETFILTER) // for AfterEach
			})
		})
		Context("Bind and Port", func() {
			It("port should be zero", func() {
				Expect(nl.Portid()).To(BeZero())
			})
			It("portid shuld be same as bind and not bind again", func() {
				err := nl.Bind(0, 65432) // may fail
				Expect(err).To(BeNil())
				Expect(nl.Portid()).To(Equal(uint32(65432)))
				err = nl.Bind(0, MNL_SOCKET_AUTOPID)
				Expect(err).To(Equal(syscall.EINVAL))
			})
		})
		Context("Send and Recv", func() {
			var nlh *Nlmsg

			BeforeEach(func() {
				_ = nl.Bind(0, MNL_SOCKET_AUTOPID)
				nlh, _ = NewNlmsg(int(MNL_NLMSG_HDRLEN))
				nlh.Type = NLMSG_NOOP
				nlh.Flags = NLM_F_ECHO | NLM_F_ACK
				nlh.Pid = nl.Portid()
				nlh.Seq = 1234
			})
			It("should sendto and recv err message", func() {
				b1, _ := nlh.MarshalBinary()
				nsent, err := nl.Sendto(b1)
				Expect(err).To(BeNil())
				Expect(nsent).To(Equal(Ssize_t(MNL_NLMSG_HDRLEN)))
				b2 := make([]byte, 256)
				nrecv, err := nl.Recvfrom(b2)
				Expect(err).To(BeNil())
				Expect(nrecv).To(Equal(Ssize_t(36))) // nlmsghdr + nlmsgerr

				// nlr := NlmsgBytes(b2[:nrecv])
				nlr, _ := NewNlmsg(int(nrecv))
				nlr.Len = uint32(nrecv)
				nlr.UnmarshalBinary(b2[:nrecv])

				Expect(nlr.Len).To(Equal(uint32(36)))
				Expect(nlr.Type).To(Equal(uint16(NLMSG_ERROR)))

				nle := (*Nlmsgerr)(nlr.Payload())
				Expect(nle.Error).To(Equal(-int32(syscall.EPERM)))
				Expect(nle.Msg.Len).To(Equal(MNL_NLMSG_HDRLEN))
				Expect(nle.Msg.Flags).To(Equal(uint16(NLM_F_ECHO | NLM_F_ACK)))
				Expect(nle.Msg.Pid).To(Equal(nl.Portid()))
				Expect(nle.Msg.Seq).To(Equal(uint32(1234)))
			})
			It("should send_nlmsg and recv err message", func() {
				nsent, err := nl.SendNlmsg(nlh)
				Expect(err).To(BeNil())
				Expect(nsent).To(Equal(Ssize_t(MNL_NLMSG_HDRLEN)))
				b2 := make([]byte, 256)
				nrecv, err := nl.Recvfrom(b2)
				Expect(err).To(BeNil())
				Expect(nrecv).To(Equal(Ssize_t(36)))

				nlr, _ := NewNlmsg(int(nrecv))
				nlr.Len = uint32(nrecv)
				nlr.UnmarshalBinary(b2[:nrecv])

				Expect(nlr.Len).To(Equal(uint32(36)))
				Expect(nlr.Type).To(Equal(uint16(NLMSG_ERROR)))

				nle := (*Nlmsgerr)(nlr.Payload())
				Expect(nle.Error).To(Equal(-int32(syscall.EPERM)))
				Expect(nle.Msg.Len).To(Equal(MNL_NLMSG_HDRLEN))
				Expect(nle.Msg.Flags).To(Equal(uint16(NLM_F_ECHO | NLM_F_ACK)))
				Expect(nle.Msg.Pid).To(Equal(nl.Portid()))
				Expect(nle.Msg.Seq).To(Equal(uint32(1234)))
			})
		})
		Context("option set and get", func() {
			It("should set/get NETLINK_BROADCAST_ERROR by Cint", func() {
				Expect(nl.SetsockoptCint(NETLINK_BROADCAST_ERROR, 1)).To(BeNil())
				ret, err := nl.Sockopt(NETLINK_BROADCAST_ERROR, SizeofCint)
				Expect(err).To(BeNil())
				Expect(len(ret)).To(Equal(SizeofCint))
				// assume at least 32bit
				Expect(Endian.Uint32(ret[0:SizeofCint])).To(Equal(uint32(1)))
			})
		})
	}
}

var _ = Describe("Socket", func() {
	fmt.Fprintf(os.Stdout, "Hello, socket tester!\n") // to import os, sys for debugging
	var (
		nl *Socket
	)

	BeforeEach(func() {
		nl, _ = NewSocket(NETLINK_NETFILTER)
	})

	AfterEach(func() {
		nl.Close()
	})

	socketContexts(nl)
})
