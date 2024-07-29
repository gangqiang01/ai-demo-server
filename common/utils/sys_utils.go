package utils

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"bytes"
	"crypto/aes"
	"encoding/base64"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/exp/rand"

	n "github.com/shirou/gopsutil/v3/net"
	"github.com/yumaojun03/dmidecode"
	"k8s.io/klog/v2"
)

var (
	SYS_1GB = uint64(1024 * 1024 * 1024)
)

type Port struct {
	PortName string `json:"port"`
}

var InitTime int64

type IpInformation struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

var TimeZone IpInformation

// process monitor
type ThreadInfos struct {
	Name       string  `json:"cmd"`
	Pid        int64   `json:"pid"`
	CpuLoading float64 `json:"cpuloading"`
	MemLoading float64 `json:"memloading"`
	Status     string  `json:"status"`
	UerName    string  `json:"username"`
}

type MonitorInfo struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"ts"`
}

// TimeZone = IpInformation{"N", "A", 0, 0}

/*
* System reboot.
 */
func SysReboot() error {
	if strings.Contains(GetOsType(), "windows") {
		_, err := Execute1("cmd", "/C", "shutdown", "/r", "/t", "0")
		return err
	}

	//default for linux.
	_, err := Execute("reboot")
	return err
}

/*
* System shutdown.
 */
func SysShutdown(command string) error {
	if strings.Contains(GetOsType(), "windows") {
		_, err := Execute1("cmd", "/C", "shutdown", "/s", "/t", "0")
		return err
	}

	if strings.Contains(command, "halt") {
		_, err := Execute("halt")
		return err
	}

	_, err := Execute("shutdown -h now")
	if err == nil {
		return nil
	}

	_, err = Execute("poweroff")
	if err == nil {
		return nil
	}

	return err
}

/*
* GetLocalMACs
* get the local host's macaddress.
 */
func GetLocalMACs() []string {

	netInterfaces, err := net.Interfaces()
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil
	}

	macAddrs := make([]string, 0)

	for _, ni := range netInterfaces {
		if !strings.HasPrefix(ni.Name, "e") &&
			!strings.HasPrefix(ni.Name, "w") &&
			!strings.HasPrefix(ni.Name, "p") {
			continue
		}
		macAddr := ni.HardwareAddr.String()
		if len(macAddr) > 6 {
			macAddrs = append(macAddrs, strings.ToUpper(macAddr))
		}
	}
	return macAddrs
}

/*
* GetBoardName
* get board name, currently, we use hostname.
 */

func IsX86() bool {
	// out, err := exec.Command("uname", "-m").Output()
	// if err != nil {
	// 	klog.Errorf("err: %v", err)
	// }
	// arch := string(out[:len(out)-1])

	// if arch == "x86_64" {
	// 	return true
	// }

	// //arm32
	// if "arm" == arch || "armv7l" == arch {
	// 	return false
	// }

	// //arm64
	// if "aarch64_be" == arch || "aarch64" == arch || "armv8b" == arch ||
	// 	"armv8l" == arch {
	// 	return false
	// }
	// return false

	//test
	return true
}

func GetBoardName() string {
	boardName := "unknown board"

	if IsX86() {
		dmi, err := dmidecode.New()
		if err != nil {
			klog.Errorf("err: %v", err)
			return boardName
		}

		infos, err := dmi.System()
		if err != nil {
			klog.Errorf("err: %v", err)
			return boardName
		}
		boardName = infos[0].ProductName
		return boardName
	}

	//for ec yocto linux.
	name, err := ioutil.ReadFile("/proc/board")
	if err == nil {
		boardName = strings.TrimSpace(string(name))
		return strings.Replace(boardName, " ", "", -1)
	}

	//we use hostname as board name by default.
	boardName, _ = os.Hostname()
	boardName = strings.TrimSpace(boardName)
	return strings.Replace(boardName, " ", "", -1)
}

/*
* GetHostName
* get host name.
 */
func GetHostName() string {
	hostName := "unknown host"

	hostName, err := os.Hostname()
	if err != nil {
		klog.Errorf("get host name with err: %v", err)
	}
	return hostName
}

// func GetKernelRelease() string {
// 	var u syscall.Utsname
// 	var buf [512]byte

// 	err := syscall.Uname(&u)
// 	if err != nil {
// 		klog.Errorf("err: %v", err)
// 		return ""
// 	}

// 	for i, b := range u.Release[:] {
// 		buf[i] = uint8(b)
// 		if b == 0 && i > 0 {
// 			return string(buf[:i])
// 		}
// 	}

// 	return ""
// }

// func GetKernelVersion() string {
// 	var u syscall.Utsname
// 	var buf [512]byte

// 	err := syscall.Uname(&u)
// 	if err != nil {
// 		klog.Errorf("err: %v", err)
// 		return ""
// 	}

// 	for i, b := range u.Version[:] {
// 		buf[i] = uint8(b)
// 		if b == 0 && i > 0 {
// 			return string(buf[:i])
// 		}
// 	}

// 	return ""
// }

// func GetKernelMachine() string {
// 	var u syscall.Utsname
// 	var buf [512]byte

// 	err := syscall.Uname(&u)
// 	if err != nil {
// 		klog.Errorf("err: %v", err)
// 		return ""
// 	}

// 	for i, b := range u.Machine[:] {
// 		buf[i] = uint8(b)
// 		if b == 0 && i > 0 {
// 			return string(buf[:i])
// 		}
// 	}

// 	return ""
// }

/*
* GetOsType:
* return lowcase ostype.
 */
func GetOsType() string {
	return runtime.GOOS
}

func RandomInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	rand.Seed(uint64(time.Now().Unix()))
	return rand.Intn(max-min) + min
}

func GetCpuPercent() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0.0
	}

	pt, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", percent[0]), 64)
	return pt
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	mp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", memInfo.UsedPercent), 64)
	return mp
}
func GetMemUsed() uint64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.Used
}
func GetTotalMemory() uint64 {
	memInfo, _ := mem.VirtualMemory()

	return memInfo.Total
}

func DetectPort(value string) bool {
	port := &Port{}
	err := json.Unmarshal([]byte(value), port)
	if err != nil {
		klog.Errorf("err: %v", err)
		return true
	}
	hostName := "127.0.0.1"
	protocol := "tcp"
	addr := net.JoinHostPort(hostName, port.PortName)
	conn, err := net.DialTimeout(protocol, addr, 3*time.Second)
	if err != nil {
		klog.Errorf("detect port with err: %v", err)
		return false
	}
	defer conn.Close()
	return true
}

func GetMemAvailable() uint64 {
	memInfo, _ := mem.VirtualMemory()

	return memInfo.Available
}

func GetDiskTotal() uint64 {
	var total uint64 = 0
	paths := GetAllPartitions()

	for _, path := range paths {
		state, err := disk.Usage(path)
		if err != nil {
			klog.Errorf("err: %v", err)
			continue
		}

		total += state.Total
	}

	return total
}

func GetDiskUsed() uint64 {
	var used uint64 = 0
	paths := GetAllPartitions()

	for _, path := range paths {
		state, err := disk.Usage(path)
		if err != nil {
			klog.Errorf("err: %v", err)
			continue
		}

		used += state.Used
	}

	return used
}

func GetDiskUsedPercent() uint64 {
	total := GetDiskTotal()
	used := GetDiskUsed()

	if 0 == total {
		return 0
	} else {
		return used * 100 / total
	}
}

// get path free space by bytes.
func GetPathFreeSpace(path string) uint64 {
	state, err := disk.Usage(path)
	if err != nil {
		klog.Errorf("get %v free space with err: %v", path, err)
		return 0
	}

	return state.Free
}

func Convert2GB(bytes float64) string {
	base := float64(1024.0)

	gBytes := bytes / base / base / base

	return fmt.Sprintf("%.2f", gBytes)
}

func GetProcessInfo() []*process.Process {
	processInfos, _ := process.Processes()
	return processInfos
}

func GetProcessByName(processName, user string) *ThreadInfos {
	infos := GetProcessInfo()
	var threadInfo *ThreadInfos
	for _, info := range infos {
		name, _ := info.Name()
		username, _ := info.Username()
		if processName == name && user == username {
			// status R: Running S: Sleep T: Stop I: Idle
			// Z: Zombie W: Wait L: Lock
			status, _ := info.Status()
			cpuLoading, _ := info.CPUPercent()
			cpuPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuLoading), 64)
			memLoading, _ := info.MemoryPercent()
			memPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", memLoading), 64)
			username, _ := info.Username()
			if threadInfo == nil {
				threadInfo = &ThreadInfos{
					Name:       processName,
					Pid:        int64(info.Pid),
					Status:     status,
					UerName:    username,
					CpuLoading: cpuPercent,
					MemLoading: memPercent,
				}
			} else {
				threadInfo.CpuLoading += cpuPercent
				threadInfo.MemLoading += memPercent
				threadInfo.CpuLoading, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", threadInfo.CpuLoading), 64)
				threadInfo.MemLoading, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", threadInfo.MemLoading), 64)
			}
		}
	}

	return threadInfo
}

func KillThreadByPid(threadpid string) bool {
	processInfos := GetProcessInfo()
	pid, _ := strconv.Atoi(threadpid)
	for _, processInfo := range processInfos {
		if int32(pid) == processInfo.Pid {
			err := processInfo.Kill()
			if err != nil {
				klog.Errorf("kill thread with err %v", err)
				return false
			}
			break
		}
	}

	return true
}

func KillThreadByName(threadname string) bool {
	processInfos := GetProcessInfo()
	for _, processInfo := range processInfos {
		name, _ := processInfo.Name()
		if name == threadname {
			err := processInfo.Kill()
			if err != nil {
				klog.Errorf("kill thread with err %v", err)
				return false
			}
			break
		}
	}

	return true
}

func FileCopy(src, dst string) error {
	isWindows := strings.Contains(GetOsType(), "windows")
	if isWindows {
		fileName := filepath.Base(src)
		targetDir := filepath.Dir(dst)

		cmd := exec.Command("xcopy", src, targetDir, "/I", "/E", "/F", "/O", "/Y")
		_, err := cmd.Output()
		if err != nil {
			klog.Errorf("err: %v", err)
			return err
		}

		return os.Rename(filepath.Join(targetDir, fileName), dst)
	}

	// default for linux.
	cmd := exec.Command("cp", "-r", src, dst)
	_, err := cmd.Output()
	return err
}

/*
* Execute:
* Execute the command without paramenta.
 */
func Execute(command string) (string, error) {
	if strings.Contains(GetOsType(), "windows") {
		cmd := exec.Command("cmd", "/C", command)
		cmd.Env = os.Environ()
		output, err := cmd.CombinedOutput()
		if err != nil {
			klog.Errorf("output: %s, err: %v", string(output), err)
			return string(output), err
		}

		return string(output), nil
	}

	// default for unix.
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Env = os.Environ()
	// cmd.Dir = "/usr/local"
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Setpgid: true,
	// }

	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("output: %s, err: %v", string(output), err)
		return string(output), err
	}

	return string(output), nil
}

/*
* Execute1:
* Execute1 the command with paramenta.
 */
func Execute1(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("output: %s, err: %v", string(output), err)
		return string(output), err
	}

	return string(output), nil
}

/*
* Execute2:
* Execute the command realtime.
 */
func Execute2(command string, callback func(io.ReadCloser, io.ReadCloser)) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}

	if callback != nil {
		callback(stdout, stderr)
	}

	return cmd.Wait()
}

func Execute3(command, path string) (string, error) {
	if strings.Contains(GetOsType(), "windows") {
		cmd := exec.Command("cmd", "/C", command)
		cmd.Env = os.Environ()
		cmd.Dir = path
		output, err := cmd.CombinedOutput()
		klog.V(4).Infof("output: %s", string(output))
		if err != nil {
			klog.Errorf("err: %v", err)
			return string(output), err
		}

		return string(output), nil
	}

	// default for unix.
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	klog.Infof("output: %s", string(output))
	if err != nil {
		klog.Errorf("err: %v", err)
		return string(output), err
	}

	return string(output), nil
}

func Execute4(command string) error {

	// default for unix.
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	klog.Infof("exec4 end")

	return nil
}

func Execute5(command, workDir string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("err: %v", err)
		return string(output), err
	}

	return string(output), nil
}

func ExecuteWithTimeOut(cmdStr string, workdir string, timeout time.Duration) (string, error) {
	var stdout, stderr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Dir = workdir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	errMsg := ""
	if len(stderr.String()) != 0 {
		errMsg = fmt.Sprintf("stderr:\n %s", stderr.String())
	}

	if len(stdout.String()) != 0 {
		if len(errMsg) != 0 {
			errMsg = fmt.Sprintf("%s \n\n; stdout:\n %s", errMsg, stdout.String())
		} else {
			errMsg = fmt.Sprintf("stdout:\n %s", stdout.String())
		}
	}

	return errMsg, err
}

/*
* XmlParser
* parse the xml file.
 */
func XmlParser(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}

	return xml.Unmarshal(data, v)
}

// read file from filepath
func GetFileContent(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return content, nil
}

func WriteFileContent(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0755)
}

/*
* JsonParser
* parse the xml file.
 */
func JsonParser(fileName string, v interface{}) error {
	data, err := GetFileContent(fileName)
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}

	return json.Unmarshal(data, v)
}

func AesECBEncrypt(serverTimestamp, srpName string) string {
	var key = []byte("androidlink/apphub%#$abc")
	nowTime, _ := strconv.ParseUint(serverTimestamp, 10, 64)
	nowTime = nowTime / 1000

	origString := strconv.FormatUint(nowTime, 10) + "#" + srpName
	origData := []byte(origString)

	block, err := aes.NewCipher(key)
	if err != nil {
		//klog.Errorf("AESEncrypt err: %v", err)
		return err.Error()
	}

	blockSize := block.BlockSize()
	padding := blockSize - len(origData)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext := append(origData, padtext...)

	if len(plaintext)%blockSize != 0 {
		klog.Errorf("Need a multiple of the blocksize")
	}
	ciphertext := make([]byte, len(plaintext))

	size := 16
	for bs, be := 0, size; bs < len(plaintext); bs, be = bs+size, be+size {
		//cipher.Decrypt(decrypted[bs:be], data[bs:be])
		block.Encrypt(ciphertext[bs:be], plaintext[bs:be])
	}

	pass64 := base64.StdEncoding.EncodeToString(ciphertext)

	return string(pass64)
}

func GetServerConfigUrl(ipAddr, port, path, srpName string) string {
	configUrl := "http://" + ipAddr + ":" + port + path
	return configUrl
}

func CopyLogFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		klog.Errorf("err: %v", err)
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		klog.Errorf("err: %v", err)
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func TimeStart() {
	InitTime = time.Now().UnixNano() / 1e6
}

func GetTime() int64 {
	return InitTime
}

func GetDeviceLocation() IpInformation {
	return TimeZone
}

func DeviceLocation() {
	resp, _ := http.Get("https://ip.gs/json")
	obj := &IpInformation{}
	if resp == nil {
		for i := 1; i <= 3; i++ {
			if resp == nil {
				resp, _ = http.Get("https://ip.gs/json")
				fmt.Printf("this is times :%d", i)
			} else {
				break
			}
		}
	}
	if resp == nil {
		obj = &IpInformation{Country: "N", City: "A", Longitude: 0, Latitude: 0}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal([]byte(body), obj)
		resp.Body.Close()
	}
	TimeZone = *obj
}

func GetSystemAllUsers() []string {
	output, err := Execute("cat /etc/passwd | grep bash | cut -f 1 -d:")
	if err != nil {
		klog.Errorf("err: %v", err)
	}
	names := strings.Split(output, "\n")
	// names := strings.TrimRight(strings.Split(output, "\n"), " ")
	var nameSlice []string = names[0 : len(names)-1]
	klog.V(4).Infof("names: %v", names)
	klog.V(4).Infof("nameSlice: %v", nameSlice)
	return nameSlice
}

func GetSystemdService(name string) bool {
	command := "systemctl list-units --all --type=service | grep " + name
	klog.V(4).Infof("command: %v", command)
	output, err := Execute(command)
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	klog.V(4).Infof("output: %v", output)
	if output != "" {
		return true
	} else {
		return false
	}

}

func saveAsjpg(img *image.RGBA, imageName string) error {
	filePath := imageName + ".jpg"
	file, err := os.Create(filePath)
	if err != nil {
		klog.Errorf("create file %v with err: %v", filePath, err)
		return err
	}
	defer file.Close()

	var opt jpeg.Options
	opt.Quality = 50

	err = jpeg.Encode(file, img, &opt) //jpeg里面三个参数
	if err != nil {
		klog.Errorf("jpg encode with err: %v", err)
	}

	pngFileInfo, err := os.Stat(filePath)
	klog.V(4).Infof("jpeg file size: %v", pngFileInfo.Size())

	return nil
}

func GetNetReceive() float64 {
	netStat, _ := n.IOCounters(false)
	time.Sleep(1 * time.Second)
	netStat1, _ := n.IOCounters(false)
	receive := 0.0
	for _, net1 := range netStat1 {
		for _, net := range netStat {
			if net1.Name == net.Name {
				if net1.BytesRecv != 0 && net.BytesRecv != 0 && net1.BytesRecv > net.BytesRecv {
					receive = float64(net1.BytesRecv-net.BytesRecv) / 1024
				}
				break
			}
		}
	}
	return receive
}

func GetNetSend() float64 {
	netStat, _ := n.IOCounters(false)
	time.Sleep(1 * time.Second)

	netStat1, _ := n.IOCounters(false)
	send := 0.0
	for _, net1 := range netStat1 {
		for _, net := range netStat {
			if net1.Name == net.Name {
				if net1.BytesSent != 0 && net.BytesSent != 0 && net1.BytesSent > net.BytesSent {
					send = float64(net1.BytesSent-net.BytesSent) / 1024
				}
				break
			}
		}
	}
	return send
}

type MonitorNetwork struct {
	Name string  `json:"name"`
	Up   float64 `gorm:"type:float" json:"up"`
	Down float64 `gorm:"type:float" json:"down"`
}

func LoadNetIO(isAll bool, interval time.Duration) []MonitorNetwork {
	netStat, _ := n.IOCounters(isAll)
	time.Sleep(interval)
	netStat2, _ := n.IOCounters(isAll)
	var netList []MonitorNetwork
	for _, net2 := range netStat2 {
		for _, net1 := range netStat {
			if net2.Name == net1.Name {
				var itemNet MonitorNetwork
				itemNet.Name = net1.Name

				if net2.BytesSent != 0 && net1.BytesSent != 0 && net2.BytesSent > net1.BytesSent {
					itemNet.Up = float64(net2.BytesSent-net1.BytesSent) / 1024 / interval.Seconds()
				}
				if net2.BytesRecv != 0 && net1.BytesRecv != 0 && net2.BytesRecv > net1.BytesRecv {
					itemNet.Down = float64(net2.BytesRecv-net1.BytesRecv) / 1024 / interval.Seconds()
				}
				klog.Infof("netIo name: %v, UP: %v, Down: %v", itemNet.Name, itemNet.Up, itemNet.Down)
				netList = append(netList, itemNet)
				break
			}
		}
	}

	return netList
}

type Slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func FindBiosVersion() string {
	// 	file, err := os.OpenFile("/dev/mem", os.O_RDONLY, 0600)
	// 	if err != nil {
	// 		klog.Errorf("err: %v", err)
	// 		return "N/A"
	// 	}

	// 	defer file.Close()

	// 	b, err := syscall.Mmap(int(file.Fd()), 0xF0000, 0xFFFF, syscall.PROT_READ, syscall.MAP_SHARED)
	// 	if err != nil {
	// 		klog.Errorf("err: %v", err)
	// 		return "N/A"
	// 	}
	// 	defer syscall.Munmap(b)
	// 	hdr := (*Slice)(unsafe.Pointer(&b))

	// 	pT := uintptr(hdr.Data)

	// 	var mt []byte
	// 	mt1 := 0
	// 	mt2 := 0

	// 	for i := 0; i < 0xffff-2; i++ {

	// 		t1 := (*byte)(unsafe.Pointer(pT))
	// 		t2 := (*byte)(unsafe.Pointer(pT + 1))
	// 		t3 := (*byte)(unsafe.Pointer(pT + 2))

	// 		if *t1 == '*' && *t2 == '*' && *t3 == '*' {
	// 			for j := 0; j < 150; j++ {
	// 				mt = append(mt, *(*byte)(unsafe.Pointer(pT)))
	// 				pT = pT + 1
	// 			}

	// 			for t := 149; t >= 0; t-- {
	// 				if mt[t] == '*' {
	// 					mt1 = t
	// 					break
	// 				}
	// 			}
	// 			for t := 149; t >= 0; t-- {
	// 				if mt[t] == 'V' {
	// 					mt2 = t
	// 				}
	// 			}
	// 			return string(mt[mt2 : mt1-4])
	// 		}
	// 		pT = pT + 1
	// 	}
	return "N/A"
}

type threadMemLoadingInfos []*ThreadInfos
type threadCpuLoadingInfos []*ThreadInfos

func (M threadMemLoadingInfos) Len() int {
	return len(M)
}

func (M threadMemLoadingInfos) Less(i, j int) bool {
	return M[i].MemLoading > M[j].MemLoading
}

func (M threadMemLoadingInfos) Swap(i, j int) {
	M[i], M[j] = M[j], M[i]
}

func GetMemLoadingTop5() []*ThreadInfos {

	processInfos := GetProcessInfo()
	var threadInfos []*ThreadInfos

	for _, processInfo := range processInfos {
		var threadInfo ThreadInfos

		threadPid := int64(processInfo.Pid)
		threadMemLoading, _ := processInfo.MemoryPercent()
		threadInfo.Pid = threadPid
		threadInfo.Name, _ = processInfo.Name()
		memLoading, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", threadMemLoading), 64)
		threadInfo.MemLoading = memLoading
		threadInfos = append(threadInfos, &threadInfo)
	}

	sort.Sort(threadMemLoadingInfos(threadInfos))
	var threadInfoSlice []*ThreadInfos = threadInfos[0:5]
	return threadInfoSlice
}
func (C threadCpuLoadingInfos) Len() int {
	return len(C)
}

func (C threadCpuLoadingInfos) Less(i, j int) bool {
	return C[i].CpuLoading > C[j].CpuLoading
}

func (C threadCpuLoadingInfos) Swap(i, j int) {
	C[i], C[j] = C[j], C[i]
}
func GetCpuLoadingTop5() []*ThreadInfos {

	processInfos := GetProcessInfo()
	var threadInfos []*ThreadInfos

	for _, processInfo := range processInfos {
		var threadInfo ThreadInfos
		threadPid := int64(processInfo.Pid)
		threadCpuLoading, _ := processInfo.CPUPercent()
		threadInfo.Pid = threadPid
		threadInfo.Name, _ = processInfo.Name()
		cpuLoading, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", threadCpuLoading), 64)
		threadInfo.CpuLoading = cpuLoading
		threadInfos = append(threadInfos, &threadInfo)
	}

	sort.Sort(threadCpuLoadingInfos(threadInfos))
	var threadInfoSlice []*ThreadInfos = threadInfos[0:5]
	return threadInfoSlice
}
