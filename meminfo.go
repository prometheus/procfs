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
	// Total usable ram (i.e. physical ram minus a few reserved
	// bits and the kernel binary code)
	MemTotal int64 `meminfo:"MemTotal"`
	// The sum of LowFree+HighFree
	MemFree int64 `meminfo:"MemFree"`
	// An estimate of how much memory is available for starting
	// new applications, without swapping. Calculated from
	// MemFree, SReclaimable, the size of the file LRU lists, and
	// the low watermarks in each zone.  The estimate takes into
	// account that the system needs some page cache to function
	// well, and that not all reclaimable slab will be
	// reclaimable, due to items being in use. The impact of those
	// factors will vary from system to system.
	MemAvailable int64 `meminfo:"MemAvailable"`
	// Relatively temporary storage for raw disk blocks shouldn't
	// get tremendously large (20MB or so)
	Buffers int64 `meminfo:"Buffers"`
	Cached  int64 `meminfo:"Cached"`
	// Memory that once was swapped out, is swapped back in but
	// still also is in the swapfile (if memory is needed it
	// doesn't need to be swapped out AGAIN because it is already
	// in the swapfile. This saves I/O)
	SwapCached int64 `meminfo:"SwapCached"`
	// Memory that has been used more recently and usually not
	// reclaimed unless absolutely necessary.
	Active int64 `meminfo:"Active"`
	// Memory which has been less recently used.  It is more
	// eligible to be reclaimed for other purposes
	Inactive     int64 `meminfo:"Inactive"`
	ActiveAnon   int64 `meminfo:"Active(anon)"`
	InactiveAnon int64 `meminfo:"Inactive(anon)"`
	ActiveFile   int64 `meminfo:"Active(file)"`
	InactiveFile int64 `meminfo:"Inactive(file)"`
	Unevictable  int64 `meminfo:"Unevictable"`
	Mlocked      int64 `meminfo:"Mlocked"`
	// total amount of swap space available
	SwapTotal int64 `meminfo:"SwapTotal"`
	// Memory which has been evicted from RAM, and is temporarily
	// on the disk
	SwapFree int64 `meminfo:"SwapFree"`
	// Memory which is waiting to get written back to the disk
	Dirty int64 `meminfo:"Dirty"`
	// Memory which is actively being written back to the disk
	Writeback int64 `meminfo:"Writeback"`
	// Non-file backed pages mapped into userspace page tables
	AnonPages int64 `meminfo:"AnonPages"`
	// files which have been mmaped, such as libraries
	Mapped int64 `meminfo:"Mapped"`
	Shmem  int64 `meminfo:"Shmem"`
	// in-kernel data structures cache
	Slab int64 `meminfo:"Slab"`
	// Part of Slab, that might be reclaimed, such as caches
	SReclaimable int64 `meminfo:"SReclaimable"`
	// Part of Slab, that cannot be reclaimed on memory pressure
	SUnreclaim  int64 `meminfo:"SUnreclaim"`
	KernelStack int64 `meminfo:"KernelStack"`
	// amount of memory dedicated to the lowest level of page
	// tables.
	PageTables int64 `meminfo:"PageTables"`
	// NFS pages sent to the server, but not yet committed to
	// stable storage
	NFSUnstable int64 `meminfo:"NFS_Unstable"`
	// Memory used for block device "bounce buffers"
	Bounce int64 `meminfo:"Bounce"`
	// Memory used by FUSE for temporary writeback buffers
	WritebackTmp int64 `meminfo:"WritebackTmp"`
	// Based on the overcommit ratio ('vm.overcommit_ratio'),
	// this is the total amount of  memory currently available to
	// be allocated on the system. This limit is only adhered to
	// if strict overcommit accounting is enabled (mode 2 in
	// 'vm.overcommit_memory').
	// The CommitLimit is calculated with the following formula:
	// CommitLimit = ([total RAM pages] - [total huge TLB pages]) *
	//                overcommit_ratio / 100 + [total swap pages]
	// For example, on a system with 1G of physical RAM and 7G
	// of swap with a `vm.overcommit_ratio` of 30 it would
	// yield a CommitLimit of 7.3G.
	// For more details, see the memory overcommit documentation
	// in vm/overcommit-accounting.
	CommitLimit int64 `meminfo:"CommitLimit"`
	// The amount of memory presently allocated on the system.
	// The committed memory is a sum of all of the memory which
	// has been allocated by processes, even if it has not been
	// "used" by them as of yet. A process which malloc()'s 1G
	// of memory, but only touches 300M of it will show up as
	// using 1G. This 1G is memory which has been "committed" to
	// by the VM and can be used at any time by the allocating
	// application. With strict overcommit enabled on the system
	// (mode 2 in 'vm.overcommit_memory'),allocations which would
	// exceed the CommitLimit (detailed above) will not be permitted.
	// This is useful if one needs to guarantee that processes will
	// not fail due to lack of memory once that memory has been
	// successfully allocated.
	CommittedAS int64 `meminfo:"Committed_AS"`
	// total size of vmalloc memory area
	VmallocTotal int64 `meminfo:"VmallocTotal"`
	// amount of vmalloc area which is used
	VmallocUsed int64 `meminfo:"VmallocUsed"`
	// largest contiguous block of vmalloc area which is free
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
	re := regexp.MustCompile(m.regex())
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

func (m Meminfo) regex() string {
	return "([A-Za-z0-9()_]*): *([0-9]*).*$"
}
