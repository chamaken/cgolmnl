package cgolmnl_test

import (
	. "cgolmnl"
	. "cgolmnl/testlib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math/rand"
	"syscall"
	"time"
	"unsafe"
	"fmt"
	"os"

	"testing"
)


//
// use testing
func TestNlattrPointer(t *testing.T) {
	ab := NewNlattrBuf(16)
	ab.SetLen(4)
	nla := NlattrPointer(*ab)
	if nla.Len != 4 {
		t.Errorf("len - want: %d, but: %d", 4, nla.Len)
	}

	ab.SetLen(10)
	if nla.Len != 10 {
		t.Errorf("len - want: %d, but: %d", 10, nla.Len)
	}
	ab.SetType(2)
	if nla.Type != 2 {
		t.Errorf("type - want: %d, but: %d", 2, nla.Type)
	}
}

func TestAttrGetType(t *testing.T) {
	ab := NewNlattrBuf(16)
	ab.SetType(2)
	nla := NlattrPointer(*ab)
	if nla.GetType() != 2 {
		t.Errorf("type - want: %d, but: %d", 2, nla.GetType())
	}
	nla.Type = 2 | NLA_F_NESTED
	if nla.GetType() != 2 {
		t.Errorf("type - want: %d, but: %d", 2, nla.GetType())
	}
}

func TestAttrParse(t *testing.T) {
	cb := func(attr *Nlattr, data interface{}) (int, syscall.Errno) {
		// fmt.Printf("attr.Len: %d, data: %d\n", (*attr).Len, data.(int))
		return MNL_CB_OK, 0
	}

	nlh, _ := NewNlmsghdr(4096)
	val := 0x12
	nlh.PutU8(uint16(MNL_TYPE_U8), 0x10)
	nlh.PutU8(uint16(MNL_TYPE_U8), 0x11)
	nlh.PutU8(uint16(MNL_TYPE_U8), 0x12)
	nlh.PutU8(uint16(MNL_TYPE_U8), 0x13)
	ret, err := nlh.Parse(0, cb, val)
	if ret != MNL_CB_OK {
		t.Errorf("type - want: %d, but: %d", MNL_CB_OK, ret)
	}
	if err != nil {
		t.Errorf("type - want: %v, but: %v", nil, err)
	}
}
// used testing
//


