package v1

import (
	"encoding/json"

	"github.com/edgehook/ithings/common/dbm/model"
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

const (
	//log handle status
	//0: processing , 1: solve, 2: solved, 3: recover
	AlertLogUnsolved int32 = 0
	AlertLogSolving  int32 = 1
	AlertLogResolved int32 = 2
	AlertLogInvalid  int32 = 3

	//log type: event, monitor
	AlertLogTypeEvent   string = "event"
	AlertLogTypeMonitor string = "monitor"
	//monitor type
	MonitorProcessType   = "process"
	MonitorCpuUsageType  = "cpuUsage"
	MonitorMemUsageType  = "memUsage"
	MonitorDiskUsageType = "diskUsage"
	MonitorCpuTempType   = "cpuTemp"
	MonitorAiType        = "ai"
)

func CheckValueByCondition(ctype string, val interface{}, monitorInfo *model.Monitor) bool {
	// klog.Infof("ctype: %s, status: %s", ctype, monitorInfo.Status)
	if monitorInfo.Status != StatusEnable {
		return false
	}
	value, err := utils.ToFloat64(val)
	if err != nil {
		klog.Errorf("Value convert error: %v", err.Error())
		return false
	}
	svalue, err := utils.ToFloat64(monitorInfo.Value)
	if err != nil {
		klog.Errorf("monitor value convert error: %v", err.Error())
		return false
	}
	klog.Infof("value: %v, svalue: %v, condition: %v", value, svalue, monitorInfo.Condition)
	switch monitorInfo.Condition {
	case ">":
		if value > svalue {
			return true
		}
	case "<":
		if value < svalue {
			return true
		}
	case "=":
		if value == svalue {
			return true
		}
	}
	return false

}

func CheckIsStopByCondition(user, name string) bool {

	MonitorProcessName := name
	MonitorProcessUserName := user
	processInfos := utils.GetProcessInfo()

	for _, processInfo := range processInfos {

		processName, _ := processInfo.Name()
		processUserName, _ := processInfo.Username()
		if MonitorProcessName == processName && MonitorProcessUserName == processUserName {
			st, _ := processInfo.Status()
			if st != "T" {
				return false
			}
			break
		}
	}
	return true
}

func HandleMonitor(ctype string, val interface{}, level int64) error {

	monitors, err := model.GetMonitorByType(ctype)
	if err != nil {
		klog.Errorf("Get monitor by type error: %v", err.Error())
		return err
	}
	monitorInfo := monitors[0]
	value := utils.ToString(val)
	details := value
	content := value
	webContent := value
	name := ctype
	// klog.Infof("ctype: %s, value: %v, level: %v", ctype, val, level)
	switch ctype {
	case MonitorCpuUsageType:
		datac := utils.GetCpuLoadingTop5()
		data, _ := json.Marshal(datac)
		details = string(data)
		content = "CPU Usage Alarm: threshold:" + monitorInfo.Value + " Value:" + value
		webContent = "threshold:" + monitorInfo.Value + " Value:" + value
		name = "CPU Usage"
	case MonitorMemUsageType:
		datac := utils.GetMemLoadingTop5()
		data, _ := json.Marshal(datac)
		details = string(data)
		content = "Memory Usage Alarm: threshold:" + monitorInfo.Value + " Value:" + value
		webContent = "threshold:" + monitorInfo.Value + " Value:" + value
		name = "Memory Usage"
	case MonitorDiskUsageType:
		content = "Disk Usage Alarm: threshold:" + monitorInfo.Value + " Value:" + value
		webContent = "threshold:" + monitorInfo.Value + " Value:" + value
		name = "Disk Usage"
	case MonitorProcessType:
		name = "Process " + value
		if level == 0 {
			details = value + " is running"
			content = "Software Monitor Alarm:" + details
			webContent = details
		} else {
			details = value + " stopped"
			content = "Software Monitor Alarm:" + details
			webContent = details
		}
	}
	// klog.Infof("web content: %s", webContent)
	SendEventMsgToWeb(ctype, webContent, level)
	alertLog, err := model.GetAlertLogByNameAndType(name, AlertLogTypeMonitor)
	status := AlertLogUnsolved
	if level == 0 {
		status = AlertLogInvalid
	}
	//add alert log
	if err != nil {
		if err := model.AddAlertLog(&model.AlertLog{
			Name:      name,
			Details:   details,
			Level:     level,
			Status:    status,
			Value:     value,
			FuncId:    ctype,
			Condition: monitorInfo.Condition + monitorInfo.Value,
			LogType:   AlertLogTypeMonitor,
		}); err != nil {
			klog.Errorf("AddAlertLog err %v", err)
			return err
		}
	} else {
		if status != alertLog.Status ||
			level != alertLog.Level {
			if err := model.SaveAlertLog(alertLog.ID, "", &status, &level, details, value); err != nil {
				klog.Errorf("SaveAlertLog err %v", err)
				return err
			}
		}
	}
	id := utils.NewUUID()
	if err := model.AddAlertHistory(&model.AlertHistory{
		ID:        id,
		Name:      name,
		Details:   details,
		FuncId:    ctype,
		Value:     value,
		Condition: monitorInfo.Condition + monitorInfo.Value,
		Level:     level,
	}); err != nil {
		klog.Errorf("Add alert history error: %s", err.Error())
		return err
	}
	config, err := model.GetIthingsConfigByType(global.NotificationEmailType)
	if err != nil {
		klog.Errorf("Get notify email error: %v", err.Error())
		return err
	}
	if !SendAlarmEmail(config.Address, "Alarm Monitor", content) {
		klog.Errorf("Failed to send alarm email")
	}
	return nil
}
