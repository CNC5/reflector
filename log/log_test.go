package log

import (
	"testing"
)

type TestWriter struct {
	testContext *testing.T
}

func (w TestWriter) Write(p []byte) (n int, err error) {
	w.testContext.Logf("%s", string(p))
	return len(p), nil
}

func testAllMessages(t *testing.T, logLevel int) {
	lgr := NewLogger("log_test.go", logLevel)
	twr := TestWriter{testContext: t}
	lgr.SetOutput(twr)
}

func TestLoggerLevelDebug(t *testing.T) {
	testAllMessages(t, DEBUG)
}
func TestLoggerLevelInfo(t *testing.T) {
	testAllMessages(t, INFO)
}
func TestLoggerLevelWarning(t *testing.T) {
	testAllMessages(t, WARNING)
}
func TestLoggerLevelError(t *testing.T) {
	testAllMessages(t, ERROR)
}
func TestLoggerLevelPanic(t *testing.T) {
	testAllMessages(t, PANIC)
}
func TestLoggerLevelFatal(t *testing.T) {
	testAllMessages(t, FATAL)
}
