package cgolmnl

/*
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>
#include <stdio.h>
*/
import "C"

import (
	"unsafe"
	"os"
)

/**
 * mnl_nlmsg_size - calculate the size of Netlink message (without alignment)
 *
 * size_t mnl_nlmsg_size(size_t len)
 */
func NlmsgSize(size Size_t) Size_t {
	return Size_t(C.mnl_nlmsg_size(C.size_t(size)))
}

/**
 * mnl_nlmsg_get_payload_len - get the length of the Netlink payload
 *
 * size_t mnl_nlmsg_get_payload_len(const struct nlmsghdr *nlh)
 */
func NlmsgGetPayloadLen(nlh *Nlmsghdr) Size_t {
	return Size_t(C.mnl_nlmsg_get_payload_len((*C.struct_nlmsghdr)(unsafe.Pointer(nlh))))
}

/**
 * mnl_nlmsg_put_header - reserve and prepare room for Netlink header
 *
 * struct nlmsghdr *mnl_nlmsg_put_header(void *buf)
 */
func NlmsgPutHeader(buf unsafe.Pointer) *Nlmsghdr {
	return (*Nlmsghdr)(unsafe.Pointer(C.mnl_nlmsg_put_header(buf)))
}

/**
 * mnl_nlmsg_put_extra_header - reserve and prepare room for an extra header
 *
 * void *
 * mnl_nlmsg_put_extra_header(struct nlmsghdr *nlh, size_t size)
 */
func NlmsgPutExtraHeader(nlh *Nlmsghdr, size Size_t) unsafe.Pointer {
	return C.mnl_nlmsg_put_extra_header((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), C.size_t(size))
}

/**
 * mnl_nlmsg_get_payload - get a pointer to the payload of the netlink message
 *
 * void *mnl_nlmsg_get_payload(const struct nlmsghdr *nlh)
 */
func NlmsgGetPayload(nlh *Nlmsghdr) unsafe.Pointer {
	return C.mnl_nlmsg_get_payload((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)))
}
func NlmsgGetPayloadBytes(nlh *Nlmsghdr) []byte {
	return SharedBytes(NlmsgGetPayload(nlh), int(NlmsgGetPayloadLen(nlh)))
}

/**
 * mnl_nlmsg_get_payload_offset - get a pointer to the payload of the message
 *
 * void *
 * mnl_nlmsg_get_payload_offset(const struct nlmsghdr *nlh, size_t offset)
 */
func NlmsgGetPayloadOffset(nlh *Nlmsghdr, offset Size_t) unsafe.Pointer {
	return C.mnl_nlmsg_get_payload_offset((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), C.size_t(offset))
}
func NlmsgGetPayloadOffsetBytes(nlh *Nlmsghdr, offset Size_t) []byte {
	return SharedBytes(NlmsgGetPayloadOffset(nlh, offset), int(NlmsgGetPayloadLen(nlh) - offset))
}

/**
 * mnl_nlmsg_ok - check a there is room for netlink message
 *
 * bool mnl_nlmsg_ok(const struct nlmsghdr *nlh, int len)
 */
func NlmsgOk(nlh *Nlmsghdr, size int) bool {
	return bool(C.mnl_nlmsg_ok((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), C.int(size)))
}

/**
 * mnl_nlmsg_next - get the next netlink message in a multipart message
 *
 * struct nlmsghdr *
 * mnl_nlmsg_next(const struct nlmsghdr *nlh, int *len)
 */
func NlmsgNext(nlh *Nlmsghdr, size int) (*Nlmsghdr, int) {
	c_size := C.int(size)
	h := (*Nlmsghdr)(unsafe.Pointer(C.mnl_nlmsg_next((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), &c_size)))
	return h, size
}

/**
 * mnl_nlmsg_get_payload_tail - get the ending of the netlink message
 *
 * void *mnl_nlmsg_get_payload_tail(const struct nlmsghdr *nlh)
 */
func NlmsgGetPayloadTail(nlh *Nlmsghdr) unsafe.Pointer {
	return C.mnl_nlmsg_get_payload_tail((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)))
}

/**
 * mnl_nlmsg_seq_ok - perform sequence tracking
 *
 * bool
 * mnl_nlmsg_seq_ok(const struct nlmsghdr *nlh, uint32_t seq)
 */
func NlmsgSeqOk(nlh *Nlmsghdr, seq uint32) bool {
	return bool(C.mnl_nlmsg_seq_ok((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), C.uint32_t(seq)))
}

/**
 * mnl_nlmsg_portid_ok - perform portID origin check
 *
 * bool
 * mnl_nlmsg_portid_ok(const struct nlmsghdr *nlh, uint32_t portid)
 */
