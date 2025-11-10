package log

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Logger struct {
	module              string
	stdout              *log.Logger
	logLevel            int
	callerDiscoverySkip int
}

var defaultLogger *Logger = &Logger{}

func NewLogger(module string, level int) *Logger {
	newLogger := Logger{
		module:              module,
		stdout:              log.New(os.Stdout, "", 1),
		callerDiscoverySkip: 3, // skip the log function and the wrapper
	}
	newLogger.SetLevel(level)
	newLogger.stdout.SetFlags(0)
	return &newLogger
}

func GetDefaultLogger() *Logger {
	if defaultLogger == (&Logger{}) {
		defaultLogger = NewLogger("", INFO)
	}
	return defaultLogger
}

func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}

func (l *Logger) logAsJSON(data map[string]any) {
	line, _ := json.Marshal(data)
	l.stdout.Println(string(line))
}

func (l *Logger) startNewChain(msgLevel int) *MsgChain {
	newChain := &MsgChain{messages: make(map[string]any), msgLevel: msgLevel, logLevel: l.logLevel, callback: l.logAsJSON}
	newChain.Update("level", logLevelToString(msgLevel))
	newChain.Update("caller", GetCallerName(l.callerDiscoverySkip))
	return newChain
}

func (l *Logger) SetLevel(level int) {
	l.logLevel = level % (MAXVERBOSELEVEL + 1)
}
func (l *Logger) GetLevel() string {
	return logLevelToString(l.logLevel)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.stdout.SetOutput(w)
}
func (l *Logger) Debug() *MsgChain {
	return l.startNewChain(DEBUG)
}
func (l *Logger) Info() *MsgChain {
	return l.startNewChain(INFO)
}
func (l *Logger) Warning() *MsgChain {
	return l.startNewChain(WARNING)
}
func (l *Logger) Error() *MsgChain {
	return l.startNewChain(ERROR)
}
func (l *Logger) Panic() *MsgChain {
	return l.startNewChain(PANIC)
}
func (l *Logger) Fatal() *MsgChain {
	return l.startNewChain(FATAL)
}
