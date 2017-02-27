package cgolmnl

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
*/
import "C"
import (
	"reflect"
	"errors"
	"os"
	"syscall"
	"unsafe"
)

// get type of netlink attribute
//
// This function returns the attribute type.
func (attr *Nlattr) GetType() uint16 {
	return attrGetType(attr)
}

// get length of netlink attribute
//
// This function returns the attribute length that is the attribute header
// plus the attribute payload.
func (attr *Nlattr) GetLen() uint16 {
	return attrGetLen(attr)
}

// get the attribute payload-value length
//
// This function returns the attribute payload-value length.
func (attr *Nlattr) PayloadLen() uint16 {
	return attrGetPayloadLen(attr)
}

// get pointer to the attribute payload
//
// This function return a pointer to the attribute payload.
func (attr *Nlattr) Payload() unsafe.Pointer {
	return attrGetPayload(attr)
}

// get pointer to the attribute payload as []byte
//
// This function wraps Payload().
func (attr *Nlattr) PayloadBytes() []byte {
	return attrGetPayloadBytes(attr)
}

// check if there is room for an attribute in a buffer
//
// This function is used to check that a buffer, which is supposed to contain
// an attribute, has enough room for the attribute that it stores, i.e. this
// function can be used to verify that an attribute is neither malformed nor
// truncated.
//
// This function does not set errno in case of error since it is intended
// for iterations. Thus, it returns true on success and false on error.
//
// The size parameter may be negative in the case of malformed messages during
// attribute iteration, that is why we use a signed integer.
func (attr *Nlattr) Ok(size int) bool {
	return attrOk(attr, size)
}

// get the next attribute in the payload of a netlink message
//
// This function returns a pointer to the next attribute after the one passed
// as parameter. You have to use Ok() to ensure that the next
// attribute is valid.
func (attr *Nlattr) Next() *Nlattr {
	return attrNext(attr)
}

// check if the attribute type is valid
//
// This function allows to check if the attribute type is higher than the
// maximum supported type. If the attribute type is invalid, this function
// returns error.
//
// Strict attribute checking in user-space is not a good idea since you may
// run an old application with a newer kernel that supports new attributes.
// This leads to backward compatibility breakages in user-space. Better check
// if you support an attribute, if not, skip it.
func (attr *Nlattr) TypeValid(max uint16) error {
	return attrTypeValid(attr, max)
}

// validate netlink attribute (simplified version)
//
// The validation is based on the data type. Specifically, it checks that
// integers (u8, u16, u32 and u64) have enough room for them. This function
// returns error.
func (attr *Nlattr) Validate(data_type AttrDataType) error {
	return attrValidate(attr, data_type)
}

// validate netlink attribute (extended version)
//
// This function allows to perform a more accurate validation for attributes
// whose size is variable. If the size of the attribute is not what we expect,
// this functions returns -1 and errno is explicitly set.
func (attr *Nlattr) Validate2(data_type AttrDataType, exp_len Size_t) error {
	return attrValidate2(attr, data_type, exp_len)
}

// parse attributes
//
// This function allows to iterate over the sequence of attributes that compose
// the Netlink message. You can then put the attribute in an array as it
// usually happens at this stage or you can use any other data structure (such
// as lists or trees).
//
// This function propagates the return value of the callback, which can be
// MNL_CB_ERROR, MNL_CB_OK or MNL_CB_STOP.
func (nlh *Nlmsg) Parse(offset Size_t, cb MnlAttrCb, data interface{}) (int, error) {
	return attrParse(nlh, offset, cb, data)
}

// parse attributes inside a nest
//
// This function allows to iterate over the sequence of attributes that compose
// the Netlink message. You can then put the attribute in an array as it
// usually happens at this stage or you can use any other data structure (such
// as lists or trees).
//
// This function propagates the return value of the callback, which can be
// MNL_CB_ERROR, MNL_CB_OK or MNL_CB_STOP.
func (attr *Nlattr) ParseNested(cb MnlAttrCb, data interface{}) (int, error) {
	return attrParseNested(attr, cb, data)
}

// returns 8-bit unsigned integer attribute payload
//
// This function returns the 8-bit value of the attribute payload.
func (attr *Nlattr) U8() uint8 {
	return attrGetU8(attr)
}

// returns 16-bit unsigned integer attribute payload
//
// This function returns the 16-bit value of the attribute payload.
func (attr *Nlattr) U16() uint16 {
	return attrGetU16(attr)
}

