package logging

import (
	"io"

	"github.com/krostar/test"
)

// NewWriter creates an io.Writer that forwards all written data to the
// testing instance's log. This is useful for capturing output from components
// that write to an io.Writer and redirecting it to the test log.
//
// Example:
//
//	logger := log.New(NewWriter(t), "PREFIX: ", 0)
//	logger.Println("This will appear in test logs")
func NewWriter(t test.TestingT) io.Writer { return loggingWriter{t} }

// loggingWriter implements io.Writer by forwarding all writes to TestingT.Logf
type loggingWriter struct{ t test.TestingT }

// Write implements io.Writer by sending data to the test log.
func (w loggingWriter) Write(p []byte) (int, error) {
	w.t.Helper()
	w.t.Logf("%s", string(p))
	return len(p), nil
}
