package inet

/*
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <net/if.h>
#include <stdlib.h>
*/
import "C"

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
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


type IPAddr []byte // big endian, verdict by len

func (ia IPAddr) String() string {
	if b, err := ia.MarshalText(); err != nil {
		return fmt.Sprintf("(None - error: %v)", err)
	} else {
		return (string)(b)
	}
}

func (ia IPAddr) MarshalBinary() ([]byte, error) {
	dst := make([]byte, len(ia))
	if len(ia) == 4 { // IPv4
		binary.BigEndian.PutUint32(dst, *(*uint32)(unsafe.Pointer(&ia[0])))
		return dst, nil
	}
	copy(dst, ia)
	return dst, nil
}

func (ia IPAddr) MarshalText() ([]byte, error) {
	switch len(ia) {
	case net.IPv4len:
		a := make([]string, net.IPv4len)
		for i := 0; i < net.IPv4len; i++ {
			a[i] = fmt.Sprintf("%d", ia[i])
		}
		return ([]byte)(strings.Join(a, ".")), nil
	case net.IPv6len:
		a := make([]string, net.IPv6len)
		for i := 0; i < net.IPv6len; i++ {
			a[i] = fmt.Sprintf("%d", ia[i])
		}
		return ([]byte)(strings.Join(a, ":")), nil
	default:
		return nil, errors.New(fmt.Sprintf("invalid address len: %d", len(ia)))
	}
}

func (ia *IPAddr) UnmarshalText(src []byte) error {
	cs := C.CString(string(src))
	defer C.free(unsafe.Pointer(cs))
	// try AF_INET first
	dst := make([]byte, net.IPv4len)
	ret, err := C.inet_pton(C.AF_INET, cs, unsafe.Pointer(&dst[0]))
	if err != nil { // ret == -1 and errno is EAFNOTSUPPORT
		return err 
	}
	if ret == 1 {
		*ia = IPAddr(dst)
	}

	// ret == 0: not a valid network address in AF_INET, try AF_INET6
	dst = make([]byte, net.IPv6len)
	ret, err = C.inet_pton(C.AF_INET6, cs, unsafe.Pointer(&dst[0]))
	if err != nil { return err }
	if ret == 0 {
		return errors.New("not a valid address string")
	}
	*ia = IPAddr(dst)
	return nil
}


func InetAddr(s string) uint32 {
	ip := net.ParseIP(s)
	d := strings.Split(ip.String(), ".")
	d0, _ := strconv.Atoi(d[0])
	d1, _ := strconv.Atoi(d[1])
	d2, _ := strconv.Atoi(d[2])
	d3, _ := strconv.Atoi(d[3])

	r := uint32(d0 << 24)
	r += uint32(d1 << 16)
	r += uint32(d2 <<  8)
	r += uint32(d3)

	return r
}

func InetNtoa(p unsafe.Pointer) string {
	in := IPAddr((*(*[net.IPv4len]byte)(p))[:])
	return in.String()
}

func Inet6Ntoa(p unsafe.Pointer) string {
	in := IPAddr((*(*[net.IPv6len]byte)(p))[:])
	return in.String()
}

func InetNtop(family int, src unsafe.Pointer) (string, error) {
	var in IPAddr
	switch family {
	case C.AF_INET:
		in = IPAddr((*(*[net.IPv4len]byte)(src))[:])
	case C.AF_INET6:
		in = IPAddr((*(*[net.IPv6len]byte)(src))[:])
	default:
		return "", errors.New("invalid address family")
	}
	return in.String(), nil
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
