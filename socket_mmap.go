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

/**
 * mnl_socket_set_ringopt - set ring opt to prepare for mnl_socket_map_ring()
 * int mnl_socket_set_ringopt(struct mnl_socket *nl, struct nl_mmap_req *req,
 *			      enum mnl_ring_types type)
 */
func SocketSetRingopt(nl *Socket, rtype RingTypes, block_size, block_nr, frame_size, frame_nr uint) error {
	_, err := C.mnl_socket_set_ringopt((*C.struct_mnl_socket)(nl), (C.enum_mnl_ring_types)(rtype),
		C.uint(block_size), C.uint(block_nr), C.uint(frame_size), C.uint(frame_nr))
	return err
}

/**
 * mnl_socket_map_ring - setup a ring for mnl_socket
 *
 * int mnl_socket_map_ring(struct mnl_socket *nl)
 */
func SocketMapRing(nl *Socket) error {
	_, err := C.mnl_socket_map_ring((*C.struct_mnl_socket)(nl))
	return err
}

/**
 * mnl_socket_unmap_ring - unmap a ring for mnl_socket
 *
 * int mnl_socket_unmap_ring(struct mnl_socket *nl)
 */
func SocketUnmapRing(nl *Socket) error {
	_, err := C.mnl_socket_unmap_ring((*C.struct_mnl_socket)(nl))
	return err
}

/**
 * mnl_socket_get_ring - get ring from mnl_socket
 *
 * struct mnl_ring *mnl_socket_get_ring(const struct mnl_socket *nl, enum mnl_ring_types type)
 */
func SocketGetRing(nl *Socket, rtype RingTypes) (*Ring, error) {
	ret, err := C.mnl_socket_get_ring((*C.struct_mnl_socket)(nl), (C.enum_mnl_ring_types)(rtype))
	return (*Ring)(unsafe.Pointer(ret)), err
}

/**
 * mnl_ring_get_frame - get current frame
 *
 * struct nl_mmap_hdr *mnl_ring_get_frame(const struct mnl_ring *ring)
 */
func RingGetFrame(ring *Ring) (*NlMmapHdr) {
	return (*NlMmapHdr)(unsafe.Pointer(C.mnl_ring_get_frame((*C.struct_mnl_ring)(ring))))
}

/**
 * mnl_ring_advance - set forward frame pointer
 *
 * void mnl_ring_advance(const struct mnl_ring *ring)
 */
func RingAdvance(ring *Ring) {
	C.mnl_ring_advance((*C.struct_mnl_ring)(ring))
}


// receivers
func (nl *Socket) SetRingopt(t RingTypes, bs, bn, fs, fn uint) error { return SocketSetRingopt(nl, t, bs, bn, fs, fn) }
func (nl *Socket) MapRing() error { return SocketMapRing(nl) }
func (nl *Socket) UnmapRing() error { return SocketUnmapRing(nl) }
func (nl *Socket) Ring(t RingTypes) (*Ring, error) { return SocketGetRing(nl, t) }
func (ring *Ring) Frame() *NlMmapHdr { return RingGetFrame(ring) }
func (ring *Ring) Advance() { RingAdvance(ring) }
