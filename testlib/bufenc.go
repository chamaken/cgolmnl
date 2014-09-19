package testlib

import (
	"C"
	"encoding/binary"
	"unsafe"
)

const SizeofCint = C.sizeof_int

// C int
func GetCint(bs []byte, start int) int {
	return int(*(*C.int)(unsafe.Pointer(&bs[start : start+SizeofCint][0])))
}
func SetCint(bs []byte, start int, val int) {
	*(*C.int)(unsafe.Pointer(&bs[start : start+SizeofCint][0])) = C.int(val)
}

// C unsigned int
func GetCuint(bs []byte, start int) uint {
	return uint(*(*C.uint)(unsafe.Pointer(&bs[start : start+SizeofCint][0])))
}
func SetCuint(bs []byte, start int, val uint) {
	*(*C.uint)(unsafe.Pointer(&bs[start : start+SizeofCint][0])) = C.uint(val)
}

// using encofing/binary
var Endian binary.ByteOrder

func init() {
	var v uint16 = 1
	if *(*byte)(unsafe.Pointer(&v)) == 1 {
		Endian = binary.LittleEndian
	} else {
		Endian = binary.BigEndian
	}
}

// 1 byte
func GetInt8(bs []byte, start int) int8 {
	return int8(bs[start])
}
func SetInt8(bs []byte, start int, val int8) {
	bs[start] = byte(val)
}
func GetUint8(bs []byte, start uint) uint8 {
	return uint8(bs[start])
}
func SetUint8(bs []byte, start uint, val uint8) {
	bs[start] = byte(val)
}

// 2 byte
func GetInt16(bs []byte, start int) int16 {
	return int16(Endian.Uint16(bs[start : start+2]))
}
func SetInt16(bs []byte, start int, val int16) {
	Endian.PutUint16(bs[start:start+2], uint16(val))
}
func GetUint16(bs []byte, start uint) uint16 {
	return Endian.Uint16(bs[start : start+2])
}
func SetUint16(bs []byte, start uint, val uint16) {
	Endian.PutUint16(bs[start:start+2], val)
}

// 4 byte
func GetInt32(bs []byte, start int) int32 {
	return int32(Endian.Uint32(bs[start : start+4]))
}
func SetInt32(bs []byte, start int, val int32) {
	Endian.PutUint32(bs[start:start+4], uint32(val))
}
func GetUint32(bs []byte, start uint) uint32 {
	return Endian.Uint32(bs[start : start+4])
}
func SetUint32(bs []byte, start uint, val uint32) {
	Endian.PutUint32(bs[start:start+4], val)
}

// 8 byte
func GetInt64(bs []byte, start int) int64 {
	return int64(Endian.Uint64(bs[start : start+8]))
}
func SetInt64(bs []byte, start int, val int64) {
	Endian.PutUint64(bs[start:start+8], uint64(val))
}
func GetUint64(bs []byte, start uint) uint64 {
	return Endian.Uint64(bs[start : start+8])
}
func SetUint64(bs []byte, start uint, val uint64) {
	Endian.PutUint64(bs[start:start+8], val)
}

// for struct nlmsghdr
type NlmsghdrBuf []byte

const (
	nlmsghdr_len_index     = 0  // __u32	nlmsg_len
	nlmsghdr_type_index    = 4  // __u16	nlmsg_type
	nlmsghdr_flags_index   = 6  // __u16	nlmsg_flags
	nlmsghdr_seq_index     = 8  // __u32	nlmsg_seq
	nlmsghdr_pid_index     = 12 // __u32	nlmsg_pid
	nlmsghdr_payload_index = 16
)

func NewNlmsghdrBuf(size int) *NlmsghdrBuf {
	nlb := NlmsghdrBuf(make([]byte, size))
	return &nlb
}
func (nlh *NlmsghdrBuf) SetLen(nlmsg_len uint32) {
	SetUint32(*nlh, nlmsghdr_len_index, nlmsg_len)
}
func (nlh *NlmsghdrBuf) Len() uint32 {
	return GetUint32(*nlh, nlmsghdr_len_index)
}
func (nlh *NlmsghdrBuf) SetType(nlmsg_type uint16) {
	SetUint16(*nlh, nlmsghdr_type_index, nlmsg_type)
}
func (nlh *NlmsghdrBuf) Type() uint16 {
	return GetUint16(*nlh, nlmsghdr_type_index)
}
func (nlh *NlmsghdrBuf) SetFlags(nlmsg_flags uint16) {
	SetUint16(*nlh, nlmsghdr_flags_index, nlmsg_flags)
}
func (nlh *NlmsghdrBuf) Flags() uint16 {
	return GetUint16(*nlh, nlmsghdr_flags_index)
}
func (nlh *NlmsghdrBuf) SetSeq(nlmsg_seq uint32) {
	SetUint32(*nlh, nlmsghdr_seq_index, nlmsg_seq)
}
func (nlh *NlmsghdrBuf) Seq() uint32 {
	return GetUint32(*nlh, nlmsghdr_seq_index)
}
func (nlh *NlmsghdrBuf) SetPid(nlmsg_pid uint32) {
	SetUint32(*nlh, nlmsghdr_pid_index, nlmsg_pid)
}
func (nlh *NlmsghdrBuf) Pid() uint32 {
	return GetUint32(*nlh, nlmsghdr_pid_index)
}
func (nlh *NlmsghdrBuf) SetPayload(payload []byte) {
	copy((*(*[]byte)(nlh))[nlmsghdr_payload_index:], payload)
}
func (nlh *NlmsghdrBuf) Payload() []byte {
	return (*(*[]byte)(nlh))[nlmsghdr_payload_index:]
}

// for struct nlattr
type NlattrBuf []byte

const (
	nla_len_index     = 0 // __u16	nla_len
	nla_type_index    = 2 // __u16	nla_type
	nla_payload_index = 4
)

func NewNlattrBuf(size int) *NlattrBuf {
	nlb := NlattrBuf(make([]byte, size))
	return &nlb
}
func (nla *NlattrBuf) SetLen(nla_len uint16) {
	SetUint16(*nla, nla_len_index, nla_len)
}
func (nla *NlattrBuf) Len() uint16 {
	return GetUint16(*nla, nla_len_index)
}
func (nla *NlattrBuf) SetType(nla_type uint16) {
	SetUint16(*nla, nla_type_index, nla_type)
}
func (nla *NlattrBuf) Type() uint16 {
	return GetUint16(*nla, nla_type_index)
}
func (nla *NlattrBuf) SetPayload(payload []byte) {
	copy((*(*[]byte)(nla))[:nla_payload_index], payload)
}
func (nla *NlattrBuf) Payload() []byte {
	return (*(*[]byte)(nla))[:nla_payload_index]
}
