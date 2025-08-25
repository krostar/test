package test

import (
	"testing"

	"github.com/krostar/test/double"
)

func Test_Assert(t *testing.T) {
	t.Run("assertion true", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			spiedT := double.NewSpy(double.NewFake())
			if result := Assert(spiedT, true, "hello from %s", t.Name()); !result {
				t.Error("Assert should return true when result is true")
			}

			spiedT.ExpectTestToPass(t)
			spiedT.ExpectNoLogs(t)
		})

		t.Run("with success message enabled", func(t *testing.T) {
			originalSuccessMessageEnabled := SuccessMessageEnabled
			t.Cleanup(func() { SuccessMessageEnabled = originalSuccessMessageEnabled })

			SuccessMessageEnabled = true

			spiedT := double.NewSpy(double.NewFake())
			if result := Assert(spiedT, true, "hello from %s", t.Name()); !result {
				t.Error("Assert should return true when result is true")
			}

			spiedT.ExpectTestToPass(t)
			spiedT.ExpectLogsToContain(t, "Success:", "[hello from Test_Assert/assertion_true/with_success_message_enabled]")
		})
	})

	t.Run("assertion false", func(t *testing.T) {
		spiedT := double.NewSpy(double.NewFake())
		if result := Assert(spiedT, false, "hello from %s", t.Name()); result {
			t.Error("Assert should return false when result is false")
		}

		spiedT.ExpectTestToFail(t)
		spiedT.ExpectLogsToContain(t, "Error:", "[hello from Test_Assert/assertion_false]")
	})
}

func Test_Require(t *testing.T) {
	t.Run("assertion true", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			spiedT := double.NewSpy(double.NewFake())
			Require(spiedT, true, "hello from %s", t.Name())
			spiedT.ExpectTestToPass(t)
			spiedT.ExpectNoLogs(t)
		})

		t.Run("with success message enabled", func(t *testing.T) {
			originalSuccessMessageEnabled := SuccessMessageEnabled
			t.Cleanup(func() { SuccessMessageEnabled = originalSuccessMessageEnabled })

			SuccessMessageEnabled = true

			spiedT := double.NewSpy(double.NewFake())
			Require(spiedT, true, "hello from %s", t.Name())
			spiedT.ExpectTestToPass(t)
			spiedT.ExpectLogsToContain(t, "Success:", "[hello from Test_Require/assertion_true/with_success_message_enabled]")
		})
	})

	t.Run("assertion false", func(t *testing.T) {
		spiedT := double.NewSpy(double.NewFake())
		Require(spiedT, false, "hello from %s", t.Name())
		spiedT.ExpectTestToFail(t)
		spiedT.ExpectRecords(t, false, double.SpyTestingTRecord{Method: "FailNow"})
		spiedT.ExpectLogsToContain(t, "Error:", "[hello from Test_Require/assertion_false]")
	})
}

func Test_logResult(t *testing.T) {
	t.Run("success without message", func(t *testing.T) {
		originalSuccessMessageEnabled := SuccessMessageEnabled
		t.Cleanup(func() { SuccessMessageEnabled = originalSuccessMessageEnabled })

		SuccessMessageEnabled = false

		spiedT := double.NewSpy(double.NewFake())
		logResult(spiedT, true, 0)
		spiedT.ExpectNoLogs(t)
	})

	t.Run("success with message", func(t *testing.T) {
		originalSuccessMessageEnabled := SuccessMessageEnabled
		t.Cleanup(func() { SuccessMessageEnabled = originalSuccessMessageEnabled })

		SuccessMessageEnabled = true

		spiedT := double.NewSpy(double.NewFake())
		logResult(spiedT, true, 0, "custom %s with %d values", "message", 42)
		spiedT.ExpectLogsToContain(t, "Success:", "custom message with 42 values")
	})

	t.Run("error with message", func(t *testing.T) {
		spiedT := double.NewSpy(double.NewFake())
		logResult(spiedT, false, 0, "failure reason")
		spiedT.ExpectLogsToContain(t, "Error:", "failure reason")
	})

	t.Run("empty message is skipped", func(t *testing.T) {
		spiedT := double.NewSpy(double.NewFake())
		logResult(spiedT, false, 0, "", "%s", "hello")
		spiedT.ExpectLogsToContain(t, "Error: literal false [hello]")
	})

	t.Run("first message is not a string", func(t *testing.T) {
		spiedT := double.NewSpy(double.NewFake())
		logResult(spiedT, false, 0, 42, "hello")
		spiedT.ExpectLogsToContain(t, "Error: literal false [42 hello]")
	})
}
