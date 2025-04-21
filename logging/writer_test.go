package logging

import (
	"testing"

	"github.com/krostar/test/double"
)

func Test_NewWriter(t *testing.T) {
	for name, input := range map[string]string{
		"empty string": "",
		"not empty":    "test message",
	} {
		t.Run(name, func(t *testing.T) {
			spiedT := double.NewSpy(double.NewFake())
			writer := NewWriter(spiedT)

			n, err := writer.Write([]byte(input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if n != len(input) {
				t.Errorf("expected to write %d bytes but wrote %d", len(input), n)
			}

			spiedT.ExpectLogsToContain(t, input)
		})
	}
}
