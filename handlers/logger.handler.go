package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"go-gerbang/database"
	"go-gerbang/models"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/datatypes"
)

var ZapLogger *zap.Logger

type logEntry struct {
	Level     string
	Message   string
	Timestamp time.Time
	Fields    []zapcore.Field
}

type AsyncPGXCore struct {
	zapcore.LevelEnabler
	encoder zapcore.Encoder
	queue   chan logEntry
}

func NewAsyncPGXCore(level zapcore.Level, queue chan logEntry) zapcore.Core {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "timestamp"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return &AsyncPGXCore{
		LevelEnabler: level,
		encoder:      zapcore.NewJSONEncoder(encCfg),
		queue:        queue,
	}
}

func (c *AsyncPGXCore) With(fields []zapcore.Field) zapcore.Core {
	return &AsyncPGXCore{
		LevelEnabler: c.LevelEnabler,
		encoder:      c.encoder,
		queue:        c.queue,
	}
}

func (c *AsyncPGXCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *AsyncPGXCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	select {
	case c.queue <- logEntry{
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Timestamp: entry.Time,
		Fields:    fields,
	}:
	default:
		// Optional: drop log if queue is full
	}
	return nil
}

func (c *AsyncPGXCore) Sync() error {
	return nil
}

func StartLogWorker(ctx context.Context, queue <-chan logEntry) {
	for {
		select {
		case <-ctx.Done():
			return

		case log := <-queue:
			fieldMap := make(map[string]interface{})
			for _, f := range log.Fields {
				switch f.Type {
				case zapcore.StringType:
					fieldMap[f.Key] = f.String
				case zapcore.Int64Type:
					fieldMap[f.Key] = f.Integer
				case zapcore.Int32Type:
					fieldMap[f.Key] = int32(f.Integer)
				case zapcore.Int16Type:
					fieldMap[f.Key] = int16(f.Integer)
				case zapcore.Int8Type:
					fieldMap[f.Key] = int8(f.Integer)
				case zapcore.Uint64Type:
					fieldMap[f.Key] = f.Integer
				case zapcore.Uint32Type:
					fieldMap[f.Key] = uint32(f.Integer)
				case zapcore.BoolType:
					fieldMap[f.Key] = f.Integer == 1
				case zapcore.Float64Type:
					fieldMap[f.Key] = math.Float64frombits(uint64(f.Integer))
				case zapcore.BinaryType:
					fieldMap[f.Key] = f.Interface
				case zapcore.DurationType:
					fieldMap[f.Key] = float64(time.Duration(f.Integer).Milliseconds())
				case zapcore.ErrorType:
					if err, ok := f.Interface.(error); ok {
						fieldMap[f.Key] = err.Error()
					} else {
						fieldMap[f.Key] = fmt.Sprint(f.Interface)
					}
				case zapcore.TimeType:
					if t, ok := f.Interface.(time.Time); ok {
						fieldMap[f.Key] = t.Format(time.RFC3339Nano)
					} else {
						fieldMap[f.Key] = fmt.Sprint(f.Interface)
					}
				default:
					fieldMap[f.Key] = fmt.Sprintf("UNSUPPORTED_TYPE_%v", f.Type)
				}
			}

			jsonFields, err := json.Marshal(fieldMap)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to marshal log fields: %v\n", err)
				jsonFields = []byte("{}")
			}

			entry := models.LogProxy{
				Level:     log.Level,
				Service:   log.Message,
				Method:    fmt.Sprint(fieldMap["method"]),
				Path:      fmt.Sprint(fieldMap["path"]),
				UserAuth:  fmt.Sprint(fieldMap["user"]),
				Status:    uint16(fieldMap["status"].(int64)),
				Duration:  0,
				Timestamp: log.Timestamp,
				Fields:    datatypes.JSON(jsonFields),
			}

			var durationVal float64
			if v, ok := fieldMap["duration"]; ok {
				switch d := v.(type) {
				case float64:
					durationVal = d
				case int64:
					durationVal = float64(d)
				default:
					durationVal = 0.0
				}
			}
			entry.Duration = durationVal

			if err := database.GDB.WithContext(ctx).Create(&entry).Error; err != nil {
				fmt.Fprintf(os.Stderr, "GORM log insert failed: %v\n", err)
			}
		}
	}
}

func InitLogger(ctx context.Context) error {
	if err := database.GDB.AutoMigrate(&models.LogProxy{}); err != nil {
		return err
	}

	logQueue := make(chan logEntry, 1000)
	go StartLogWorker(ctx, logQueue)

	pgxCore := NewAsyncPGXCore(zapcore.InfoLevel, logQueue)
	ZapLogger = zap.New(pgxCore)

	return nil
}
