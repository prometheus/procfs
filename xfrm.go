package procfs

import (
	"fmt"
	"io/ioutil"
	"os"
)

// XfrmStat models the content of /proc/net/xfrm_stat
type XfrmStat struct {
	// All errors which are not matched by other
	XfrmInError int
	// No buffer is left
	XfrmInBufferError int
	// Header Error
	XfrmInHdrError int
	// No state found
	// i.e. either inbound SPI, address, or IPSEC protocol at SA is wrong
	XfrmInNoStates int
	// Transformation protocol specific error
	// e.g. SA Key is wrong
	XfrmInStateProtoError int
	// Transformation mode specific error
	XfrmInStateModeError int
	// Sequence error
	// e.g. sequence number is out of window
	XfrmInStateSeqError int
	// State is expired
	XfrmInStateExpired int
	// State has mismatch option
	// e.g. UDP encapsulation type is mismatched
	XfrmInStateMismatch int
	// State is invalid
	XfrmInStateInvalid int
	// No matching template for states
	// e.g. Inbound SAs are correct but SP rule is wrong
	XfrmInTmplMismatch int
	// No policy is found for states
	// e.g. Inbound SAs are correct but no SP is found
	XfrmInNoPols int
	// Policy discards
	XfrmInPolBlock int
	// Policy error
	XfrmInPolError int
	// All errors which are not matched by others
	XfrmOutError int
	// Bundle generation error
	XfrmOutBundleGenError int
	// Bundle check error
	XfrmOutBundleCheckError int
	// No state was found
	XfrmOutNoStates int
	// Transformation protocol specific error
	XfrmOutStateProtoError int
	// Transportation mode specific error
	XfrmOutStateModeError int
	// Sequence error
	// i.e sequence number overflow
	XfrmOutStateSeqError int
	// State is expired
	XfrmOutStateExpired int
	// Policy discads
	XfrmOutPolBlock int
	// Policy is dead
	XfrmOutPolDead int
	// Policy Error
	XfrmOutPolError     int
	XfrmFwdHdrError     int
	XfrmOutStateInvalid int
	XfrmAcquireError    int
}

// NewXfrmStat reads the xfrm_stat info
func NewXfrmStat() (XfrmStat, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return XfrmStat{}, err
	}

	return fs.NewXfrmStat()
}

// NewXfrmStat reads the xfrm_stat statistics from the specified `proc` filesystem
func (fs FS) NewXfrmStat() (XfrmStat, error) {
	x := XfrmStat{}

	file, err := os.Open(fs.Path("net/xfrm_stat"))
	if err != nil {
		return x, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return x, err
	}

	ioFormat := "XfrmInError %d\nXfrmInBufferError %d\nXfrmInHdrError %d\n" +
		"XfrmInNoStates %d\nXfrmInStateProtoError %d\nXfrmInStateModeError %d\n" +
		"XfrmInStateSeqError %d\nXfrmInStateExpired %d\nXfrmInStateMismatch %d\n" +
		"XfrmInStateInvalid %d\nXfrmInTmplMismatch %d\nXfrmInNoPols %d\n" +
		"XfrmInPolBlock %d\nXfrmInPolError %d\n" +
		"XfrmOutError %d\nXfrmOutBundleGenError %d\nXfrmOutBundleCheckError %d\n" +
		"XfrmOutNoStates %d\nXfrmOutStateProtoError %d\nXfrmOutStateModeError %d\n" +
		"XfrmOutStateSeqError %d\nXfrmOutStateExpired %d\nXfrmOutPolBlock %d\n" +
		"XfrmOutPolDead %d\nXfrmOutPolError %d\nXfrmFwdHdrError %d\nXfrmOutStateInvalid %d\n" +
		"XfrmAcquireError %d\n"

	_, err = fmt.Sscanf(string(data), ioFormat, &x.XfrmInError, &x.XfrmInBufferError,
		&x.XfrmInHdrError, &x.XfrmInNoStates, &x.XfrmInStateProtoError,
		&x.XfrmInStateModeError, &x.XfrmInStateSeqError, &x.XfrmInStateExpired,
		&x.XfrmInStateMismatch, &x.XfrmInStateInvalid, &x.XfrmInTmplMismatch, &x.XfrmInNoPols,
		&x.XfrmInPolBlock, &x.XfrmInPolError, &x.XfrmOutError, &x.XfrmOutBundleGenError, &x.XfrmOutBundleCheckError,
		&x.XfrmOutNoStates, &x.XfrmOutStateProtoError, &x.XfrmOutStateModeError,
		&x.XfrmOutStateSeqError, &x.XfrmOutStateExpired, &x.XfrmOutPolBlock, &x.XfrmOutPolDead,
		&x.XfrmOutPolError, &x.XfrmFwdHdrError, &x.XfrmOutStateInvalid,
		&x.XfrmAcquireError)

	if err != nil {
		return x, err
	}

	return x, nil
}
