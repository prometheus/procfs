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

func TestParseNFNetLinkQueueLine(t *testing.T) {
	tests := []struct {
		name           string
		s              string
		shouldErr      bool
		nfNetLinkQueue *NFNetLinkQueue
	}{
		{
			name:      "nf_net_link_queue simple line",
			s:         "  230  44306     1 2 65531     3     4        5  6",
			shouldErr: false,
			nfNetLinkQueue: &NFNetLinkQueue{
				QueueID:          230,
				PeerPID:          44306,
				QueueTotal:       1,
				CopyMode:         2,
				CopyRange:        65531,
				QueueDropped:     3,
				QueueUserDropped: 4,
				SequenceID:       5,
				Use:              6,
			},
		},
		{
			name:           "empty line",
			s:              "",
			shouldErr:      true,
			nfNetLinkQueue: nil,
		},
		{
			name:           "incorrect parameters count in line",
			s:              " 1 2 3 4 55555 ",
			shouldErr:      true,
			nfNetLinkQueue: nil,
		},
	}

	for i, test := range tests {
		t.Logf("[%02d] test %q", i, test.name)

		nfNetLinkQueue, err := parseNFNetLinkQueueLine(test.s)

		if test.shouldErr && err == nil {
			t.Errorf("%s: expected an error, but none occurred", test.name)
		}
		if !test.shouldErr && err != nil {
			t.Errorf("%s: unexpected error: %v", test.name, err)
		}

		if diff := cmp.Diff(test.nfNetLinkQueue, nfNetLinkQueue); diff != "" {
			t.Fatalf("unexpected nfNetLinkQueue (-want +got):\n%s", diff)
		}
	}
}
