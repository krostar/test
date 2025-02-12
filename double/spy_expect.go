package double

import (
	"strings"

	gocmp "github.com/google/go-cmp/cmp"
)

// ExpectRecords verifies that the spy contains the expected method call records.
//
// Two matching modes are supported:
//   - Strict mode (strict=true): The actual records must exactly match the expected records
//     in the same order.
//   - Non-strict mode (strict=false): The expected records can appear in any order within the
//     actual records, and additional unexpected records are allowed.
//
// In non-strict mode, all expected records must still exist in the actual records, but:
//   - The order doesn't matter (records can be matched in any order)
//   - Extra records in the spy beyond what's expected are ignored
//   - All expected records must be present in the actual records
//
// The method fails the test if the verification doesn't pass.
func (spy *Spy) ExpectRecords(t TestingT, strict bool, expected ...SpyTestingTRecord) {
	spy.m.RLock()
	defer spy.m.RUnlock()

	gocmpOpts := []gocmp.Option{gocmp.Comparer(func(a, b SpyTestingTRecord) bool { return a.seemsEqualTo(b) })}

	t.Helper()

	if strict {
		if diff := gocmp.Diff(spy.records, expected, gocmpOpts...); diff != "" {
			t.Logf("Expected provided records to match\n%s", diff)
			t.Fail()
		}
		return
	}

	foundExpected := make(map[int]bool, len(expected))
	for idx := range expected {
		foundExpected[idx] = false
	}

	for _, actual := range spy.records {
		for i, exp := range expected {
			if foundExpected[i] {
				continue
			}

			if gocmp.Equal(actual, exp, gocmpOpts...) {
				foundExpected[i] = true
				break
			}
		}
	}

	var missingExpectedRecords []SpyTestingTRecord
	for k, v := range foundExpected {
		if !v {
			missingExpectedRecords = append(missingExpectedRecords, expected[k])
		}
	}

	if len(missingExpectedRecords) > 0 {
		t.Logf("Missing expected records:\n%s", gocmp.Diff([]SpyTestingTRecord{}, missingExpectedRecords, gocmpOpts...))
		t.Fail()
	}
}

// ExpectNoLogs verifies that no logs were captured by the spy.
// Fails the test if any logs were captured.
// This is useful for ensuring that no messages were logged during the test.
func (spy *Spy) ExpectNoLogs(t TestingT) {
	spy.m.RLock()
	defer spy.m.RUnlock()

	t.Helper()

	if len(spy.logs) > 0 {
		t.Logf("Expected no logs, got:\n\t%s", strings.Join(spy.logs, "\n"))
		t.Fail()
	}
}

// ExpectLogsToContain verifies that all the provided strings are contained within the spy's logs.
// Fails the test if any of the strings are not found in the concatenated logs.
// This is useful for verifying that specific messages were logged during the test.
func (spy *Spy) ExpectLogsToContain(t TestingT, expect string, more ...string) {
	spy.m.RLock()
	defer spy.m.RUnlock()

	t.Helper()

	log := strings.Join(spy.logs, "\n")

	for _, str := range append([]string{expect}, more...) {
		if !strings.Contains(log, str) {
			t.Logf("Expected log to contain message:\nexpected: %s\nlog: %s", str, log)
			t.Fail()
		}
	}
}

// ExpectTestToFail verifies that the test failed.
// Fails the test if no failure was recorded.
// This is useful for testing assertion functions that should fail tests.
func (spy *Spy) ExpectTestToFail(t TestingT) {
	spy.m.RLock()
	defer spy.m.RUnlock()

	t.Helper()

	if !spy.failed {
		t.Logf("Expected test to fail but test succeeded")
		t.Fail()
	}
}

// ExpectTestToPass verifies that the test passed.
// Fails the test if any failure was recorded.
// This is useful for testing assertion functions that should pass tests.
func (spy *Spy) ExpectTestToPass(t TestingT) {
	spy.m.RLock()
	defer spy.m.RUnlock()

	t.Helper()

	if spy.failed {
		t.Logf("Expected test to succeed but test failed")
		t.Fail()
	}
}
