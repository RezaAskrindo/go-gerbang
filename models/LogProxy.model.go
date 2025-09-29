package models

import (
	"go-gerbang/database"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LogProxy struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;type:bigint" json:"id"`
	Level     string         `gorm:"size:16;index" json:"level"`
	Service   string         `gorm:"size:32" json:"service"`
	Method    string         `gorm:"size:16" json:"method"`
	Path      string         `gorm:"not null;size:512" json:"path"`
	UserAuth  string         `gorm:"default:null;size:255" json:"user_auth"`
	Status    uint16         `gorm:"not null" json:"status"`
	Duration  float64        `gorm:"not null" json:"duration"`
	Fields    datatypes.JSON `gorm:"type:jsonb" json:"fields"`
	Timestamp time.Time      `gorm:"type:timestamptz" json:"timestamp"`
}

func FindLogProxy(dest *[]LogProxy, service, method, path, status string, from, to time.Time) *gorm.DB {
	query := `
			SELECT id, level, service, method, path, user_auth, status, duration, fields, timestamp FROM log_proxies 
			WHERE service = ?
			AND method = ?
			AND path = ?
			AND status = ?
			AND timestamp BETWEEN ? AND ?
			ORDER BY timestamp DESC
	`
	return database.GDB.Raw(query, service, method, path, status, from, to).Scan(dest)
}

type PathStats struct {
	Service      string  `json:"service"`
	Method       string  `json:"method"`
	Path         string  `json:"path"`
	Status       uint16  `json:"status"`
	AvgDuration  float64 `json:"avg_duration"`
	RequestCount int64   `json:"request_count"`
}

func FindStatsLogProxy(dest *[]PathStats, from, to time.Time) *gorm.DB {
	query := `
			SELECT service,
							method,
							path,
							status,
							AVG(duration)   AS avg_duration,
							COUNT(*)        AS request_count
			FROM log_proxies
			WHERE timestamp BETWEEN ? AND ?
			GROUP BY service, method, path, status
			ORDER BY request_count DESC
	`
	return database.GDB.Raw(query, from, to).Scan(dest)
}
