package models

import (
	"go-gerbang/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Config struct {
	IdConfiguration    uuid.UUID      `gorm:"type:uuid;primaryKey" json:"idConfiguration"`
	ConfigurationGroup string         `gorm:"uniqueIndex;not null;size:64" json:"configurationGroup" validate:"required"`
	Data               datatypes.JSON `gorm:"type:jsonb;not null" json:"data"`
	CreatedAt          int            `gorm:"autoCreateTime:true" json:"created_at"`
	UpdatedAt          int            `gorm:"default:0;autoCreateTime:false" json:"updated_at"`
}

func (c *Config) BeforeCreate(tx *gorm.DB) error {
	if c.IdConfiguration == uuid.Nil {
		c.IdConfiguration = uuid.New()
	}
	return nil
}

func CreatConfig(config *Config) *gorm.DB {
	// return database.GDB.Create(config)
	now := time.Now().Unix()
	return database.GDB.Exec(`
		INSERT INTO configs (id_configuration, configuration_group, data, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (configuration_group)
		DO UPDATE SET 
			data = EXCLUDED.data, 
			updated_at = EXCLUDED.updated_at
	`, uuid.New(), config.ConfigurationGroup, config.Data, now)
}

func FindConfig(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.GDB.Model(&Config{}).Find(dest, conds...)
}
