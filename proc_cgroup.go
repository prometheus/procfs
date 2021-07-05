package procfs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"log"
	"strconv"
	"strings"

	"github.com/prometheus/procfs/internal/util"
)

// github.com/prometheus/procfs@v0.2.0/proc_cgroup.go
type Cgroup struct {
	HierarchyID int
	Controllers []string
	Path string
	CgroupMemMax int64
}


// parseCgroupString parses each line of the /proc/[pid]/cgroup file
// Line format is hierarchyID:[controller1,controller2]:path
func parseCgroupString(cgroupStr string) (*Cgroup, error) {
	var err error

	fields := strings.Split(cgroupStr, ":")
	if len(fields) < 3 {
		return nil, fmt.Errorf("at least 3 fields required, found %d fields in cgroup string: %s", len(fields), cgroupStr)
	}
	cgroup := &Cgroup{
		Path:        fields[2],
		Controllers: nil,
	}
	if fields[1] == "memory" {
		cgroupfile := "/sys/fs/cgroup/memory" + fields[2]
		myfile := cgroupfile + "/memory.limit_in_bytes"
		_, err := os.Stat(myfile)
		log.Printf("file file file file %v", myfile)
		log.Printf("file err file err file err %v", err)
		if err == nil {
			//data, _ := ioutil.ReadFile(myfile)
			data, _ := util.ReadFileNoStat(fmt.Sprintf("%v", myfile))
			//var CgroupMemMax int64 = 100
			log.Printf("str str str str is %v", string(data))
			log.Printf("str str str str is %T", string(data))
			trimdata := strings.TrimSpace(string(data))
			CgroupMemMax, err := strconv.ParseInt(trimdata, 10, 64)
			log.Printf("data data data is%v", CgroupMemMax)
			log.Printf("err err err is%v", err)
			cgroup.CgroupMemMax = CgroupMemMax
			log.Printf("%v", cgroup)
		}
	}
	cgroup.HierarchyID, err = strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse hierarchy ID")
	}
	if fields[1] != "" {
		ssNames := strings.Split(fields[1], ",")
		cgroup.Controllers = append(cgroup.Controllers, ssNames...)
	}
	return cgroup, nil
}

// parseCgroups reads each line of the /proc/[pid]/cgroup file
func parseCgroups(data []byte) ([]Cgroup, error) {
	var cgroups []Cgroup
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		mountString := scanner.Text()
		parsedMounts, err := parseCgroupString(mountString)
		if err != nil {
			return nil, err
		}
		if parsedMounts.Controllers[0] != "memory" {
			continue
		}
		cgroups = append(cgroups, *parsedMounts)
	}

	err := scanner.Err()
	return cgroups, err
}

// Cgroups reads from /proc/<pid>/cgroups and returns a []*Cgroup struct locating this PID in each process
// control hierarchy running on this system. On every system (v1 and v2), all hierarchies contain all processes,
// so the len of the returned struct is equal to the number of active hierarchies on this system
func (p Proc) Cgroups() ([]Cgroup, error) {
	data, err := util.ReadFileNoStat(fmt.Sprintf("/proc/%d/cgroup", p.PID))
	if err != nil {
		return nil, err
	}
	return parseCgroups(data)
}


//func (p Proc) MyCgroups() (Cgroup, error) {
//	var clist []Cgroup
//	data, err := util.ReadFileNoStat(fmt.Sprintf("/proc/%d/cgroup", p.PID))
//	if err != nil {
//		return Cgroup{}, err
//	}
//	scanner := bufio.NewScanner(bytes.NewReader(data))
//	for scanner.Scan() {
//		mountString := scanner.Text()
//		parsedMounts, err := parseCgroupString(mountString)
//		if err != nil {
//			return Cgroup{}, err
//		}
//		clist = append(clist, *parsedMounts)
//		cgroups := *parsedMounts
//	}
//	return Cgroup{}, err
//}

func (p Proc) NewCgroup() ([]Cgroup, error) {
   aa, _ := p.Cgroups()
   log.Printf("cgroup cgroup cgroup cgroup is %+v", aa)
   return p.Cgroups()
}


//func (c Cgroup) CgroupMemMax() int64 {
//	return s.VSize
//}