package xray

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"reflector/log"
	"reflector/utils"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-cmp/cmp"
)

type portableXray struct {
	version        *xrayVersion
	binaryLocation string
	xrayjson       *xrayJSON
	currentProcess *exec.Cmd
}

func (c *portableXray) xrayVersion() *xrayVersion {
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

func (c *portableXray) EnsureBinary() {
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

func (c *portableXray) uploadConfig() error {
	panic("not implemented")
	req, err := http.NewRequest(
		"POST",
		"http://localhost:2019/load",
		bytes.NewBuffer(c.xrayjson.Marshal()))
	log.GetDefaultLogger().Debug().
		Update("config", c.xrayjson.Marshal()).Msg("Uploading new xray config")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.GetDefaultLogger().Error().
			Update("status_code", resp.StatusCode).
			Update("body", resp.Body).
			Msg("failed to update xray config via API")
		return errors.New("non 200 status code")
	}
	log.GetDefaultLogger().Debug().
		Update("status_code", resp.StatusCode).Msg("upload config succeeded")
	return nil
}

func forwardXrayLogs(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.GetDefaultLogger().Debug().
			UpdateWithJSON(scanner.Text()).Done()
	}
}

func (c *portableXray) Start() {
	c.currentProcess = exec.Command(c.binaryLocation, "run")
	stdout, _ := c.currentProcess.StdoutPipe()
	stderr, _ := c.currentProcess.StderrPipe()

	if err := c.currentProcess.Start(); err != nil {
		log.GetDefaultLogger().Fatal().
			Update("error", err).Msg("failed to start:")
	}

	go forwardXrayLogs(stdout)
	go forwardXrayLogs(stderr)

	c.currentProcess.Start()
}

func (c *portableXray) Reload() {
	backoffDelay := time.Second
	maxBackoffDelay := 60 * time.Second
	maxRetry := 100
	for range maxRetry {
		err := c.uploadConfig()
		if err == nil {
			break
		}
		log.GetDefaultLogger().Error().
			Update("retry_delay_seconds", backoffDelay.Seconds()).Msg("config upload failed")
		time.Sleep(backoffDelay)
		if backoffDelay < maxBackoffDelay {
			backoffDelay *= 2
		} else {
			backoffDelay = 60 * time.Second
		}
	}
}

func (c *portableXray) Stop() {
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

func NewPortableXray(version string) *portableXray {
	newXray :=
		&portableXray{
			binaryLocation: "./xray-bin",
			xrayjson:       NewXrayJSON()}
	newXray.version = LoadXrayVersion(&version)
	newXray.EnsureBinary()
	log.GetDefaultLogger().Info().Msg("portable xray ready")
	return newXray
}
