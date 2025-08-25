package logging

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/krostar/test"
)

// NewSlogHandler creates a new slog.Handler that forwards logs to a testing instance.
// By default, it uses slog.LevelInfo as the minimum log level.
func NewSlogHandler(t test.TestingT) slog.Handler {
	return &slogHandler{t: t}
}

// slogHandler is a slog.Handler implementation that forwards all log
// records to a TestingT instance. This allows capturing structured logs emitted
// by the code under test directly in the test output.
//
// The handler supports level filtering, attribute collection, and group nesting.
// Log messages will be formatted as "group.subgroup.level=LEVEL group.subgroup.attr=value message".
type slogHandler struct {
	m sync.Mutex
	t test.TestingT

	attrs  []slog.Attr
	groups []string
}

// Enabled checks if the provided log level meets or exceeds the handler's configured minimum level.
func (*slogHandler) Enabled(context.Context, slog.Level) bool { return true }

// Handle formats the log record and its attributes, then forwarding it to the test log.
//
//nolint:gocritic // record is huge to be passed by copy, but its slog's decision
func (h *slogHandler) Handle(_ context.Context, record slog.Record) error {
	h.m.Lock()
	defer h.m.Unlock()

	var attrs []string

	attrs = append(attrs, fmt.Sprintf("%s=%s", strings.Join(append(h.groups, "level"), "."), record.Level.String()))
	for _, attr := range h.attrs {
		attrs = append(attrs, fmt.Sprintf("%s=%s", strings.Join(append(h.groups, attr.Key), "."), attr.Value.Any()))
	}

	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%s", strings.Join(append(h.groups, attr.Key), "."), attr.Value.Any()))
		return true
	})

	h.t.Logf("%s %s", strings.Join(attrs, " "), record.Message)

	return nil
}

// WithAttrs creates a new handler with the combined attributes from this handler and the provided attributes.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.m.Lock()
	defer h.m.Unlock()

	return &slogHandler{
		t:      h.t,
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
	}
}

// WithGroup creates a new handler with the provided group name appended to the existing group path.
func (h *slogHandler) WithGroup(name string) slog.Handler {
	h.m.Lock()
	defer h.m.Unlock()

	return &slogHandler{
		t:      h.t,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
