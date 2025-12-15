package camo

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
	"reflector/log"
	"reflector/utils"
	"strings"
)

type CamoController struct {
	camoDir               string
	containerController   *utils.ContainerImageController
	preparedCamoLocations map[string]string
}

func NewCamoController() *CamoController {
	return &CamoController{
		camoDir:               "/tmp/camo/",
		containerController:   utils.NewContainerImageController(),
		preparedCamoLocations: make(map[string]string),
	}
}

// template can be any of
//   - `oci://docker.io/library/image:example`
//   - `./image`
//   - `/image`
//   - `image`
func (cc *CamoController) camoInternalNameFromTemplate(template string) string {
	untaggedName := template
	if strings.Count(template, ":") > 0 {
		// >...image< :example
		spl := strings.Split(template, ":")
		untaggedName = spl[len(spl)-2]
	}

	internalName := untaggedName
	if strings.Count(untaggedName, "/") > 0 {
		// ...library/ >image<
		spl := strings.Split(untaggedName, "/")
		internalName = spl[len(spl)-1]
	}

	// many templates can result in the same name, add the template hash
	templateHashBytes := sha1.Sum([]byte(template))
	templateHash := hex.EncodeToString(templateHashBytes[:])
	internalName += templateHash
	return internalName
}

func (cc *CamoController) camoFullPathForTemplate(template string) string {
	return strings.TrimSuffix(cc.camoDir, "/") +
		"/" +
		cc.camoInternalNameFromTemplate(template)
}

// Default camo loading, pull policy if not present
func (cc *CamoController) PreLoadCamo(template string) error {
	return cc.preLoadCamo(template, false)
}

// Camo loading with controlled pull policy
func (cc *CamoController) PreLoadCamoWithPullPolicy(template string, alwaysPull bool) error {
	return cc.preLoadCamo(template, alwaysPull)
}

func (cc *CamoController) preLoadCamo(template string, alwaysPull bool) error {
	if existingLocation, exists := cc.preparedCamoLocations[template]; exists {
		// already loaded and accounted
		// ignore alwaysPull, it was respected when the accounting happened
		log.GetDefaultLogger().
			Info().
			Update("template", template).
			Update("location", existingLocation).
			Msg("camo is already loaded")
		return nil
	}
	camoLocation := cc.camoFullPathForTemplate(template)

	if stat, err := os.Stat(camoLocation); os.IsExist(err) || (stat != nil && stat.IsDir()) {
		// already downloaded but not accounted
		// repull if alwaysPull
		if !alwaysPull {
			cc.preparedCamoLocations[template] = camoLocation
			log.GetDefaultLogger().
				Info().
				Update("template", template).
				Update("location", camoLocation).
				Msg("camo is already downloaded")
			return nil
		}
		log.GetDefaultLogger().
			Debug().
			Msg("camo dir already exists, cleaning and redownloading")
		err := os.RemoveAll(camoLocation)
		if err != nil {
			return err
		}
	}
	err := cc.containerController.UnpackImage(template, camoLocation)
	if err != nil {
		return err
	}
	cc.preparedCamoLocations[template] = camoLocation
	log.GetDefaultLogger().
		Info().
		Update("template", template).
		Update("location", camoLocation).
		Msg("loaded camo")
	return nil
}

func (cc *CamoController) CamoLocation(template string) (string, error) {
	if location, exists := cc.preparedCamoLocations[template]; exists {
		return location, nil
	}
	return "", errors.New("this camo was never loaded")
}
