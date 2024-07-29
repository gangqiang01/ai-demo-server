package model

import (
	"github.com/edgehook/ithings/common/global"
	"github.com/edgehook/ithings/common/utils"
	"k8s.io/klog/v2"
)

type AiModel struct {
	Id       string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name     string `gorm:"column:name; not null; type:varchar(256);" json:"name"`
	CType    string `gorm:"column:type; not null; type:varchar(256);" json:"type"`
	FileName string `gorm:"column:filename; type:text;" json:"filename"`
	//detect python path
	Path            string `gorm:"column:path; type:text;" json:"path"`
	Labels          string `gorm:"column:labels; type:text;" json:"labels"`
	Description     string `gorm:"column:description; default:null; type:varchar(512);" json:"description"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (d *AiModel) TableName() string {
	return "ai_model"
}

func GetAiModelsByPageAndKeywords(page int, limit int, keywords string) ([]*AiModel, error) {
	var ai_models []*AiModel
	tx := global.DBAccess.Model(&AiModel{}).Offset((page - 1) * limit).Limit(limit)
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Order("create_time_stamp desc").Find(&ai_models).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ai_models, err
}
func GetAiModels() ([]*AiModel, error) {
	var ai_models []*AiModel
	err := global.DBAccess.Model(&AiModel{}).Find(&ai_models).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return ai_models, nil
}

func GetAiModelsByKeywords(keywords string) ([]*AiModel, error) {
	var ai_models []*AiModel
	tx := global.DBAccess.Model(&AiModel{})
	if keywords != "" {
		tx = tx.Where("name LIKE ?", "%"+keywords+"%")
	}
	err := tx.Order("create_time_stamp desc").Find(&ai_models).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return ai_models, err
}

func GetAiModelCount(keywords string) (int64, error) {
	var count int64
	tx := global.DBAccess.Model(&AiModel{})
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

func GetAiModelById(id string) (*AiModel, error) {
	ai_model := &AiModel{}
	err := global.DBAccess.Where("id = ?", id).First(ai_model).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return ai_model, nil
}

func IsExistAiModelByName(name string) bool {
	var count int64
	err := global.DBAccess.Model(&AiModel{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func AddAiModel(AiModel *AiModel) error {
	AiModel.CreateTimeStamp = utils.GetNowTimeStamp()
	err := global.DBAccess.Create(&AiModel).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveAiModel(id, name, labels, description string) error {
	vals := make(map[string]interface{})
	if name != "" {
		vals["name"] = name
	}
	if labels != "" {
		vals["labels"] = labels
	}

	if description != "" {
		vals["description"] = description
	}

	err := global.DBAccess.Model(&AiModel{}).Where("id = ?", id).Updates(vals).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteAiModel(id string) error {
	err := global.DBAccess.Where("id = ?", id).Delete(&AiModel{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
