package log

import (
	"runtime"
	"strings"
)

func GetCallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		fullName := details.Name()
		if strings.Contains(fullName, "/") {
			parts := strings.Split(fullName, "/")
			return parts[len(parts)-1]
		} else {
			return fullName
		}
	}
	return ""
}