func NlmsgPortidOk(nlh *Nlmsghdr, portid uint32) bool {
	return bool(C.mnl_nlmsg_portid_ok((*C.struct_nlmsghdr)(unsafe.Pointer(nlh)), C.uint32_t(portid)))
}


/**
 * mnl_nlmsg_fprintf - print netlink message to file
 *
 * void
 * mnl_nlmsg_fprintf(FILE *fd, const void *data, size_t datalen,
 *		     size_t extra_header_size)
 */
func NlmsgFprint(fd *os.File, data []byte, extra_header_size Size_t) {
	C.mnl_nlmsg_fprintf(C.fdopen(C.int(fd.Fd()), C.CString("w")),
		unsafe.Pointer(&data[0]), C.size_t(len(data)), C.size_t(extra_header_size))
}

/*
 * Netlink message batch helpers
 */

// struct mnl_nlmsg_batch
type NlmsgBatchDescriptor [0]byte // unsafe.Pointer

/**
 * mnl_nlmsg_batch_start - initialize a batch
 *
 * struct mnl_nlmsg_batch *mnl_nlmsg_batch_start(void *buf, size_t limit)
 */
func NlmsgBatchStart(buf []byte, limit Size_t) (*NlmsgBatchDescriptor, error) {
	ret, err := C.mnl_nlmsg_batch_start(unsafe.Pointer(&buf[0]), C.size_t(limit))
	return  (*NlmsgBatchDescriptor)(ret), err
}

/**
 * mnl_nlmsg_batch_stop - release a batch
 *
 * void mnl_nlmsg_batch_stop(struct mnl_nlmsg_batch *b)
 */
func NlmsgBatchStop(b *NlmsgBatchDescriptor) {
	C.mnl_nlmsg_batch_stop((*[0]byte)(b))
}

/**
 * mnl_nlmsg_batch_next - get room for the next message in the batch
 *
 * bool mnl_nlmsg_batch_next(struct mnl_nlmsg_batch *b)
 */
func NlmsgBatchNext(b *NlmsgBatchDescriptor) bool {
	return bool(C.mnl_nlmsg_batch_next((*[0]byte)(b)))
}

/**
 * mnl_nlmsg_batch_reset - reset the batch
 *
 * void mnl_nlmsg_batch_reset(struct mnl_nlmsg_batch *b)
 */
func NlmsgBatchReset(b *NlmsgBatchDescriptor) {
	C.mnl_nlmsg_batch_reset((*[0]byte)(b))
}

/**
 * mnl_nlmsg_batch_size - get current size of the batch
 *
 * size_t mnl_nlmsg_batch_size(struct mnl_nlmsg_batch *b)
 */
func NlmsgBatchSize(b *NlmsgBatchDescriptor) Size_t {
	return Size_t(C.mnl_nlmsg_batch_size((*[0]byte)(b)))
}

/**
 * mnl_nlmsg_batch_head - get head of this batch
 *
 * void *mnl_nlmsg_batch_head(struct mnl_nlmsg_batch *b)
 */
func NlmsgBatchHead(b *NlmsgBatchDescriptor) unsafe.Pointer {
	return C.mnl_nlmsg_batch_head((*[0]byte)(b))
}
func NlmsgBatchHeadBytes(b *NlmsgBatchDescriptor) []byte {
	return SharedBytes(NlmsgBatchHead(b), int(NlmsgBatchSize(b)))
}

/**
 * mnl_nlmsg_batch_current - returns current position in the batch
 *
 * void *mnl_nlmsg_batch_current(struct mnl_nlmsg_batch *b)
 */
func NlmsgBatchCurrent(b *NlmsgBatchDescriptor) unsafe.Pointer {
	return C.mnl_nlmsg_batch_current((*[0]byte)(b))
}

/**
 * mnl_nlmsg_batch_is_empty - check if there is any message in the batch
 *
 * bool mnl_nlmsg_batch_is_empty(struct mnl_nlmsg_batch *b)
 */
func NlmsgBatchIsEmpty(b *NlmsgBatchDescriptor) bool {
	return bool(C.mnl_nlmsg_batch_is_empty((*[0]byte)(b)))
}

/**
 * mnl_nlmsg_batch_reset_buffer - reset the new buffer
 *
 * void mnl_nlmsg_batch_reset_buffer(struct mnl_nlmsg_batch *b, void *buf, size_t limit)
func NlmsgBatchResetBuffer(b NlmsgBatchDescriptor, buf []byte) {
	C.mnl_nlmsg_batch_reset_buffer((*[0]byte)(b), (unsafe.Pointer)(&buf[0]), C.size_t(len(buf)))
}
 */
