package double

import (
	"testing"
)

func TestSpyTestingTRecord_seemsEqualTo(t *testing.T) {
	for name, tt := range map[string]struct {
		a    SpyTestingTRecord
		b    SpyTestingTRecord
		want bool
	}{
		"identical records": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			want: true,
		},
		"different methods": {
			a: SpyTestingTRecord{
				Method:  "TestMethod1",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod2",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			want: false,
		},
		"different inputs length": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1"},
				Outputs: []any{"output1", true},
			},
			want: false,
		},
		"different outputs length": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1"},
			},
			want: false,
		},
		"different input types": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", "42"},
				Outputs: []any{"output1", true},
			},
			want: false,
		},
		"different output types": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", 1},
			},
			want: false,
		},
		"with ignore parameter in inputs": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", SpyTestingTRecordIgnoreParam},
				Outputs: []any{"output1", true},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", true},
			},
			want: true,
		},
		"with ignore parameter in outputs": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", SpyTestingTRecordIgnoreParam},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", 42},
				Outputs: []any{"output1", false},
			},
			want: true,
		},
		"comparing nil and non-nil functions": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", (func())(nil)},
				Outputs: []any{"output1"},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", func() {}},
				Outputs: []any{"output1"},
			},
			want: false,
		},
		"comparing nil functions": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", (func())(nil)},
				Outputs: []any{"output1"},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", (func())(nil)},
				Outputs: []any{"output1"},
			},
			want: true,
		},
		"comparing non-nil functions": {
			a: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", func() {}},
				Outputs: []any{"output1"},
			},
			b: SpyTestingTRecord{
				Method:  "TestMethod",
				Inputs:  []any{"input1", func() {}},
				Outputs: []any{"output1"},
			},
			want: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := tt.a.seemsEqualTo(tt.b); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
