package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

const versionFileDir = "/etc/pjversion/"
const versionFilePath = "/etc/pjversion/%s.version"

const (
	versionUnknown = "unknown"
	versionError   = "errorInstalled"
)

var commonConfig struct {
	Version string
}

type versionStu struct {
	ServiceName    string `json:"servicename"`
	PackageVersion string `json:"packageversion"`
	ConfigVersion  string `json:"configversion"`
}

func LoadVersionV2(servicename string, configpath string, buildversion string) (err error) {
	if _, err = toml.DecodeFile(configpath, &commonConfig); err != nil {
		return
	}

	if commonConfig.Version == "" {
		commonConfig.Version = versionUnknown
	}

	var versionData versionStu
	versionData.ServiceName = servicename
	versionData.PackageVersion = buildversion
	if versionData.PackageVersion == "" {
		versionData.PackageVersion = getPackageVersion(servicename)
	}
	versionData.ConfigVersion = commonConfig.Version

	jsbyte, err := json.Marshal(versionData)
	if err != nil {
		return
	}

	return saveToFile(fmt.Sprintf(versionFilePath, servicename), jsbyte)
}

func getPackageVersion(servicename string) (version string) {
	out, err := execCMD("/bin/sh", "-c", "dpkg -s "+servicename)
	if err != nil {
		version = versionUnknown
		return
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Status:") && !strings.Contains(line, "ok") {
			version = versionError
			return
		} else if strings.Contains(line, "Version:") {
			versionSlice := strings.Fields(line)
			if len(versionSlice) < 2 {
				version = versionUnknown
				return
			}
			version = versionSlice[1]
			return
		}
	}

	version = versionUnknown
	return
}

func execCMD(pendCMD string, keys ...string) (bytes.Buffer, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	var cmd = exec.Command(pendCMD, keys...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return bytes.Buffer{}, err
	}

	return out, nil
}

func saveToFile(filepath string, data []byte) (err error) {
	_, err = os.Stat(filepath)
	if !os.IsNotExist(err) {
		os.Remove(filepath)
	}

	out, err := os.Create(filepath)
	if err != nil {
		err = os.MkdirAll(path.Clean(versionFileDir), 0666)
		if err != nil {
			return
		}
		out, err = os.Create(filepath)
		if err != nil {
			return
		}
	}
	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(data))
	return
}
