package cgolmnl

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
*/
import "C"
import (
	"unsafe"
	"errors"
	"os"
)

// attr.go
func (attr *Nlattr) GetType() uint16 { return attrGetType(attr) }
func (attr *Nlattr) GetLen() uint16 { return attrGetLen(attr) }
func (attr *Nlattr) PayloadLen() uint16 { return attrGetPayloadLen(attr) }
func (attr *Nlattr) Payload() unsafe.Pointer { return attrGetPayload(attr) }
func (attr *Nlattr) PayloadBytes() []byte { return attrGetPayloadBytes(attr) } // added
func (attr *Nlattr) Ok(size int) bool { return attrOk(attr, size) }
func (attr *Nlattr) Next() *Nlattr { return attrNext(attr) }
func (attr *Nlattr) TypeValid(max uint16) (int, error) { return attrTypeValid(attr, max) }
func (attr *Nlattr) Validate(data_type AttrDataType) (int, error) { return attrValidate(attr, data_type) }
func (attr *Nlattr) Validate2(data_type AttrDataType, exp_len Size_t) (int, error) { return attrValidate2(attr, data_type, exp_len) } 
func (nlh *Nlmsghdr) Parse(offset Size_t, cb MnlAttrCb, data interface{}) (int, error) { return attrParse(nlh, offset, cb, data) }
func (attr *Nlattr) ParseNested(cb MnlAttrCb, data interface{}) (int, error) { return attrParseNested(attr, cb, data) }
func (attr *Nlattr) U8() uint8 { return attrGetU8(attr) }
func (attr *Nlattr) U16() uint16 { return attrGetU16(attr) }
func (attr *Nlattr) U32() uint32 { return attrGetU32(attr) }
func (attr *Nlattr) U64() uint64 { return attrGetU64(attr) }
func (attr *Nlattr) Str() string { return attrGetStr(attr) }
func (nlh *Nlmsghdr) Put(attr_type uint16, size Size_t, p unsafe.Pointer) { attrPut(nlh, attr_type, size, p) }
func (nlh *Nlmsghdr) PutPtr(attr_type uint16, data interface{}) { attrPutPtr(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutBytes(attr_type uint16, data []byte) { attrPutBytes(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU8(attr_type uint16, data uint8) { attrPutU8(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU16(attr_type uint16, data uint16) { attrPutU16(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU32(attr_type uint16, data uint32) { attrPutU32(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU64(attr_type uint16, data uint64) { attrPutU64(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutStr(attr_type uint16, data string) { attrPutStr(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutStrz(attr_type uint16, data string) { attrPutStrz(nlh, attr_type, data) }
func (nlh *Nlmsghdr) NestStart(attr_type uint16) *Nlattr { return attrNestStart(nlh, attr_type) }
func (nlh *Nlmsghdr) PutCheck(buflen Size_t, attr_type uint16, data []byte) bool { return attrPutCheck(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU8Check(buflen Size_t, attr_type uint16, data uint8) bool { return attrPutU8Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU16Check(buflen Size_t, attr_type uint16, data uint16) bool { return attrPutU16Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU32Check(buflen Size_t, attr_type uint16, data uint32) bool { return attrPutU32Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU64Check(buflen Size_t, attr_type uint16, data uint64) bool { return attrPutU64Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutStrCheck(buflen Size_t, attr_type uint16, data string) bool { return attrPutStrCheck(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutStrzCheck(buflen Size_t, attr_type uint16, data string) bool { return attrPutStrzCheck(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) NestEnd(start *Nlattr) { attrNestEnd(nlh, start) }
func (nlh *Nlmsghdr) NestCancel(start *Nlattr) { attrNestCancel(nlh, start) }

// mnl_attr_for_each macro in libmnl.h
func (nlh *Nlmsghdr) Attributes(offset Size_t) <-chan *Nlattr {
	c := make(chan *Nlattr)
	go func() {
		attr := (*Nlattr)(nlh.PayloadOffset(offset))
		for attr.Ok(int(uintptr(nlh.PayloadTail()) - uintptr(unsafe.Pointer(attr)))) {
			c <- attr
			attr = attr.Next()
		}
		close(c)
	}()
	return c
}

// mnl_attr_for_each_nested macro in libmnl.h
func (nest *Nlattr) Nesteds() <-chan *Nlattr {
	c := make(chan *Nlattr)
	go func() {
		attr := (*Nlattr)(nest.Payload())
		for attr.Ok(int(uintptr(nest.Payload()) + uintptr(nest.PayloadLen()) - uintptr(unsafe.Pointer(attr)))) {
			c <- attr
			attr = attr.Next()
		}
		close(c)
	}()
	return c
}

// mnl_attr_for_each_payload macro in libmnl.h
func PayloadAttributes(payload []byte) <-chan *Nlattr {
	p := unsafe.Pointer(&payload[0])
	attr := (*Nlattr)(p)
	c := make(chan *Nlattr)

	go func() {
		for attr.Ok(int(uintptr(p) + uintptr(len(payload)) - uintptr(p))) {
			c <- attr
			attr = attr.Next()
		}
		close(c)
	}()

	return c
}

// helper function
func NewNlattr(size int) (*Nlattr, error) {
	if size < SizeofNlattr {
		return nil, errors.New("too short size")
	}
	b := make([]byte, size)
	return (*Nlattr)(unsafe.Pointer(&b[0])), nil
}

func NlattrPointer(b []byte) *Nlattr {
	// ???: check buf len
	//      len(b) >= SizeofNlattr
	//      nla.len <= len(b)
	return (*Nlattr)(unsafe.Pointer(&b[0]))
}

// nlmsg.go
func (nlh *Nlmsghdr) PayloadLen() Size_t { return nlmsgGetPayloadLen(nlh) }
func (nlh *Nlmsghdr) PutExtraHeader(size Size_t) unsafe.Pointer { return nlmsgPutExtraHeader(nlh, size) }
func (nlh *Nlmsghdr) Payload() unsafe.Pointer { return nlmsgGetPayload(nlh) }
func (nlh *Nlmsghdr) PayloadBytes() []byte { return nlmsgGetPayloadBytes(nlh) }
func (nlh *Nlmsghdr) PayloadOffset(offset Size_t) unsafe.Pointer { return nlmsgGetPayloadOffset(nlh, offset) }
func (nlh *Nlmsghdr) PayloadOffsetBytes(offset Size_t) []byte { return nlmsgGetPayloadOffsetBytes(nlh, offset) }
func (nlh *Nlmsghdr) Ok(size int) bool { return nlmsgOk(nlh, size) }
func (nlh *Nlmsghdr) Next(size int) (*Nlmsghdr, int) { return nlmsgNext(nlh, size) }
func (nlh *Nlmsghdr) PayloadTail() unsafe.Pointer { return nlmsgGetPayloadTail(nlh) }
func (nlh *Nlmsghdr) SeqOk(seq uint32) bool { return nlmsgSeqOk(nlh, seq) }
func (nlh *Nlmsghdr) PortidOk(portid uint32) bool { return nlmsgPortidOk(nlh, portid) }
func (nlh *Nlmsghdr) Fprint(fd *os.File, extra_header_size Size_t) { nlmsgFprintNlmsg(fd, nlh, extra_header_size) }
func NewNlmsgBatch(buf []byte, limit Size_t) (*NlmsgBatch, error) { return nlmsgBatchStart(buf, limit) }
func (b *NlmsgBatch) Stop() { nlmsgBatchStop(b) }
func (b *NlmsgBatch) Next() bool { return nlmsgBatchNext(b) }
func (b *NlmsgBatch) Reset() { nlmsgBatchReset(b) }
func (b *NlmsgBatch) Size() Size_t { return nlmsgBatchSize(b) }
func (b *NlmsgBatch) Head() unsafe.Pointer { return nlmsgBatchHead(b) }
func (b *NlmsgBatch) HeadBytes() []byte { return nlmsgBatchHeadBytes(b) }
func (b *NlmsgBatch) Current() unsafe.Pointer { return nlmsgBatchCurrent(b) }
func (b *NlmsgBatch) IsEmpty() bool { return nlmsgBatchIsEmpty(b) }

// helper function
func NewNlmsghdr(size int) (*Nlmsghdr, error) {
	if size < int(MNL_NLMSG_HDRLEN) {
		return nil, errors.New("too short size")
	}
	b := make([]byte, size)
	return (*Nlmsghdr)(unsafe.Pointer(&b[0])), nil
}

func NlmsghdrBytes(b []byte) *Nlmsghdr {
	return (*Nlmsghdr)(unsafe.Pointer(&b[0]))
}

func (nlh *Nlmsghdr) PutHeader() {
	// ???: check buf len
	//      len(b) >= SizeofNlmsghdr
	//      nlh.len <= len(buf)
	C.mnl_nlmsg_put_header(unsafe.Pointer(nlh))
}

func PutNewNlmsghdr(size int) (*Nlmsghdr, error) {
	nlh, err := NewNlmsghdr(size)
	if err != nil {
		return nil, err
	}
	nlh.PutHeader()
	return nlh, nil
}

// socket.go
func NewSocket(bus int) (*Socket, error) { return socketOpen(bus) }
func (nl *Socket) Fd() int { return socketGetFd(nl) }
func (nl *Socket) Portid() uint32 { return socketGetPortid(nl) }
func (nl *Socket) Bind(groups uint, pid Pid_t) error { return socketBind(nl, groups, pid) }
func (nl *Socket) Sendto(buf []byte) (Ssize_t, error) { return socketSendto(nl, buf) }
func (nl *Socket) SendNlmsg(nlh *Nlmsghdr) (Ssize_t, error) { return socketSendNlmsg(nl, nlh) }
func (nl *Socket) Recvfrom(buf []byte) (Ssize_t, error) { return socketRecvfrom(nl, buf) }
func (nl *Socket) Close() error { return socketClose(nl) }
func (nl *Socket) Setsockopt(t int, v unsafe.Pointer, l Socklen_t) error { return socketSetsockopt(nl, t, v, l) }
func (nl *Socket) SetsockoptBytes(optype int, buf []byte) error { return socketSetsockoptBytes(nl, optype, buf) }
func (nl *Socket) SetsockoptByte(optype int, v byte) error { return socketSetsockoptByte(nl, optype, v) }
func (nl *Socket) SetsockoptCint(optype int, v int) error { return socketSetsockoptCint(nl, optype, v) }
func (nl *Socket) Sockopt(optype int, size Socklen_t) ([]byte, error) { return socketGetsockopt(nl, optype, size) }

