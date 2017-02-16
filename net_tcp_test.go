package procfs

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNetTCP_regex(t *testing.T) {

	var expected = []struct {
		Line     string
		Expected NetTCPLine
	}{
		{
			Line: "   0: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 16071 1 0000000000000000 100 0 0 10 0",
			Expected: NetTCPLine{
				Sl:            "0",
				LocalAddress:  "00000000:0016",
				RemoteAddress: "00000000:0000",
				St:            "0A",
				TxQueue:       "00000000",
				RxQueue:       "00000000",
				Tr:            "00",
				TmWhen:        "00000000",
				Retrnsmt:      "00000000",
				UID:           "0",
				Timeout:       "0",
				Inode:         "16071",
				RefCount:      "1",
				MemoryAddress: "0000000000000000",
			},
		},
		{
			Line: "  8: 260FC90A:01BB AE6857C8:C747 03 00000000:00000000 01:000000DC 00000005     0        0 0 2 0000000000000000",
			Expected: NetTCPLine{
				Sl:            "8",
				LocalAddress:  "260FC90A:01BB",
				RemoteAddress: "AE6857C8:C747",
				St:            "03",
				TxQueue:       "00000000",
				RxQueue:       "00000000",
				Tr:            "01",
				TmWhen:        "000000DC",
				Retrnsmt:      "00000005",
				UID:           "0",
				Timeout:       "0",
				Inode:         "0",
				RefCount:      "2",
				MemoryAddress: "0000000000000000",
			},
		},
	}

	m := NetTCP{}
	re := regexp.MustCompile(m.regex())

	for _, i := range expected {
		have, _ := m.netTCPLine(re, i.Line)

		if !reflect.DeepEqual(i.Expected, have) {
			t.Errorf("NetTCP Regex is invalid")
		}
	}

}

func TestNetTCP(t *testing.T) {
	have, err := FS("fixtures").NewNetTCP()
	if err != nil {
		t.Fatal(err)
	}

	if len(have) != 19 {
		t.Errorf("incorrect number of NetTCPLine items.")
	}
}
