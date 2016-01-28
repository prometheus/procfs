package procfs

import (
	"bufio"
	"os"
	"reflect"
	"regexp"
)

type ProcStatus struct {
	Name                     string `proc_status:"Name"`
	State                    string `proc_status:"State"`
	Tgid                     string `proc_status:"Tgid"`
	Ngid                     string `proc_status:"Ngid"`
	Pid                      string `proc_status:"Pid"`
	PPid                     string `proc_status:"PPid"`
	TracerPid                string `proc_status:"TracerPid"`
	Uid                      string `proc_status:"Uid"`
	Gid                      string `proc_status:"Gid"`
	FDSize                   string `proc_status:"FDSize"`
	Groups                   string `proc_status:"Groups"`
	VmPeak                   string `proc_status:"VmPeak"`
	VmSize                   string `proc_status:"VmSize"`
	VmLck                    string `proc_status:"VmLck"`
	VmPin                    string `proc_status:"VmPin"`
	VmHWM                    string `proc_status:"VmHWM"`
	VmRSS                    string `proc_status:"VmRSS"`
	VmData                   string `proc_status:"VmData"`
	VmStk                    string `proc_status:"VmStk"`
	VmExe                    string `proc_status:"VmExe"`
	VmLib                    string `proc_status:"VmLib"`
	VmPTE                    string `proc_status:"VmPTE"`
	VmSwap                   string `proc_status:"VmSwap"`
	Threads                  string `proc_status:"Threads"`
	SigQ                     string `proc_status:"SigQ"`
	SigPnd                   string `proc_status:"SigPnd"`
	ShdPnd                   string `proc_status:"ShdPnd"`
	SigBlk                   string `proc_status:"SigBlk"`
	SigIgn                   string `proc_status:"SigIgn"`
	SigCgt                   string `proc_status:"SigCgt"`
	CapInh                   string `proc_status:"CapInh"`
	CapPrm                   string `proc_status:"CapPrm"`
	CapEff                   string `proc_status:"CapEff"`
	CapBnd                   string `proc_status:"CapBnd"`
	Seccomp                  string `proc_status:"Seccomp"`
	CpusAllowed              string `proc_status:"Cpus_allowed"`
	CpusAllowedList          string `proc_status:"Cpus_allowed_list"`
	MemsAllowed              string `proc_status:"Mems_allowed"`
	MemsAllowedList          string `proc_status:"Mems_allowed_list"`
	VoluntaryCtxtSwitches    string `proc_status:"voluntary_ctxt_switches"`
	NonvoluntaryCtxtSwitches string `proc_status:"nonvoluntary_ctxt_switches"`

	fs FS
}

func (ps ProcStatus) regex() string {
	return "([A-Za-z0-9()_]*):[ \t]*([A-Za-z0-9]*).*$"
}

// NewStatus returns the current status information of the process.
func (p Proc) NewStatus() (ps ProcStatus, err error) {
	f, err := os.Open(p.path("status"))
	if err != nil {
		return ProcStatus{}, err
	}
	defer f.Close()

	st := reflect.TypeOf(ps)
	re := regexp.MustCompile(ps.regex())
	s := bufio.NewScanner(f)

	for s.Scan() {
		line := s.Text()

		submatch := re.FindAllStringSubmatch(line, 1)
		if submatch == nil {
			continue
		}

		key := submatch[0][1]
		val := submatch[0][2]

		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			if field.Tag.Get("proc_status") == key {
				// v, err := strconv.ParseInt(val, 10, 64)
				// if err != nil {
				// 	// no op
				// }
				reflect.ValueOf(&ps).Elem().Field(i).SetString(val)
			}
		}
	}

	return ps, nil
}
