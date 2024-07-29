package model

import (
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

type Monitor struct {
	Id    string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Ctype string `gorm:"column:type; not null; type:varchar(256);" json:"type"`
	Name  string `gorm:"column:name; type:varchar(256);" json:"name"`
	User  string `gorm:"column:user; type:varchar(256);" json:"user"`
	Pid   string `gorm:"column:pid; type:varchar(256);" json:"pid"`
	//1min
	Delay string `gorm:"column:delay; type:varchar(256);" json:"delay"`
	//< > =
	Condition string `gorm:"column:condition; type:varchar(256);" json:"condition"`
	Value     string `gorm:"column:value; type:varchar(256);" json:"value"`
	//enable disable
	Status          string `gorm:"column:status; type:varchar(256);" json:"status"`
	Level           string `gorm:"column:level; type:varchar(256);" json:"level"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (d *Monitor) TableName() string {
	return "monitor"
}

func GetMonitorByPageAndKeywords(ctype string, page int, limit int, keywords string) ([]*Monitor, error) {
	var monitors []*Monitor
	tx := global.DBAccess.Model(&Monitor{}).Offset((page - 1) * limit).Limit(limit)
	if ctype != "" {
		tx = tx.Where("type = ?", ctype)
	}
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Order("create_time_stamp desc").Find(&monitors).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return monitors, err
}

func GetMonitorByKeywords(ctype, keywords string) ([]*Monitor, error) {
	var monitors []*Monitor
	tx := global.DBAccess.Model(&Monitor{})
	if ctype != "" {
		tx = tx.Where("type = ?", ctype)
	}
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Order("create_time_stamp desc").Find(&monitors).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return monitors, err
}

func GetMonitorCount(ctype, keywords string) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&Monitor{})
	if ctype != "" {
		tx = tx.Where("type = ?", ctype)
	}
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetMonitorByType(ctype string) ([]*Monitor, error) {
	var monitors []*Monitor
	tx := global.DBAccess.Model(&Monitor{})
	if ctype != "" {
		tx = tx.Where("type = ?", ctype)
	}
	err := tx.Order("create_time_stamp desc").Find(&monitors).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return monitors, err
}
func GetMonitorByCondition(ctype, name, user string) (*Monitor, error) {
	var monitor *Monitor
	tx := global.DBAccess.Model(&Monitor{})
	if ctype != "" {
		tx = tx.Where("type = ?", ctype)
	}
	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if user != "" {
		tx = tx.Where("user = ?", user)
	}
	err := tx.Order("create_time_stamp desc").First(&monitor).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return monitor, err
}
func GetMonitorById(id string) (*Monitor, error) {
	monitor := &Monitor{}
	err := global.DBAccess.Where("id = ?", id).First(monitor).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return monitor, nil
}
func IsExistMonitorByNameAndUser(ctype, name, user string) bool {
	var count int64
	tx := global.DBAccess.Model(&Monitor{})
	if ctype != "" {
		tx = tx.Where("type = ?", ctype)
	}
	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if user != "" {
		tx = tx.Where("user = ?", user)
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
func AddMonitor(Monitor *Monitor) error {
	klog.Infof("AddMonitor: %v", *Monitor)
	Monitor.CreateTimeStamp = utils.GetNowTimeStamp()
	err := global.DBAccess.Create(&Monitor).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveMonitor(id, status, value, delay string) error {
	klog.Infof("SaveMonitor: id: %v, status: %v, value: %v, delay: %v", id, status, value, delay)
	vals := make(map[string]interface{})
	if status != "" {
		vals["status"] = status
	}

	if value != "" {
		vals["value"] = value
	}
	if delay != "" {
		vals["delay"] = delay
	}

	err := global.DBAccess.Model(&Monitor{}).Where("id = ?", id).Updates(vals).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveMonitorLevel(id, level string) error {
	vals := make(map[string]interface{})
	if level != "" {
		vals["level"] = level
	}

	err := global.DBAccess.Model(&Monitor{}).Where("id = ?", id).Updates(vals).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteMonitor(id string) error {
	klog.Infof("DeleteMonitor: %v", id)
	err := global.DBAccess.Where("id = ?", id).Delete(&Monitor{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
