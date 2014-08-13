package cgolmnl_test

import (
	. "github.com/chamaken/cgolmnl"
	. "github.com/chamaken/cgolmnl/testlib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"syscall"
	"fmt"
	"os"
)

var _ = Describe("callback", func() {
	fmt.Fprintf(os.Stdout, "Hello, callback tester!\n")
	var (
		nlmsghdr_noop,
		nlmsghdr_done,
		nlmsghdr_overrun,
		nlmsghdr_mintype,
		nlmsghdr_error,
		nlmsghdr_typeFF,
		nlmsghdr_pid2,
		nlmsghdr_seq2,
		nlmsghdr_intr	*NlmsghdrBuf
	)

	BeforeEach(func() {
		nlmsghdr_noop = NewNlmsghdrBuf(16)
		nlmsghdr_noop.SetLen(16)
		nlmsghdr_noop.SetType(NLMSG_NOOP)
		nlmsghdr_noop.SetFlags(NLM_F_REQUEST)
		nlmsghdr_noop.SetSeq(1)
		nlmsghdr_noop.SetPid(1)

		nlmsghdr_done = NewNlmsghdrBuf(16)
		nlmsghdr_done.SetLen(16)
		nlmsghdr_done.SetType(NLMSG_DONE)
		nlmsghdr_done.SetFlags(NLM_F_REQUEST)
		nlmsghdr_done.SetSeq(1)
		nlmsghdr_done.SetPid(1)

		nlmsghdr_overrun = NewNlmsghdrBuf(16)
		nlmsghdr_overrun.SetLen(16)
		nlmsghdr_overrun.SetType(NLMSG_OVERRUN)
		nlmsghdr_overrun.SetFlags(NLM_F_REQUEST)
		nlmsghdr_overrun.SetSeq(1)
		nlmsghdr_overrun.SetPid(1)

		mintype_msg := NewNlmsghdrBuf(16)
		mintype_msg.SetLen(16)
		mintype_msg.SetType(NLMSG_MIN_TYPE)
		mintype_msg.SetFlags(NLM_F_REQUEST)
		mintype_msg.SetSeq(1)
		mintype_msg.SetPid(1)
		nlmsghdr_mintype = NewNlmsghdrBuf(16)
		copy(([]byte)(*nlmsghdr_mintype), ([]byte)(*nlmsghdr_noop))
		b1 := append(([]byte)(*nlmsghdr_mintype), ([]byte)(*mintype_msg)...)
		nlmsghdr_mintype = (*NlmsghdrBuf)(&b1)

		nlmsghdr_error = NewNlmsghdrBuf(16)
		nlmsghdr_error.SetType(NLMSG_ERROR)
		nlmsghdr_error.SetFlags(NLM_F_REQUEST)
		nlmsghdr_error.SetSeq(1)
		nlmsghdr_error.SetPid(1)
		errno := make([]byte, SizeofCint)
		SetCint(errno, 0, 1) // EPERM
		b2 := append(([]byte)(*nlmsghdr_error), errno...)
		error_msg := NewNlmsghdrBuf(16)
		error_msg.SetLen(16)
		error_msg.SetType(NLMSG_ERROR)
		error_msg.SetFlags(NLM_F_REQUEST)
		error_msg.SetSeq(1)
		error_msg.SetPid(1)
		b2 = append(b2, ([]byte)(*error_msg)...)
		nlmsghdr_error = (*NlmsghdrBuf)(&b2)
		nlmsghdr_error.SetLen(uint32(len(b2)))

		typeFF_msg := NewNlmsghdrBuf(16)
		typeFF_msg.SetLen(16)
		typeFF_msg.SetType(0xff)
		typeFF_msg.SetFlags(NLM_F_REQUEST)
		typeFF_msg.SetSeq(1)
		typeFF_msg.SetPid(1)
		b3 := append(([]byte)(*nlmsghdr_mintype), ([]byte)(*typeFF_msg)...)
		nlmsghdr_typeFF = (*NlmsghdrBuf)(&b3)

		pid2_msg := NewNlmsghdrBuf(16)
		pid2_msg.SetLen(16)
		pid2_msg.SetType(0xff)
		pid2_msg.SetFlags(NLM_F_REQUEST)
		pid2_msg.SetSeq(1)
		pid2_msg.SetPid(2)
		b4 := append(([]byte)(*nlmsghdr_mintype), ([]byte)(*pid2_msg)...)
		nlmsghdr_pid2 = (*NlmsghdrBuf)(&b4)

		seq2_msg := NewNlmsghdrBuf(16)
		seq2_msg.SetLen(16)
		seq2_msg.SetType(0xff)
		seq2_msg.SetFlags(NLM_F_REQUEST)
		seq2_msg.SetSeq(2)
		seq2_msg.SetPid(1)
		b5 := append(([]byte)(*nlmsghdr_mintype), ([]byte)(*seq2_msg)...)
		nlmsghdr_seq2 = (*NlmsghdrBuf)(&b5)

		intr_msg := NewNlmsghdrBuf(16)
		intr_msg.SetLen(16)
		intr_msg.SetType(0xff)
		intr_msg.SetFlags(NLM_F_REQUEST | NLM_F_DUMP_INTR)
		intr_msg.SetSeq(1)
		intr_msg.SetPid(1)
		b6 := append(([]byte)(*nlmsghdr_mintype), ([]byte)(*intr_msg)...)
		nlmsghdr_intr = (*NlmsghdrBuf)(&b6)
	})

	Context("CbRun", func() {
		cb := func(nlh *Nlmsghdr, data interface{}) (int, syscall.Errno) {
			if data != nil {
				l := data.(*[]uint16)
				*l = append(*l, nlh.Type)
			}
			if nlh.Type == 0xff {
				return MNL_CB_ERROR, 0
			}
			return MNL_CB_OK, 0
		}

		It("should return MNL_CB_OK for NOOP", func() {
			ret, err := CbRun(([]byte)(*nlmsghdr_noop), 1, 1, nil, nil)
			Expect(err).To(BeNil())
			Expect(ret).To(Equal(MNL_CB_OK))
		})
		It("should return MNL_CB_OK for MIN_TYPE", func() {
			ret, err := CbRun(([]byte)(*nlmsghdr_mintype), 1, 1, nil, nil)
			Expect(err).To(BeNil())
			Expect(ret).To(Equal(MNL_CB_OK))
		})
		It("should return MNL_CB_ERROR and EPERM(1) for NLMSG_ERROR", func() {
			ret, err := CbRun(([]byte)(*nlmsghdr_error), 1, 1, nil, nil)
			Expect(err).To(Equal(syscall.Errno(1)))
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
		It("should return MNL_CB_ERROR and type is set", func() {
			l := make([]uint16, 0)
			ret, err := CbRun(([]byte)(*nlmsghdr_typeFF), 1, 1, cb, &l)
			Expect(err).To(BeNil())
			Expect(ret).To(Equal(MNL_CB_ERROR))
			Expect(l[0]).To(Equal(uint16(NLMSG_MIN_TYPE)))
			Expect(l[1]).To(Equal(uint16(0xff)))
		})
		It("should return ESRCH", func() {
			ret, err := CbRun(([]byte)(*nlmsghdr_pid2), 1, 1, nil, nil)
			Expect(err).To(Equal(syscall.ESRCH))
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
		It("should return EPROTO", func() {
			ret, err := CbRun(([]byte)(*nlmsghdr_seq2), 1, 1, nil, nil)
			Expect(err).To(Equal(syscall.EPROTO))
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
		It("should return EINTR", func() {
			ret, err := CbRun(([]byte)(*nlmsghdr_intr), 1, 1, nil, nil)
			Expect(err).To(Equal(syscall.EINTR))
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
	})
	Context("CbRun2", func() {
		// NLMSG_NOOP		0x1
		// NLMSG_ERROR		0x2
		// NLMSG_DONE		0x3
		// NLMSG_OVERRUN	0x4
		// NLMSG_MIN_TYPE	0x10

		ctl_cb := func(nlh *Nlmsghdr, msgtype uint16, data interface{}) (int, syscall.Errno) {
			switch (msgtype) {
			case NLMSG_NOOP:
				return MNL_CB_ERROR, 0
			case NLMSG_OVERRUN:
				return MNL_CB_ERROR, syscall.ENOBUFS
			case NLMSG_DONE:
				return MNL_CB_STOP, 0
			case NLMSG_ERROR:
				// see original mnl_cb_error()
				var errno syscall.Errno
				err := (*Nlmsgerr)(nlh.Payload())
				if nlh.Len < uint32(NlmsgSize(SizeofNlmsgerr)) {
					return MNL_CB_ERROR, syscall.EBADMSG
				}
				if err.Error < 0 {
					errno = -syscall.Errno(err.Error)
				} else {
					errno = syscall.Errno(err.Error)
				}
				if err.Error == 0 {
					return MNL_CB_STOP, 0
				} else {
					return MNL_CB_ERROR, errno
				}
			default:
				return MNL_CB_OK, 0
			}
		}
		var ctltypes []uint16
		BeforeEach(func() {
			ctltypes = []uint16{NLMSG_OVERRUN, NLMSG_NOOP, NLMSG_ERROR, NLMSG_DONE}
		})
		It("should return MNL_CB_ERROR", func() {
			ret, err := CbRun2(([]byte)(*nlmsghdr_noop), 1, 1, nil, nil,
				ctl_cb, ctltypes)
			Expect(err).To(BeNil())
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
		It("should return (default) MNL_CB_OK", func() {
			ctltypes := []uint16{NLMSG_OVERRUN, NLMSG_ERROR, NLMSG_DONE}
			ret, err := CbRun2(([]byte)(*nlmsghdr_noop), 1, 1, nil, nil,
				ctl_cb, ctltypes)
			Expect(err).To(BeNil())
			Expect(ret).To(Equal(MNL_CB_OK))
		})
		It("should return MNL_CB_STOP", func() {
			ret, err := CbRun2(([]byte)(*nlmsghdr_done), 1, 1, nil, nil,
				ctl_cb, ctltypes)
			Expect(err).To(BeNil())
			Expect(ret).To(Equal(MNL_CB_STOP))
		})
		It("should return ENOBUFS", func() {
			ret, err := CbRun2(([]byte)(*nlmsghdr_overrun), 1, 1, nil, nil,
				ctl_cb, ctltypes)
			Expect(err).To(Equal(syscall.ENOBUFS))
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
		It("should return EPERM", func() {
			ret, err := CbRun2(([]byte)(*nlmsghdr_error), 1, 1, nil, nil,
				ctl_cb, ctltypes)
			Expect(err).To(Equal(syscall.EPERM))
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
		It("should return EOPNOTSUPP", func() {
			ret, err := CbRun2(([]byte)(*nlmsghdr_noop), 1, 1, nil, nil,
				nil, ctltypes)
			Expect(err).To(Equal(syscall.EOPNOTSUPP))
			Expect(ret).To(Equal(MNL_CB_ERROR))
		})
	})
})
