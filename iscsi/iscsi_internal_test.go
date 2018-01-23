// Copyright 2018 The Prometheus Authors
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

package iscsi

import (
	"reflect"
	"testing"
    "path/filepath"
)

func TestGetStats(t *testing.T) {
	tests := []struct {
		invalid bool
		stat    *Stats
	}{
		{
			stat: &Stats{
				Name: "iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.8888bbbbddd0",
				Tpgt: []TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.8888bbbbddd0/tpgt_1",
						IsEnable: true,
						Luns: []LUN{
							{
								Name:       "lun_0",
								LunPath:    "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.8888bbbbddd0/tpgt_1/lun/lun_0",
								Backstore:  "rd_mcp",
								ObjectName: "ramdisk_lio_1G",
								TypeNumber: "0",
							},
						},
					},
				},
			},
		},
		{
			stat: &Stats{
				Name: "iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.abcd1abcd2ab",
				Tpgt: []TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.abcd1abcd2ab/tpgt_1",
						IsEnable: true,
						Luns: []LUN{
							{
								Name:       "lun_0",
								LunPath:    "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.abcd1abcd2ab/tpgt_1/lun/lun_0",
								Backstore:  "iblock",
								ObjectName: "block_lio_rbd1",
								TypeNumber: "0",
							},
						},
					},
				},
			},
		},
		{
			stat: &Stats{
				Name: "iqn.2016-11.org.linux-iscsi.igw.x86:dev.rbd0",
				Tpgt: []TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2016-11.org.linux-iscsi.igw.x86:dev.rbd0/tpgt_1",
						IsEnable: true,
						Luns: []LUN{
							{
								Name:       "lun_0",
								LunPath:    "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2016-11.org.linux-iscsi.igw.x86:dev.rbd0/tpgt_1/lun/lun_0",
								Backstore:  "fileio",
								ObjectName: "file_lio_1G",
								TypeNumber: "1",
							},
						},
					},
				},
			},
		},
		{
			stat: &Stats{
				Name: "iqn.2016-11.org.linux-iscsi.igw.x86:sn.ramdemo",
				Tpgt: []TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2016-11.org.linux-iscsi.igw.x86:sn.ramdemo/tpgt_1",
						IsEnable: true,
						Luns: []LUN{
							{
								Name:       "lun_0",
								LunPath:    "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2016-11.org.linux-iscsi.igw.x86:sn.ramdemo/tpgt_1/lun/lun_0",
								Backstore:  "rbd",
								ObjectName: "iscsi-images-demo",
								TypeNumber: "0",
							},
						},
					},
				},
			},
		},
	}

    SetPath("../sysfs/fixtures")
    matches, err := filepath.Glob(filepath.Join(sysPath, TARGETPATH, "/iqn*"))
	if err != nil {
		t.Errorf("unexpected test fixtures")
	}
    fstats := make([]*Stats, 0, len(matches))
    for _, iqnPath := range matches {
        name := filepath.Base(iqnPath)
        s, err := GetStats(iqnPath)
        if err != nil {
			t.Errorf("unexpected iSCSI stats:%s", iqnPath)
        }
        s.Name = name
        fstats = append(fstats, s)
    }

	for i, stat := range fstats {
		if want, have := tests[i].stat, stat; !reflect.DeepEqual(want, have) {
			t.Errorf("unexpected iSCSI stats:\nwant:\n%v\nhave:\n%v", want, have)
		}
	}
}
