package test

import (
	"context"
	"testing"
	"time"

	"github.com/krostar/test/double"
)

func Test_Context(t *testing.T) {
	t.Run("with deadline", func(t *testing.T) {
		for _, tt := range []struct {
			deadline time.Time
			lower    time.Duration
			upper    time.Duration
		}{
			{
				deadline: time.Now().Add(time.Second),
				lower:    9 * time.Millisecond,
				upper:    10 * time.Millisecond,
			},
			{
				deadline: time.Now().Add(time.Minute),
				lower:    599 * time.Millisecond,
				upper:    600 * time.Millisecond,
			},
			{
				deadline: time.Now().Add(time.Hour),
				lower:    time.Second,
				upper:    time.Second,
			},
		} {
			spiedT := double.NewSpy(double.NewFake(
				double.FakeWithDeadline(tt.deadline),
				double.FakeWithContext(func() context.Context {
					ctx, cancel := context.WithDeadline(t.Context(), tt.deadline)
					t.Cleanup(cancel)
					return ctx
				}()),
			))

			newDeadline, hasDeadline := Context(spiedT).Deadline()
			if !hasDeadline {
				t.Error("expected context to have a deadline")
			}

			if !newDeadline.Before(tt.deadline) {
				t.Error("expected new deadline to be before old deadline")
			}

			if d := tt.deadline.Sub(newDeadline); d < tt.lower || d > tt.upper {
				t.Errorf("expected delta between new and old deadline to be between %d;%d, got %d", tt.lower, tt.upper, d)
			}

			spiedT.ExpectRecords(t, true,
				double.SpyTestingTRecord{Method: "Context", Outputs: []any{double.SpyTestingTRecordIgnoreParam}},
				double.SpyTestingTRecord{Method: "Deadline", Outputs: []any{tt.deadline, true}},
				double.SpyTestingTRecord{Method: "Cleanup", Inputs: []any{double.SpyTestingTRecordIgnoreParam}},
			)
		}
	})

	t.Run("with long deadline ", func(t *testing.T) {
		deadline := time.Now().Add(time.Hour)

		spiedT := double.NewSpy(double.NewFake(
			double.FakeWithDeadline(deadline),
			double.FakeWithContext(func() context.Context {
				ctx, cancel := context.WithDeadline(t.Context(), deadline)
				t.Cleanup(cancel)
				return ctx
			}()),
		))

		newDeadline, hasDeadline := Context(spiedT).Deadline()
		if !hasDeadline {
			t.Error("expected context to have a deadline")
		}

		if !newDeadline.Before(deadline) {
			t.Error("expected new deadline to be before old deadline")
		}

		spiedT.ExpectRecords(t, true,
			double.SpyTestingTRecord{Method: "Context", Outputs: []any{double.SpyTestingTRecordIgnoreParam}},
			double.SpyTestingTRecord{Method: "Deadline", Outputs: []any{deadline, true}},
			double.SpyTestingTRecord{Method: "Cleanup", Inputs: []any{double.SpyTestingTRecordIgnoreParam}},
		)
	})

	t.Run("without deadline", func(t *testing.T) {
		spiedT := double.NewSpy(double.NewFake())

		_, hasDeadline := Context(spiedT).Deadline()
		if hasDeadline {
			t.Error("expected context to have no deadline")
		}

		spiedT.ExpectRecords(t, true,
			double.SpyTestingTRecord{Method: "Context", Outputs: []any{double.SpyTestingTRecordIgnoreParam}},
			double.SpyTestingTRecord{Method: "Deadline", Outputs: []any{time.Time{}, false}},
		)
	})
}
