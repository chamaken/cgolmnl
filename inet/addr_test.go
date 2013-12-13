package inet

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	var ia IPAddr

	ia.UnmarshalText(([]byte)("192.168.1.1"))
	if ia.String() != "192.168.1.1" {
		t.Errorf("want: 192.168.1.1, but: %s", ia.String())
	}
}
