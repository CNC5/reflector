package xray

import (
	"fmt"
	"reflector/log"
	"regexp"
	"strconv"
)

type xrayVersion struct {
	Major int
	Minor int
	Patch int
}

func LoadXrayVersion(version *string) *xrayVersion {
	if *version == "" {
		return nil
	}
	validSemVerRegExStr := "^v{0,1}(?P<major>[0-9]*)\\.(?P<minor>[0-9]*)\\.(?P<patch>[0-9]*)"
	re := regexp.MustCompile(validSemVerRegExStr)
	semVer := re.FindStringSubmatch(*version)
	if len(semVer) == 0 {
		log.GetDefaultLogger().Error().
			Update("string", &version).
			Msg("failed to get any version from the string, fallback to default")
		return nil
	}
	major, err := strconv.Atoi(semVer[re.SubexpIndex("major")])
	if err != nil {
		log.GetDefaultLogger().Error().
			Update("string", &version).
			Update("error", err).
			Update("substr", semVer[re.SubexpIndex("major")]).
			Msg("failed to get major version from string, fallback to default")
		return nil
	}
	minor, err := strconv.Atoi(semVer[re.SubexpIndex("minor")])
	if err != nil {
		log.GetDefaultLogger().Error().
			Update("string", &version).
			Update("error", err).
			Update("substr", semVer[re.SubexpIndex("minor")]).
			Msg("failed to get minor version from string, fallback to default")
		return nil
	}
	patch, err := strconv.Atoi(semVer[re.SubexpIndex("patch")])
	if err != nil {
		log.GetDefaultLogger().Error().
			Update("string", &version).
			Update("error", err).
			Update("substr", semVer[re.SubexpIndex("patch")]).
			Msg("failed to get patch version from string, fallback to default")
		return nil
	}
	return &xrayVersion{Major: major, Minor: minor, Patch: patch}
}

func (cv *xrayVersion) Repr() string {
	return fmt.Sprintf("%d.%d.%d", cv.Major, cv.Minor, cv.Patch)
}

func (cv *xrayVersion) ReprV() string {
	return "v" + cv.Repr()
}
