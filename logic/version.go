package logic

import "fmt"

var (
	// ld should change these during build
	Version   = "undefined!"
	Commit    = "undefined!"
	BuildDate = "undefined!"
)

func VersionString() string {
	return fmt.Sprintf(
		"Reflector %s %s",
		Version, Commit)
}
