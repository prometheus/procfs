package procfs

import (
	"testing"
)

func TestSys(t *testing.T) {
	have, err := FS("fixtures").NewSys()
	if err != nil {
		t.Fatal(err)
	}

	if have["fixtures/sys/vm/user_reserve_kbytes"] != "29155" {
		t.Errorf("doesn't have the user_reserve_kbytes key '%s'", have["fixtures/sys/vm/user_reserve_kbytes"])
	}
}
