package procfs

import (
	"testing"
)

func TestXfrmStats(t *testing.T) {
	xfrmStats, err := FS("fixtures").NewXfrmStat()

	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		want int
		got  int
	}{
		{name: "XfrmInError", want: 1, got: xfrmStats.XfrmInError},
		{name: "XfrmInBufferError", want: 2, got: xfrmStats.XfrmInBufferError},
		{name: "XfrmInHdrError", want: 4, got: xfrmStats.XfrmInHdrError},
		{name: "XfrmInNoStates", want: 3, got: xfrmStats.XfrmInNoStates},
		{name: "XfrmInStateProtoError", want: 40, got: xfrmStats.XfrmInStateProtoError},
		{name: "XfrmInStateModeError", want: 100, got: xfrmStats.XfrmInStateModeError},
		{name: "XfrmInStateSeqError", want: 6000, got: xfrmStats.XfrmInStateSeqError},
		{name: "XfrmInStateExpired", want: 4, got: xfrmStats.XfrmInStateExpired},
		{name: "XfrmInStateMismatch", want: 23451, got: xfrmStats.XfrmInStateMismatch},
		{name: "XfrmInStateInvalid", want: 55555, got: xfrmStats.XfrmInStateInvalid},
		{name: "XfrmInTmplMismatch", want: 51, got: xfrmStats.XfrmInTmplMismatch},
		{name: "XfrmInNoPols", want: 65432, got: xfrmStats.XfrmInNoPols},
		{name: "XfrmInPolBlock", want: 100, got: xfrmStats.XfrmInPolBlock},
		{name: "XfrmInPolError", want: 10000, got: xfrmStats.XfrmInPolError},
		{name: "XfrmOutError", want: 1000000, got: xfrmStats.XfrmOutError},
		{name: "XfrmOutBundleGenError", want: 43321, got: xfrmStats.XfrmOutBundleGenError},
		{name: "XfrmOutBundleCheckError", want: 555, got: xfrmStats.XfrmOutBundleCheckError},
		{name: "XfrmOutNoStates", want: 869, got: xfrmStats.XfrmOutNoStates},
		{name: "XfrmOutStateProtoError", want: 4542, got: xfrmStats.XfrmOutStateProtoError},
		{name: "XfrmOutStateModeError", want: 4, got: xfrmStats.XfrmOutStateModeError},
		{name: "XfrmOutStateSeqError", want: 543, got: xfrmStats.XfrmOutStateSeqError},
		{name: "XfrmOutStateExpired", want: 565, got: xfrmStats.XfrmOutStateExpired},
		{name: "XfrmOutPolBlock", want: 43456, got: xfrmStats.XfrmOutPolBlock},
		{name: "XfrmOutPolDead", want: 7656, got: xfrmStats.XfrmOutPolDead},
		{name: "XfrmOutPolError", want: 1454, got: xfrmStats.XfrmOutPolError},
		{name: "XfrmFwdHdrError", want: 6654, got: xfrmStats.XfrmFwdHdrError},
		{name: "XfrmOutStateInvaliad", want: 28765, got: xfrmStats.XfrmOutStateInvalid},
		{name: "XfrmAcquireError", want: 24532, got: xfrmStats.XfrmAcquireError},
		{name: "XfrmInStateInvalid", want: 55555, got: xfrmStats.XfrmInStateInvalid},
		{name: "XfrmOutError", want: 1000000, got: xfrmStats.XfrmOutError},
	} {
		if test.want != test.got {
			t.Errorf("Want %s %d, have %d", test.name, test.want, test.got)
		}
	}
}
