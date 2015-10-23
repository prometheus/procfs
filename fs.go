package procfs

import (
	"fmt"
	"os"
	"path"
	"strings"
	"strconv"
)

// FS represents the pseudo-filesystem proc, which provides an interface to
// kernel data structures.
type FS string

// DefaultMountPoint is the common mount point of the proc filesystem.
const DefaultMountPoint = "/proc"

// NewFS returns a new FS mounted under the given mountPoint. It will error
// if the mount point can't be read.
func NewFS(mountPoint string) (FS, error) {
	info, err := os.Stat(mountPoint)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %s", mountPoint, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("mount point %s is not a directory", mountPoint)
	}

	return FS(mountPoint), nil
}

func (fs FS) stat(p string) (os.FileInfo, error) {
	return os.Stat(path.Join(string(fs), p))
}

func (fs FS) open(p string) (*os.File, error) {
	return os.Open(path.Join(string(fs), p))
}

func (fs FS) readlink(p string) (string, error) {
	return os.Readlink(path.Join(string(fs), p))
}

// Self returns a process for the current process.
func (fs FS) Self() (Proc, error) {
	p, err := fs.readlink("self")
	if err != nil {
		return Proc{}, err
	}
	pid, err := strconv.Atoi(strings.Replace(p, string(fs), "", -1))
	if err != nil {
		return Proc{}, err
	}
	return fs.NewProc(pid)
}

// NewProc returns a process for the given pid.
func (fs FS) NewProc(pid int) (Proc, error) {
	if _, err := fs.stat(strconv.Itoa(pid)); err != nil {
		return Proc{}, err
	}
	return Proc{PID: pid, fs: fs}, nil
}

// AllProcs returns a list of all currently avaible processes.
func (fs FS) AllProcs() (Procs, error) {
	d, err := fs.open("")
	if err != nil {
		return Procs{}, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return Procs{}, fmt.Errorf("could not read %s: %s", d.Name(), err)
	}

	p := Procs{}
	for _, n := range names {
		pid, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			continue
		}
		p = append(p, Proc{PID: int(pid), fs: fs})
	}
	return p, nil
}
