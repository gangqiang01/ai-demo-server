package model

import (
	"time"

	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type AlertHistory struct {
	ID              string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	FuncId          string `gorm:"column:funcId; type:varchar(256);" json:"funcId"`
	Level           int64  `gorm:"column:level;" json:"level"`
	Details         string `gorm:"column:details; type:text;" json:"details"`
	Value           string `gorm:"column:value; type:varchar(256);" json:"value"`
	Condition       string `gorm:"column:condition; type:varchar(256);" json:"condition"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (AlertHistory) TableName() string {
	return "alert_history"
}

func GetAlertHistory() ([]*AlertHistory, error) {
	var alertHistorys []*AlertHistory
	err := global.DBAccess.Order("update_time_stamp desc").Find(&alertHistorys).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertHistorys, err
}

func GetAlertHistoryById(id string) (*AlertHistory, error) {
	var alertHistory *AlertHistory
	err := global.DBAccess.Where(&AlertHistory{ID: id}).Find(&alertHistory).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertHistory, err
}

func GetAlertHistoryByPage(page int, limit int) ([]*AlertHistory, error) {
	var alertHistorys []*AlertHistory
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&alertHistorys).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertHistorys, err
}

func GetAlertHistoryCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertHistory{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetAlertHistoryByPageAndCondition(page int, limit int, name string, level *int64, beginTs *int64, endTs *int64) ([]*AlertHistory, error) {
	var alertHistorys []*AlertHistory
	tx := global.DBAccess.Model(&AlertHistory{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}

	if level != nil {
		tx = tx.Where("level = ?", *level)
	}
	if beginTs != nil && endTs != nil {
		tx = tx.Where("create_time_stamp between ? and ?", *beginTs, *endTs)
	}
	err := tx.Offset((page - 1) * limit).Limit(limit).Order("update_time_stamp desc").Find(&alertHistorys).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertHistorys, err
}

func GetAlertHistoryByCondition(name, edgeId, deviceId string, level *int64, beginTs *int64, endTs *int64, count int) ([]*AlertHistory, error) {
	var alertHistorys []*AlertHistory
	tx := global.DBAccess.Model(&AlertHistory{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if level != nil {
		tx = tx.Where("level = ?", *level)
	}
	if beginTs != nil && endTs != nil {
		tx = tx.Where("create_time_stamp between ? and ?", *beginTs, *endTs)
	}
	if count > 0 {
		tx = tx.Limit(count)
	}
	err := tx.Order("update_time_stamp desc").Find(&alertHistorys).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertHistorys, err
}

func GetAlertHistoryCountByCondition(name string, level *int64, beginTs *int64, endTs *int64) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&AlertHistory{})

	if name != "" {
		tx = tx.Where("name = ?", name)
	}
	if level != nil {
		tx = tx.Where("level = ?", *level)
	}
	if beginTs != nil && endTs != nil {
		tx = tx.Where("update_time_stamp between ? and ?", *beginTs, *endTs)
	}
	err := tx.Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0, err
	}
	return count, err
}

func GetAlertHistoryByPageAndKeywords(page int, limit int, keywords string) ([]*AlertHistory, error) {
	var alertHistorys []*AlertHistory
	err := global.DBAccess.Where("edge_name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("update_time_stamp desc").Limit(limit).Find(&alertHistorys).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return alertHistorys, err
}
func GetAlertHistoryCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&AlertHistory{}).Where("edge_name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return 0, err
	}
	return count, err
}

func GetAlertHistoryByName(name string) ([]*AlertHistory, error) {
	var alertHistorys []*AlertHistory
	err := global.DBAccess.Model(&AlertHistory{}).Where("name = ?", name).First(&alertHistorys).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return alertHistorys, err
	}
	return alertHistorys, err
}

func AddAlertHistory(alertHistory *AlertHistory) error {
	alertHistory.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&alertHistory).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAlertHistory(id string) error {
	if err := global.DBAccess.Where("id = ?", id).Delete(&AlertHistory{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAllAlertHistory() error {
	if err := global.DBAccess.Exec("DELETE FROM alert_history").Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAlertHistoryByName(name string) error {
	if err := global.DBAccess.Where("name = ?", name).Delete(&AlertHistory{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAlertHistoryBeforeCreateTs(ts int64) error {
	if err := global.DBAccess.Where("create_time_stamp < ?", ts).Delete(&AlertHistory{}).Error; err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func BatchDeleteAlertHistory(ids []string) error {
	err := global.DBAccess.Where("id in ?", ids).Delete(&AlertHistory{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
