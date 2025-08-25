package double

import (
	"testing"
	"time"
)

func Test_SpyTestingT_Helper(t *testing.T) {
	spiedT := NewSpy(NewFake())
	spiedT.Helper()
	spiedT.ExpectRecords(t, true, SpyTestingTRecord{Method: "Helper"})
}

func Test_SpyTestingT_Cleanup(t *testing.T) {
	var cleanup func()

	spiedT := NewSpy(NewFake(FakeWithRegisterCleanup(func(f func()) { cleanup = f })))

	cleanupCalled := false
	spiedT.Cleanup(func() { cleanupCalled = true })

	spiedT.ExpectRecords(t, true, SpyTestingTRecord{
		Method: "Cleanup",
		Inputs: []any{SpyTestingTRecordIgnoreParam},
	})

	if cleanupCalled {
		t.Error("cleanup function should not be called yet")
	}

	cleanup()

	if !cleanupCalled {
		t.Error("cleanup function not called")
	}
}

func Test_SpyTestingT_Fail(t *testing.T) {
	spiedT := NewSpy(NewFake())
	spiedT.Fail()
	spiedT.ExpectTestToFail(t)
	spiedT.ExpectRecords(t, true, SpyTestingTRecord{Method: "Fail"})
}

func Test_SpyTestingT_FailNow(t *testing.T) {
	spiedT := NewSpy(NewFake())
	spiedT.FailNow()
	spiedT.ExpectTestToFail(t)
	spiedT.ExpectRecords(t, true, SpyTestingTRecord{Method: "FailNow"})
}

func Test_SpyTestingT_Log(t *testing.T) {
	spiedT := NewSpy(NewFake())
	spiedT.Log("hello", "world")
	spiedT.ExpectLogsToContain(t, "hello", "world")
	spiedT.ExpectRecords(t, true, SpyTestingTRecord{
		Method: "Log",
		Inputs: []any{"hello", "world"},
	})
}

func Test_SpyTestingT_Logf(t *testing.T) {
	spiedT := NewSpy(NewFake())
	spiedT.Logf("hello %s", "world")
	spiedT.ExpectLogsToContain(t, "hello world")
	spiedT.ExpectRecords(t, true, SpyTestingTRecord{
		Method: "Logf",
		Inputs: []any{"hello %s", []any{"world"}},
	})
}

func Test_SpyTestingT_Context(t *testing.T) {
	ctx := t.Context()

	spiedT := NewSpy(NewFake(FakeWithContext(ctx)))

	if spiedT.Context() != ctx {
		t.Error("Context did not return the expected context")
	}

	spiedT.ExpectRecords(t, true, SpyTestingTRecord{
		Method:  "Context",
		Outputs: []any{ctx},
	})
}

func Test_SpyTestingT_Deadline(t *testing.T) {
	deadline := time.Now().Add(time.Hour)

	spiedT := NewSpy(NewFake(FakeWithDeadline(deadline)))
	returnedDeadline, hasDeadline := spiedT.Deadline()

	if !returnedDeadline.Equal(deadline) {
		t.Error("Deadline did not return the expected deadline")
	}

	if !hasDeadline {
		t.Error("Deadline should have returned hasDeadline=true")
	}

	spiedT.ExpectRecords(t, true, SpyTestingTRecord{
		Method:  "Deadline",
		Outputs: []any{deadline, true},
	})
}
