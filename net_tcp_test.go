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
	/*
		expected := NetTCP{
			MemTotal:          15666184,
			MemFree:           440324,
			Buffers:           1020128,
			Cached:            12007640,
			SwapCached:        0,
			Active:            6761276,
			Inactive:          6532708,
			ActiveAnon:        267256,
			InactiveAnon:      268,
			ActiveFile:        6494020,
			InactiveFile:      6532440,
			Unevictable:       0,
			Mlocked:           0,
			SwapTotal:         0,
			SwapFree:          0,
			Dirty:             768,
			Writeback:         0,
			AnonPages:         266216,
			Mapped:            44204,
			Shmem:             1308,
			Slab:              1807264,
			SReclaimable:      1738124,
			SUnreclaim:        69140,
			KernelStack:       1616,
			PageTables:        5288,
			NFSUnstable:       0,
			Bounce:            0,
			WritebackTmp:      0,
			CommitLimit:       7833092,
			CommittedAS:       530844,
			VmallocTotal:      34359738367,
			VmallocUsed:       36596,
			VmallocChunk:      34359637840,
			HardwareCorrupted: 0,
			AnonHugePages:     12288,
			HugePagesTotal:    0,
			HugePagesFree:     0,
			HugePagesRsvd:     0,
			HugePagesSurp:     0,
			Hugepagesize:      2048,
			DirectMap4k:       91136,
			DirectMap2M:       16039936,
		}

		have, err := FS("fixtures").NewNetTCP()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(have, expected) {
			t.Errorf("structs are not equal")
		}
	*/
}