// returns 32-bit unsigned integer attribute payload
//
// This function returns the 32-bit value of the attribute payload.
func (attr *Nlattr) U32() uint32 {
	return attrGetU32(attr)
}

// returns 64-bit unsigned integer attribute.
//
// This function returns the 64-bit value of the attribute payload.
func (attr *Nlattr) U64() uint64 {
	return attrGetU64(attr)
}

// returns string attribute.
//
// This function returns the payload of string attribute value.
func (attr *Nlattr) Str() string {
	return attrGetStr(attr)
}

// add an attribute to netlink message
//
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) Put(attr_type uint16, size Size_t, p unsafe.Pointer) {
	attrPut(nlh, attr_type, size, p)
}

// add an attribute to netlink message
//
// This function wraps Put().
func (nlh *Nlmsg) PutPtr(attr_type uint16, data interface{}) {
	attrPutPtr(nlh, attr_type, data)
}

// add an attribute to netlink message
//
// This function wraps Put().
func (nlh *Nlmsg) PutBytes(attr_type uint16, data []byte) {
	attrPutBytes(nlh, attr_type, data)
}

// add 8-bit unsigned integer attribute to netlink message
//
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutU8(attr_type uint16, data uint8) {
	attrPutU8(nlh, attr_type, data)
}

// add 16-bit unsigned integer attribute to netlink message
//
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutU16(attr_type uint16, data uint16) {
	attrPutU16(nlh, attr_type, data)
}

// add 32-bit unsigned integer attribute to netlink message
//
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutU32(attr_type uint16, data uint32) {
	attrPutU32(nlh, attr_type, data)
}

// add 64-bit unsigned integer attribute to netlink message
//
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutU64(attr_type uint16, data uint64) {
	attrPutU64(nlh, attr_type, data)
}

// add string attribute to netlink message
//
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutStr(attr_type uint16, data string) {
	attrPutStr(nlh, attr_type, data)
}

// add string attribute to netlink message
//
// This function is similar to mnl_attr_put_str, but it includes the
// NUL/zero ('\0') terminator at the end of the string.
//
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutStrz(attr_type uint16, data string) {
	attrPutStrz(nlh, attr_type, data)
}

// start an attribute nest
//
// This function adds the attribute header that identifies the beginning of
// an attribute nest. This function always returns a valid pointer to the
// beginning of the nest.
func (nlh *Nlmsg) NestStart(attr_type uint16) *Nlattr {
	return attrNestStart(nlh, attr_type)
}

// add an attribute to netlink message
//
// This function first checks that the data can be added to the message
// (fits into the buffer) and then updates the length field of the Netlink
// message (nlmsg_len) by adding the size (header + payload) of the new
// attribute. The function returns true if the attribute could be added
// to the message, otherwise false is returned.
func (nlh *Nlmsg) PutCheck(attr_type uint16, size Size_t, data unsafe.Pointer) bool {
	return attrPutCheck(nlh, attr_type, size, data)
}

// add an attribute to netlink message
//
// This function wraps PutCheck().
func (nlh *Nlmsg) PutCheckPtr(attr_type uint16, data interface{}) bool {
	return attrPutCheckPtr(nlh, attr_type, data)
}

// add an attribute to netlink message
//
// This function wraps PutCheck().
func (nlh *Nlmsg) PutCheckBytes(attr_type uint16, data []byte) bool {
	return attrPutCheckBytes(nlh, attr_type, data)
}

// add 8-bit unsigned int attribute to netlink message
//
// This function first checks that the data can be added to the message
// (fits into the buffer) and then updates the length field of the Netlink
// message (nlmsg_len) by adding the size (header + payload) of the new
// attribute. The function returns true if the attribute could be added
// to the message, otherwise false is returned.
func (nlh *Nlmsg) PutU8Check(attr_type uint16, data uint8) bool {
	return attrPutU8Check(nlh, attr_type, data)
}

// add 16-bit unsigned int attribute to netlink message
//
// This function first checks that the data can be added to the message
// (fits into the buffer) and then updates the length field of the Netlink
// message (nlmsg_len) by adding the size (header + payload) of the new
// attribute. The function returns true if the attribute could be added
// to the message, otherwise false is returned.
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutU16Check(attr_type uint16, data uint16) bool {
	return attrPutU16Check(nlh, attr_type, data)
}

// add 32-bit unsigned int attribute to netlink message
//
// This function first checks that the data can be added to the message
// (fits into the buffer) and then updates the length field of the Netlink
// message (nlmsg_len) by adding the size (header + payload) of the new
// attribute. The function returns true if the attribute could be added
// to the message, otherwise false is returned.
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutU32Check(attr_type uint16, data uint32) bool {
	return attrPutU32Check(nlh, attr_type, data)
}