var _ = Describe("Attr", func() {
	fmt.Fprintf(os.Stdout, "Hello, attr tester!\n") // to import os, sys for debugging
	var (
		BUFLEN		= 512
		r		*rand.Rand

		// Nlmsghdr
		hbuf		*NlmsghdrBuf
		nlh		*Nlmsghdr
		rand_hbuf	*NlmsghdrBuf
		rand_nlh	*Nlmsghdr

		// Nlattr
		abuf		*NlattrBuf
		nla		*Nlattr
		rand_abuf	*NlattrBuf
		rand_nla	*Nlattr

		// for nlattr validation
		valid_len 	= map[AttrDataType][2]uint16 {
			MNL_TYPE_UNSPEC		: {0, 0},
			MNL_TYPE_U8		: {1, 1},
			MNL_TYPE_U16		: {2, 2},
			MNL_TYPE_U32		: {4, 4},
			MNL_TYPE_U64		: {8, 8},
			MNL_TYPE_STRING		: {64, 64},
			MNL_TYPE_FLAG		: {0, 0},
			MNL_TYPE_MSECS		: {8, 8},
			MNL_TYPE_NESTED		: {32, 32},
			MNL_TYPE_NESTED_COMPAT	: {32, 32},
			MNL_TYPE_NUL_STRING	: {64, 64},
			MNL_TYPE_BINARY		: {64, 64},
			// mnl.TYPE_MAX		: {, ),
		}
		invalid_len	= map[AttrDataType][2]uint16 {
			MNL_TYPE_U8		: {2, 3},
			MNL_TYPE_U16		: {3, 4},
			MNL_TYPE_U32		: {5, 6},
			MNL_TYPE_U64		: {9, 10},
		}
	)

	BeforeEach(func() {
		r = rand.New(rand.NewSource(time.Now().Unix()))
		hbuf = NewNlmsghdrBuf(BUFLEN)
		nlh = NlmsghdrBytes(*hbuf)
		rand_hbuf = NewNlmsghdrBuf(BUFLEN)
		for i := 0; i < BUFLEN; i++ {
			(*(*[]byte)(rand_hbuf))[i] = byte(r.Int() % 256)
		}
		rand_nlh = NlmsghdrBytes(*rand_hbuf)

		abuf = NewNlattrBuf(BUFLEN)
		nla = NlattrPointer(*abuf)
		rand_abuf = NewNlattrBuf(BUFLEN)
		for i := 0; i < BUFLEN; i++ {
			(*(*[]byte)(rand_abuf))[i] = byte(r.Int() % 256)
		}
		rand_nla = NlattrPointer(*rand_abuf)
	})

	Context("NlattrPointer", func() {
		It("should share uint16 4 len", func() {
			abuf.SetLen(4)
			Expect(nla.Len).To(Equal(uint16(4)))
		})
		It("should share uint16 2 type", func() {
			abuf.SetType(2)
			Expect(nla.Type).To(Equal(uint16(2)))
		})
	})

	Context("AttrGetType", func() {
		It("should share uint16 2 len", func() {
			abuf.SetType(2)
			Expect(nla.GetType()).To(Equal(uint16(2)))
		})
		It("should be 2 even NLA_F_NESTED", func() {
			abuf.SetType(2 | NLA_F_NESTED)
			Expect(nla.GetType()).To(Equal(uint16(2)))
		})
	})

	Context("AttrGetLen", func() {
		It("should be 10", func() {
			abuf.SetLen(16)
			Expect(nla.GetLen()).To(Equal(uint16(16)))
		})
	})

	Context("AttrGetPayloadLen", func() {
		It("should be 123 - MNL_ATTR_HDRLEN (not aligned)", func() {
			rand_abuf.SetLen(123)
			Expect(rand_nla.PayloadLen()).To(Equal(uint16(123 - MNL_ATTR_HDRLEN)))
		})
	})

	Context("AttrGetPayload", func() {
		It("should have same pointer", func() {
			abuf_value := *((*[]byte)(abuf))
			Expect(nla.Payload()).To(Equal(unsafe.Pointer(&abuf_value[MNL_ATTR_HDRLEN])))
		})
	})

	Context("AttrGetPayloadBytes", func() {
		It("should have same value", func() {
			rand_abuf_value := *((*[]byte)(rand_abuf))
			alen := uint16(91)
			rand_abuf.SetLen(alen)
			Expect(rand_nla.PayloadBytes()).To(Equal(rand_abuf_value[MNL_ATTR_HDRLEN:alen]))
		})
	})

	Context("AttrOk", func() {
		It("has invalid len and all false", func() {
			abuf.SetLen(3)
			Expect(nla.Ok(3)).To(Equal(false))
			Expect(nla.Ok(4)).To(Equal(false))
			Expect(nla.Ok(5)).To(Equal(false))
		})
		It("have valid len and should fail only shorten", func() {
			abuf.SetLen(8)
			Expect(nla.Ok(7)).To(Equal(false))
			Expect(nla.Ok(8)).To(Equal(true))
			Expect(nla.Ok(9)).To(Equal(true))
		})
	})

	Context("AttrNext", func() {
		It("next should share the same buf", func() {
			rand_abuf.SetLen(256)
			next_buf := (*(*[]byte)(rand_abuf))[256:]
			next_abuf := NlattrBuf(next_buf)
			next_abuf.SetLen(128)

			Expect(rand_nla.Next().Len).To(Equal(uint16(128)))
			next_nla := rand_nla.Next()
			Expect(next_nla.PayloadBytes()).To(Equal(next_buf[MNL_ATTR_HDRLEN:128]))
		})
	})

	Context("AttrTypeValid", func() {
		It("should be valid if type is le param", func() {
			var i uint16
			for i = 0; i < uint16(MNL_TYPE_MAX); i++ {
				abuf.SetType(i)
				ret, err := nla.TypeValid(uint16(MNL_TYPE_MAX))
				Expect(ret).To(Equal(MNL_CB_OK))
				Expect(err).To(BeNil())
			}
		})
		It("should error gt param", func() {
			abuf.SetType(uint16(MNL_TYPE_MAX + 1))
			ret, err := nla.TypeValid(uint16(MNL_TYPE_MAX))
			Expect(ret).To(Equal(MNL_CB_ERROR))
			Expect(err.(syscall.Errno)).To(Equal(syscall.EOPNOTSUPP))
		})
	})

	Context("AttrValidate", func() {
		It("should be invalid because of type", func() {
			ret, err := nla.Validate(MNL_TYPE_MAX)
			Expect(ret).To(Equal(-1))
			Expect(err.(syscall.Errno)).To(Equal(syscall.EINVAL))

			ret, err = nla.Validate2(MNL_TYPE_MAX, 1)
			Expect(ret).To(Equal(-1))
			Expect(err.(syscall.Errno)).To(Equal(syscall.EINVAL))
		})
		It("should be valid", func() {
			for t := range valid_len {
				abuf.SetLen(uint16(MNL_ATTR_HDRLEN) + valid_len[t][0])

				ret, err := nla.Validate( t)
				Expect(ret).To(Equal(0))
				Expect(err).To(BeNil())

				ret, err = nla.Validate2(t, Size_t(valid_len[t][1]))
				Expect(ret).To(Equal(0))
				Expect(err).To(BeNil())
			}
		})
		It("should be invalid by mnl_attr_data_type_len, ERANGE", func() {
			for t := range invalid_len {
				abuf.SetLen(uint16(MNL_ATTR_HDRLEN) + invalid_len[t][0])
				ret, err := nla.Validate(t)
				Expect(ret).To(Equal(-1))
				Expect(err.(syscall.Errno)).To(Equal(syscall.ERANGE))
			}
		})
		It("should be invalid MNL_TYPE_FLAG", func() {
			abuf.SetLen(uint16(MNL_ATTR_HDRLEN + 1))
			ret, err := nla.Validate(MNL_TYPE_FLAG)
			Expect(ret).To(Equal(-1))
			Expect(err.(syscall.Errno)).To(Equal(syscall.ERANGE))
		})
		It("should be invalid MNL_TYPE_NUL_STRING", func() {
			abuf.SetLen(256)
			(*abuf)[abuf.Len() - 1] = 1
			ret, err := nla.Validate(MNL_TYPE_NUL_STRING)
			Expect(ret).To(Equal(-1))
			Expect(err.(syscall.Errno)).To(Equal(syscall.EINVAL))

			abuf.SetLen(uint16(MNL_ATTR_HDRLEN))
			ret, err = nla.Validate(MNL_TYPE_NUL_STRING)
			Expect(ret).To(Equal(-1))
			Expect(err.(syscall.Errno)).To(Equal(syscall.ERANGE))
		})
		It("should be invalid MNL_TYPE_STRING", func() {
			abuf.SetLen(uint16(MNL_ATTR_HDRLEN))
			ret, err := nla.Validate(MNL_TYPE_STRING)
			Expect(ret).To(Equal(-1))
			Expect(err.(syscall.Errno)).To(Equal(syscall.ERANGE))
		})
		It("should be valid MNL_TYPE_NESTED with 0 payload", func() {
			abuf.SetLen(uint16(MNL_ATTR_HDRLEN))
			ret, err := nla.Validate(MNL_TYPE_NESTED)
			Expect(ret).To(Equal(0))
			Expect(err).To(BeNil())
		})
		It("should be invalid MNL_TYPE_NESTED", func() {
			abuf.SetLen(uint16(MNL_ATTR_HDRLEN * 2 - 1))
			ret, err := nla.Validate(MNL_TYPE_NESTED)
			Expect(ret).To(Equal(-1))
			Expect(err.(syscall.Errno)).To(Equal(syscall.ERANGE))
		})
	})

	Context("AttrParse", func() {
		cb := func(val uint8) (func(*Nlattr, interface{}) (int, syscall.Errno)) {
			return func(attr *Nlattr, data interface{}) (int, syscall.Errno) {
				if data != nil {
					return MNL_CB_ERROR, data.(syscall.Errno)
				}
				preval := val
				val += 1
				if preval == attr.U8() && data == nil {
					return MNL_CB_OK, 0
				}
				return MNL_CB_STOP, 0
			}
		}(0x10)

		// XXX: using functions defined here nlmsg.go
		It("should return MNL_CB_OK, nil", func() {
			nlh, _ := PutNewNlmsghdr(512)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x10)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x11)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x12)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x13)
			ret, err := nlh.Parse(0, cb, nil)
			Expect(ret).To(Equal(MNL_CB_OK))
			Expect(err).To(BeNil())

			nlh, _ = PutNewNlmsghdr(512)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x14)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x15)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x16)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x17)
			ret, err = nlh.Parse(0, cb, nil)
			Expect(ret).To(Equal(MNL_CB_OK))
			Expect(err).To(BeNil())
		})
		It("should return MNL_CB_STOP, nil", func() {
			nlh, _ := PutNewNlmsghdr(512)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x18)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x19)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x1a)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x00)
			ret, err := nlh.Parse(0, cb, nil)
			Expect(ret).To(Equal(MNL_CB_STOP))
			Expect(err).To(BeNil())
		})
		It("should return MNL_CB_ERROR, nil", func() {
			nlh, _ := PutNewNlmsghdr(512)
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x00)
			ret, err := nlh.Parse(0, cb, syscall.Errno(3))
			Expect(ret).To(Equal(MNL_CB_ERROR))
			Expect(err).To(Equal(syscall.Errno(3)))
		})
	})

	Describe("Attribute callback", func() {
		cb_f := func(val uint16) (func(*Nlattr, interface{}) (int, syscall.Errno)) {
			return func(attr *Nlattr, data interface{}) (int, syscall.Errno) {
				if !data.(bool) {
					return MNL_CB_STOP, 0
				}
				if val != attr.Type {
					return MNL_CB_ERROR, 123
				}
				val += 1
				return MNL_CB_OK, 0
			}
		}

		Context("AttrParseNested", func() {
			// XXX: using functions defined here nlmsg.go
			nlh, _ := PutNewNlmsghdr(512)
			nested := nlh.NestStart(1)
			nlh.PutU8(uint16(2), 10)
			nlh.PutU8(uint16(3), 20)
			nlh.PutU8(uint16(4), 30)
			nlh.NestEnd(nested)

			It("should return MNL_CB_OK, nil", func() {
				ret, err := nested.ParseNested(cb_f(2), true)
				Expect(ret).To(Equal(MNL_CB_OK))
				Expect(err).To(BeNil())
			})
			It("should return MNL_CB_STOP, nil", func() {
				ret, err := nested.ParseNested(cb_f(2), false)
				Expect(ret).To(Equal(MNL_CB_STOP))
				Expect(err).To(BeNil())
			})
			It("should return MNL_CB_ERROR, nil", func() {
				ret, err := nested.ParseNested(cb_f(0), true)
				Expect(ret).To(Equal(MNL_CB_ERROR))
				Expect(err).To(Equal(syscall.Errno(123)))
			})
		})

		Context("AttrParsePayload", func() {
			// again, using functions defined here, nlmsg.go
			nlh, _ := PutNewNlmsghdr(512)
			nlh.PutU8(uint16(2), 10)
			nlh.PutU8(uint16(3), 20)
			nlh.PutU8(uint16(4), 30)

			It("should return MNL_CB_OK, nil", func() {
				ret, err := AttrParsePayload(nlh.PayloadBytes(), cb_f(2), true)
				Expect(ret).To(Equal(MNL_CB_OK))
				Expect(err).To(BeNil())
			})

			It("should return MNL_CB_STOP, nil", func() {
				ret, err := AttrParsePayload(nlh.PayloadBytes(), cb_f(2), false)
				Expect(ret).To(Equal(MNL_CB_STOP))
				Expect(err).To(BeNil())
			})
			It("should return MNL_CB_ERROR, nil", func() {
				ret, err := AttrParsePayload(nlh.PayloadBytes(), cb_f(0), true)
				Expect(ret).To(Equal(MNL_CB_ERROR))
				Expect(err).To(Equal(syscall.Errno(123)))
			})
		})
	})

	Context("AttrGetU8", func() {
		It("should be 0x11", func() {
			abuf.SetLen(5)
			abuf.SetType(uint16(MNL_TYPE_U8))
			SetUint8(*abuf, 4, 0x11)
			Expect(nla.U8()).To(Equal(uint8(0x11)))
		})
	})

	Context("AttrGetU16", func() {
		It("should be 0x1234", func() {
			abuf.SetLen(6)
			abuf.SetType(uint16(MNL_TYPE_U16))
			SetUint16(*abuf, 4, 0x1234)
			Expect(nla.U16()).To(Equal(uint16(0x1234)))
		})
	})

	Context("AttrGetU32", func() {
		It("should be 0x12345678", func() {
			abuf.SetLen(8)
			abuf.SetType(uint16(MNL_TYPE_U32))
			SetUint32(*abuf, 4, 0x12345678)
			Expect(nla.U32()).To(Equal(uint32(0x12345678)))
		})
	})

	Context("AttrGetU64", func() {
		It("should be 0x123456789abcdef", func() {
			abuf.SetLen(12)
			abuf.SetType(uint16(MNL_TYPE_U64))
			SetUint64(*abuf, 4, 0x123456789abcdef)
			Expect(nla.U64()).To(Equal(uint64(0x123456789abcdef)))
		})
	})

	Context("AttrGetStr", func() {
		It("should be abcDEF", func() {
			abuf.SetLen(11)
			abuf.SetType(uint16(MNL_TYPE_STRING))
			for i, c := range []byte("abcDEF") {
				(*abuf)[i + 4] = c
			}
			(*abuf)[11] = 0
			Expect(nla.Str()).To(Equal("abcDEF"))
		})
	})

	Context("AttrPut", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.Put(3, SizeofNlmsghdr, unsafe.Pointer(rand_nlh))
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + 16", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MNL_NLMSG_HDRLEN))
		})
		It("attr type should be 3", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(3)))
		})
		It("attr len should be 4 + MNL_NLMSG_HDRLEN", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + MNL_NLMSG_HDRLEN)))
		})
		It("attr contents should be equal", func() {
			Expect(([]byte)(_tbuf)[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN + MNL_NLMSG_HDRLEN)]).To(Equal((*(*[]byte)(rand_hbuf))[:MNL_NLMSG_HDRLEN]))
		})
	})

	Context("AttrPutPtr", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutPtr(3, rand_nlh)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + 16", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MNL_NLMSG_HDRLEN))
		})
		It("attr type should be 3", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(3)))
		})
		It("attr len should be 4 + MNL_NLMSG_HDRLEN", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + MNL_NLMSG_HDRLEN)))
		})
		It("attr contents should be equal", func() {
			Expect(([]byte)(_tbuf)[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN + MNL_NLMSG_HDRLEN)]).To(Equal((*(*[]byte)(rand_hbuf))[:MNL_NLMSG_HDRLEN]))
		})
	})

	Context("AttrPutBytes", func() {
		b := []byte{1, 2, 3}
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutBytes(1, b)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(3)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(len(b)))))
		})
		It("attr type should be 1", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(1)))
		})
		It("attr len should be 4 + 3", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + uint32(len(b)))))
		})
		It("attr contents should be equal", func() {
			Expect(([]byte)(_tbuf)[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN) + len(b)]).To(Equal(b))
		})
	})

	Context("AttrPutU8", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutU8(uint16(MNL_TYPE_U8), 7)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(1)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(1))))
		})
		It("attr type should be MNL_TYPE_U8", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U8)))
		})
		It("attr len should be 4 + 1", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 1)))
		})
		It("value should be 7", func() {
			Expect(_tbuf[MNL_ATTR_HDRLEN]).To(Equal(uint8(7)))
		})
	})

	Context("AttrPutU16", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutU16(uint16(MNL_TYPE_U16), 12345)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(2)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(2))))
		})
		It("attr type should be MNL_TYPE_U16", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U16)))
		})
		It("attr len should be 4 + 2", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 2)))
		})
		It("valud should be 12345", func() {
			Expect(GetUint16(_tbuf, uint(MNL_ATTR_HDRLEN))).To(Equal(uint16(12345)))
		})
	})

	Context("AttrPutU32", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutU32(uint16(MNL_TYPE_U32), 0x12345678)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(4)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(4))))
		})
		It("attr type should be MNL_TYPE_U32", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U32)))
		})
		It("attr len should be 4 + 4", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 4)))
		})
		It("valud should be 0x12345678", func() {
			Expect(GetUint32(_tbuf, uint(MNL_ATTR_HDRLEN))).To(Equal(uint32(0x12345678)))
		})
	})

	Context("AttrPutU64", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutU64(uint16(MNL_TYPE_U64), 0x123456789abcdef0)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(8)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(8))))
		})
		It("attr type should be MNL_TYPE_U64", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U64)))
		})
		It("attr len should be 4 + 8", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 8)))
		})
		It("value should be 0x123456789abcdef0", func() {
			Expect(GetUint64(_tbuf, uint(MNL_ATTR_HDRLEN))).To(Equal(uint64(0x123456789abcdef0)))
		})
	})

	Context("AttrPutStr", func() {
		var _tbuf NlattrBuf
		s := "abcdEFGH"
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutStr(uint16(MNL_TYPE_STRING), s)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(len(s))", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(len(s)))))
		})
		It("attr type should be MNL_TYPE_STRING", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_STRING)))
		})
		It("attr len should be 4 + len(s)", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(int(MNL_ATTR_HDRLEN) + len(s))))
		})
		It("value should be equal to s", func() {
			Expect(([]byte)(_tbuf[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN) + len(s)])).To(Equal(([]byte)(s)))
		})
		It("next byte should not be NULL", func() {
			t := [333]byte{} // 256 <= size  <= 512
			for i, _ := range t {
				t[i]= 'a'
			}
			nlh.PutStr(uint16(MNL_TYPE_STRING), string(t[:]))
			Expect(_tbuf[int(MNL_ATTR_HDRLEN) + len(s) + 1]).NotTo(BeZero())
		})
	})

	Context("AttrPutStrz", func() {
		var _tbuf NlattrBuf
		s := "abcdEFGH"
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.PutStrz(uint16(MNL_TYPE_NUL_STRING), s)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(len(s) + 1)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(len(s) + 1))))
		})
		It("attr type should be MNL_TYPE_NUL_STRING", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_NUL_STRING)))
		})
		It("attr len should be 4 + len(s) + 1", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(int(MNL_ATTR_HDRLEN) + len(s) + 1)))
		})
		It("value should be equal to s", func() {
			Expect(([]byte)(_tbuf[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN) + len(s)])).To(Equal(([]byte)(s)))
		})
		It("next byte should not be NULL", func() {
			t := [333]byte{} // 256 <= size  <= 512
			for i, _ := range t {
				t[i]= 'a'
			}
			nlh.PutStr(uint16(MNL_TYPE_STRING), string(t[:]))
			Expect(_tbuf[int(MNL_ATTR_HDRLEN) + len(s) + 1]).To(BeZero())
		})
	})

	Context("AttrNestStart", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			nlh.NestStart(1)
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("header len should be MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN", func() {
			Expect(hbuf.Len()).To(Equal(uint32(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN)))
		})
		It("attribute type has NLA_F_NESTED", func() {
			Expect(_tbuf.Type() & NLA_F_NESTED).To(Equal(uint16(NLA_F_NESTED)))
		})
		It("attribute type has set value", func() {
			Expect(_tbuf.Type() & 1).To(Equal(uint16(1)))
		})
	})

	Context("AttrPutCheck", func() {
		b := []byte{1, 2, 3}
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			if nlh.PutCheck(Size_t(len(*hbuf)), 1, 3, unsafe.Pointer(&b[0])) == false {
				panic("invalid test assumption")
			}
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(3)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(len(b)))))
		})
		It("attr type should be 1", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(1)))
		})
		It("attr len should be 4 + 3", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + uint32(len(b)))))
		})
		It("attr contents should be equal", func() {
			Expect(([]byte)(_tbuf)[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN) + len(b)]).To(Equal(b))
		})
		It("buflen just MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN should return false and nothing has changed", func() {
			_hbuf := NewNlmsghdrBuf(int(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN))
			_nlh := NlmsghdrBytes(*_hbuf)
			_nlh.PutHeader()
			pb := make([]byte, len(*_hbuf))
			copy(pb, *_hbuf)
			Expect(nlh.PutCheck(Size_t(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN), 1, 3, unsafe.Pointer(&b[0]))).To(Equal(false))
			Expect(*(*[]byte)(_hbuf)).To(BeEquivalentTo(pb))
		})
	})

	Context("AttrPutU8Check", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			if nlh.PutU8Check(Size_t(len(*hbuf)), uint16(MNL_TYPE_U8), 7) == false {
				panic("invalid test assumption")
			}
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(1)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(1))))
		})
		It("attr type should be MNL_TYPE_U8", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U8)))
		})
		It("attr len should be 4 + 1", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 1)))
		})
		It("value should be 7", func() {
			Expect(_tbuf[MNL_ATTR_HDRLEN]).To(Equal(uint8(7)))
		})
		It("buflen just MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN should return false and nothing has changed", func() {
			_hbuf := NewNlmsghdrBuf(int(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN))
			_nlh := NlmsghdrBytes(*_hbuf)
			_nlh.PutHeader()
			pb := make([]byte, len(*_hbuf))
			copy(pb, *_hbuf)
			Expect(nlh.PutU8Check(Size_t(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN), uint16(MNL_TYPE_U8), 7)).To(Equal(false))
			Expect(*(*[]byte)(_hbuf)).To(BeEquivalentTo(pb))
		})
	})

	Context("AttrPutU16Check", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			if nlh.PutU16Check(Size_t(len(*hbuf)), uint16(MNL_TYPE_U16), 12345) == false {
				panic("invalid test assumption")
			}
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(2)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(2))))
		})
		It("attr type should be MNL_TYPE_U16", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U16)))
		})
		It("attr len should be 4 + 2", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 2)))
		})
		It("valud should be 12345", func() {
			Expect(GetUint16(_tbuf, uint(MNL_ATTR_HDRLEN))).To(Equal(uint16(12345)))
		})
		It("buflen just MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN should return false and nothing has changed", func() {
			_hbuf := NewNlmsghdrBuf(int(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN))
			_nlh := NlmsghdrBytes(*_hbuf)
			_nlh.PutHeader()
			pb := make([]byte, len(*_hbuf))
			copy(pb, *_hbuf)
			Expect(nlh.PutU16Check(Size_t(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN), uint16(MNL_TYPE_U16), 12345)).To(Equal(false))
			Expect(*(*[]byte)(_hbuf)).To(BeEquivalentTo(pb))
		})
	})

	Context("AttrPutU32Check", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			if nlh.PutU32Check(Size_t(len(*hbuf)), uint16(MNL_TYPE_U32), 0x12345678) == false {
				panic("invalid test assumption")
			}
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(4)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(4))))
		})
		It("attr type should be MNL_TYPE_U32", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U32)))
		})
		It("attr len should be 4 + 4", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 4)))
		})
		It("valud should be 0x12345678", func() {
			Expect(GetUint32(_tbuf, uint(MNL_ATTR_HDRLEN))).To(Equal(uint32(0x12345678)))
		})
		It("buflen just MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN should return false and nothing has changed", func() {
			_hbuf := NewNlmsghdrBuf(int(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN))
			_nlh := NlmsghdrBytes(*_hbuf)
			_nlh.PutHeader()
			pb := make([]byte, len(*_hbuf))
			copy(pb, *_hbuf)
			Expect(nlh.PutU32Check(Size_t(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN), uint16(MNL_TYPE_U32), 0x12345678)).To(Equal(false))
			Expect(*(*[]byte)(_hbuf)).To(BeEquivalentTo(pb))
		})
	})

	Context("AttrPutU64Check", func() {
		var _tbuf NlattrBuf
		BeforeEach(func() {
			nlh.PutHeader()
			if nlh.PutU64Check(Size_t(len(*hbuf)), uint16(MNL_TYPE_U64), 0x123456789abcdef0) == false {
				panic("invalid test assumption")
			}
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(8)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(8))))
		})
		It("attr type should be MNL_TYPE_U64", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_U64)))
		})
		It("attr len should be 4 + 8", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(MNL_ATTR_HDRLEN + 8)))
		})
		It("value should be 0x123456789abcdef0", func() {
			Expect(GetUint64(_tbuf, uint(MNL_ATTR_HDRLEN))).To(Equal(uint64(0x123456789abcdef0)))
		})
		It("buflen just MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN should return false and nothing has changed", func() {
			_hbuf := NewNlmsghdrBuf(int(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN))
			_nlh := NlmsghdrBytes(*_hbuf)
			_nlh.PutHeader()
			pb := make([]byte, len(*_hbuf))
			copy(pb, *_hbuf)
			Expect(nlh.PutU64Check(Size_t(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN), uint16(MNL_TYPE_U64), 0x123456789abcdef0)).To(Equal(false))
			Expect(*(*[]byte)(_hbuf)).To(BeEquivalentTo(pb))
		})
	})

	Context("AttrPutStrCheck", func() {
		var _tbuf NlattrBuf
		s := "abcdEFGH"
		BeforeEach(func() {
			nlh.PutHeader()
			if nlh.PutStrCheck(Size_t(len(*hbuf)), uint16(MNL_TYPE_STRING), s) == false {
				panic("invalid test assumption")
			}
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(len(s))", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(len(s)))))
		})
		It("attr type should be MNL_TYPE_STRING", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_STRING)))
		})
		It("attr len should be 4 + len(s)", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(int(MNL_ATTR_HDRLEN) + len(s))))
		})
		It("value should be equal to s", func() {
			Expect(([]byte)(_tbuf[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN) + len(s)])).To(Equal(([]byte)(s)))
		})
		It("next byte should not be NULL", func() {
			t := [333]byte{} // 256 <= size  <= 512
			for i, _ := range t {
				t[i]= 'a'
			}
			nlh.PutStr(uint16(MNL_TYPE_STRING), string(t[:]))
			Expect(_tbuf[int(MNL_ATTR_HDRLEN) + len(s) + 1]).NotTo(BeZero())
		})
		It("buflen just MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN should return false and nothing has changed", func() {
			_hbuf := NewNlmsghdrBuf(int(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN))
			_nlh := NlmsghdrBytes(*_hbuf)
			_nlh.PutHeader()
			pb := make([]byte, len(*_hbuf))
			copy(pb, *_hbuf)
			Expect(nlh.PutStrCheck(Size_t(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN), uint16(MNL_TYPE_STRING), s)).To(Equal(false))
			Expect(*(*[]byte)(_hbuf)).To(BeEquivalentTo(pb))
		})
	})

	Context("AttrPutStrzCheck", func() {
		var _tbuf NlattrBuf
		s := "abcdEFGH"
		BeforeEach(func() {
			nlh.PutHeader()
			if nlh.PutStrzCheck(Size_t(len(*hbuf)), uint16(MNL_TYPE_NUL_STRING), s) == false {
				panic("invalid test assumption")
			}
			_tbuf = NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
		})
		It("nlh len should be 16 + 4 + align(len(s) + 1)", func() {
			Expect(nlh.Len).To(Equal(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN + MnlAlign(uint32(len(s) + 1))))
		})
		It("attr type should be MNL_TYPE_NUL_STRING", func() {
			Expect(_tbuf.Type()).To(Equal(uint16(MNL_TYPE_NUL_STRING)))
		})
		It("attr len should be 4 + len(s) + 1", func() {
			Expect(_tbuf.Len()).To(Equal(uint16(int(MNL_ATTR_HDRLEN) + len(s) + 1)))
		})
		It("value should be equal to s", func() {
			Expect(([]byte)(_tbuf[MNL_ATTR_HDRLEN:int(MNL_ATTR_HDRLEN) + len(s)])).To(Equal(([]byte)(s)))
		})
		It("next byte should not be NULL", func() {
			t := [333]byte{} // 256 <= size  <= 512
			for i, _ := range t {
				t[i]= 'a'
			}
			nlh.PutStr(uint16(MNL_TYPE_STRING), string(t[:]))
			Expect(_tbuf[int(MNL_ATTR_HDRLEN) + len(s) + 1]).To(BeZero())
		})
		It("buflen just MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN should return false and nothing has changed", func() {
			_hbuf := NewNlmsghdrBuf(int(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN))
			_nlh := NlmsghdrBytes(*_hbuf)
			_nlh.PutHeader()
			pb := make([]byte, len(*_hbuf))
			copy(pb, *_hbuf)
			Expect(nlh.PutStrzCheck(Size_t(MNL_NLMSG_HDRLEN + MNL_ATTR_HDRLEN), uint16(MNL_TYPE_NUL_STRING), s)).To(Equal(false))
			Expect(*(*[]byte)(_hbuf)).To(BeEquivalentTo(pb))
		})
	})

	Context("AttrNestEnd", func() {
		It("attr len should be updated", func() {
			nlh.PutHeader()
			_nla := nlh.NestStart(1)				// payload len: aligned
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x12)			// 1: 4
			nlh.PutU16(uint16(MNL_TYPE_U16), 0x3456)		// 2: 4
			nlh.PutU32(uint16(MNL_TYPE_U32), 0x3456789a)		// 4: 4
			nlh.PutU64(uint16(MNL_TYPE_U64), 0xbcdef0123456789a)	// 8: 8
			nlh.PutStr(uint16(MNL_TYPE_STRING), "bcdef")		// 5: 8
			nlh.PutStrz(uint16(MNL_TYPE_NUL_STRING), "01234567")	// 9: 12 = 40 + MNL_ATTR_HDRLEN * 7
			nlh.NestEnd(_nla)
			_abuf := NlattrBuf((*hbuf)[MNL_NLMSG_HDRLEN:])
			Expect(_abuf.Len()).To(Equal(uint16(68)))
		})
	})

	Context("AttrNestCancel", func() {
		It("nlh len should be updated", func() {
			nlh.PutHeader()
			_nla := nlh.NestStart(1)				// payload len: aligned
			nlh.PutU8(uint16(MNL_TYPE_U8), 0x12)			// 1: 4
			nlh.PutU16(uint16(MNL_TYPE_U16), 0x3456)		// 2: 4
			nlh.PutU32(uint16(MNL_TYPE_U32), 0x3456789a)		// 4: 4
			nlh.PutU64(uint16(MNL_TYPE_U64), 0xbcdef0123456789a)	// 8: 8
			nlh.PutStr(uint16(MNL_TYPE_STRING), "bcdef")		// 5: 8
			nlh.PutStrz(uint16(MNL_TYPE_NUL_STRING), "01234567")	// 9: 12 = 40 + MNL_ATTR_HDRLEN * 7
			Expect(hbuf.Len()).To(Equal(uint32(MNL_NLMSG_HDRLEN + 68)))
			nlh.NestCancel(_nla)
			Expect(hbuf.Len()).To(Equal(uint32(MNL_NLMSG_HDRLEN)))
		})
	})

	Context("Attributes", func() {
		It("has valid 4 attrs", func() {
			nlh.PutHeader()
			nlh.PutU8(uint16(0), 0x10)
			nlh.PutU8(uint16(1), 0x11)
			nlh.PutU8(uint16(2), 0x12)
			nlh.PutU8(uint16(3), 0x13)
			i := 0
			for attr := range(nlh.Attributes(0)) {
				Expect(attr.Type).To(Equal(uint16(i)))
				Expect(attr.U8()).To(Equal(uint8(0x10 + i)))
				i += 1
			}
		})
	})
	Context("Nesteds", func() {
		It("nested 4 valid attributes", func() {
			// XXX: using functions defined here nlmsg.go
			nlh.PutHeader()
			nested := nlh.NestStart(1)
			nlh.PutU8(uint16(0), 0)
			nlh.PutU8(uint16(1), 10)
			nlh.PutU8(uint16(2), 20)
			nlh.PutU8(uint16(3), 30)
			nlh.NestEnd(nested)
			i := 0
			for attr := range(nested.Nesteds()) {
				Expect(attr.Type).To(Equal(uint16(i)))
				Expect(attr.U8()).To(Equal(uint8(10 * i)))
				i += 1
			}
		})
	})

	Context("PayloadAttributes", func() {
		It("has 4 valid attributes", func() {
			nlh.PutHeader()
			nlh.PutU8(uint16(1), 10)
			nlh.PutU8(uint16(2), 20)
			nlh.PutU8(uint16(3), 30)
			nlh.PutU8(uint16(4), 40)

			i := 0
			for attr := range(PayloadAttributes(*hbuf)) {
				Expect(attr.Type).To(Equal(uint16(i)))
				Expect(attr.U8()).To(Equal(uint8(10 * i)))
				i += 1
			}
		})
	})
})
