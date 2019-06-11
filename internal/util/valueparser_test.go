package util_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/procfs/internal/util"
)

func TestValueParser(t *testing.T) {
	tests := []struct {
		name string
		v    string
		ok   bool
		fn   func(t *testing.T, vp *util.ValueParser)
	}{
		{
			name: "bad PInt64",
			v:    "hello",
			fn: func(_ *testing.T, vp *util.ValueParser) {
				_ = vp.PInt64()
			},
		},
		{
			name: "bad hex PInt64",
			v:    "0xhello",
			fn: func(_ *testing.T, vp *util.ValueParser) {
				_ = vp.PInt64()
			},
		},
		{
			name: "ok PInt64",
			v:    "1",
			ok:   true,
			fn: func(t *testing.T, vp *util.ValueParser) {
				want := int64(1)
				got := vp.PInt64()

				if diff := cmp.Diff(&want, got); diff != "" {
					t.Fatalf("unexpected integer (-want +got):\n%s", diff)
				}
			},
		},
		{
			name: "ok hex PInt64",
			v:    "0xff",
			ok:   true,
			fn: func(t *testing.T, vp *util.ValueParser) {
				want := int64(255)
				got := vp.PInt64()

				if diff := cmp.Diff(&want, got); diff != "" {
					t.Fatalf("unexpected integer (-want +got):\n%s", diff)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vp := util.NewValueParser(tt.v)
			tt.fn(t, vp)

			err := vp.Err()
			if err != nil {
				if tt.ok {
					t.Fatalf("unexpected error: %v", err)
				}

				t.Logf("OK err: %v", err)
				return
			}

			if err == nil && !tt.ok {
				t.Fatal("expected an error, but none occurred")
			}
		})
	}
}
