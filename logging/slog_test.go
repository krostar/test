package logging

import (
	"log/slog"
	"testing"
	"time"

	"github.com/krostar/test/double"
)

func Test_NewSlogHandler(t *testing.T) {
	tests := map[string]struct {
		name        string
		level       slog.Level
		attrs       []slog.Attr
		recordAttrs []slog.Attr
		groups      []string
		message     string

		expected string
	}{
		"simple message": {
			level:    slog.LevelInfo,
			message:  "test message",
			expected: "level=INFO test message",
		},
		"message with attributes": {
			level:    slog.LevelWarn,
			attrs:    []slog.Attr{slog.String("key", "value")},
			message:  "with attributes",
			expected: "level=WARN key=value with attributes",
		},
		"message with group": {
			level:    slog.LevelError,
			groups:   []string{"group1"},
			message:  "with group",
			expected: "group1.level=ERROR with group",
		},
		"message with groups and attributes": {
			level:    slog.LevelInfo,
			groups:   []string{"group1", "group2"},
			attrs:    []slog.Attr{slog.String("key", "value")},
			message:  "complex message",
			expected: "group1.group2.level=INFO group1.group2.key=value complex message",
		},
		"message with record attributes": {
			level:       slog.LevelDebug,
			recordAttrs: []slog.Attr{slog.String("key", "value")},
			message:     "with record attributes",
			expected:    "level=DEBUG key=value with record attributes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spiedT := double.NewSpy(double.NewFake())
			handler := NewSlogHandler(spiedT)

			for _, group := range tt.groups {
				handler = handler.WithGroup(group)
			}

			if len(tt.attrs) > 0 {
				handler = handler.WithAttrs(tt.attrs)
			}

			record := slog.Record{
				Time:    time.Now(),
				Message: tt.message,
				Level:   tt.level,
			}
			record.AddAttrs(tt.recordAttrs...)

			if !handler.Enabled(t.Context(), record.Level) {
				t.Errorf("expected handler to be enabled for all levels, false for %s", record.Level.String())
			}

			err := handler.Handle(t.Context(), record)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			spiedT.ExpectLogsToContain(t, tt.expected)
		})
	}
}
