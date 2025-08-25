package double

import (
	"testing"
)

func Test_SpyTestingT_ExpectRecords(t *testing.T) {
	t.Run("strict matching", func(t *testing.T) {
		t.Run("exact match", func(t *testing.T) {
			testedT := NewSpy(NewFake())
			testedT.Helper()
			testedT.Logf("hello %s, %d is the only way", "world", 42)
			testedT.Fail()

			testedT.ExpectRecords(t, true,
				SpyTestingTRecord{Method: "Helper"},
				SpyTestingTRecord{Method: "Logf", Inputs: []any{"hello %s, %d is the only way", []any{"world", 42}}},
				SpyTestingTRecord{Method: "Fail"},
			)
		})

		t.Run("any order", func(t *testing.T) {
			testedT := NewSpy(NewFake())
			testedT.Helper()
			testedT.Fail()

			spiedT := NewSpy(NewFake())
			testedT.ExpectRecords(spiedT, true,
				SpyTestingTRecord{Method: "Fail"},
				SpyTestingTRecord{Method: "Helper"},
			)
			spiedT.ExpectTestToFail(t)
			spiedT.ExpectLogsToContain(t, "Expected provided records to match")
		})

		t.Run("extra records", func(t *testing.T) {
			testedT := NewSpy(NewFake())
			testedT.Fail()
			testedT.Helper()

			spiedT := NewSpy(NewFake())
			testedT.ExpectRecords(spiedT, true, SpyTestingTRecord{Method: "Fail"})
			spiedT.ExpectTestToFail(t)
			spiedT.ExpectLogsToContain(t, "Expected provided records to match")
		})
	})

	t.Run("non-strict matching", func(t *testing.T) {
		t.Run("exact match", func(t *testing.T) {
			testedT := NewSpy(NewFake())
			testedT.Helper()
			testedT.Logf("hello %s, %d is the only way", "world", 42)
			testedT.Fail()

			testedT.ExpectRecords(t, false,
				SpyTestingTRecord{Method: "Helper"},
				SpyTestingTRecord{Method: "Logf", Inputs: []any{"hello %s, %d is the only way", []any{"world", 42}}},
				SpyTestingTRecord{Method: "Fail"},
			)
		})

		t.Run("any order", func(t *testing.T) {
			testedT := NewSpy(NewFake())
			testedT.Helper()
			testedT.Fail()

			testedT.ExpectRecords(t, false,
				SpyTestingTRecord{Method: "Fail"},
				SpyTestingTRecord{Method: "Helper"},
			)
		})

		t.Run("extra records", func(t *testing.T) {
			testedT := NewSpy(NewFake())
			testedT.Fail()
			testedT.Helper()

			testedT.ExpectRecords(t, false, SpyTestingTRecord{Method: "Fail"})
		})

		t.Run("missing expected", func(t *testing.T) {
			testedT := NewSpy(NewFake())
			testedT.Helper()
			testedT.Fail()

			spiedT := NewSpy(NewFake())
			testedT.ExpectRecords(spiedT, false, SpyTestingTRecord{Method: "Logf"})
			spiedT.ExpectTestToFail(t)
			spiedT.ExpectLogsToContain(t, "Missing expected records")
		})
	})
}

func Test_SpyTestingT_ExpectNoLogs(t *testing.T) {
	testedT := NewSpy(NewFake())
	testedT.ExpectNoLogs(t)
	testedT.Logf("hello world")

	spiedT := NewSpy(NewFake())
	testedT.ExpectNoLogs(spiedT)
	spiedT.ExpectTestToFail(t)
	spiedT.ExpectLogsToContain(t, "Expected no logs, got:\n\thello world")
}

func Test_SpyTestingT_ExpectLogsToContain(t *testing.T) {
	testedT := NewSpy(NewFake())
	spiedT := NewSpy(NewFake())

	testedT.ExpectLogsToContain(spiedT, "foo")
	spiedT.ExpectTestToFail(t)
	spiedT.ExpectLogsToContain(t, "Expected log to contain message")

	spiedT = NewSpy(NewFake())

	testedT.Logf("hello world")
	testedT.ExpectLogsToContain(t, "hello world")
	testedT.ExpectLogsToContain(spiedT, "foo")
	spiedT.ExpectTestToFail(t)
	spiedT.ExpectLogsToContain(t, "Expected log to contain message")

	testedT.Logf("another log")
	testedT.ExpectLogsToContain(t, "hello world")
	testedT.ExpectLogsToContain(t, "another log")
	testedT.ExpectLogsToContain(t, "another log", "hello world")
}

func Test_SpyTestingT_ExpectTestToFail(t *testing.T) {
	testedT := NewSpy(NewFake())
	spiedT := NewSpy(NewFake())

	testedT.ExpectTestToFail(spiedT)
	spiedT.ExpectTestToFail(t)
	spiedT.ExpectLogsToContain(t, "Expected test to fail but test succeeded")
	testedT.Fail()
	testedT.ExpectTestToFail(t)
}

func Test_SpyTestingT_ExpectTestToPass(t *testing.T) {
	testedT := NewSpy(NewFake())
	testedT.ExpectTestToPass(t)
	testedT.Fail()

	spiedT := NewSpy(NewFake())
	testedT.ExpectTestToPass(spiedT)
	spiedT.ExpectTestToFail(t)
	spiedT.ExpectLogsToContain(t, "Expected test to succeed but test failed")
}
