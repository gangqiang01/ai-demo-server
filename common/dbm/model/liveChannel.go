package model

import (
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

type LiveChannel struct {
	Id              string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name            string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	AccessType      string `gorm:"column:access_type; type:varchar(256);" json:"accessType"`
	Status          string `gorm:"column:status; type:varchar(256);" json:"status"`
	Config          string `gorm:"column:config; type:text;" json:"config"`
	Duration        int    `gorm:"column:duration; default:null; type:int;" json:"duration"`
	Segment         int    `gorm:"column:segment; default:null; type:int;" json:"segment"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (d *LiveChannel) TableName() string {
	return "live_channel"
}
func GetLiveChannels() ([]*LiveChannel, error) {
	var liveChannels []*LiveChannel
	err := global.DBAccess.Order("create_time_stamp desc").Find(&liveChannels).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return liveChannels, nil
}

func GetLiveChannelByStatus(status string) ([]*LiveChannel, error) {
	var liveChannels []*LiveChannel
	err := global.DBAccess.Where("status = ?", status).Order("create_time_stamp desc").Find(&liveChannels).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return liveChannels, nil
}

func GetLiveChannelByPageAndKeywords(page int, limit int, keywords string, status string) ([]*LiveChannel, error) {
	var liveChannels []*LiveChannel
	tx := global.DBAccess.Model(&LiveChannel{}).Offset((page - 1) * limit).Limit(limit)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Order("create_time_stamp desc").Debug().Find(&liveChannels).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return liveChannels, err
}

func GetLiveChannelByName(name string) (*LiveChannel, error) {
	liveChannel := &LiveChannel{}
	err := global.DBAccess.Where("name = ?", name).First(liveChannel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return liveChannel, err
}

func GetLiveChannelCount(keywords string) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&LiveChannel{})
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

func GetLiveChannelById(id string) (*LiveChannel, error) {
	liveChannel := &LiveChannel{}
	err := global.DBAccess.Where("id = ?", id).First(liveChannel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return liveChannel, nil
}

func IsExistLiveChannelByName(name string) bool {
	var count int64
	err := global.DBAccess.Model(&LiveChannel{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func AddLiveChannel(LiveChannel *LiveChannel) error {
	LiveChannel.CreateTimeStamp = utils.GetNowTimeStamp()
	err := global.DBAccess.Create(&LiveChannel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveLiveChannel(id, name, accessType, status, config string, duration, segment *int) error {
	vals := make(map[string]interface{})
	if name != "" {
		vals["Name"] = name
	}
	if status != "" {
		vals["status"] = status
	}

	if config != "" {
		vals["config"] = config
	}

	if duration != nil {
		vals["duration"] = duration
	}

	if segment != nil {
		vals["segment"] = segment
	}

	if accessType != "" {
		vals["access_type"] = accessType
	}
	err := global.DBAccess.Model(&LiveChannel{}).Where("id = ?", id).Updates(vals).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteLiveChannel(id string) error {
	err := global.DBAccess.Where("id = ?", id).Delete(&LiveChannel{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
