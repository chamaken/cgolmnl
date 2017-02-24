package cgolmnl

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <stdlib.h>
#include <stdio.h>
#include <libmnl/libmnl.h>
#include "helper.h"
*/
import "C"

import (
	"os"
	"unsafe"
)

// mnl_nlmsg_size - calculate the size of Netlink message (without alignment)
//
// This function returns the size of a netlink message (header plus payload)
// without alignment.
func NlmsgSize(size Size_t) Size_t {
	return Size_t(C.mnl_nlmsg_size(C.size_t(size)))
}

// size_t mnl_nlmsg_get_payload_len(const struct nlmsghdr *nlh)
func nlmsgGetPayloadLen(nlh *Nlmsg) Size_t {
	return Size_t(
		C.mnl_nlmsg_get_payload_len(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr))))
}

// void *
// mnl_nlmsg_put_extra_header(struct nlmsghdr *nlh, size_t size)
func nlmsgPutExtraHeader(nlh *Nlmsg, size Size_t) unsafe.Pointer {
	return C.mnl_nlmsg_put_extra_header(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.size_t(size))
}

// void *mnl_nlmsg_get_payload(const struct nlmsghdr *nlh)
func nlmsgGetPayload(nlh *Nlmsg) unsafe.Pointer {
	return C.mnl_nlmsg_get_payload(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)))
}
// XXX
func nlmsgGetPayloadBytes(nlh *Nlmsg) []byte {
	return SharedBytes(nlmsgGetPayload(nlh), int(nlmsgGetPayloadLen(nlh)))
}

// void *
// mnl_nlmsg_get_payload_offset(const struct nlmsghdr *nlh, size_t offset)
func nlmsgGetPayloadOffset(nlh *Nlmsg, offset Size_t) unsafe.Pointer {
	return C.mnl_nlmsg_get_payload_offset(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
		C.size_t(offset))
}
// XXX
func nlmsgGetPayloadOffsetBytes(nlh *Nlmsg, offset Size_t) []byte {
	return SharedBytes(nlmsgGetPayloadOffset(nlh, offset), int(nlmsgGetPayloadLen(nlh)-Size_t(MnlAlign(uint32(offset)))))
}

// bool mnl_nlmsg_ok(const struct nlmsghdr *nlh, int len)
func nlmsgOk(nlh *Nlmsg, size int) bool {
	// test fails
	//   unexpected fault address 0x--------
	//   fatal error: fault
	// sometimes without below
	if size < SizeofNlmsg {
		return false
	}
	return bool(
		C.mnl_nlmsg_ok(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.int(size)))
}

// struct nlmsghdr *
// mnl_nlmsg_next(const struct nlmsghdr *nlh, int *len)
func nlmsgNext(nlh *Nlmsg, size int) (*Nlmsg, int) {
	c_size := C.int(size)
	h := (*C.struct_nlmsghdr)(unsafe.Pointer(
		C.mnl_nlmsg_next(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			&c_size)))
	return nlmsgPointer(h), int(c_size)
}

// void *mnl_nlmsg_get_payload_tail(const struct nlmsghdr *nlh)
func nlmsgGetPayloadTail(nlh *Nlmsg) unsafe.Pointer {
	return C.mnl_nlmsg_get_payload_tail(
		(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)))
}

// bool
// mnl_nlmsg_seq_ok(const struct nlmsghdr *nlh, uint32_t seq)
func nlmsgSeqOk(nlh *Nlmsg, seq uint32) bool {
	return bool(
		C.mnl_nlmsg_seq_ok(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.uint(seq)))
}

// bool
// mnl_nlmsg_portid_ok(const struct nlmsghdr *nlh, uint32_t portid)
func nlmsgPortidOk(nlh *Nlmsg, portid uint32) bool {
	return bool(
		C.mnl_nlmsg_portid_ok(
			(*C.struct_nlmsghdr)(unsafe.Pointer(nlh.Nlmsghdr)),
			C.uint(portid)))
}

// void
// mnl_nlmsg_fprintf(FILE *fd, const void *data, size_t datalen,
//		     size_t extra_header_size)
func nlmsgFprint(fd *os.File, data unsafe.Pointer, size Size_t, extra_header_size Size_t) {
	mode := C.CString("w")
	defer C.free(unsafe.Pointer(mode))
	C.mnl_nlmsg_fprintf(C.fdopen(C.int(fd.Fd()), mode),
		data, C.size_t(size), C.size_t(extra_header_size))
}
func nlmsgFprintBytes(fd *os.File, data []byte, extra_header_size Size_t) {
	nlmsgFprint(fd, unsafe.Pointer(&data[0]), Size_t(len(data)), extra_header_size)
}
func nlmsgFprintNlmsg(fd *os.File, nlh *Nlmsg, extra_header_size Size_t) {
	nlmsgFprint(
		fd, unsafe.Pointer(nlh.Nlmsghdr),
		Size_t(nlh.Len), extra_header_size)
}

