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

package cifs_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/prometheus/procfs/cifs"
)

func TestNewCifsRPCStats(t *testing.T) {
	tests := []struct {
		name    string
		content string
		stats   *cifs.ClientStats
		invalid bool
	}{
		{
			name:    "invalid file",
			content: "invalid",
			invalid: true,
		}, {
			name: "SMB1 statistics",
			content: `Resources in use
CIFS Session: 1
Share (unique mount targets): 2
SMB Request/Response Buffer: 1 Pool size: 5
SMB Small Req/Resp Buffer: 1 Pool size: 30
Operations (MIDs): 0

0 session 0 share reconnects
Total vfs operations: 16 maximum at one time: 2

1) \\server\share
SMBs: 9 Oplocks breaks: 0
Reads:  0 Bytes: 0
Writes: 0 Bytes: 0
Flushes: 0
Locks: 0 HardLinks: 0 Symlinks: 0
Opens: 0 Closes: 0 Deletes: 0
Posix Opens: 0 Posix Mkdirs: 0
Mkdirs: 0 Rmdirs: 0
Renames: 0 T2 Renames 0
FindFirst: 1 FNext 0 FClose 0`,
			stats: &cifs.ClientStats{
				Header: map[string]uint64{
					"operations":         0,
					"sessionCount":       0,
					"sessions":           1,
					"shareReconnects":    0,
					"shares":             2,
					"smbBuffer":          1,
					"smbPoolSize":        5,
					"smbSmallBuffer":     1,
					"smbSmallPoolSize":   30,
					"totalMaxOperations": 2,
					"totalOperations":    16,
				},
				SMBStatsList: []*cifs.SMBStats{
					&cifs.SMBStats{
						SessionIDs: cifs.SessionIDs{
							SessionID: 1,
							Server:    "server",
							Share:     "\\share",
						},
						Stats: map[string]uint64{
							"breaks":      0,
							"closes":      0,
							"deletes":     0,
							"fClose":      0,
							"fNext":       0,
							"findFirst":   1,
							"flushes":     0,
							"hardlinks":   0,
							"locks":       0,
							"mkdirs":      0,
							"opens":       0,
							"posixMkdirs": 0,
							"posixOpens":  0,
							"reads":       0,
							"readsBytes":  0,
							"renames":     0,
							"rmdirs":      0,
							"smbs":        9,
							"symlinks":    0,
							"t2Renames":   0,
							"writes":      0,
							"writesBytes": 0,
						},
					},
				},
			},
		}, {
			name: "SMB2 statistics",
			content: `Resources in use
CIFS Session: 2
Share (unique mount targets): 4
SMB Request/Response Buffer: 2 Pool size: 6
SMB Small Req/Resp Buffer: 2 Pool size: 30
Operations (MIDs): 0

0 session 0 share reconnects
Total vfs operations: 90 maximum at one time: 2

1) \\server\share1
SMBs: 20
Negotiates: 0 sent 0 failed
SessionSetups: 0 sent 0 failed
Logoffs: 0 sent 0 failed
TreeConnects: 0 sent 0 failed
TreeDisconnects: 0 sent 0 failed
Creates: 0 sent 2 failed
Closes: 0 sent 0 failed
Flushes: 0 sent 0 failed
Reads: 0 sent 0 failed
Writes: 0 sent 0 failed
Locks: 0 sent 0 failed
IOCTLs: 0 sent 0 failed
Cancels: 0 sent 0 failed
Echos: 0 sent 0 failed
QueryDirectories: 0 sent 0 failed
ChangeNotifies: 0 sent 0 failed
QueryInfos: 0 sent 0 failed
SetInfos: 0 sent 0 failed
OplockBreaks: 0 sent 0 failed`,
			stats: &cifs.ClientStats{
				Header: map[string]uint64{
					"operations":         0,
					"sessionCount":       0,
					"sessions":           2,
					"shareReconnects":    0,
					"shares":             4,
					"smbBuffer":          2,
					"smbPoolSize":        6,
					"smbSmallBuffer":     2,
					"smbSmallPoolSize":   30,
					"totalMaxOperations": 2,
					"totalOperations":    90,
				},
				SMBStatsList: []*cifs.SMBStats{
					&cifs.SMBStats{
						SessionIDs: cifs.SessionIDs{
							SessionID: 1,
							Server:    "server",
							Share:     "\\share1",
						},
						Stats: map[string]uint64{
							"cancelsSent":            0,
							"cancelsFailed":          0,
							"changeNotifiesSent":     0,
							"changeNotifiesFailed":   0,
							"closesSent":             0,
							"closesFailed":           0,
							"createsSent":            0,
							"createsFailed":          2,
							"echosSent":              0,
							"echosFailed":            0,
							"flushesSent":            0,
							"flushesFailed":          0,
							"ioCTLsSent":             0,
							"ioCTLsFailed":           0,
							"locksSent":              0,
							"locksFailed":            0,
							"logoffsSent":            0,
							"logoffsFailed":          0,
							"negotiatesSent":         0,
							"negotiatesFailed":       0,
							"oplockBreaksSent":       0,
							"oplockBreaksFailed":     0,
							"queryDirectoriesSent":   0,
							"queryDirectoriesFailed": 0,
							"queryInfosSent":         0,
							"queryInfosFailed":       0,
							"readsSent":              0,
							"readsFailed":            0,
							"sessionSetupsSent":      0,
							"sessionSetupsFailed":    0,
							"setInfosSent":           0,
							"setInfosFailed":         0,
							"treeConnectsSent":       0,
							"treeConnectsFailed":     0,
							"treeDisconnectsSent":    0,
							"treeDisconnectsFailed":  0,
							"writesSent":             0,
							"writesFailed":           0,
							"smbs":                   20,
						},
					},
				},
			},
		}, {
			name: "Mixed statistics (SMB1 then SMB2)",
			content: `Resources in use
CIFS Session: 1
Share (unique mount targets): 2
SMB Request/Response Buffer: 1 Pool size: 5
SMB Small Req/Resp Buffer: 1 Pool size: 30
Operations (MIDs): 0

0 session 0 share reconnects
Total vfs operations: 16 maximum at one time: 2

1) \\server1\share1
SMBs: 9 Oplocks breaks: 0
Reads:  0 Bytes: 0
Writes: 0 Bytes: 0
Flushes: 0
Locks: 0 HardLinks: 0 Symlinks: 0
Opens: 0 Closes: 0 Deletes: 0
Posix Opens: 0 Posix Mkdirs: 0
Mkdirs: 0 Rmdirs: 0
Renames: 0 T2 Renames 0
FindFirst: 1 FNext 0 FClose 0

2) \\server2\share2
SMBs: 20
Negotiates: 1 sent 2 failed
SessionSetups: 3 sent 4 failed
Logoffs: 5 sent 6 failed
TreeConnects: 7 sent 8 failed
TreeDisconnects: 9 sent 10 failed
Creates: 11 sent 12 failed
Closes: 13 sent 14 failed
Flushes: 15 sent 16 failed
Reads: 17 sent 18 failed
Writes: 99 sent 98 failed
Locks: 97 sent 96 failed
IOCTLs: 95 sent 94 failed
Cancels: 10 sent 110 failed
Echos: 20 sent 220 failed
QueryDirectories: 30 sent 330 failed
ChangeNotifies: 40 sent 440 failed
QueryInfos: 50 sent 550 failed
SetInfos: 60 sent 660 failed
OplockBreaks: 770 sent 70 failed`,
			stats: &cifs.ClientStats{
				Header: map[string]uint64{
					"operations":         0,
					"sessionCount":       0,
					"sessions":           1,
					"shareReconnects":    0,
					"shares":             2,
					"smbBuffer":          1,
					"smbPoolSize":        5,
					"smbSmallBuffer":     1,
					"smbSmallPoolSize":   30,
					"totalMaxOperations": 2,
					"totalOperations":    16,
				},
				SMBStatsList: []*cifs.SMBStats{
					&cifs.SMBStats{
						SessionIDs: cifs.SessionIDs{
							SessionID: 1,
							Server:    "server1",
							Share:     "\\share1",
						},
						Stats: map[string]uint64{
							"breaks":      0,
							"closes":      0,
							"deletes":     0,
							"fClose":      0,
							"fNext":       0,
							"findFirst":   1,
							"flushes":     0,
							"hardlinks":   0,
							"locks":       0,
							"mkdirs":      0,
							"opens":       0,
							"posixMkdirs": 0,
							"posixOpens":  0,
							"reads":       0,
							"readsBytes":  0,
							"renames":     0,
							"rmdirs":      0,
							"smbs":        9,
							"symlinks":    0,
							"t2Renames":   0,
							"writes":      0,
							"writesBytes": 0,
						},
					},
					&cifs.SMBStats{
						SessionIDs: cifs.SessionIDs{
							SessionID: 2,
							Server:    "server2",
							Share:     "\\share2",
						},
						Stats: map[string]uint64{
							"cancelsSent":            10,
							"cancelsFailed":          110,
							"changeNotifiesSent":     40,
							"changeNotifiesFailed":   440,
							"closesSent":             13,
							"closesFailed":           14,
							"createsSent":            11,
							"createsFailed":          12,
							"echosSent":              20,
							"echosFailed":            220,
							"flushesSent":            15,
							"flushesFailed":          16,
							"ioCTLsSent":             95,
							"ioCTLsFailed":           94,
							"locksSent":              97,
							"locksFailed":            96,
							"logoffsSent":            5,
							"logoffsFailed":          6,
							"negotiatesSent":         1,
							"negotiatesFailed":       2,
							"oplockBreaksSent":       70,
							"oplockBreaksFailed":     770,
							"queryDirectoriesSent":   30,
							"queryDirectoriesFailed": 330,
							"queryInfosSent":         50,
							"queryInfosFailed":       550,
							"readsSent":              17,
							"readsFailed":            18,
							"sessionSetupsSent":      3,
							"sessionSetupsFailed":    4,
							"setInfosSent":           60,
							"setInfosFailed":         660,
							"treeConnectsSent":       7,
							"treeConnectsFailed":     8,
							"treeDisconnectsSent":    9,
							"treeDisconnectsFailed":  10,
							"writesSent":             99,
							"writesFailed":           98,
							"smbs":                   20,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := cifs.ParseClientStats(strings.NewReader(tt.content))

			if tt.invalid && nil == err {
				t.Fatal("expected an error, but none occured")
			}
			if !tt.invalid && nil != err {
				t.Fatalf("unexpected error: %v", err)
			}
			if want, have := tt.stats, stats; !reflect.DeepEqual(want, have) {
				t.Fatalf("unexpected CIFS Stats:\nwant:\n%v\nhave:\n%v", want, have)
			}
		})
	}
}
