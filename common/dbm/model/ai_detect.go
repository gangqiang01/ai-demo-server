package model

import (
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

type AiDetect struct {
	Id              string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	CType           string `gorm:"column:type; type:varchar(256);" json:"type"`
	Source          string `gorm:"column:source; type:varchar(256);" json:"source"`
	Config          string `gorm:"column:config; type:varchar(256);" json:"config"`
	AiModelId       string `gorm:"column:ai_model_id; type:varchar(36);" json:"aiModelId"`
	IsShow          string `gorm:"column:is_show; type:varchar(36);" json:"is_show"`
	Notification    string `gorm:"column:notification; type:text;" json:"notification"`
	Status          string `gorm:"column:status; type:varchar(256);" json:"status"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (d *AiDetect) TableName() string {
	return "ai_detect"
}

func GetAiDetectsByPageAndKeywords(page int, limit int, keywords string) ([]*AiDetect, error) {
	var ai_detects []*AiDetect
	tx := global.DBAccess.Model(&AiDetect{}).Offset((page - 1) * limit).Limit(limit)
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Order("create_time_stamp desc").Find(&ai_detects).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ai_detects, err
}

func GetAiDetectsByKeywords(keywords string) ([]*AiDetect, error) {
	var ai_detects []*AiDetect
	tx := global.DBAccess.Model(&AiDetect{})
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Order("create_time_stamp desc").Find(&ai_detects).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ai_detects, err
}

func GetAiDetectCount(keywords string) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&AiDetect{})
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

func GetAiDetectById(id string) (*AiDetect, error) {
	ai_detect := &AiDetect{}
	err := global.DBAccess.Where("id = ?", id).First(ai_detect).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return ai_detect, nil
}

func GetAiDetectsByStatus(status string) ([]*AiDetect, error) {
	var ai_detects []*AiDetect
	tx := global.DBAccess.Model(&AiDetect{})
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	err := tx.Order("create_time_stamp desc").Find(&ai_detects).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ai_detects, err
}

func IsExistAiDetectByAiModelId(id string) bool {
	var count int64
	err := global.DBAccess.Model(&AiDetect{}).Where("ai_model_id = ?", id).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}

	if count > 0 {
		return true
	}
	return false
}

func IsExistAiDetectByChannelId(id string) bool {
	var count int64
	err := global.DBAccess.Model(&AiDetect{}).Where("channel_id = ?", id).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func IsExistAiDetectByName(name string) bool {
	var count int64
	err := global.DBAccess.Model(&AiDetect{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func GetAiDetectsByChannelId(id string) ([]*AiDetect, error) {
	var ai_detects []*AiDetect
	err := global.DBAccess.Model(&AiDetect{}).Where("channel_id = ?", id).Find(&ai_detects).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return ai_detects, err
	}
	return ai_detects, nil
}

func AddAiDetect(AiDetect *AiDetect) error {
	AiDetect.CreateTimeStamp = utils.GetNowTimeStamp()
	err := global.DBAccess.Debug().Create(&AiDetect).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveAiDetect(id, name, notification, status string) error {

	vals := make(map[string]interface{})
	if name != "" {
		vals["name"] = name
	}

	if notification != "" {
		vals["notification"] = notification
	}

	if status != "" {
		vals["status"] = status
	}
	err := global.DBAccess.Model(&AiDetect{}).Where("id = ?", id).Updates(vals).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAiDetect(id string) error {
	err := global.DBAccess.Where("id = ?", id).Delete(&AiDetect{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
