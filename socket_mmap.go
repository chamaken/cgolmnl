// +build nlmmap

package cgolmnl

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -lmnl
#include <libmnl/libmnl.h>

void *mnl_ring_msghdr(void *hdr)
{
	return MNL_RING_MSGHDR(hdr);
}
*/
import "C"

import "unsafe"

type MnlRingTypes C.enum_mnl_ring_types
const (
	MNL_RING_RX MnlRingTypes = C.MNL_RING_RX
	MNL_RING_TX MnlRingTypes = C.MNL_RING_TX
)

func RingMsghdr(hdr *NlMmapHdr) *Nlmsghdr {
	return (*Nlmsghdr)(C.mnl_ring_msghdr(unsafe.Pointer(hdr)))
}

/**
 * mnl_socket_set_ringopt - set ring opt to prepare for mnl_socket_map_ring()
 * int mnl_socket_set_ringopt(struct mnl_socket *nl, struct nl_mmap_req *req,
 *			      enum mnl_ring_types type)
 */
func SocketSetRingopt(nl *SocketDescriptor, req *NlMmapReq, rtype MnlRingTypes) (int, error) {
	ret, err := C.mnl_socket_set_ringopt((*[0]byte)(nl), (*C.struct_nl_mmap_req)(unsafe.Pointer(req)),
		(C.enum_mnl_ring_types)(rtype))
	return int(ret), err
}

/**
 * mnl_socket_map_ring - setup a ring for mnl_socket
 *
 * int mnl_socket_map_ring(struct mnl_socket *nl)
 */
func SocketMapRing(nl *SocketDescriptor) (int, error) {
	ret, err := C.mnl_socket_map_ring((*[0]byte)(nl))
	return int(ret), err
}

/**
 * mnl_socket_get_frame - get current frame
 *
 * struct nl_mmap_hdr *mnl_socket_get_frame(const struct mnl_socket *nl,
 *					    enum mnl_ring_types type)
 */
func SocketGetFrame(nl *SocketDescriptor, rtype MnlRingTypes) (*NlMmapHdr) {
	return (*NlMmapHdr)(unsafe.Pointer(C.mnl_socket_get_frame((*[0]byte)(nl), (C.enum_mnl_ring_types)(rtype))))
}

/**
 * mnl_socket_advance_ring - set forward frame pointer
 *
 * int mnl_socket_advance_ring(const struct mnl_socket *nl, enum mnl_ring_types type)
 */
func SocketAdvanceRing(nl *SocketDescriptor, rtype MnlRingTypes) int {
	return int(C.mnl_socket_advance_ring((*[0]byte)(nl), (C.enum_mnl_ring_types)(rtype)))
}


// receivers
func (nl *SocketDescriptor) SetRingopt(r *NlMmapReq, t MnlRingTypes) (int, error) { return SocketSetRingopt(nl, r, t) }
func (nl *SocketDescriptor) MapRing() (int, error) { return SocketMapRing(nl) }
func (nl *SocketDescriptor) Frame(t MnlRingTypes) *NlMmapHdr { return SocketGetFrame(nl, t) }
func (nl *SocketDescriptor) AdvanceRing(t MnlRingTypes) int { return SocketAdvanceRing(nl, t) }
