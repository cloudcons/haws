package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput captures stdout for testing log output
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	return <-outC
}

func TestLogLevels(t *testing.T) {
	testCases := []struct {
		name          string
		level         LogLevel
		logFunc       func(string, ...interface{})
		message       string
		expectedLevel string
		shouldContain bool
	}{
		{"Debug when Debug enabled", LevelDebug, Debug, "test debug message", "DEBUG", true},
		{"Info when Debug enabled", LevelDebug, Info, "test info message", "INFO", true},
		{"Warn when Debug enabled", LevelDebug, Warn, "test warn message", "WARN", true},
		{"Error when Debug enabled", LevelDebug, Error, "test error message", "ERROR", true},
		
		{"Debug when Info enabled", LevelInfo, Debug, "test debug message", "DEBUG", false},
		{"Info when Info enabled", LevelInfo, Info, "test info message", "INFO", true},
		{"Warn when Info enabled", LevelInfo, Warn, "test warn message", "WARN", true},
		{"Error when Info enabled", LevelInfo, Error, "test error message", "ERROR", true},
		
		{"Debug when Warn enabled", LevelWarn, Debug, "test debug message", "DEBUG", false},
		{"Info when Warn enabled", LevelWarn, Info, "test info message", "INFO", false},
		{"Warn when Warn enabled", LevelWarn, Warn, "test warn message", "WARN", true},
		{"Error when Warn enabled", LevelWarn, Error, "test error message", "ERROR", true},
		
		{"Debug when Error enabled", LevelError, Debug, "test debug message", "DEBUG", false},
		{"Info when Error enabled", LevelError, Info, "test info message", "INFO", false},
		{"Warn when Error enabled", LevelError, Warn, "test warn message", "WARN", false},
		{"Error when Error enabled", LevelError, Error, "test error message", "ERROR", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the log level for this test
			oldLevel := CurrentLevel
			SetLevel(tc.level)
			defer SetLevel(oldLevel)

			output := captureOutput(func() {
				tc.logFunc(tc.message)
			})

			// Check if the message was logged (we don't check for the level since it's colorized)
			messageContains := strings.Contains(output, tc.message)

			if tc.shouldContain && !messageContains {
				t.Errorf("Expected output to contain '%s', got: %s", tc.message, output)
			} else if !tc.shouldContain && messageContains {
				t.Errorf("Expected output to NOT contain '%s', but it did: %s", tc.message, output)
			}
		})
	}
}

func TestSetLevel(t *testing.T) {
	oldLevel := CurrentLevel
	defer SetLevel(oldLevel)

	levels := []LogLevel{LevelDebug, LevelInfo, LevelWarn, LevelError}
	
	for _, level := range levels {
		SetLevel(level)
		if CurrentLevel != level {
			t.Errorf("SetLevel(%d) didn't set CurrentLevel to %d, it's %d", level, level, CurrentLevel)
		}
	}
}
