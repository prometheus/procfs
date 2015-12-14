package procfs

import (
	"reflect"
	"testing"
)

func TestNetSockstat(t *testing.T) {
	expected := NetSockstat{}
	expected.Sockets.Used = 3014
	expected.TCP.InUse = 3168
	expected.TCP.Orphan = 612
	expected.TCP.Tw = 13701
	expected.TCP.Alloc = 3169
	expected.TCP.Mem = 2734
	expected.UDP.InUse = 9
	expected.UDP.Mem = 2
	expected.UDPLite.InUse = 0
	expected.RAW.InUse = 1
	expected.FRAG.InUse = 2
	expected.FRAG.Memory = 2

	have, err := FS("fixtures").NewNetSockstat()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, have) {
		t.Errorf("structs are not equal")
	}
}
