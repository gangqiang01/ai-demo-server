package model

import (
	"gorm.io/gorm"
)

/*
* RegisterTables create all database tables in this function.
* Notice! you should create tables at here!
 */
func RegisterTables(db *gorm.DB) error {
	err := db.AutoMigrate(
		&ProtocolTypes{},
		&AlertConfig{},
		&AlertLog{},
		&AlertHistory{},
		&IthingsConfig{},
		&LiveChannel{},
		&User{},
		&Monitor{},
		&AiModel{},
		&AiDetect{})

	if err != nil {
		return err
	}

	return nil
}
