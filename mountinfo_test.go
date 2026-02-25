// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package procfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMountInfo(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		mount   *MountInfo
		invalid bool
	}{
		{
			name:    "Regular sysfs mounted at /sys",
			s:       "16 21 0:16 / /sys rw,nosuid,nodev,noexec,relatime shared:7 - sysfs sysfs rw",
			invalid: false,
			mount: &MountInfo{
				MountID:        16,
				ParentID:       21,
				MajorMinorVer:  "0:16",
				Root:           "/",
				MountPoint:     "/sys",
				Options:        map[string]string{"rw": "", "nosuid": "", "nodev": "", "noexec": "", "relatime": ""},
				OptionalFields: map[string]string{"shared": "7"},
				FSType:         "sysfs",
				Source:         "sysfs",
				SuperOptions:   map[string]string{"rw": ""},
			},
		},
		{
			name:    "Not enough information",
			s:       "hello",
			invalid: true,
		},
		{
			name: "Tmpfs mounted at /run",
			s:    "225 20 0:39 / /run/user/112 rw,nosuid,nodev,relatime shared:177 - tmpfs tmpfs rw,size=405096k,mode=700,uid=112,gid=116",
			mount: &MountInfo{
				MountID:        225,
				ParentID:       20,
				MajorMinorVer:  "0:39",
				Root:           "/",
				MountPoint:     "/run/user/112",
				Options:        map[string]string{"rw": "", "nosuid": "", "nodev": "", "relatime": ""},
				OptionalFields: map[string]string{"shared": "177"},
				FSType:         "tmpfs",
				Source:         "tmpfs",
				SuperOptions:   map[string]string{"rw": "", "size": "405096k", "mode": "700", "uid": "112", "gid": "116"},
			},
			invalid: false,
		},
		{
			name: "Tmpfs mounted at /run, but no optional values",
			s:    "225 20 0:39 / /run/user/112 rw,nosuid,nodev,relatime  - tmpfs tmpfs rw,size=405096k,mode=700,uid=112,gid=116",
			mount: &MountInfo{
				MountID:        225,
				ParentID:       20,
				MajorMinorVer:  "0:39",
				Root:           "/",
				MountPoint:     "/run/user/112",
				Options:        map[string]string{"rw": "", "nosuid": "", "nodev": "", "relatime": ""},
				OptionalFields: nil,
				FSType:         "tmpfs",
				Source:         "tmpfs",
				SuperOptions:   map[string]string{"rw": "", "size": "405096k", "mode": "700", "uid": "112", "gid": "116"},
			},
			invalid: false,
		},
		{
			name: "Tmpfs mounted at /run, with multiple optional values",
			s:    "225 20 0:39 / /run/user/112 rw,nosuid,nodev,relatime shared:177 master:8 - tmpfs tmpfs rw,size=405096k,mode=700,uid=112,gid=116",
			mount: &MountInfo{
				MountID:        225,
				ParentID:       20,
				MajorMinorVer:  "0:39",
				Root:           "/",
				MountPoint:     "/run/user/112",
				Options:        map[string]string{"rw": "", "nosuid": "", "nodev": "", "relatime": ""},
				OptionalFields: map[string]string{"shared": "177", "master": "8"},
				FSType:         "tmpfs",
				Source:         "tmpfs",
				SuperOptions:   map[string]string{"rw": "", "size": "405096k", "mode": "700", "uid": "112", "gid": "116"},
			},
			invalid: false,
		},
		{
			name: "Tmpfs mounted at /run, with a mixture of valid and invalid optional values",
			s:    "225 20 0:39 / /run/user/112 rw,nosuid,nodev,relatime shared:177 master:8 foo:bar - tmpfs tmpfs rw,size=405096k,mode=700,uid=112,gid=116",
			mount: &MountInfo{
				MountID:        225,
				ParentID:       20,
				MajorMinorVer:  "0:39",
				Root:           "/",
				MountPoint:     "/run/user/112",
				Options:        map[string]string{"rw": "", "nosuid": "", "nodev": "", "relatime": ""},
				OptionalFields: map[string]string{"shared": "177", "master": "8"},
				FSType:         "tmpfs",
				Source:         "tmpfs",
				SuperOptions:   map[string]string{"rw": "", "size": "405096k", "mode": "700", "uid": "112", "gid": "116"},
			},
			invalid: false,
		},
		{
			name: "CIFS mounted at /with/a-hyphen",
			s:    "454 29 0:87 / /with/a-hyphen rw,relatime shared:255 - cifs //remote-storage/Path rw,vers=3.1.1,cache=strict,username=user,uid=1000,forceuid,gid=0,noforcegid,addr=127.0.0.1,file_mode=0755,dir_mode=0755,soft,nounix,serverino,mapposix,echo_interval=60,actimeo=1",
			mount: &MountInfo{
				MountID:        454,
				ParentID:       29,
				MajorMinorVer:  "0:87",
				Root:           "/",
				MountPoint:     "/with/a-hyphen",
				Options:        map[string]string{"rw": "", "relatime": ""},
				OptionalFields: map[string]string{"shared": "255"},
				FSType:         "cifs",
				Source:         "//remote-storage/Path",
				SuperOptions:   map[string]string{"rw": "", "vers": "3.1.1", "cache": "strict", "username": "user", "uid": "1000", "forceuid": "", "gid": "0", "noforcegid": "", "addr": "127.0.0.1", "file_mode": "0755", "dir_mode": "0755", "soft": "", "nounix": "", "serverino": "", "mapposix": "", "echo_interval": "60", "actimeo": "1"},
			},
			invalid: false,
		},
		{
			name: "Docker overlay with 10 fields (no optional fields)",
			s:    "137 45 253:2 /lib/docker/overlay2 /var/lib/docker/overlay2 rw,relatime - ext4 /dev/mapper/vg0-lv_var rw,data=ordered",
			mount: &MountInfo{
				MountID:        137,
				ParentID:       45,
				MajorMinorVer:  "253:2",
				Root:           "/lib/docker/overlay2",
				MountPoint:     "/var/lib/docker/overlay2",
				Options:        map[string]string{"rw": "", "relatime": ""},
				OptionalFields: map[string]string{},
				FSType:         "ext4",
				Source:         "/dev/mapper/vg0-lv_var",
				SuperOptions:   map[string]string{"rw": "", "data": "ordered"},
			},
		},
		{
			name: "bind chroot bind mount with 10 fields (no optional fields)",
			s:    "157 47 253:2 /etc/named /var/named/chroot/etc/named rw,relatime - ext4 /dev/mapper/vg0-lv_root rw,data=ordered",
			mount: &MountInfo{
				MountID:        157,
				ParentID:       47,
				MajorMinorVer:  "253:2",
				Root:           "/etc/named",
				MountPoint:     "/var/named/chroot/etc/named",
				Options:        map[string]string{"rw": "", "relatime": ""},
				OptionalFields: map[string]string{},
				FSType:         "ext4",
				Source:         "/dev/mapper/vg0-lv_root",
				SuperOptions:   map[string]string{"rw": "", "data": "ordered"},
			},
		},
	}

	for i, test := range tests {
		t.Logf("[%02d] test %q", i, test.name)

		mount, err := parseMountInfoString(test.s)

		if test.invalid && err == nil {
			t.Error("expected an error, but none occurred")
		}
		if !test.invalid && err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if diff := cmp.Diff(test.mount, mount); diff != "" {
			t.Fatalf("unexpected diff (-want +got):\n%s", diff)
		}
	}
}

