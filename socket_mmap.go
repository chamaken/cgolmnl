// +build nlmmap

package cgolmnl

import "unsafe"

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>

void *mnl_frame_payload(void *hdr)
{
	return MNL_FRAME_PAYLOAD(hdr);
}
*/
import "C"

type Ring C.struct_mnl_ring
type RingTypes C.enum_mnl_ring_types
const (
	MNL_RING_RX RingTypes = C.MNL_RING_RX
	MNL_RING_TX RingTypes = C.MNL_RING_TX
)

func FramePayload(hdr *NlMmapHdr) []byte {
	return SharedBytes(C.mnl_frame_payload(unsafe.Pointer(hdr)), int(hdr.Len))
}

// int mnl_socket_set_ringopt(struct mnl_socket *nl, struct nl_mmap_req *req,
//			      enum mnl_ring_types type)
func socketSetRingopt(nl *Socket, rtype RingTypes, block_size, block_nr, frame_size, frame_nr uint) error {
	_, err := C.mnl_socket_set_ringopt((*C.struct_mnl_socket)(nl), (C.enum_mnl_ring_types)(rtype),
		C.uint(block_size), C.uint(block_nr), C.uint(frame_size), C.uint(frame_nr))
	return err
}

// int mnl_socket_map_ring(struct mnl_socket *nl)
func socketMapRing(nl *Socket) error {
	_, err := C.mnl_socket_map_ring((*C.struct_mnl_socket)(nl))
	return err
}

// int mnl_socket_unmap_ring(struct mnl_socket *nl)
func socketUnmapRing(nl *Socket) error {
	_, err := C.mnl_socket_unmap_ring((*C.struct_mnl_socket)(nl))
	return err
}

// struct mnl_ring *mnl_socket_get_ring(const struct mnl_socket *nl, enum mnl_ring_types type)
func socketGetRing(nl *Socket, rtype RingTypes) (*Ring, error) {
	ret, err := C.mnl_socket_get_ring((*C.struct_mnl_socket)(nl), (C.enum_mnl_ring_types)(rtype))
	return (*Ring)(unsafe.Pointer(ret)), err
}

// struct nl_mmap_hdr *mnl_ring_get_frame(const struct mnl_ring *ring)
func ringGetFrame(ring *Ring) (*NlMmapHdr) {
	return (*NlMmapHdr)(unsafe.Pointer(C.mnl_ring_get_frame((*C.struct_mnl_ring)(ring))))
}

// void mnl_ring_advance(const struct mnl_ring *ring)
func ringAdvance(ring *Ring) {
	C.mnl_ring_advance((*C.struct_mnl_ring)(ring))
}

// void mnl_nlmsg_batch_reset_buffer(struct mnl_nlmsg_batch *b, void *buf, size_t limit)
func nlmsgBatchResetBuffer(b *NlmsgBatch, buf []byte) {
	C.mnl_nlmsg_batch_reset_buffer((*C.struct_mnl_nlmsg_batch)(b), (unsafe.Pointer)(&buf[0]), C.size_t(len(buf)))
}

// receivers

func (nl *Socket) SetRingopt(t RingTypes, bs, bn, fs, fn uint) error {
	return socketSetRingopt(nl, t, bs, bn, fs, fn)
}

func (nl *Socket) MapRing() error {
	return socketMapRing(nl)
}

func (nl *Socket) UnmapRing() error {
	return socketUnmapRing(nl)
}

func (nl *Socket) Ring(t RingTypes) (*Ring, error) {
	return socketGetRing(nl, t)
}

func (ring *Ring) Frame() *NlMmapHdr {
	return ringGetFrame(ring)
}

func (ring *Ring) Advance() {
	ringAdvance(ring)
}

func (b *NlmsgBatch) ResetBuffer(buf []byte) {
	nlmsgBatchResetBuffer(b, buf)
}
