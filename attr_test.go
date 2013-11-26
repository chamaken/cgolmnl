package cgolmnl

import (
	"testing"
	"fmt"
)

func TestAttrParse(t *testing.T) {
	cb := func(attr *Nlattr, data interface{}) int {
		fmt.Printf("attr.Len: %d, data: %d\n", (*attr).Len, data.(int))
		return 1
	}

	b := make([]byte, 4096)
	nlh := NlmsgPutHeader(b)
	val := 0x12
	AttrPutU8(nlh, MNL_TYPE_U8, 0x10)
	AttrPutU8(nlh, MNL_TYPE_U8, 0x11)
	AttrPutU8(nlh, MNL_TYPE_U8, 0x12)
	AttrPutU8(nlh, MNL_TYPE_U8, 0x13)
	AttrParse(nlh, 0, cb, val)
}

func TestCAttrParse2(t *testing.T) {
	cb := func(attr *Nlattr, data interface{}) int {
		fmt.Printf("attr.Len: %d, data: %d\n", (*attr).Len, data.(int))
		return 1
	}

	b := make([]byte, 4096)
	nlh := NlmsgPutHeader(b)
	val := 0x12
	AttrPutU8(nlh, MNL_TYPE_U8, 0x10)
	AttrPutU8(nlh, MNL_TYPE_U8, 0x11)
	AttrPutU8(nlh, MNL_TYPE_U8, 0x12)
	AttrPutU8(nlh, MNL_TYPE_U8, 0x13)
	CAttrParse2(nlh, 0, cb, val)
}
