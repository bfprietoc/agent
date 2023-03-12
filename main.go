package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Process struct {
	Pid  string `json:"pid"`
	Name string `json:"name"`
}
type Data struct {
	CpuInfo     string `json:"cpu_info"`
	ProcessInfo string `json:"process_info,omitempty"`
	UsersInfo   string `json:"users_info"`
	OsInfo      string `json:"os_info"`
}

type MacOsInformation struct {
	SPSoftwareDataType []SPSoftwareDataType
}

type SPSoftwareDataType struct {
	OsVersion     string `json:"os_version"`
	KernelVersion string `json:"kernel_version"`
}

func CpuMacInfo() string {
	cpuCmd, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(cpuCmd)
}

func CpuLinuxInfo() string {
	cpuCmd, err := exec.Command("bash", "-c", "lscpu | egrep 'Model name'").Output()
	if err != nil {
		panic(err)
	}
	return string(cpuCmd)
}

func ProccesInfo() string {
	psCmd, err := exec.Command("ps", "-e").Output()
	if err != nil {
		panic(err)
	}

	psLines := strings.Split(string(psCmd), "\n")

	arrProcess := []Process{}

	for i, psLine := range psLines {
		if i == 0 || len(psLine) == 0 {
			continue
		}
		fields := strings.Fields(psLine)
		process := Process{Pid: fields[0], Name: fields[3]}
		arrProcess = append(arrProcess, process)
	}

	processArray, err := json.Marshal(arrProcess)
	if err != nil {
		log.Fatal(err)
	}
	return string(processArray)
}

func UsersInfo() string {
	usersCmd, err := exec.Command("who").Output()
	if err != nil {
		log.Fatal(err)
	}
	usersSplited := strings.Split(string(usersCmd), "\n")
	usersArray, err := json.Marshal(usersSplited)
	if err != nil {
		log.Fatal(err)
	}
	return string(usersArray)
}

func OsMacInfo() string {
	var macInfo MacOsInformation
	osCmd, err := exec.Command("system_profiler", "SPSoftwareDataType", "-json").Output()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(osCmd, &macInfo)
	if err != nil {
		log.Fatal(err)
	}
	macVersion := macInfo.SPSoftwareDataType[0].KernelVersion
	macName := macInfo.SPSoftwareDataType[0].OsVersion
	macOsInfo := macName + macVersion

	return macOsInfo
}

func OsLinuxInfo() string {
	osCmd, err := exec.Command("lsb_release", "-d").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(osCmd)
}

func main() {
	cpuInfo := CpuLinuxInfo()
	processInfo := ProccesInfo()
	usersInfo := UsersInfo()
	osInfo := OsLinuxInfo()

	content := "Content-Type: application/json"
	url := "http://ec2-18-188-229-228.us-east-2.compute.amazonaws.com:8080/data"

	body := Data{CpuInfo: cpuInfo, ProcessInfo: processInfo, UsersInfo: usersInfo, OsInfo: osInfo}

	info, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	_, err = exec.Command("curl", "-X", "POST", "-H", content, "-d", string(info), url, "-v").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Succesfull")

}
