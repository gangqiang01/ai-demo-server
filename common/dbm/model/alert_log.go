package model

import (
	"time"

	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type AlertLog struct {
	ID              int64  `gorm:"primary_key; auto_increment" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	LogType         string `gorm:"column:log_type; type:varchar(256);" json:"logType"`
	FuncId          string `gorm:"column:funcId; type:varchar(256);" json:"funcId"`
	Level           int64  `gorm:"column:level;" json:"level"`
	DeviceName      string `gorm:"column:device_name; type:varchar(256);" json:"deviceName"`
	DeviceId        string `gorm:"column:device_id; type:varchar(256);" json:"deviceId"`
	Details         string `gorm:"column:details; type:text;" json:"details"`
	Value           string `gorm:"column:value; type:varchar(256);" json:"value"`
	Condition       string `gorm:"column:condition; type:varchar(256);" json:"condition"`
	HandleStatus    int32  `gorm:"column:handleStatus;" json:"handleStatus"`
	Status          int32  `gorm:"column:status;" json:"status"`
	Record          string `gorm:"column:record; type:text;" json:"record"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (AlertLog) TableName() string {
	return "alert_log"
}

func GetAlertLog() ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	err := global.DBAccess.Order("update_time_stamp desc").Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetAlertLogByType(logType string) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	err := global.DBAccess.Where("log_type = ?", logType).Order("update_time_stamp desc").Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetUnresolvedAlertLogByTypeAndDeviceId(logType, deviceId string, status []int32) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	err := global.DBAccess.Where("log_type = ? and device_id = ? and status in ?", logType, deviceId, status).Order("update_time_stamp desc").Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetAlertLogById(id int64) (*AlertLog, error) {
	var alertLog *AlertLog
	err := global.DBAccess.Where(&AlertLog{ID: id}).Find(&alertLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLog, err
}

func GetAlertLogByPage(page int, limit int) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetAlertLogCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertLog{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0, err
	}
	return count, err
}

func GetAlertLogByPageAndType(page int, limit int, logType string) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	err := global.DBAccess.Where("log_type = ?", logType).Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}
func GetAlertLogCountByType(logType string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertLog{}).Where("log_type = ?", logType).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0, err
	}
	return count, err
}
func GetAlertLogByPageAndCondition(page int, limit int, name string, status *int32, level *int64, beginTs *int64, endTs *int64, logType string) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	tx := global.DBAccess.Model(&AlertLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}

	if status != nil {
		tx = tx.Where("status = ?", *status)
	}
	if level != nil {
		tx = tx.Where("level = ?", *level)
	}
	if logType != "" {
		tx = tx.Where("log_type = ?", logType)
	}
	if beginTs != nil && endTs != nil {
		tx = tx.Where("create_time_stamp between ? and ?", *beginTs, *endTs)
	}
	err := tx.Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetAlertLogByCondition(name string, status *int32, level *int64, beginTs *int64, endTs *int64, logType string) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	tx := global.DBAccess.Model(&AlertLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}

	if status != nil {
		tx = tx.Where("status = ?", *status)
	}
	if level != nil {
		tx = tx.Where("level = ?", *level)
	}
	if logType != "" {
		tx = tx.Where("log_type = ?", logType)
	}
	if beginTs != nil && endTs != nil {
		tx = tx.Where("create_time_stamp between ? and ?", *beginTs, *endTs)
	}
	err := tx.Order("update_time_stamp desc").Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetAlertLogCountByCondition(name string, status *int32, level *int64, beginTs *int64, endTs *int64, logType string) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&AlertLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}

	if status != nil {
		tx = tx.Where("status = ?", *status)
	}
	if level != nil {
		tx = tx.Where("level = ?", *level)
	}
	if logType != "" {
		tx = tx.Where("log_type = ?", logType)
	}
	if beginTs != nil && endTs != nil {
		tx = tx.Where("create_time_stamp between ? and ?", *beginTs, *endTs)
	}
	err := tx.Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0, err
	}
	return count, err
}

func GetAlertLogByPageAndKeywords(page int, limit int, keywords string) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	err := global.DBAccess.Where("edge_name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("update_time_stamp desc").Limit(limit).Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetAlertLogCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertLog{}).Where("edge_name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0, err
	}
	return count, err
}
func GetAlertLogByPageAndKeywordsAndType(page int, limit int, keywords string, logType string) ([]*AlertLog, error) {
	var alertLogs []*AlertLog
	err := global.DBAccess.Where("log_type = ? and edge_name LIKE ?", logType, "%"+keywords+"%").Offset((page - 1) * limit).Order("update_time_stamp desc").Limit(limit).Find(&alertLogs).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertLogs, err
}

func GetAlertLogCountByKeywordsAndType(keywords, logType string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertLog{}).Where("log_type = ? and edge_name LIKE ?", logType, "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0, err
	}
	return count, err
}
func GetAlertLogByName(name, deviceId string) (*AlertLog, error) {
	var alertLog *AlertLog
	tx := global.DBAccess.Model(&AlertLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if deviceId != "" {
		tx = tx.Where("device_id = ?", deviceId)
	}
	err := tx.First(&alertLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return alertLog, err
	}
	return alertLog, err
}

func IsExistAlertLogByDeviceIdAndLevelAndStatus(deviceId string, level int64, status []int32) bool {
	var count int64
	err := global.DBAccess.Model(&AlertLog{}).Where("device_id = ? and level = ? and status in ?", deviceId, level, status).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func GetAlertLogByNameAndType(name, logType string) (*AlertLog, error) {
	var alertLog *AlertLog
	tx := global.DBAccess.Model(&AlertLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if logType != "" {
		tx = tx.Where("log_type = ?", logType)
	}
	err := tx.First(&alertLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return alertLog, err
	}
	return alertLog, err
}

func IsExistAlertLogByNameAndDevice(name, deviceId string) bool {
	var count int64
	tx := global.DBAccess.Model(&AlertLog{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if deviceId != "" {
		tx = tx.Where("device_id = ?", deviceId)
	}
	err := tx.Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}

	if count > 0 {
		return true
	}
	return false
}

func AddAlertLog(alertLog *AlertLog) error {
	alertLog.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&alertLog).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
func SaveAlertLog(id int64, record string, status *int32, level *int64, details, value string) error {
	vals := make(map[string]interface{})
	if record != "" {
		vals["record"] = record
	}
	if status != nil {
		vals["status"] = status
	}
	if level != nil {
		vals["level"] = level
	}
	if details != "" {
		vals["details"] = details
	}
	if value != "" {
		vals["value"] = value
	}
	err := global.DBAccess.Model(&AlertLog{}).Where("id = ?", id).Updates(vals).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAlertLog(id int64) error {
	if err := global.DBAccess.Where("id = ?", id).Delete(&AlertLog{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAlertLogByName(name string) error {
	if err := global.DBAccess.Where("name = ?", name).Delete(&AlertLog{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAllAlertLog() error {
	if err := global.DBAccess.Exec("DELETE FROM alert_log").Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func BatchDeleteAlertLog(ids []int64) error {
	err := global.DBAccess.Where("id in ?", ids).Delete(&AlertLog{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
