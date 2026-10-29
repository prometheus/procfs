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

//go:build !windows

package procfs

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestZoneinfo(t *testing.T) {
	fs := getProcFixtures(t)

	protectionDMA := []*int64{new(int64(0)), new(int64(2877)), new(int64(7826)), new(int64(7826)), new(int64(7826))}
	protectionDMA32 := []*int64{new(int64(0)), new(int64(0)), new(int64(4949)), new(int64(4949)), new(int64(4949))}
	protectionNormal := []*int64{new(int64(0)), new(int64(0)), new(int64(0)), new(int64(0)), new(int64(0))}
	protectionMovable := []*int64{new(int64(0)), new(int64(0)), new(int64(0)), new(int64(0)), new(int64(0))}
	protectionDevice := []*int64{new(int64(0)), new(int64(0)), new(int64(0)), new(int64(0)), new(int64(0))}

	refs := []Zoneinfo{
		{
			Node:                       "0",
			Zone:                       "DMA",
			NrFreePages:                new(int64(3952)),
			Min:                        new(int64(33)),
			Low:                        new(int64(41)),
			High:                       new(int64(49)),
			Spanned:                    new(int64(4095)),
			Present:                    new(int64(3975)),
			Managed:                    new(int64(3956)),
			NrActiveAnon:               new(int64(547580)),
			NrInactiveAnon:             new(int64(230981)),
			NrIsolatedAnon:             new(int64(0)),
			NrAnonPages:                new(int64(795576)),
			NrAnonTransparentHugepages: new(int64(0)),
			NrActiveFile:               new(int64(346282)),
			NrInactiveFile:             new(int64(316904)),
			NrIsolatedFile:             new(int64(0)),
			NrFilePages:                new(int64(761874)),
			NrSlabReclaimable:          new(int64(131220)),
			NrSlabUnreclaimable:        new(int64(47320)),
			NrKernelStack:              new(int64(0)),
			NrMapped:                   new(int64(215483)),
			NrDirty:                    new(int64(908)),
			NrWriteback:                new(int64(0)),
			NrUnevictable:              new(int64(115467)),
			NrShmem:                    new(int64(224925)),
			NrDirtied:                  new(int64(8007423)),
			NrWritten:                  new(int64(7752121)),
			NumaHit:                    new(int64(1)),
			NumaMiss:                   new(int64(0)),
			NumaForeign:                new(int64(0)),
			NumaInterleave:             new(int64(0)),
			NumaLocal:                  new(int64(1)),
			NumaOther:                  new(int64(0)),
			Protection:                 protectionDMA,
		},
		{
			Node:           "0",
			Zone:           "DMA32",
			NrFreePages:    new(int64(204252)),
			Min:            new(int64(19510)),
			Low:            new(int64(21059)),
			High:           new(int64(22608)),
			Spanned:        new(int64(1044480)),
			Present:        new(int64(759231)),
			Managed:        new(int64(742806)),
			NrKernelStack:  new(int64(2208)),
			NumaHit:        new(int64(113952967)),
			NumaMiss:       new(int64(0)),
			NumaForeign:    new(int64(0)),
			NumaInterleave: new(int64(0)),
			NumaLocal:      new(int64(113952967)),
			NumaOther:      new(int64(0)),
			Protection:     protectionDMA32,
		},
		{
			Node:           "0",
			Zone:           "Normal",
			NrFreePages:    new(int64(18553)),
			Min:            new(int64(11176)),
			Low:            new(int64(13842)),
			High:           new(int64(16508)),
			Spanned:        new(int64(1308160)),
			Present:        new(int64(1308160)),
			Managed:        new(int64(1268711)),
			NrKernelStack:  new(int64(15136)),
			NumaHit:        new(int64(162718019)),
			NumaMiss:       new(int64(0)),
			NumaForeign:    new(int64(0)),
			NumaInterleave: new(int64(26812)),
			NumaLocal:      new(int64(162718019)),
			NumaOther:      new(int64(0)),
			Protection:     protectionNormal,
		},
		{
			Node:       "0",
			Zone:       "Movable",
			Min:        new(int64(0)),
			Low:        new(int64(0)),
			High:       new(int64(0)),
			Spanned:    new(int64(0)),
			Present:    new(int64(0)),
			Managed:    new(int64(0)),
			Protection: protectionMovable,
		},

		{
			Node:       "0",
			Zone:       "Device",
			Min:        new(int64(0)),
			Low:        new(int64(0)),
			High:       new(int64(0)),
			Spanned:    new(int64(0)),
			Present:    new(int64(0)),
			Managed:    new(int64(0)),
			Protection: protectionDevice,
		},
	}

	data, err := fs.Zoneinfo()
	if err != nil {
		t.Fatalf("failed to parse zoneinfo: %v", err)
	}
	if len(data) != 5 {
		t.Fatal("failed to get all all node info")
	}

	for index, ref := range refs {
		want, got := ref, data[index]
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("unexpected crypto entry (-want +got):\n%s", diff)
		}

	}
}
