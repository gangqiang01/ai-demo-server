//go:build windows
// +build windows

package utils

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"k8s.io/klog/v2"
)

var (
	winRing0, _ = syscall.LoadLibrary("WinRing0x64.dll")
	initOls, _  = syscall.GetProcAddress(winRing0, "InitializeOls")
	rdmsr, _    = syscall.GetProcAddress(winRing0, "Rdmsr")

	psapi           = syscall.NewLazyDLL("Psapi.dll")
	QueryWorkingSet = psapi.NewProc("QueryWorkingSet")
)

type ethernetInfo struct {
	GUID            string
	Index           uint32
	MACAddress      string
	Name            string
	NetEnabled      bool
	PhysicalAdapter bool
	PNPDeviceID     string
	Description     string
}

func GetCpuTemp() string {
	defer func() {
		if err := recover(); err != nil {
			klog.Errorf("Get CPU Temp failed with err: %v", err)
		}
	}()

	// klog.V(4).Infof("GetCpuTemp")
	// h, err := syscall.LoadLibrary("WinRing0x64.dll")
	// if err != nil {
	// 	klog.Errorf("LoadLibrary: %v", err)
	// 	return "N/A"
	// }
	// defer syscall.FreeLibrary(h)

	// init, err := syscall.GetProcAddress(h, "InitializeOls")
	// if err != nil {
	// 	klog.Errorf("InitializeOls: %v", err)
	// 	return "N/A"
	// }
	syscall.Syscall(uintptr(initOls), 0, 0, 0, 0)

	//cpu temp
	var cpuTemp string
	var a uint32
	var b uint32
	cpuTemp = "N/A"
	a = 0
	b = 0
	// rdmsr, err := syscall.GetProcAddress(h, "Rdmsr")
	// if err != nil {
	// 	klog.Errorf("Rdmsr: %v", err)
	// 	return "N/A"
	// }
	syscall.Syscall(uintptr(rdmsr), 3, 0x1A2, uintptr(unsafe.Pointer(&a)), uintptr(unsafe.Pointer(&b)))

	tjMax := fmt.Sprintf("%x", unsafe.Pointer((*uint32)(unsafe.Pointer(uintptr(a)))))
	tjMaxValue, _ := strconv.ParseInt(tjMax, 16, 64)
	tjMaxValue1 := (uint64(tjMaxValue) & 0x7F0000) >> 16

	a = 0
	b = 0

	syscall.Syscall(uintptr(rdmsr), 3, 0x19C, uintptr(unsafe.Pointer(&a)), uintptr(unsafe.Pointer(&b)))

	value := fmt.Sprintf("%x", unsafe.Pointer((*uint32)(unsafe.Pointer(uintptr(a)))))
	tempValue, _ := strconv.ParseInt(value, 16, 64)
	tempValue1 := (uint64(tempValue) & 0xFF0000) >> 16

	temp := tjMaxValue1 - tempValue1
	klog.Infof("cpuTemp: %v, tjMaxValue1: %v, tempValue1: %v", temp, tjMaxValue1, tempValue1)
	cpuTemp = strconv.FormatUint(temp, 10)
	return cpuTemp
}

func adapterAddresses() ([]*windows.IpAdapterAddresses, error) {
	var b []byte
	l := uint32(15000) // recommended initial size
	for {
		b = make([]byte, l)
		err := windows.GetAdaptersAddresses(syscall.AF_UNSPEC, windows.GAA_FLAG_INCLUDE_PREFIX, 0, (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])), &l)
		if err == nil {
			if l == 0 {
				return nil, nil
			}
			break
		}
		if err.(syscall.Errno) != syscall.ERROR_BUFFER_OVERFLOW {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
		if l <= uint32(len(b)) {
			return nil, os.NewSyscallError("getadaptersaddresses", err)
		}
	}
	var aas []*windows.IpAdapterAddresses
	for aa := (*windows.IpAdapterAddresses)(unsafe.Pointer(&b[0])); aa != nil; aa = aa.Next {
		aas = append(aas, aa)
	}
	return aas, nil
}

func GetNetInterfaceType(index int) int {
	aas, err := adapterAddresses()
	if err != nil {
		return -1
	}

	for _, aa := range aas {
		ifIndex := aa.IfIndex
		ifType := aa.IfType

		if index != int(ifIndex) {
			continue
		}

		switch ifType {
		case windows.IF_TYPE_ETHERNET_CSMACD: //ethernet
			return 1
		case windows.IF_TYPE_IEEE80211: //wifi
			return 2
		case windows.IF_TYPE_PPP: //ppp
			return 3
		default:
			return 0
		}
	}

	return 0
}

