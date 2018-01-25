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
//	"reflect"
	"testing"

    "github.com/prometheus/procfs/sysfs"
    "github.com/prometheus/procfs/iscsi"
    "github.com/kr/pretty"
)

func TestGetStats(t *testing.T) {
	tests := []struct {
		invalid bool
		stat    *iscsi.Stats
	}{
		{
			stat: &iscsi.Stats{
				Name: "iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.8888bbbbddd0",
				Tpgt: []iscsi.TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.8888bbbbddd0/tpgt_1",
						IsEnable: true,
						Luns: []iscsi.LUN{
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
			stat: &iscsi.Stats{
				Name: "iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.abcd1abcd2ab",
				Tpgt: []iscsi.TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2003-01.org.linux-iscsi.osd1.x8664:sn.abcd1abcd2ab/tpgt_1",
						IsEnable: true,
						Luns: []iscsi.LUN{
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
			stat: &iscsi.Stats{
				Name: "iqn.2016-11.org.linux-iscsi.igw.x86:dev.rbd0",
				Tpgt: []iscsi.TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2016-11.org.linux-iscsi.igw.x86:dev.rbd0/tpgt_1",
						IsEnable: true,
						Luns: []iscsi.LUN{
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
			stat: &iscsi.Stats{
				Name: "iqn.2016-11.org.linux-iscsi.igw.x86:sn.ramdemo",
				Tpgt: []iscsi.TPGT{
					{
						Name:     "tpgt_1",
						TpgtPath: "../sysfs/fixtures/kernel/config/target/iscsi/iqn.2016-11.org.linux-iscsi.igw.x86:sn.ramdemo/tpgt_1",
						IsEnable: true,
						Luns: []iscsi.LUN{
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

    readTests := []struct {
        read uint64
        write uint64
        iops uint64
    }{
        { 10325, 40325, 204950 },
        { 20095, 71235, 104950 },
        { 10195, 30195, 301950 },
        { 1504 , 4733 , 1234   }, 
    }

    sysfsStat, err := sysfs.FS("../sysfs/fixtures").ISCSIStats()
	if err != nil {
		t.Errorf("unexpected test fixtures")
	}

	for i, stat := range sysfsStat {
		want, have := tests[i].stat, stat;
        diff := pretty.Diff(want, have)
        if diff != nil {
			t.Errorf("unexpected iSCSI stats:\ndiff:\n%v", diff)
			t.Errorf("\nwant:\n%v\nhave:\n%v", want, have)
        } else {
            readMB, writeMB, iops, err := iscsi.ReadWriteOPS(stat.Name,
            stat.Tpgt[0].Name, stat.Tpgt[0].Luns[0].Name )
            if err != nil { 
			    t.Errorf("unexpected iSCSI ReadWriteOPS path %s %s %s",
                stat.Name, stat.Tpgt[0].Name, stat.Tpgt[0].Luns[0].Name)
			    t.Errorf("%v", err);
            }
            // datawant, datahave := readTests[i], { readMB, writeMB, iops };

            diff = pretty.Diff(readTests[i].read, readMB)
            if diff != nil {
                t.Errorf("unexpected iSCSI read data :\ndiff:\n%v", diff)
                t.Errorf("\nwant:\n%v\nhave:\n%v", readTests[i].read, readMB)
            }
            diff = pretty.Diff(readTests[i].write, writeMB)
            if diff != nil {
                t.Errorf("unexpected iSCSI write data :\ndiff:\n%v", diff)
                t.Errorf("\nwant:\n%v\nhave:\n%v", readTests[i].write, writeMB)
            }
            diff = pretty.Diff(readTests[i].iops, iops)
            if diff != nil {
                t.Errorf("unexpected iSCSI data :\ndiff:\n%v", diff)
                t.Errorf("\nwant:\n%v\nhave:\n%v", readTests[i].iops, iops)
            }
        }
	}
}