// add 64-bit unsigned int attribute to netlink message
//
// This function first checks that the data can be added to the message
// (fits into the buffer) and then updates the length field of the Netlink
// message (nlmsg_len) by adding the size (header + payload) of the new
// attribute. The function returns true if the attribute could be added
// to the message, otherwise false is returned.
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutU64Check(attr_type uint16, data uint64) bool {
	return attrPutU64Check(nlh, attr_type, data)
}

// add string attribute to netlink message
//
// This function first checks that the data can be added to the message
// (fits into the buffer) and then updates the length field of the Netlink
// message (nlmsg_len) by adding the size (header + payload) of the new
// attribute. The function returns true if the attribute could be added
// to the message, otherwise false is returned.
// This function updates the length field of the Netlink message (nlmsg_len)
// by adding the size (header + payload) of the new attribute.
func (nlh *Nlmsg) PutStrCheck(attr_type uint16, data string) bool {
	return attrPutStrCheck(nlh, attr_type, data)
}

// add string attribute to netlink message
//
// This function is similar to mnl_attr_put_str, but it includes the
// NUL/zero ('\0') terminator at the end of the string.
//
// This function first checks that the data can be added to the message
// (fits into the buffer) and then updates the length field of the Netlink
// message (nlmsg_len) by adding the size (header + payload) of the new
// attribute. The function returns true if the attribute could be added
// to the message, otherwise false is returned.
func (nlh *Nlmsg) PutStrzCheck(attr_type uint16, data string) bool {
	return attrPutStrzCheck(nlh, attr_type, data)
}

// start an attribute nest
//
// This function adds the attribute header that identifies the beginning of
// an attribute nest. If the nested attribute cannot be added then nil,
// otherwise valid pointer to the beginning of the nest is returned.
func (nlh *Nlmsg) NestStartCheck(attr_type uint16) *Nlattr {
	return attrNestStartCheck(nlh, attr_type)
}

// end an attribute nest
//
// This function updates the attribute header that identifies the nest.
func (nlh *Nlmsg) NestEnd(start *Nlattr) {
	attrNestEnd(nlh, start)
}

// cancel an attribute nest
//
// This function updates the attribute header that identifies the nest.
func (nlh *Nlmsg) NestCancel(start *Nlattr) {
	attrNestCancel(nlh, start)
}

