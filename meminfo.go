package procfs

import (
	"bufio"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

// Meminfo represents memory statistics.
type Meminfo struct {
	MemTotal          int64 `meminfo:"MemTotal"`
	MemFree           int64 `meminfo:"MemFree"`
	Buffers           int64 `meminfo:"Buffers"`
	Cached            int64 `meminfo:"Cached"`
	SwapCached        int64 `meminfo:"SwapCached"`
	Active            int64 `meminfo:"Active"`
	Inactive          int64 `meminfo:"Inactive"`
	ActiveAnon        int64 `meminfo:"Active(anon)"`
	InactiveAnon      int64 `meminfo:"Inactive(anon)"`
	ActiveFile        int64 `meminfo:"Active(file)"`
	InactiveFile      int64 `meminfo:"Inactive(file)"`
	Unevictable       int64 `meminfo:"Unevictable"`
	Mlocked           int64 `meminfo:"Mlocked"`
	SwapTotal         int64 `meminfo:"SwapTotal"`
	SwapFree          int64 `meminfo:"SwapFree"`
	Dirty             int64 `meminfo:"Dirty"`
	Writeback         int64 `meminfo:"Writeback"`
	AnonPages         int64 `meminfo:"AnonPages"`
	Mapped            int64 `meminfo:"Mapped"`
	Shmem             int64 `meminfo:"Shmem"`
	Slab              int64 `meminfo:"Slab"`
	SReclaimable      int64 `meminfo:"SReclaimable"`
	SUnreclaim        int64 `meminfo:"SUnreclaim"`
	KernelStack       int64 `meminfo:"KernelStack"`
	PageTables        int64 `meminfo:"PageTables"`
	NFSUnstable       int64 `meminfo:"NFS_Unstable"`
	Bounce            int64 `meminfo:"Bounce"`
	WritebackTmp      int64 `meminfo:"WritebackTmp"`
	CommitLimit       int64 `meminfo:"CommitLimit"`
	CommittedAS       int64 `meminfo:"Committed_AS"`
	VmallocTotal      int64 `meminfo:"VmallocTotal"`
	VmallocUsed       int64 `meminfo:"VmallocUsed"`
	VmallocChunk      int64 `meminfo:"VmallocChunk"`
	HardwareCorrupted int64 `meminfo:"HardwareCorrupted"`
	AnonHugePages     int64 `meminfo:"AnonHugePages"`
	HugePagesTotal    int64 `meminfo:"HugePages_Total"`
	HugePagesFree     int64 `meminfo:"HugePages_Free"`
	HugePagesRsvd     int64 `meminfo:"HugePages_Rsvd"`
	HugePagesSurp     int64 `meminfo:"HugePages_Surp"`
	Hugepagesize      int64 `meminfo:"Hugepagesize"`
	DirectMap4k       int64 `meminfo:"DirectMap4k"`
	DirectMap2M       int64 `meminfo:"DirectMap2M"`
}

// NewMeminfo returns kernel/system statistics read from /proc/stat.
func NewMeminfo() (Meminfo, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Meminfo{}, err
	}

	return fs.NewMeminfo()
}

// NewMeminfo returns an information about current kernel/system statistics.
func (fs FS) NewMeminfo() (m Meminfo, err error) {
	f, err := os.Open(fs.Path("meminfo"))
	if err != nil {
		return Meminfo{}, err
	}
	defer f.Close()

	st := reflect.TypeOf(m)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		re := regexp.MustCompile(m.Regex())
		submatch := re.FindAllStringSubmatch(line, 1)
		if submatch == nil {
			continue
		}

		key := submatch[0][1]
		val := submatch[0][2]

		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			if field.Tag.Get("meminfo") == key {
				v, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					// no op
				}
				reflect.ValueOf(&m).Elem().Field(i).SetInt(v)
			}
		}
	}

	return m, nil
}

func (m Meminfo) Regex() string {
	return "([A-Za-z()]*): *([0-9]*).*$"
}
