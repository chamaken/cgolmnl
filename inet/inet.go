package inet

/*
#include <arpa/inet.h>
#include <net/if.h>
#include <stdlib.h>
*/
import "C"

import (
	"encoding/binary"
	"net"
	"unsafe"
)

func Ntohl(i uint32) uint32 {
	return binary.BigEndian.Uint32((*(*[4]byte)(unsafe.Pointer(&i)))[:])
}
func Htonl(i uint32) uint32 {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return *(*uint32)(unsafe.Pointer(&b[0]))
}

func Ntohs(i uint16) uint16 {
	return binary.BigEndian.Uint16((*(*[2]byte)(unsafe.Pointer(&i)))[:])
}
func Htons(i uint16) uint16 {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return *(*uint16)(unsafe.Pointer(&b[0]))
}

func Be64toh(i uint64) uint64 {
	return binary.BigEndian.Uint64((*(*[8]byte)(unsafe.Pointer(&i)))[:])
}

func InetAddr(s string) uint32 {
	return binary.BigEndian.Uint32(net.ParseIP(s).To4())
}

func InetNtoa(p unsafe.Pointer) string {
	return net.IP((*(*[net.IPv4len]byte)(p))[:]).String()
}

func Inet6Ntoa(p unsafe.Pointer) string {
	return net.IP((*(*[net.IPv6len]byte)(p))[:]).String()
}

func InetNtop(family int, src unsafe.Pointer) string {
	switch family {
	case C.AF_INET:
		return net.IP((*(*[net.IPv4len]byte)(src))[:]).String()
	case C.AF_INET6:
		return net.IP((*(*[net.IPv6len]byte)(src))[:]).String()
	default:
		return "(unknown family)"
	}
	return "(unknown family)"
}

func InetPton(family int, src string) ([]byte, error) {
	var dst []byte
	switch family {
	case C.AF_INET:
		dst = make([]byte, net.IPv4len)
	case C.AF_INET6:
		dst = make([]byte, net.IPv6len)
	}
	cs := C.CString(src)
	defer C.free(unsafe.Pointer(cs))
	if _, err := C.inet_pton(C.int(family), cs, unsafe.Pointer(&dst[0])); err != nil {
		return nil, err
	}
	return dst, nil
}

func IfNametoindex(name string) (int, error) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))
	ret, err := C.if_nametoindex(cs)
	if err != nil {
		return 0, err
	}
	return int(ret), nil
}
