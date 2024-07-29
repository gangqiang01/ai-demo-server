//go:build linux
// +build linux

package utils

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/disk"
	"k8s.io/klog/v2"
)

func GetCpuTemp() string {
	hw := GetHvMonitor()
	if hw.TempCPU == 0 {
		tempPath := "/sys/class/thermal/thermal_zone0/temp"
		if _, err1 := os.Stat(tempPath); err1 != nil {
			tempPath = "/sys/class/hwmon/hwmon1/device/temp1_input"
			if _, err2 := os.Stat(tempPath); err2 != nil {
				klog.Errorf("err: %v", err1)
				klog.Errorf("err: %v", err2)
				return "N/A"
			}
		}
		temp, err := os.ReadFile(tempPath)
		if err != nil {
			klog.Errorf("err: %v", err)
			return "N/A"
		}
		t, err := strconv.Atoi(strings.TrimSpace(string(temp)))
		if err != nil {
			klog.Errorf("err: %v", err)
			return "N/A"
		}

		return fmt.Sprintf("%d", t/1000)
	} else {
		return fmt.Sprintf("%d", int64(hw.TempCPU))
	}
}

/*
* GetLocalIPs
* get all IP of all network interfaces.
 */
func GetLocalIPs() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil
	}

	ips := make([]string, 0)

	for _, addr := range addrs {
		ipNet, isThisType := addr.(*net.IPNet)
		if isThisType && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}

	return ips
}

func GetOSVersion() string {

	cmd := exec.Command("cat", "/etc/issue")

	cmdout, _ := cmd.StdoutPipe()

	cmd.Start()

	bytes, _ := io.ReadAll(cmdout)
	osversion := strings.Split(string(bytes), "\\n")

	cmd.Wait()

	return string(osversion[0])

}

func GetAllPartitions() []string {
	paths := make([]string, 0)
	parts := make([]disk.PartitionStat, 0)

	infos, _ := disk.Partitions(true)

	findDev := func(name string) string {
		for _, pi := range parts {
			if pi.Device == name {
				return name
			}
		}

		return ""
	}

	//filter all duplicate device
	for _, info := range infos {
		if findDev(info.Device) == "" {
			if strings.Contains(info.Device, "/dev/loop") {
				backFile, _ := Execute1("losetup", "-nO", "BACK-FILE", info.Device)
				backFile = strings.TrimSpace(backFile)
				if findDev(backFile) == "" {
					info.Device = backFile
					parts = append(parts, info)
					continue
				}
			}

			parts = append(parts, info)
		}
	}

	for _, info := range parts {
		if !strings.Contains(info.Device, "/dev/") {
			continue
		}

		paths = append(paths, info.Mountpoint)
	}

	return paths
}
