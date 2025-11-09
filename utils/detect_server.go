package utils

import (
	"reflector/log"
	"regexp"
)

func DetectExistingServer() string {
	procs := PSStrings()
	re := regexp.MustCompile("nginx|caddy")
	for _, p := range *procs {
		matches := re.FindStringSubmatch(p.Cmdline)
		if len(matches) > 0 {
			log.GetDefaultLogger().Debug().
				Update("matches", matches).
				Update("pid", p.PID).Done()
			return matches[0]
		}
	}
	return "unknown"
}