func TestFSMountInfo(t *testing.T) {
	fs, err := NewFS(procTestFixtures)
	if err != nil {
		t.Fatalf("failed to open procfs: %v", err)
	}

	want := []*MountInfo{
		{
			MountID:        1,
			ParentID:       1,
			MajorMinorVer:  "0:5",
			Root:           "/",
			Options:        map[string]string{"/root": ""},
			OptionalFields: map[string]string{"shared": "8"},
			FSType:         "rootfs",
			Source:         "rootfs",
			SuperOptions:   map[string]string{"rw": ""},
		},
		{
			MountID:        16,
			ParentID:       21,
			MajorMinorVer:  "0:16",
			Root:           "/",
			MountPoint:     "/sys",
			Options:        map[string]string{"nodev": "", "noexec": "", "nosuid": "", "relatime": "", "rw": ""},
			OptionalFields: map[string]string{"shared": "7"},
			FSType:         "sysfs",
			Source:         "sysfs",
			SuperOptions:   map[string]string{"rw": ""},
		},
		{
			MountID:        17,
			ParentID:       21,
			MajorMinorVer:  "0:4",
			Root:           "/",
			MountPoint:     "/proc",
			Options:        map[string]string{"nodev": "", "noexec": "", "nosuid": "", "relatime": "", "rw": ""},
			OptionalFields: map[string]string{"shared": "12"},
			FSType:         "proc",
			Source:         "proc",
			SuperOptions:   map[string]string{"rw": ""},
		},
		{
			MountID:        21,
			MajorMinorVer:  "8:1",
			Root:           "/",
			MountPoint:     "/",
			Options:        map[string]string{"relatime": "", "rw": ""},
			OptionalFields: map[string]string{"shared": "1"},
			FSType:         "ext4",
			Source:         "/dev/sda1",
			SuperOptions:   map[string]string{"data": "ordered", "errors": "remount-ro", "rw": ""},
		},
		{
			MountID:        194,
			ParentID:       21,
			MajorMinorVer:  "0:42",
			Root:           "/",
			MountPoint:     "/mnt/nfs/test",
			Options:        map[string]string{"rw": ""},
			OptionalFields: map[string]string{"shared": "144"},
			FSType:         "nfs4",
			Source:         "192.168.1.1:/srv/test",
			SuperOptions: map[string]string{
				"acdirmax":   "60",
				"acdirmin":   "30",
				"acregmax":   "60",
				"acregmin":   "3",
				"addr":       "192.168.1.1",
				"clientaddr": "192.168.1.5",
				"hard":       "",
				"local_lock": "none",
				"namlen":     "255",
				"port":       "0",
				"proto":      "tcp",
				"retrans":    "2",
				"rsize":      "1048576",
				"rw":         "",
				"sec":        "sys",
				"timeo":      "600",
				"vers":       "4.0",
				"wsize":      "1048576",
			},
		},
		{
			MountID:        177,
			ParentID:       21,
			MajorMinorVer:  "0:42",
			Root:           "/",
			MountPoint:     "/mnt/nfs/test",
			Options:        map[string]string{"rw": ""},
			OptionalFields: map[string]string{"shared": "130"},
			FSType:         "nfs4",
			Source:         "192.168.1.1:/srv/test",
			SuperOptions: map[string]string{
				"acdirmax":   "60",
				"acdirmin":   "30",
				"acregmax":   "60",
				"acregmin":   "3",
				"addr":       "192.168.1.1",
				"clientaddr": "192.168.1.5",
				"hard":       "",
				"local_lock": "none",
				"namlen":     "255",
				"port":       "0",
				"proto":      "tcp",
				"retrans":    "2",
				"rsize":      "1048576",
				"rw":         "",
				"sec":        "sys",
				"timeo":      "600",
				"vers":       "4.0",
				"wsize":      "1048576",
			},
		},
		{
			MountID:        1398,
			ParentID:       798,
			MajorMinorVer:  "0:44",
			Root:           "/",
			MountPoint:     "/mnt/nfs/test",
			Options:        map[string]string{"relatime": "", "rw": ""},
			OptionalFields: map[string]string{"shared": "1154"},
			FSType:         "nfs",
			Source:         "192.168.1.1:/srv/test",
			SuperOptions: map[string]string{
				"addr":       "192.168.1.1",
				"hard":       "",
				"local_lock": "none",
				"mountaddr":  "192.168.1.1",
				"mountport":  "49602",
				"mountproto": "udp",
				"mountvers":  "3",
				"namlen":     "255",
				"proto":      "udp",
				"retrans":    "3",
				"rsize":      "32768",
				"rw":         "",
				"sec":        "sys",
				"timeo":      "11",
				"vers":       "3",
				"wsize":      "32768",
			},
		},
		{
			MountID:        1128,
			ParentID:       67,
			MajorMinorVer:  "253:0",
			Root:           "/var/lib/containers/storage/overlay",
			MountPoint:     "/var/lib/containers/storage/overlay",
			Options:        map[string]string{"relatime": "", "rw": ""},
			OptionalFields: map[string]string{},
			FSType:         "xfs",
			Source:         "/dev/mapper/rhel-root",
			SuperOptions: map[string]string{
				"attr2":    "",
				"inode64":  "",
				"logbsize": "32k",
				"logbufs":  "8",
				"noquota":  "",
				"rw":       "",
				"seclabel": "",
			},
		},
	}

	got, err := fs.GetMounts()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected mountpoints (-want +got):\n%s", diff)
	}

	got, err = fs.GetProcMounts(26231)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected mountpoints (-want +got):\n%s", diff)
	}
}
