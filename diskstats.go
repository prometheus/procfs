package procfs

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type DiskstatLine struct {
	// 1 - major number
	MajorVersion int64
	// 2 - minor mumber
	MinorNumber int64
	// 3 - device name
	DeviceName int64
	// 4 - reads completed successfully
	ReadsCompleted int64
	// 5 - reads merged
	ReadsMerged int64
	// 6 - sectors read
	SectorsRead int64
	// 7 - time spent reading (ms)
	TimeSpentReading int64
	// 8 - writes completed
	WritesCompleted int64
	// 9 - writes merged
	WritesMerged int64
	// 10 - sectors written
	SectorsWritten int64
	// 11 - time spent writing (ms)
	TimeSpentWriting int64
	// 12 - I/Os currently in progress
	IOInProgress int64
	// 13 - time spent doing I/Os (ms)
	TimeDoingIO int64
	// 14 - weighted time spent doing I/Os (ms)
	WeightedTimeDoingIO int64
}

type Diskstats []DiskstatLine

func NewDiskstats() (Diskstats, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return Diskstats{}, err
	}

	return fs.NewDiskstats()
}

// NewDiskstats returns an information about current kernel/system statistics.
func (fs FS) NewDiskstats() (m Diskstats, err error) {
	f, err := os.Open(fs.Path("diskstats"))
	if err != nil {
		return Diskstats{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		var fieldsInt []int64

		line := s.Text()

		fields := strings.Fields(line)
		for i := range fields {
			var err error

			fieldsInt[i], err = strconv.ParseInt(fields[i], 10, 64)
			if err != nil {
				fieldsInt[i] = -1
			}
		}

		m = append(m, DiskstatLine{
			fieldsInt[0],
			fieldsInt[1],
			fieldsInt[2],
			fieldsInt[3],
			fieldsInt[4],
			fieldsInt[5],
			fieldsInt[6],
			fieldsInt[7],
			fieldsInt[8],
			fieldsInt[9],
			fieldsInt[10],
			fieldsInt[11],
			fieldsInt[12],
			fieldsInt[13],
		})
	}

	return m, nil
}
