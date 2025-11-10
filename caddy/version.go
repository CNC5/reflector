package caddy

import (
	"fmt"
	"reflector/log"
	"regexp"
	"strconv"
)

type caddyVersion struct {
	Major int
	Minor int
	Patch int
}

func LoadCaddyVersion(version *string) *caddyVersion {
	defaultCaddyVersion := caddyVersion{Major: 0, Minor: 0, Patch: 0}
	if *version == "" {
		return &defaultCaddyVersion
	}
	validSemVerRegExStr := "^v{0,1}(?P<major>[0-9]*)\\.(?P<minor>[0-9]*)\\.(?P<patch>[0-9]*)"
	re := regexp.MustCompile(validSemVerRegExStr)
	semVer := re.FindStringSubmatch(*version)
	major, err := strconv.Atoi(semVer[re.SubexpIndex("major")])
	if err != nil {
		log.GetDefaultLogger().Error().
			Update("string", &version).
			Update("error", err).
			Update("substr", semVer[re.SubexpIndex("major")]).
			Msg("failed to get major version from string, fallback to default")
		return &defaultCaddyVersion
	}
	minor, err := strconv.Atoi(semVer[re.SubexpIndex("minor")])
	if err != nil {
		log.GetDefaultLogger().Error().
			Update("string", &version).
			Update("error", err).
			Update("substr", semVer[re.SubexpIndex("minor")]).
			Msg("failed to get minor version from string, fallback to default")
		return &defaultCaddyVersion
	}
	patch, err := strconv.Atoi(semVer[re.SubexpIndex("patch")])
	if err != nil {
		log.GetDefaultLogger().Error().
			Update("string", &version).
			Update("error", err).
			Update("substr", semVer[re.SubexpIndex("patch")]).
			Msg("failed to get patch version from string, fallback to default")
		return &defaultCaddyVersion
	}
	return &caddyVersion{Major: major, Minor: minor, Patch: patch}
}

func (cv *caddyVersion) Repr() string {
	return fmt.Sprintf("%d.%d.%d", cv.Major, cv.Minor, cv.Patch)
}

func (cv *caddyVersion) ReprV() string {
	return "v" + cv.Repr()
}