// mnl_attr_for_each() macro in libmnl.h
func (nlh *Nlmsg) Attributes(offset Size_t) <-chan *Nlattr {
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

// mnl_attr_for_each_nested() macro in libmnl.h
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

// mnl_attr_for_each_payload() macro in libmnl.h
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

// This function creates a new Nlattr
func NewNlattr(size int) (*Nlattr, error) {
	if size < SizeofNlattr {
		return nil, errors.New("too short size")
	}
	b := make([]byte, size)
	return (*Nlattr)(unsafe.Pointer(&b[0])), nil
}

// This function points []byte as pointer to a Nlattr
func NlattrPointer(b []byte) *Nlattr {
	// ???: check buf len
	//      len(b) >= SizeofNlattr
	//      nla.len <= len(b)
	return (*Nlattr)(unsafe.Pointer(&b[0]))
}

// get the length of the Netlink payload
//
// This function returns the Length of the netlink payload, ie. the length
// of the full message minus the size of the Netlink header.
func (nlh *Nlmsg) PayloadLen() Size_t {
	return nlmsgGetPayloadLen(nlh)
}

// reserve and prepare room for an extra header
//
// This function sets to zero the room that is required to put the extra
// header after the initial Netlink header. This function also increases
// the nlmsg_len field. You have to invoke mnl_nlmsg_put_header() before
// you call this function. This function returns a pointer to the extra
// header.
func (nlh *Nlmsg) PutExtraHeader(size Size_t) unsafe.Pointer {
	return nlmsgPutExtraHeader(nlh, size)
}

// get a pointer to the payload of the netlink message
//
// This function returns a pointer to the payload of the netlink message.
func (nlh *Nlmsg) Payload() unsafe.Pointer {
	return nlmsgGetPayload(nlh)
}

// get the payload of the netlink message as a []byte
//
// This function returns the payload of the netlink message as a []byte.
func (nlh *Nlmsg) PayloadBytes() []byte {
	return nlmsgGetPayloadBytes(nlh)
}

// get a pointer to the payload of the message
//
// This function returns a pointer to the payload of the netlink message plus
// a given offset.
func (nlh *Nlmsg) PayloadOffset(offset Size_t) unsafe.Pointer {
	return nlmsgGetPayloadOffset(nlh, offset)
}

// get the payload of the message as a []byte
//
// This function returns the payload of the netlink message as a []byte.
func (nlh *Nlmsg) PayloadOffsetBytes(offset Size_t) []byte {
	return nlmsgGetPayloadOffsetBytes(nlh, offset)
}

// check a there is room for netlink message
//
// This function is used to check that a buffer that contains a netlink
// message has enough room for the netlink message that it stores, ie. this
// function can be used to verify that a netlink message is not malformed nor
// truncated.
//
// This function does not set errno in case of error since it is intended
// for iterations. Thus, it returns true on success and false on error.
//
// The size parameter may become negative in malformed messages during message
// iteration, that is why we use a signed integer.
func (nlh *Nlmsg) Ok(size int) bool {
	return nlmsgOk(nlh, size)
}

// get the next netlink message in a multipart message
//
// This function returns a pointer to the next netlink message that is part
// of a multi-part netlink message. Netlink can batch several messages into
// one buffer so that the receiver has to iterate over the whole set of
// Netlink messages.
//
// You have to use Ok() to check if the next Netlink message is
// valid.
func (nlh *Nlmsg) Next(size int) (*Nlmsg, int) {
	return nlmsgNext(nlh, size)
}

// get the ending of the netlink message
//
// This function returns a pointer to the netlink message tail. This is useful
// to build a message since we continue adding attributes at the end of the
// message.
func (nlh *Nlmsg) PayloadTail() unsafe.Pointer {
	return nlmsgGetPayloadTail(nlh)
}

// perform sequence tracking
//
// This functions returns true if the sequence tracking is fulfilled, otherwise
// false is returned. We skip the tracking for netlink messages whose sequence
// number is zero since it is usually reserved for event-based kernel
// notifications. On the other hand, if seq is set but the message sequence
// number is not set (i.e. this is an event message coming from kernel-space),
// then we also skip the tracking. This approach is good if we use the same
// socket to send commands to kernel-space (that we want to track) and to
// listen to events (that we do not track).
func (nlh *Nlmsg) SeqOk(seq uint32) bool {
	return nlmsgSeqOk(nlh, seq)
}

// perform portID origin check
//
// This functions returns true if the origin is fulfilled, otherwise
// false is returned. We skip the tracking for netlink message whose portID
// is zero since it is reserved for event-based kernel notifications. On the
// other hand, if portid is set but the message PortID is not (i.e. this
// is an event message coming from kernel-space), then we also skip the
// tracking. This approach is good if we use the same socket to send commands
// to kernel-space (that we want to track) and to listen to events (that we
// do not track).
func (nlh *Nlmsg) PortidOk(portid uint32) bool {
	return nlmsgPortidOk(nlh, portid)
}

// print netlink message to file
//
// This function prints the netlink header to a file handle.
// It may be useful for debugging purposes. One example of the output
// is the following:
//
//  ----------------        ------------------
//  |  0000000040  |        | message length |
//  | 00016 | R-A- |        |  type | flags  |
//  |  1289148991  |        | sequence number|
//  |  0000000000  |        |     port ID    |
//  ----------------        ------------------
//  | 00 00 00 00  |        |  extra header  |
//  | 00 00 00 00  |        |  extra header  |
//  | 01 00 00 00  |        |  extra header  |
//  | 01 00 00 00  |        |  extra header  |
//  |00008|--|00003|        |len |flags| type|
//  | 65 74 68 30  |        |      data      |       e t h 0
//  ----------------        ------------------
//
// This example above shows the netlink message that is send to kernel-space
// to set up the link interface eth0. The netlink and attribute header data
// are displayed in base 10 whereas the extra header and the attribute payload
// are expressed in base 16. The possible flags in the netlink header are:
//
// - R, that indicates that NLM_F_REQUEST is set.
// - M, that indicates that NLM_F_MULTI is set.
// - A, that indicates that NLM_F_ACK is set.
// - E, that indicates that NLM_F_ECHO is set.
//
// The lack of one flag is displayed with '-'. On the other hand, the possible
// attribute flags available are:
//
// - N, that indicates that NLA_F_NESTED is set.
// - B, that indicates that NLA_F_NET_BYTEORDER is set.
func (nlh *Nlmsg) Fprint(fd *os.File, extra_header_size Size_t) {
	nlmsgFprintNlmsg(fd, nlh, extra_header_size)
}

// initialize a batch
//
// The buffer that you pass must be double of MNL_SOCKET_BUFFER_SIZE. The
// limit must be half of the buffer size, otherwise expect funny memory
// corruptions 8-).
//
// You can allocate the buffer that you use to store the batch in the stack or
// the heap, no restrictions in this regard. This function returns nil on
// error.
func NewNlmsgBatch(buf []byte, limit Size_t) (*NlmsgBatch, error) {
	return nlmsgBatchStart(buf, limit)
}

// release a batch
//
// This function releases the batch.
func (b *NlmsgBatch) Stop() {
	nlmsgBatchStop(b)
}

// get room for the next message in the batch
//
// This function returns false if the last message did not fit into the
// batch. Otherwise, it prepares the batch to provide room for the new
// Netlink message in the batch and returns true.
//
// You have to put at least one message in the batch before calling this
// function, otherwise your application is likely to crash.
func (b *NlmsgBatch) Next() bool {
	return nlmsgBatchNext(b)
}

// reset the batch
//
// This function allows to reset a batch, so you can reuse it to create a
// new one. This function moves the last message which does not fit the
// batch to the head of the buffer, if any.
func (b *NlmsgBatch) Reset() {
	nlmsgBatchReset(b)
}

// get current size of the batch
//
// This function returns the current size of the batch.
func (b *NlmsgBatch) Size() Size_t {
	return nlmsgBatchSize(b)
}

// get head of this batch
//
// This function returns a pointer to the head of the batch, which is the
// beginning of the buffer that is used.
func (b *NlmsgBatch) Head() unsafe.Pointer {
	return nlmsgBatchHead(b)
}

// get head of this batch
//
// This function wraps Head, returns as a []byte
func (b *NlmsgBatch) HeadBytes() []byte {
	return nlmsgBatchHeadBytes(b)
}

// returns current position in the batch
//
// This function returns a pointer to the current position in the buffer
// that is used to store the batch.
func (b *NlmsgBatch) Current() unsafe.Pointer {
	return nlmsgBatchCurrent(b)
}

func (b *NlmsgBatch) CurrentNlmsg() *Nlmsg {
	var d []byte
	buf := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	len := int(nlmsgBatchLimit(b) - nlmsgBatchSize(b))
	buf.Cap = len
	buf.Len = len
	buf.Data = uintptr(nlmsgBatchCurrent(b))
	return RawNlmsgBytes(d)
}

// check if there is any message in the batch
//
// This function returns true if the batch is empty.
func (b *NlmsgBatch) IsEmpty() bool {
	return nlmsgBatchIsEmpty(b)
}

// reserve and prepare room for Netlink header
//
// This function sets to zero the room that is required to put the Netlink
// header in the memory buffer passed as parameter. This function also
// initializes the nlmsg_len field to the size of the Netlink header. This
// function returns a pointer to the Netlink header structure.
func (nlh *Nlmsg) PutHeader() {
	C.mnl_nlmsg_put_header(unsafe.Pointer(nlh.Nlmsghdr))
}

// create a new Nlmsg from []byte
func RawNlmsgBytes(b []byte) *Nlmsg {
	var d []byte
	nlh := (*Nlmsghdr)(unsafe.Pointer(&b[0]))
	s := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	s.Cap = int(nlh.Len)
	s.Len = int(nlh.Len)
	s.Data = uintptr(unsafe.Pointer(&b[0]))
	return &Nlmsg{nlh, d}
}

// create a new Nlmsg from raw nlmsg pointer
func nlmsgPointer(nlh *C.struct_nlmsghdr) *Nlmsg {
	var b []byte
	p := (*Nlmsghdr)(unsafe.Pointer(nlh))
	s := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	s.Cap = int(p.Len)
	s.Len = int(p.Len)
	s.Data = uintptr(unsafe.Pointer(nlh))
	return &Nlmsg{p, b}
}

// reserve and prepare room for Netlink header
//
// This function sets to zero the room that is required to put the Netlink
// header in the memory buffer passed as parameter. This function also
// initializes the nlmsg_len field to the size of the Netlink header. This
// function returns a pointer to the Netlink header structure.
func NewNlmsgBytes(buf []byte) (*Nlmsg, error) {
	if len(buf) < int(MNL_NLMSG_HDRLEN) {
		return nil, syscall.EINVAL
	}
	nlh := Nlmsg{(*Nlmsghdr)(unsafe.Pointer(&buf[0])), buf}
	C.mnl_nlmsg_put_header(unsafe.Pointer(&buf[0]))
	return &nlh, nil
}

// create and reserve room for Netlink header
func NewNlmsg(size int) (*Nlmsg, error) {
	b := make([]byte, size)
	return NewNlmsgBytes(b)
}

// open a netlink socket
//
// On error, it returns nil and errno is appropriately set. Otherwise, it
// returns a valid pointer to the mnl_socket structure.
func NewSocket(bus int) (*Socket, error) {
	return socketOpen(bus)
}

// obtain file descriptor from netlink socket
//
// This function returns the file descriptor of a given netlink socket.
func (nl *Socket) Fd() int {
	return socketGetFd(nl)
}

// obtain Netlink PortID from netlink socket
//
// This function returns the Netlink PortID of a given netlink socket.
// It's a common mistake to assume that this PortID equals the process ID
// which is not always true. This is the case if you open more than one
// socket that is binded to the same Netlink subsystem from the same process.
func (nl *Socket) Portid() uint32 {
	return socketGetPortid(nl)
}

// bind netlink socket
//
// On error, this function returns error. On
// success, 0 is returned. You can use MNL_SOCKET_AUTOPID which is 0 for
// automatic port ID selection.
func (nl *Socket) Bind(groups uint, pid Pid_t) error {
	return socketBind(nl, groups, pid)
}

// send a netlink message of a certain size
//
// On error, it returns -1 and errno is appropriately set. Otherwise, it
// returns the number of bytes sent.
func (nl *Socket) Sendto(buf []byte) (Ssize_t, error) {
	return socketSendto(nl, buf)
}

// send a netlink message
//
// This function wraps Sendto().
func (nl *Socket) SendNlmsg(nlh *Nlmsg) (Ssize_t, error) {
	return socketSendNlmsg(nl, nlh)
}

// receive a netlink message
//
// On error, it returns -1 and errno is appropriately set. If errno is set
// to ENOSPC, it means that the buffer that you have passed to store the
// netlink message is too small, so you have received a truncated message.
// To avoid this, you have to allocate a buffer of MNL_SOCKET_BUFFER_SIZE
// (which is 8KB, see linux/netlink.h for more information). Using this
// buffer size ensures that your buffer is big enough to store the netlink
// message without truncating it.
func (nl *Socket) Recvfrom(buf []byte) (Ssize_t, error) {
	return socketRecvfrom(nl, buf)
}

// close a given netlink socket
//
// On error, this function error.
func (nl *Socket) Close() error {
	return socketClose(nl)
}

// set Netlink socket option
//
// This function allows you to set some Netlink socket option. As of writing
// this (see linux/netlink.h), the existing options are:
//
//	- #define NETLINK_ADD_MEMBERSHIP  1
//	- #define NETLINK_DROP_MEMBERSHIP 2
//	- #define NETLINK_PKTINFO         3
//	- #define NETLINK_BROADCAST_ERROR 4
//	- #define NETLINK_NO_ENOBUFS      5
//
// In the early days, Netlink only supported 32 groups expressed in a
// 32-bits mask. However, since 2.6.14, Netlink may have up to 2^32 multicast
// groups but you have to use setsockopt() with NETLINK_ADD_MEMBERSHIP to
// join a given multicast group. This function internally calls setsockopt()
// to join a given netlink multicast group. You can still use mnl_bind()
// and the 32-bit mask to join a set of Netlink multicast groups.
//
// On error, this function error.
func (nl *Socket) Setsockopt(t int, v unsafe.Pointer, l Socklen_t) error {
	return socketSetsockopt(nl, t, v, l)
}

// set Netlink socket option
//
// This function wraps Setsockopt()
func (nl *Socket) SetsockoptBytes(optype int, buf []byte) error {
	return socketSetsockoptBytes(nl, optype, buf)
}

// set Netlink socket option
//
// This function wraps Setsockopt()
func (nl *Socket) SetsockoptByte(optype int, v byte) error {
	return socketSetsockoptByte(nl, optype, v)
}

// set Netlink socket option
//
// This function wraps Setsockopt()
func (nl *Socket) SetsockoptCint(optype int, v int) error {
	return socketSetsockoptCint(nl, optype, v)
}

// get a Netlink socket option
//
// On error, this function returns nil and error.
func (nl *Socket) Sockopt(optype int, size Socklen_t) ([]byte, error) {
	return socketGetsockopt(nl, optype, size)
}
