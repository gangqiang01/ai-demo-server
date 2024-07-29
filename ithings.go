package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/edgehook/ithings/cmd"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	os.Chdir(dir)
	execPath, _ := os.Getwd()

	var logFileName string
	if !strings.Contains(utils.GetOsType(), "windows") {
		log := strings.Split(os.Args[0], "/")
		logFileName = log[1] + ".log"
	} else {
		log := strings.Split(os.Args[0], ".")
		logFileName = log[0] + ".log"
	}

	// if os.Getenv("LOG_PATH") != "" {
	// 	logFile = os.Getenv("LOG_PATH")
	// }
	logFileName = "AppHub-Agent.log"
	flag.Set("log_file", logFileName)
	flag.Set("log_file_max_size", "5") //in MB, default as 1800MB
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")

	args := flag.Args()
	klog.Infof("flag args: %v", args)
	// cmd.ConfigureAndRun()
	klog.Infof("This is AppHub-Agent log")
	// klog.Infof("DeviceOn-iService version is: %s", constants.M2mClientVersion)
	klog.Infof("AppHub-Agent exec path is %s", execPath)
	klog.Infof("AppHub-Agent dir is %s", dir)
	klog.Infof("AppHub-Agent path is %s", os.Args[0])

	logs.InitLogs()
	defer logs.FlushLogs()

	if !strings.Contains(utils.GetOsType(), "windows") {
		if err := cmd.Execute(); err != nil {
			os.Exit(1)
		}
	} else {
		cmd.NewAppService()
	}
}
