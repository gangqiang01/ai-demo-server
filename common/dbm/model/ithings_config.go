package model

import (
	"time"

	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type IthingsConfig struct {
	ID int64 `gorm:"primary_key; auto_increment" json:"id"`
	// Name     		string `gorm:"column:name; type:varchar(256);" json:"name"`
	Type            string `gorm:"column:type; type:varchar(256);" json:"type"`
	Address         string `gorm:"column:address; type:varchar(256);" json:"address"`
	Port            int    `gorm:"column:port;" json:"port"`
	Username        string `gorm:"column:username; type:varchar(256);" json:"username"`
	Password        string `gorm:"column:password; type:varchar(256);" json:"password"`
	Content         string `gorm:"column:content; type:varchar(256);" json:"content"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (IthingsConfig) TableName() string {
	return "ithings_config"
}

func GetIthingsConfigByType(ctype string) (*IthingsConfig, error) {
	config := &IthingsConfig{}
	err := global.DBAccess.Where("type = ?", ctype).First(config).Error
	if err != nil {
		return nil, err
	}

	return config, err
}

func IsExistIthingsConfigByType(ctype string) bool {
	var count int64
	err := global.DBAccess.Model(&IthingsConfig{}).Where("type = ?", ctype).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func AddIthingsConfig(config *IthingsConfig) error {
	config.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&config).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveIthingsConfig(ctype, address, username, password, content string, port *int) error {
	configMap := make(map[string]interface{}, 4)
	if address != "" {
		configMap["address"] = address
	}
	if username != "" {
		configMap["username"] = username
	}
	if password != "" {
		configMap["password"] = password
	}
	if port != nil {
		configMap["port"] = port
	}
	if content != "" {
		configMap["content"] = content
	}
	err := global.DBAccess.Model(&IthingsConfig{}).Where("type = ?", ctype).Updates(configMap).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
func DeleteIthingsConfig(ctype string) error {
	err := global.DBAccess.Where("type = ?", ctype).Delete(&IthingsConfig{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)

		return err
	}
	return nil

}
