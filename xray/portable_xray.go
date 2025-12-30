package xray

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflector/log"
	"reflector/utils"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-cmp/cmp"
)

type PortableXray struct {
	version        *xrayVersion
	binaryLocation string
	configLocation string
	XrayConfig     *XrayConfig
	currentProcess *exec.Cmd
}

func (c *PortableXray) xrayVersion() *xrayVersion {
	outputBytes, err := exec.Command(c.binaryLocation, "version").Output()
	if err != nil {
		return &xrayVersion{Major: 0, Minor: 0, Patch: 0}
	}
	splitOutput := strings.Split(string(outputBytes), " ")
	// from 'Xray X.Y.Z' split and select version
	if len(splitOutput) < 2 {
		return &xrayVersion{Major: 0, Minor: 0, Patch: 0}
	}
	output := splitOutput[1]
	return LoadXrayVersion(&output)
}

func (c *PortableXray) EnsureBinary() {
	zipURL := "https://github.com/XTLS/Xray-core/releases/download/" + c.version.ReprV() + "/Xray-linux-64.zip"
	zipLocation := "/tmp/Xray-linux-64.zip"
	zipBinPath := "xray"
	installedXrayVersion := c.xrayVersion()
	if _, err :=
		os.Stat(zipLocation); errors.Is(err, os.ErrNotExist) ||
		!cmp.Equal(installedXrayVersion, c.version) {
		log.
			GetDefaultLogger().Info().
			Update("current_version", installedXrayVersion).
			Update("desired_version", c.version).
			Msg("new xray binary required, downloading")
		utils.DownloadFile(zipLocation, zipURL)
	} else {
		log.
			GetDefaultLogger().Debug().
			Update("current_version", installedXrayVersion).
			Update("desired_version", c.version).
			Msg("current xray is good")
		return
	}
	if _, err :=
		os.Stat(c.binaryLocation); errors.Is(err, os.ErrNotExist) ||
		!cmp.Equal(installedXrayVersion, c.version) {
		utils.UnpackZipSubpath(zipLocation, zipBinPath, c.binaryLocation)
	}
	os.Chmod(c.binaryLocation, 0o755)
}

func (c *PortableXray) updateConfig() error {
	// TODO load config first for crypto persistence
	configBytes, err := json.Marshal(c.XrayConfig)
	if err != nil {
		return fmt.Errorf("error marshaling config: %s", err.Error())
	}
	configFile, err := os.Create(c.configLocation)
	if err != nil {
		return fmt.Errorf("error opening config file: %s", err.Error())
	}
	_, err = configFile.Write(configBytes)
	if err != nil {
		return fmt.Errorf("error writing config file: %s", err.Error())
	}
	return nil
}

func forwardXrayLogs(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.GetDefaultLogger().Debug().
			UpdateWithJSON(scanner.Text()).Done()
	}
}

func (c *PortableXray) Start() {
	if err := c.updateConfig(); err != nil {
		panic(err)
	}
	c.currentProcess = exec.Command(
		c.binaryLocation, "run",
		"-c", c.configLocation,
	)
	stdout, _ := c.currentProcess.StdoutPipe()
	stderr, _ := c.currentProcess.StderrPipe()

	if err := c.currentProcess.Start(); err != nil {
		panic(err)
	}

	go forwardXrayLogs(stdout)
	go forwardXrayLogs(stderr)

	c.currentProcess.Start()
}

func (c *PortableXray) Reload() {
	backoffDelay := time.Second
	maxBackoffDelay := 60 * time.Second
	maxRetry := 100
	for range maxRetry {
		err := c.updateConfig()
		if err == nil {
			break
		}
		log.GetDefaultLogger().Error().
			Update("retry_delay_seconds", backoffDelay.Seconds()).Msg("config update failed")
		time.Sleep(backoffDelay)
		if backoffDelay < maxBackoffDelay {
			backoffDelay *= 2
		} else {
			backoffDelay = 60 * time.Second
		}
	}
}

func (c *PortableXray) Stop() {
	syscall.Kill(c.currentProcess.Process.Pid, syscall.SIGTERM)
	c.currentProcess.Wait()
	// for range 30 {
	// 	time.Sleep(time.Second)
	// 	if c.currentProcess.ProcessState.Exited() {
	// 		return
	// 	}
	// }
	//syscall.Kill(c.currentProcess.Process.Pid, syscall.SIGKILL)
}

func NewPortableXray(version string) *PortableXray {
	newXray :=
		&PortableXray{
			binaryLocation: "./xray-bin",
			XrayConfig:     NewXrayConfig(),
			configLocation: "./xray-config.json",
		}
	newXray.version = LoadXrayVersion(&version)
	newXray.EnsureBinary()
	log.GetDefaultLogger().Info().Msg("portable xray ready")
	return newXray
}
