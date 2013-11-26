package cgolmnl

import (
	"testing"
	// "fmt"
)

func TestNlmsgPutHeader(t *testing.T) {
	b := make([]byte, 32)
	nlh := *NlmsgPutHeader(b)
	if nlh.Len != MNL_NLMSG_HDRLEN {
		t.Errorf("want: %d, returns: %d", MNL_NLMSG_HDRLEN, nlh.Len)
	}
}
