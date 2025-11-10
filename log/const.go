package log

const (
	DEBUG           = 5
	INFO            = 4
	WARNING         = 3
	ERROR           = 2
	PANIC           = 1
	FATAL           = 0
	MAXVERBOSELEVEL = DEBUG
)

func logLevelToString(level int) string {
	switch i := level; i {
	case 0:
		return "fatal"
	case 1:
		return "panic"
	case 2:
		return "error"
	case 3:
		return "warning"
	case 4:
		return "info"
	case 5:
		return "debug"
	default:
		return "unknown"
	}
}
