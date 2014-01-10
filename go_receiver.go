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

/*
 * attr.go
 */
func (attr *Nlattr) GetType() uint16 { return AttrGetType(attr) }
func (attr *Nlattr) GetLen() uint16 { return AttrGetLen(attr) }
func (attr *Nlattr) PayloadLen() uint16 { return AttrGetPayloadLen(attr) }
func (attr *Nlattr) Payload() unsafe.Pointer { return AttrGetPayload(attr) }
func (attr *Nlattr) PayloadBytes() []byte { return AttrGetPayloadBytes(attr) } // added
func (attr *Nlattr) Ok(size int) bool { return AttrOk(attr, size) }
func (attr *Nlattr) Next() *Nlattr { return AttrNext(attr) }
func (attr *Nlattr) TypeValid(max uint16) (int, error) { return AttrTypeValid(attr, max) }
func (attr *Nlattr) Validate(data_type AttrDataType) (int, error) { return AttrValidate(attr, data_type) }
func (attr *Nlattr) Validate2(data_type AttrDataType, exp_len Size_t) (int, error) { return AttrValidate2(attr, data_type, exp_len) } 
func (nlh *Nlmsghdr) Parse(offset Size_t, cb MnlAttrCb, data interface{}) (int, error) { return AttrParse(nlh, offset, cb, data) }
func (attr *Nlattr) ParseNested(cb MnlAttrCb, data interface{}) (int, error) { return AttrParseNested(attr, cb, data) }
func (attr *Nlattr) U8() uint8 { return AttrGetU8(attr) }
func (attr *Nlattr) U16() uint16 { return AttrGetU16(attr) }
func (attr *Nlattr) U32() uint32 { return AttrGetU32(attr) }
func (attr *Nlattr) U64() uint64 { return AttrGetU64(attr) }
func (attr *Nlattr) Str() string { return AttrGetStr(attr) }
func (nlh *Nlmsghdr) Put(attr_type uint16, size Size_t, p unsafe.Pointer) { AttrPut(nlh, attr_type, size, p) }
func (nlh *Nlmsghdr) PutData(attr_type uint16, data interface{}) { AttrPutData(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutBytes(attr_type uint16, data []byte) { AttrPutBytes(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU8(attr_type uint16, data uint8) { AttrPutU8(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU16(attr_type uint16, data uint16) { AttrPutU16(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU32(attr_type uint16, data uint32) { AttrPutU32(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutU64(attr_type uint16, data uint64) { AttrPutU64(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutStr(attr_type uint16, data string) { AttrPutStr(nlh, attr_type, data) }
func (nlh *Nlmsghdr) PutStrz(attr_type uint16, data string) { AttrPutStrz(nlh, attr_type, data) }
func (nlh *Nlmsghdr) NestStart(attr_type uint16) *Nlattr { return AttrNestStart(nlh, attr_type) }
func (nlh *Nlmsghdr) PutCheck(buflen Size_t, attr_type uint16, data []byte) bool { return AttrPutCheck(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU8Check(buflen Size_t, attr_type uint16, data uint8) bool { return AttrPutU8Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU16Check(buflen Size_t, attr_type uint16, data uint16) bool { return AttrPutU16Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU32Check(buflen Size_t, attr_type uint16, data uint32) bool { return AttrPutU32Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutU64Check(buflen Size_t, attr_type uint16, data uint64) bool { return AttrPutU64Check(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutStrCheck(buflen Size_t, attr_type uint16, data string) bool { return AttrPutStrCheck(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) PutStrzCheck(buflen Size_t, attr_type uint16, data string) bool { return AttrPutStrzCheck(nlh, buflen, attr_type, data) }
func (nlh *Nlmsghdr) NestEnd(start *Nlattr) { AttrNestEnd(nlh, start) }
func (nlh *Nlmsghdr) NestCancel(start *Nlattr) { AttrNestCancel(nlh, start) }

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

/*
 * nlmsg.go
 */
func (nlh *Nlmsghdr) PayloadLen() Size_t { return NlmsgGetPayloadLen(nlh) }
func (nlh *Nlmsghdr) PutExtraHeader(size Size_t) unsafe.Pointer { return NlmsgPutExtraHeader(nlh, size) }
func (nlh *Nlmsghdr) Payload() unsafe.Pointer { return NlmsgGetPayload(nlh) }
func (nlh *Nlmsghdr) PayloadBytes() []byte { return NlmsgGetPayloadBytes(nlh) }
func (nlh *Nlmsghdr) PayloadOffset(offset Size_t) unsafe.Pointer { return NlmsgGetPayloadOffset(nlh, offset) }
func (nlh *Nlmsghdr) PayloadOffsetBytes(offset Size_t) []byte { return NlmsgGetPayloadOffsetBytes(nlh, offset) }
func (nlh *Nlmsghdr) Ok(size int) bool { return NlmsgOk(nlh, size) }
func (nlh *Nlmsghdr) Next(size int) (*Nlmsghdr, int) { return NlmsgNext(nlh, size) }
func (nlh *Nlmsghdr) PayloadTail() unsafe.Pointer { return NlmsgGetPayloadTail(nlh) }
func (nlh *Nlmsghdr) SeqOk(seq uint32) bool { return NlmsgSeqOk(nlh, seq) }
func (nlh *Nlmsghdr) PortidOk(portid uint32) bool { return NlmsgPortidOk(nlh, portid) }
func (nlh *Nlmsghdr) Fprint(fd *os.File, extra_header_size Size_t) { NlmsgFprintNlmsg(fd, nlh, extra_header_size) }
func (b *MnlNlmsgBatch) Stop() { NlmsgBatchStop(b) }
func (b *MnlNlmsgBatch) Next() bool { return NlmsgBatchNext(b) }
func (b *MnlNlmsgBatch) Reset() { NlmsgBatchReset(b) }
func (b *MnlNlmsgBatch) Size() Size_t { return NlmsgBatchSize(b) }
func (b *MnlNlmsgBatch) Head() unsafe.Pointer { return NlmsgBatchHead(b) }
func (b *MnlNlmsgBatch) HeadBytes() []byte { return NlmsgBatchHeadBytes(b) }
func (b *MnlNlmsgBatch) Current() unsafe.Pointer { return NlmsgBatchCurrent(b) }
func (b *MnlNlmsgBatch) IsEmpty() bool { return NlmsgBatchIsEmpty(b) }

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

/*
 * socket.go
 */
func (nl *MnlSocket) Fd() int { return SocketGetFd(nl) }
func (nl *MnlSocket) Portid() uint32 { return SocketGetPortid(nl) }
func (nl *MnlSocket) Bind(groups uint, pid Pid_t) error { return SocketBind(nl, groups, pid) }
func (nl *MnlSocket) Sendto(buf []byte) (Ssize_t, error) { return SocketSendto(nl, buf) }
func (nl *MnlSocket) SendNlmsg(nlh *Nlmsghdr) (Ssize_t, error) { return SocketSendNlmsg(nl, nlh) }
func (nl *MnlSocket) Recvfrom(buf []byte) (Ssize_t, error) { return SocketRecvfrom(nl, buf) }
func (nl *MnlSocket) Close() error { return SocketClose(nl) }
func (nl *MnlSocket) Setsockopt(t int, v unsafe.Pointer, l Socklen_t) error { return SocketSetsockopt(nl, t, v, l) }
func (nl *MnlSocket) SetsockoptBytes(optype int, buf []byte) error { return SocketSetsockoptBytes(nl, optype, buf) }
func (nl *MnlSocket) SetsockoptByte(optype int, v byte) error { return SocketSetsockoptByte(nl, optype, v) }
func (nl *MnlSocket) SetsockoptCint(optype int, v int) error { return SocketSetsockoptCint(nl, optype, v) }
func (nl *MnlSocket) Sockopt(optype int, size Socklen_t) ([]byte, error) { return SocketGetsockopt(nl, optype, size) }

