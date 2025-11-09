package utils

import (
	"bytes"
	"io"
	"os"
	"reflector/log"
	"strconv"
	"strings"
)

func DropEmpty(s []string) []string {
	var d []string
	for _, str := range s {
		if str != "" {
			d = append(d, str)
		}
	}
	return d
}

func dismissError(a any, _ error) any {
	return a
}

type process struct {
	PID     int
	Cmdline []string
}

type processStringCmdline struct {
	PID     int
	Cmdline string
}

func nullTermBytesToStrings(b []byte) []string {
	x := bytes.Split(b, []byte{0})
	l := make([]string, len(x))
	for i, b := range x {
		l[i] = string(b)
	}
	return DropEmpty(l)
}

func ps(onlyCommands bool) *[]process {
	proc_dir := "/proc/"
	processes := new([]process)
	entries, err := os.ReadDir(proc_dir)
	if err != nil {
		log.GetDefaultLogger().Fatal().Update("error", err.Error()).Done()
	}
	for _, e := range entries {
		pname := e.Name()

		// filter only numeric (PIDs)
		if _, err := strconv.Atoi(pname); err != nil {
			continue
		}
		cmdlineFile, err := os.Open(proc_dir + pname + "/cmdline")
		if err != nil {
			log.GetDefaultLogger().Fatal().Update("error", err.Error()).Done()
		}
		cmdlineBytes := bytes.NewBuffer(nil)
		io.Copy(cmdlineBytes, cmdlineFile)
		p := process{
			PID:     dismissError(strconv.Atoi(pname)).(int),
			Cmdline: nullTermBytesToStrings(cmdlineBytes.Bytes()),
		}
		if len(p.Cmdline) == 0 { // empty cmdline, this likely is a kworker
			continue
		}
		if onlyCommands {
			p.Cmdline = p.Cmdline[:1]
		}
		*processes = append(*processes, p)
	}
	return processes
}

func PS() *[]process {
	return ps(false)
}

func PSComms() *[]process {
	return ps(true)
}

func PSStrings() *[]processStringCmdline {
	r := []processStringCmdline{}
	for _, procStrings := range *ps(false) {
		r = append(r, processStringCmdline{PID: procStrings.PID, Cmdline: strings.Join(procStrings.Cmdline, " ")})
	}
	return &r
}