// Netlink message batch helpers
//
// This library provides helpers to batch several messages into one single
// datagram. These helpers do not perform strict memory boundary checkings.
//
// The following figure represents a Netlink message batch:
//
//   |<-------------- MNL_SOCKET_BUFFER_SIZE ------------->|
//   |<-------------------- batch ------------------>|     |
//   |-----------|-----------|-----------|-----------|-----------|
//   |<- nlmsg ->|<- nlmsg ->|<- nlmsg ->|<- nlmsg ->|<- nlmsg ->|
//   |-----------|-----------|-----------|-----------|-----------|
//                                             ^           ^
//                                             |           |
//                                        message N   message N+1
//
// To start the batch, you have to call mnl_nlmsg_batch_start() and you can
// use mnl_nlmsg_batch_stop() to release it.
//
// You have to invoke mnl_nlmsg_batch_next() to get room for a new message
// in the batch. If this function returns NULL, it means that the last
// message that was added (message N+1 in the figure above) does not fit the
// batch. Thus, you have to send the batch (which includes until message N)
// and, then, you have to call mnl_nlmsg_batch_reset() to re-initialize
// the batch (this moves message N+1 to the head of the buffer). For that
// reason, the buffer that you have to use to store the batch must be double
// of MNL_SOCKET_BUFFER_SIZE to ensure that the last message (message N+1)
// that did not fit into the batch is written inside valid memory boundaries.
type NlmsgBatch struct {
	c *C.struct_mnl_nlmsg_batch	// [0]byte
	buf []byte			// holder to prevent gc
}

// struct mnl_nlmsg_batch *mnl_nlmsg_batch_start(void *buf, size_t limit)
func nlmsgBatchStart(buf []byte, limit Size_t) (*NlmsgBatch, error) {
	rs := NlmsgBatch{nil, buf}
	var err error
	rs.c, err = C.mnl_nlmsg_batch_start(unsafe.Pointer(&buf[0]), C.size_t(limit))
	// return (*NlmsgBatch)(ret), err
	return &rs, err
}

// void mnl_nlmsg_batch_stop(struct mnl_nlmsg_batch *b)
func nlmsgBatchStop(b *NlmsgBatch) {
	C.mnl_nlmsg_batch_stop(b.c)
}

// bool mnl_nlmsg_batch_next(struct mnl_nlmsg_batch *b)
func nlmsgBatchNext(b *NlmsgBatch) bool {
	return bool(C.mnl_nlmsg_batch_next(b.c))
}

// void mnl_nlmsg_batch_reset(struct mnl_nlmsg_batch *b)
func nlmsgBatchReset(b *NlmsgBatch) {
	C.mnl_nlmsg_batch_reset(b.c)
}

// size_t mnl_nlmsg_batch_size(struct mnl_nlmsg_batch *b)
func nlmsgBatchSize(b *NlmsgBatch) Size_t {
	return Size_t(C.mnl_nlmsg_batch_size(b.c))
}

// void *mnl_nlmsg_batch_head(struct mnl_nlmsg_batch *b)
func nlmsgBatchHead(b *NlmsgBatch) unsafe.Pointer {
	return C.mnl_nlmsg_batch_head(b.c)
}
func nlmsgBatchHeadBytes(b *NlmsgBatch) []byte {
	return SharedBytes(nlmsgBatchHead(b), int(nlmsgBatchSize(b)))
}

// void *mnl_nlmsg_batch_current(struct mnl_nlmsg_batch *b)
func nlmsgBatchCurrent(b *NlmsgBatch) unsafe.Pointer {
	return C.mnl_nlmsg_batch_current(b.c)
}

// bool mnl_nlmsg_batch_is_empty(struct mnl_nlmsg_batch *b)
func nlmsgBatchIsEmpty(b *NlmsgBatch) bool {
	return bool(C.mnl_nlmsg_batch_is_empty(b.c))
}

func nlmsgBatchLimit(b *NlmsgBatch) Size_t {
	return Size_t(C.mnl_nlmsg_batch_limit(b.c))
}
