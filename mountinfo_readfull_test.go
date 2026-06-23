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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/procfs/internal/util"
)

// TestReadMountInfoExceeds1MiB verifies that a mountinfo file larger than 1 MiB
// is read and parsed in full. The util.ReadFileNoStat helper caps reads at
// 1 MiB, which truncates and corrupts mountinfo on hosts with a very large
// number of mounts (e.g. busy container hosts), so the mountinfo helpers must
// not use it.
func TestReadMountInfoExceeds1MiB(t *testing.T) {
	const n = 20000

	var b strings.Builder
	for i := range n {
		fmt.Fprintf(&b, "%d 35 98:0 /mnt%d /mnt%d rw,noatime shared:1 - ext3 /dev/root rw,errors=continue\n", i, i, i)
	}
	content := b.String()
	if len(content) <= 1024*1024 {
		t.Fatalf("test fixture must exceed 1 MiB, got %d bytes", len(content))
	}

	path := filepath.Join(t.TempDir(), "mountinfo")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := readMountInfo(path)
	if err != nil {
		t.Fatalf("readMountInfo: %v", err)
	}
	if len(data) != len(content) {
		t.Fatalf("readMountInfo truncated the file: got %d bytes, want %d", len(data), len(content))
	}

	mounts, err := parseMountInfo(data)
	if err != nil {
		t.Fatalf("parseMountInfo returned error: %v", err)
	}
	if len(mounts) != n {
		t.Fatalf("parsed %d mounts, want %d", len(mounts), n)
	}

	// Guard the regression: util.ReadFileNoStat caps at 1 MiB and would truncate.
	capped, err := util.ReadFileNoStat(path)
	if err != nil {
		t.Fatalf("ReadFileNoStat: %v", err)
	}
	if len(capped) >= len(content) {
		t.Fatalf("expected ReadFileNoStat to cap the read below %d bytes, got %d", len(content), len(capped))
	}
}