func IsPhysicalAdapter(name string) bool {
	controlKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Class`, registry.ALL_ACCESS)
	controlSubkeyNames, _ := controlKey.ReadSubKeyNames(0)
	for _, controlSubkeyName := range controlSubkeyNames {
		controlSubKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Class`+"\\"+controlSubkeyName, registry.ALL_ACCESS)
		controlClass, _, _ := controlSubKey.GetStringValue("Class")
		if "Net" == controlClass {
			netKeyNames, _ := controlSubKey.ReadSubKeyNames(0)
			for _, netKeyName := range netKeyNames {
				netSubKey, _ := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Class`+"\\"+controlSubkeyName+"\\"+netKeyName, registry.ALL_ACCESS)
				netName, _, _ := netSubKey.GetStringValue("DriverDesc")
				if netName == name {
					netCharacteristics, _, _ := netSubKey.GetIntegerValue("Characteristics")
					if 0x01 == netCharacteristics&0x01 {
						return false
					} else if 0x04 == netCharacteristics&0x04 {
						return true
					}
				}
			}
		}
	}
	return false
}

/*
* GetLocalIPs
* get all IP of all network interfaces.
 */
func GetLocalIPs() []string {
	idxs := make([]int, 0)
	ips := make([]string, 0)
	var ethernetInfos []ethernetInfo

	klog.V(4).Infof("GetLocalIPs")

	err := wmi.Query("SELECT * FROM Win32_NetworkAdapter WHERE (MACAddress IS NOT NULL) AND (NOT (PNPDeviceID LIKE 'ROOT%'))", &ethernetInfos)
	if err != nil {
		klog.Errorf("GetLocalIPs wmi get networkAdapter infos failed and error is %s", err.Error())
	}

	klog.V(4).Infof("ethernetInfos: %v", ethernetInfos)

	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		if len(ips) < 1 {
			netInterfaces, err := net.Interfaces()
			if err != nil {
				klog.Errorf("GetLocalIPs netInterfaces error is %s", err.Error())
				ips = append(ips, "127.0.0.1")
				return ips
			}

			for _, ni := range netInterfaces {
				if GetNetInterfaceType(ni.Index) <= 0 {
					continue
				}

				if ni.Flags&net.FlagUp == 0 {
					continue
				}

				addrs, err := ni.Addrs()
				if err != nil || len(addrs) <= 0 {
					continue
				}

				for _, addr := range addrs {
					ipNet, isThisType := addr.(*net.IPNet)
					if !isThisType || ipNet.IP.IsLoopback() {
						continue
					}

					if ipNet.IP.To4() == nil {
						continue
					}

					idxs = append(idxs, ni.Index)
				}
			}

			sort.Sort(sort.IntSlice(idxs))

			for _, id := range idxs {
				for _, ni := range netInterfaces {
					if ni.Index == id {
						addrs, err := ni.Addrs()
						if err != nil || len(addrs) <= 0 {
							continue
						}

						for _, addr := range addrs {
							ipNet, isThisType := addr.(*net.IPNet)
							if !isThisType || ipNet.IP.IsLoopback() {
								continue
							}

							if ipNet.IP.To4() == nil {
								continue
							}

							macAddr := ni.HardwareAddr.String()
							for ix := range ethernetInfos {
								if true == ethernetInfos[ix].PhysicalAdapter && macAddr == strings.ToLower(ethernetInfos[ix].MACAddress) && true == ethernetInfos[ix].NetEnabled {
									if IsPhysicalAdapter(ethernetInfos[ix].Description) {
										ips = append(ips, ipNet.IP.String())
										break
									}
								}
							}
						}
					}
				}
			}
		} else {
			break
		}
	}

	klog.V(4).Infof("ips: %v", ips)
	return ips
}

func GetOSVersion() string {
	var osVersion string
	hostInfo, _ := host.Info()
	klog.V(4).Infof("os info: %v", hostInfo)
	if strings.Contains(hostInfo.Platform, "Microsoft") {
		osVersion = strings.ReplaceAll(hostInfo.Platform, "Microsoft", "")
	}
	osVersion = strings.TrimLeft(osVersion, " ")
	return osVersion
}

func GetAllPartitions() []string {
	paths := make([]string, 0)

	infos, _ := disk.Partitions(true)

	for _, info := range infos {
		paths = append(paths, info.Mountpoint)
	}

	return paths
}
