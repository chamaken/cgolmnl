package cgolmnl_test

import (
	. "github.com/chamaken/cgolmnl"
	. "github.com/chamaken/cgolmnl/testlib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math/rand"
	// "syscall"
	"fmt"
	"os"
	"time"
	"unsafe"
)

var _ = Describe("Attr", func() {
	fmt.Fprintf(os.Stdout, "Hello, nlmsg tester!\n") // to import os, sys for debugging
	var (
		BUFLEN = 512
		r      *rand.Rand

		// Nlmsg
		hbuf      *NlmsgBuf
		nlh       *Nlmsg
		rand_hbuf *NlmsgBuf
		rand_nlh  *Nlmsg

		// Nlattr
		abuf      *NlattrBuf
		nla       *Nlattr
		rand_abuf *NlattrBuf
		rand_nla  *Nlattr
	)

	BeforeEach(func() {
		r = rand.New(rand.NewSource(time.Now().Unix()))
		hbuf = NewNlmsgBuf(BUFLEN)
		nlh = NlmsgBytes(*hbuf)
		rand_hbuf = NewNlmsgBuf(BUFLEN)
		for i := 0; i < BUFLEN; i++ {
			(*(*[]byte)(rand_hbuf))[i] = byte(r.Int() % 256)
		}
		rand_nlh = NlmsgBytes(*rand_hbuf)
		rand_hbuf.SetLen(uint32(BUFLEN))

		abuf = NewNlattrBuf(BUFLEN)
		nla = NlattrPointer(*abuf)
		rand_abuf = NewNlattrBuf(BUFLEN)
		for i := 0; i < BUFLEN; i++ {
			(*(*[]byte)(rand_abuf))[i] = byte(r.Int() % 256)
		}
		rand_nla = NlattrPointer(*rand_abuf)
		rand_abuf.SetLen(uint16(BUFLEN))
	})

	Context("NlmsgBytes", func() {
		It("should share uint32 0x12345678 len", func() {
			hbuf.SetLen(0x12345678)
			Expect(nlh.Len).To(Equal(uint32(0x12345678)))
		})
		It("should share uint16 0x9abc type", func() {
			hbuf.SetType(0x9abc)
			Expect(nlh.Type).To(Equal(uint16(0x9abc)))
		})
		It("should share uint16 0xdef0 flags", func() {
			hbuf.SetFlags(0xdef0)
			Expect(nlh.Flags).To(Equal(uint16(0xdef0)))
		})
		It("should share uint32 0x23456789 flags", func() {
			hbuf.SetSeq(0x23456789)
			Expect(nlh.Seq).To(Equal(uint32(0x23456789)))
		})
		It("should share uint32 0xabcdef01 pid", func() {
			hbuf.SetPid(0xabcdef01)
			Expect(nlh.Pid).To(Equal(uint32(0xabcdef01)))
		})
	})

	Context("NlmsgSize", func() {
		It("should return 19", func() {
			Expect(NlmsgSize(3)).To(Equal(Size_t(19)))
		})
	})

	Context("NlmsgPayloadLen", func() {
		It("should return aligned 123 - MNL_NLMSG_HDRLEN", func() {
			hbuf.SetLen(MnlAlign(123))
			Expect(nlh.PayloadLen()).To(Equal(Size_t((MnlAlign(123) - MNL_NLMSG_HDRLEN))))
		})
	})

	Context("NlmsgPutExtraHeader", func() {
		var p unsafe.Pointer
		BeforeEach(func() {
			rand_hbuf.SetLen(MnlAlign(256))
			p = rand_nlh.PutExtraHeader(123)
		})
		It("nlmsg_len should be added aligned 123", func() {
			Expect(rand_nlh.Len).To(Equal(uint32(256 + MnlAlign(123))))
		})
		It("contents should be all 0", func() {
			for i := 256; i < 256+int(MnlAlign(123)); i++ {
				Expect((*rand_hbuf)[i]).To(BeZero())
			}
		})
	})

	Context("NlmsgGetPayload", func() {
		It("points buffer[MNL_NLMSG_HDRLEN]", func() {
			rand_hbuf.SetLen(MnlAlign(384))
			Expect(rand_nlh.Payload()).To(Equal(unsafe.Pointer(&(*rand_hbuf)[MNL_NLMSG_HDRLEN])))
		})
	})

	Context("NlmsgGetPayloadBytes", func() {
		It("length should be len - MNL_NLMSG_HDRLEN", func() {
			rand_hbuf.SetLen(MnlAlign(384))
			Expect(len(rand_nlh.PayloadBytes())).To(Equal(int(MnlAlign(384) - MNL_NLMSG_HDRLEN)))
		})
		It("contents should be the same", func() {
			rand_hbuf.SetLen(MnlAlign(384))
			Expect(rand_nlh.PayloadBytes()).To(Equal((*(*[]byte)(rand_hbuf))[MNL_NLMSG_HDRLEN:MnlAlign(384)]))
		})
	})

	Context("NlmsgGetPayloadOffset", func() {
		It("points buffer[MNL_NLMSG_HDRLEN + offset]", func() {
			Expect(rand_nlh.PayloadOffset(191)).To(Equal(unsafe.Pointer(&(*rand_hbuf)[MNL_NLMSG_HDRLEN+MnlAlign(191)])))
		})
	})

	Context("NlmsgGetPayloadOffsetBytes", func() {
		It("length should be len - MNL_NLMSG_HDRLEN - offset", func() {
			Expect(len(rand_nlh.PayloadOffsetBytes(191))).To(Equal(BUFLEN - int(MNL_NLMSG_HDRLEN+MnlAlign(191))))
		})
		It("contets should be the same", func() {
			Expect(rand_nlh.PayloadOffsetBytes(191)).To(Equal((*(*[]byte)(rand_hbuf))[int(MNL_NLMSG_HDRLEN+MnlAlign(191)):]))
		})
	})

	Context("NlmsgOk", func() {
		Describe("length: 16", func() {
			BeforeEach(func() {
				hbuf.SetLen(16)
			})
			It("param 15 should be false", func() {
				Expect(nlh.Ok(15)).To(Equal(false))
			})
			It("param 16 shoule be true", func() {
				Expect(nlh.Ok(16)).To(Equal(true))
			})
			It("param 17 shoule be true", func() {
				Expect(nlh.Ok(17)).To(Equal(true))
			})
		})
		Describe("length: 8", func() {
			BeforeEach(func() {
				hbuf.SetLen(8)
			})
			It("param 7 should be false", func() {
				Expect(nlh.Ok(7)).To(Equal(false))
			})
			It("param 8 shoule be false", func() {
				Expect(nlh.Ok(8)).To(Equal(false))
			})
			It("param 9 shoule be false", func() {
				Expect(nlh.Ok(9)).To(Equal(false))
			})
		})
		Describe("length: 32", func() {
			BeforeEach(func() {
				hbuf.SetLen(32)
			})
			It("param 31 should be false", func() {
				Expect(nlh.Ok(31)).To(Equal(false))
			})
			It("param 32 shoule be true", func() {
				Expect(nlh.Ok(32)).To(Equal(true))
			})
			It("param 33 shoule be true", func() {
				Expect(nlh.Ok(33)).To(Equal(true))
			})
		})
	})

	Context("NlmsgNext", func() {
		It("should have 3 valid (empty) messages", func() {
			hbuf.SetLen(MnlAlign(256))
			SetUint32(*hbuf, uint(MnlAlign(256)), 128)
			SetUint32(*hbuf, uint(MnlAlign(256)+MnlAlign(128)), 64)

			next_nlh, rest := nlh.Next(BUFLEN)
			Expect(rest).To(Equal(BUFLEN - 256))
			Expect(next_nlh.Len).To(Equal(MnlAlign(128)))
			Expect(next_nlh.Ok(rest)).To(BeTrue())

			next_nlh, rest = next_nlh.Next(rest)
			Expect(rest).To(Equal(BUFLEN - 256 - 128))
			Expect(next_nlh.Len).To(Equal(MnlAlign(64)))
			Expect(next_nlh.Ok(rest)).To(BeTrue())

			next_nlh, rest = next_nlh.Next(rest)
			Expect(rest).To(Equal(BUFLEN - 256 - 128 - 64))
			Expect(next_nlh.Ok(rest)).To(BeFalse())
		})
	})

	Context("NlmsgGetPayloadTail", func() {
		It("address should be the same", func() {
			hbuf.SetLen(MnlAlign(323))
			Expect(nlh.PayloadTail()).To(Equal(unsafe.Pointer(&(*hbuf)[MnlAlign(323)])))
		})
	})

	Context("NlmsgSeqOk", func() {
		It("should be equal to but accept 0", func() {
			hbuf.SetSeq(0x12345678)
			Expect(nlh.SeqOk(0x12345678)).To(BeTrue())
			Expect(nlh.SeqOk(0x12345)).To(BeFalse())
			Expect(nlh.SeqOk(0)).To(BeTrue())

			hbuf.SetSeq(0)
			Expect(nlh.SeqOk(0x12345678)).To(BeTrue())
		})
	})

	Context("NlmsgPortidOk", func() {
		It("should be equal to but accept 0", func() {
			hbuf.SetPid(0x12345678)
			Expect(nlh.PortidOk(0x12345678)).To(BeTrue())
			Expect(nlh.PortidOk(0x12345)).To(BeFalse())
			Expect(nlh.PortidOk(0)).To(BeTrue())

			hbuf.SetPid(0)
			Expect(nlh.PortidOk(0x12345678)).To(BeTrue())
		})
	})

	Context("NlmsgBatches", func() {
		// init
		var b *NlmsgBatch

		BeforeEach(func() {
			buf := make([]byte, 301)
			b, _ = NewNlmsgBatch(buf, 163)
		})

		It("should indicate initial, empty states", func() {
			Expect(b.Size()).To(BeZero())
			Expect(len(b.HeadBytes())).To(BeZero())
			Expect(b.IsEmpty()).To(BeTrue())
		})
		It("has apropriate length in filling buffer", func() {
			// filling buf
			for i := 1; i < 11; i++ {
				_nlh := (*Nlmsg)(b.Current())
				_nlh.PutHeader()
				Expect(b.Next()).To(BeTrue())
				Expect(b.Size()).To(Equal(Size_t(int(MNL_NLMSG_HDRLEN) * i)))
				Expect(len(b.HeadBytes())).To(Equal(int(MNL_NLMSG_HDRLEN) * i))
				Expect(b.IsEmpty()).To(BeFalse())
			}
		})
		It("next should indicate false after filling up", func() {
			for i := 0; i < 11; i++ {
				_nlh := (*Nlmsg)(b.Current())
				_nlh.PutHeader()
				b.Next()
			}
			Expect(b.Next()).To(BeFalse())
			Expect(b.Size()).To(Equal(Size_t(int(MNL_NLMSG_HDRLEN) * 10)))
			Expect(len(b.HeadBytes())).To(Equal(int(MNL_NLMSG_HDRLEN) * 10))
			Expect(b.IsEmpty()).To(BeFalse())
		})
		It("should one header after filling up and reset", func() {
			for i := 0; i < 11; i++ {
				_nlh := (*Nlmsg)(b.Current())
				_nlh.PutHeader()
				b.Next()
			}
			b.Reset()
			Expect(b.Size()).To(Equal(Size_t(MNL_NLMSG_HDRLEN)))
			Expect(len(b.HeadBytes())).To(Equal(int(MNL_NLMSG_HDRLEN)))
			Expect(b.IsEmpty()).To(BeFalse())
			Expect(b.Next()).To(BeTrue())

			b.Reset()
			Expect(b.Size()).To(BeZero())
			Expect(len(b.HeadBytes())).To(BeZero())
			Expect(b.IsEmpty()).To(BeTrue())
		})
	})
})
