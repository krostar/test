package double

import "reflect"

// spyTestingTRecordIgnoreParam is a special type used as a marker for parameters
// that should be ignored during comparison in Spy expectations.
type spyTestingTRecordIgnoreParam uint

// SpyTestingTRecordIgnoreParam is a sentinel value used in expected records
// to indicate that a particular parameter should be ignored during comparison.
// This is useful for testing calls with parameters that are unpredictable or irrelevant.
//
// Example:
//
//	spy.ExpectRecords(t, true, SpyTestingTRecord{
//		Method: "Logf",
//		Inputs: []any{"Error: %v", SpyTestingTRecordIgnoreParam}, // Don't care about the exact error
//	})
const SpyTestingTRecordIgnoreParam = spyTestingTRecordIgnoreParam(42)

// SpyTestingTRecord represents a single method call recorded by Spy.
// It captures the method name along with its inputs and outputs.
type SpyTestingTRecord struct {
	Method  string // Name of the method called
	Inputs  []any  // Arguments passed to the method (if any)
	Outputs []any  // Return values from the method (if any)
}

// seemsEqualTo compares two SpyTestingTRecord instances for practical equality.
// It implements a special comparison that:
// - Requires method names to match exactly
// - Handles function values specially (only comparing nil status)
// - Ignores parameters marked with SpyTestingTRecordIgnoreParam
// - Requires the same number of parameters in the same positions
//
// This method is used by ExpectRecords to determine if an expected record
// matches an actual record.
func (a SpyTestingTRecord) seemsEqualTo(b SpyTestingTRecord) bool {
	if a.Method != b.Method {
		return false
	}

	ignore := reflect.TypeOf(SpyTestingTRecordIgnoreParam)

	assertParams := func(x, y []any) bool {
		lenX, lenY := len(x), len(y)
		if lenX != lenY {
			return false
		}

		for i := range lenX {
			switch ia, ib := reflect.ValueOf(x[i]), reflect.ValueOf(y[i]); {
			case ia.Type() == ignore || ib.Type() == ignore:
			case ia.Type() != ib.Type():
				return false
			case ia.Kind() == reflect.Func && ia.IsNil() != ib.IsNil():
				return false
			}
		}

		return true
	}

	if !assertParams(a.Inputs, b.Inputs) {
		return false
	}

	if !assertParams(a.Outputs, b.Outputs) {
		return false
	}

	return true
}
