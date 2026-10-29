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
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFS_Crypto(t *testing.T) {
	fs := getProcFixtures(t)
	crypto, err := fs.Crypto()
	if err != nil {
		t.Fatalf("parsing of reference-file failed entirely: %s", err)
	}

	refs := []Crypto{
		{
			Name:        "ccm(aes)",
			Driver:      "ccm_base(ctr(aes-aesni),cbcmac(aes-aesni))",
			Module:      "ccm",
			Priority:    new(int64(300)),
			Refcnt:      new(int64(4)),
			Selftest:    "passed",
			Internal:    "no",
			Type:        "aead",
			Async:       false,
			Blocksize:   new(uint64(1)),
			Ivsize:      new(uint64(16)),
			Maxauthsize: new(uint64(16)),
			Geniv:       "<none>",
		},
		{
			Name:       "cbcmac(aes)",
			Driver:     "cbcmac(aes-aesni)",
			Module:     "ccm",
			Priority:   new(int64(300)),
			Refcnt:     new(int64(7)),
			Selftest:   "passed",
			Internal:   "no",
			Type:       "shash",
			Blocksize:  new(uint64(1)),
			Digestsize: new(uint64(16)),
		},
		{
			Name:     "ecdh",
			Driver:   "ecdh-generic",
			Module:   "ecdh_generic",
			Priority: new(int64(100)),
			Refcnt:   new(int64(1)),
			Selftest: "passed",
			Internal: "no",
			Type:     "kpp",
			Async:    true,
		},
		{
			Name:       "ecb(arc4)",
			Driver:     "ecb(arc4)-generic",
			Module:     "arc4",
			Priority:   new(int64(100)),
			Refcnt:     new(int64(1)),
			Selftest:   "passed",
			Internal:   "no",
			Type:       "skcipher",
			Async:      false,
			Blocksize:  new(uint64(1)),
			MinKeysize: new(uint64(1)),
			MaxKeysize: new(uint64(256)),
			Ivsize:     new(uint64(0)),
			Chunksize:  new(uint64(1)),
			Walksize:   new(uint64(1)),
		},
		{
			Name:       "arc4",
			Driver:     "arc4-generic",
			Module:     "arc4",
			Priority:   new(int64(0)),
			Refcnt:     new(int64(3)),
			Selftest:   "passed",
			Internal:   "no",
			Type:       "cipher",
			Blocksize:  new(uint64(1)),
			MinKeysize: new(uint64(1)),
			MaxKeysize: new(uint64(256)),
		},
		{
			Name:       "crct10dif",
			Driver:     "crct10dif-pclmul",
			Module:     "crct10dif_pclmul",
			Priority:   new(int64(200)),
			Refcnt:     new(int64(2)),
			Selftest:   "passed",
			Internal:   "no",
			Type:       "shash",
			Blocksize:  new(uint64(1)),
			Digestsize: new(uint64(2)),
		},
	}

	if want, have := len(refs), len(crypto); want > have {
		t.Errorf("want at least %d parsed crypto-entries, have %d", want, have)
	}
	for index, ref := range refs {
		want, got := ref, crypto[index]
		if diff := cmp.Diff(want, got); diff != "" {
			t.Fatalf("unexpected crypto entry (-want +got):\n%s", diff)
		}
	}
}

func TestFS_CryptoCorrupted(t *testing.T) {
	cryptoFile = "crypto_corrupted"
	fs := getProcFixtures(t)
	crypto, err := fs.Crypto()

	if !errors.Is(err, ErrFileParse) {
		t.Fatalf("expected ErrFileParse error, got: %s", err)
	}

	if crypto != nil {
		t.Fatalf("expected empty crypto result")
	}
}
