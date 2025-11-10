package caddy

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

type portableCaddy struct {
	version        *caddyVersion
	binaryLocation string
	caddyjson      *caddyJSON
	currentProcess *exec.Cmd
}

func (c *portableCaddy) caddyVersion() *caddyVersion {
	outputBytes, _ := exec.Command(c.binaryLocation, "version").Output()
	output := strings.Split(string(outputBytes), " ")[0] // from 'vX.Y.Z <HASH>' split and select version
	return LoadCaddyVersion(&output)
}

func (c *portableCaddy) EnsureBinary() {
	debURL := "https://github.com/caddyserver/caddy/releases/download/" + c.version.ReprV() + "/caddy_" + c.version.Repr() + "_linux_amd64.deb"
	debLocation := "/tmp/caddy.deb"
	debBinPath := "./usr/bin/caddy"
	installedCaddyVersion := c.caddyVersion()
	if _, err :=
		os.Stat(debLocation); errors.Is(err, os.ErrNotExist) ||
		!cmp.Equal(installedCaddyVersion, c.version) {
		log.
			GetDefaultLogger().Info().
			Update("current_version", installedCaddyVersion).
			Update("desired_version", c.version).
			Msg("new caddy binary required, downloading")
		utils.DownloadFile(debLocation, debURL)
	} else {
		log.
			GetDefaultLogger().Debug().
			Update("current_version", installedCaddyVersion).
			Update("desired_version", c.version).
			Msg("current caddy is good")
		return
	}
	if _, err :=
		os.Stat(c.binaryLocation); errors.Is(err, os.ErrNotExist) ||
		!cmp.Equal(installedCaddyVersion, c.version) {
		utils.UnpackDebSubpath(debLocation, debBinPath, c.binaryLocation)
	}
	os.Chmod(c.binaryLocation, 0o755)
}

func (c *portableCaddy) uploadConfig() error {
	req, err := http.NewRequest(
		"POST",
		"http://localhost:2019/load",
		bytes.NewBuffer(c.caddyjson.Marshal()))
	log.GetDefaultLogger().Debug().
		Update("config", c.caddyjson.Marshal()).Msg("Uploading new caddy config")
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
			Msg("failed to update caddy config via API")
		return errors.New("non 200 status code")
	}
	log.GetDefaultLogger().Debug().
		Update("status_code", resp.StatusCode).Msg("upload config succeeded")
	return nil
}

func forwardCaddyLogs(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log.GetDefaultLogger().Debug().
			UpdateWithJSON(scanner.Text()).Done()
	}
}

func (c *portableCaddy) Start() {
	c.currentProcess = exec.Command(c.binaryLocation, "run")
	stdout, _ := c.currentProcess.StdoutPipe()
	stderr, _ := c.currentProcess.StderrPipe()

	if err := c.currentProcess.Start(); err != nil {
		log.GetDefaultLogger().Fatal().
			Update("error", err).Msg("failed to start:")
	}

	go forwardCaddyLogs(stdout)
	go forwardCaddyLogs(stderr)

	c.currentProcess.Start()
}

func (c *portableCaddy) Reload() {
	backoffInitDelay := time.Second
	maxRetry := 100
	for range maxRetry {
		err := c.uploadConfig()
		if err == nil {
			break
		}
		log.GetDefaultLogger().Error().
			Update("retry_delay_seconds", backoffInitDelay.Seconds()).Msg("config upload failed")
		time.Sleep(backoffInitDelay)
		if backoffInitDelay < 60*time.Second {
			backoffInitDelay *= 2
		} else {
			backoffInitDelay = 60 * time.Second
		}
	}
}

func (c *portableCaddy) AddRootStaticLocation(domain string, staticDir string) {
	c.caddyjson.AddRootStaticLocation(domain, staticDir)
}

func (c *portableCaddy) AddProxyLocation(domain string, url string, proxyTarget string) {
	c.caddyjson.AddProxyLocation(domain, url, proxyTarget)
}

func (c *portableCaddy) Stop() {
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

func NewPortableCaddy(version string) *portableCaddy {
	newCaddy :=
		&portableCaddy{
			binaryLocation: "./caddy-bin",
			caddyjson:      NewCaddyJSON([]string{":8443"})}
	newCaddy.caddyjson.Apps.Http.HTTPPort = 8008
	newCaddy.version = LoadCaddyVersion(&version)
	newCaddy.EnsureBinary()
	log.GetDefaultLogger().Info().Msg("portable caddy ready")
	return newCaddy
}
