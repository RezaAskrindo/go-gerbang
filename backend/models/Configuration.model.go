package models

import (
	"fmt"
	"go-gerbang/database"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Configuration struct {
	IdConfiguration    uuid.UUID `gorm:"type:uuid;primaryKey" json:"idConfiguration"`
	ConfigurationGroup string    `gorm:"not null;size:64" json:"configurationGroup" validate:"required"`
	ConfigurationName  *string   `gorm:"default:null;size:128" json:"configurationName"`
	ConfigurationIndex *int      `gorm:"default:null" json:"configurationIndex"`
	ConfigurationKey   string    `gorm:"not null;size:128" json:"configurationKey" validate:"required"`
	ConfigurationValue *string   `gorm:"default:null;size:512" json:"configurationValue"`
}

type GroupConfiguration struct {
	ConfigurationGroup string
	ConfigurationName  *string
}

func (c *Configuration) BeforeCreate(tx *gorm.DB) error {
	if c.IdConfiguration == uuid.Nil {
		c.IdConfiguration = uuid.New()
	}
	return nil
}

func CreateConfiguration(configuration []Configuration) *gorm.DB {
	tx := database.GDB.Begin()
	if tx.Error != nil {
		return tx
	}

	err := tx.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id_configuration"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"configuration_group",
				"configuration_index",
				"configuration_key",
				"configuration_name",
				"configuration_value",
			}),
		}).
		CreateInBatches(configuration, 100).
		Error

	if err != nil {
		tx.Rollback()
		return tx
	}

	return tx.Commit()
}

func UpdateConfigurationIndex(group interface{}, name interface{}, index int) *gorm.DB {
	return database.GDB.Model(&Configuration{}).Where("configuration_group = ? and configuration_name = ?", group, name).Update("configuration_index", index)
}

func FindConfiguration(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.GDB.Model(&Configuration{}).Find(dest, conds...)
}

func FindGroupConfiguration(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.GDB.Model(&Configuration{}).Group("configuration_group, configuration_name").Find(dest, conds...)
}

func ParseConfiguration(configs *[]Configuration) interface{} {
	hasIndex := false

	for _, c := range *configs {
		if c.ConfigurationIndex != nil {
			hasIndex = true
			break
		}
	}

	if !hasIndex {
		result := make(map[string]string)

		for _, c := range *configs {
			result[c.ConfigurationKey] = *c.ConfigurationValue
		}

		return result
	}

	grouped := make(map[int]map[string]string)

	for _, c := range *configs {
		idx := *c.ConfigurationIndex

		if _, exists := grouped[idx]; !exists {
			grouped[idx] = make(map[string]string)
		}

		grouped[idx][c.ConfigurationKey] = *c.ConfigurationValue
	}

	indexes := make([]int, 0, len(grouped))
	for idx := range grouped {
		indexes = append(indexes, idx)
	}
	sort.Ints(indexes)

	result := make([]map[string]string, 0, len(indexes))
	for _, idx := range indexes {
		result = append(result, grouped[idx])
	}

	return result
}

func DeleteConfigurationByConfName(configurationGroup string, configurationName string) error {
	tx := database.GDB.Unscoped().Where("configuration_group = ? AND configuration_name = ?", configurationGroup, configurationName).Delete(&Configuration{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("no configuration found with group = %s and name = %s", configurationGroup, configurationName)
	}

	return nil
}
