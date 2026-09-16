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

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/datatypes"
)

var (
	ZapLogger *zap.Logger
	// LogQueue is package-level so both the local zap core (Write) and the
	// NATS subscriber (SubscribeLogEvents) can feed the same single DB-writer
	// goroutine, avoiding concurrent writes from multiple sources.
	LogQueue chan logEntry
)

type logEntry struct {
	Level     string
	Message   string
	Timestamp time.Time
	FieldMap  map[string]interface{} // pre-converted, so the worker doesn't care about origin (zap vs NATS)
}

type AsyncPGXCore struct {
	zapcore.LevelEnabler
	encoder zapcore.Encoder
	queue   chan logEntry
	fields  []zapcore.Field
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
	merged := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	merged = append(merged, c.fields...)
	merged = append(merged, fields...)
	return &AsyncPGXCore{
		LevelEnabler: c.LevelEnabler,
		encoder:      c.encoder,
		queue:        c.queue,
		fields:       merged,
	}
}

func (c *AsyncPGXCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *AsyncPGXCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	allFields = append(allFields, c.fields...)
	allFields = append(allFields, fields...)

	select {
	case c.queue <- logEntry{
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Timestamp: entry.Time,
		FieldMap:  fieldsToMap(allFields),
	}:
	default:
	}
	return nil
}

func (c *AsyncPGXCore) Sync() error { return nil }

// fieldsToMap converts zapcore fields to a plain map — used only by the local
// zap path. NATS-sourced events already arrive as JSON, so they skip this.
func fieldsToMap(fields []zapcore.Field) map[string]interface{} {
	fieldMap := make(map[string]interface{})
	for _, f := range fields {
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
	return fieldMap
}

// LogEventPayload is the JSON shape any service publishes to the "logger"
// subject for centralized DB logging over NATS. Service is required here —
// unlike a per-service subject, a single shared subject can't imply it.
type LogEventPayload struct {
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error,omitempty"`
	Duration  string                 `json:"duration,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
}

// HandleLogger is a plain nats.MsgHandler — wire it into your central
// SubscribeEvent() map under the "logger" subject:
//
//	handlers := map[string]nats.MsgHandler{
//	    "user.notification": handleMsg,
//	    "logger":            handlers.HandleLogger,
//	}
//
// It does NOT call GORM directly. It pushes onto LogQueue, and the single
// goroutine started by InitLogger (StartLogWorker) is what actually calls
// database.GDB.Create(&entry) — keeping exactly one writer regardless of how
// many services are publishing logs concurrently.
func HandleLogger(msg *nats.Msg) {
	var payload LogEventPayload
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		fmt.Fprintf(os.Stderr, "invalid log event: %v\n", err)
		return
	}

	fieldMap := payload.Fields
	if fieldMap == nil {
		fieldMap = make(map[string]interface{})
	}
	fieldMap["service"] = payload.Service
	if payload.Error != "" {
		fieldMap["error"] = payload.Error
	}
	if payload.Duration != "" {
		fieldMap["duration"] = payload.Duration
	}

	ts := time.Now()
	if payload.Timestamp != "" {
		if parsed, err := time.Parse(time.RFC3339, payload.Timestamp); err == nil {
			ts = parsed
		}
	}

	entry := logEntry{
		Level:     payload.Level,
		Message:   payload.Message,
		Timestamp: ts,
		FieldMap:  fieldMap,
	}

	if LogQueue == nil {
		// InitLogger was never called (queue not set up) — write directly
		// so the event isn't silently lost instead of panicking on a nil channel.
		writeLogEntry(context.Background(), entry)
		return
	}

	select {
	case LogQueue <- entry:
	default:
		fmt.Fprintf(os.Stderr, "log queue full, dropping event from %s\n", payload.Service)
	}
}

func StartLogWorker(ctx context.Context, queue <-chan logEntry) {
	for {
		select {
		case <-ctx.Done():
			return

		case log := <-queue:
			writeLogEntry(ctx, log)
		}
	}
}

func writeLogEntry(ctx context.Context, log logEntry) {
	fieldMap := log.FieldMap
	if fieldMap == nil {
		fieldMap = make(map[string]interface{})
	}
	if _, exists := fieldMap["message"]; !exists {
		fieldMap["message"] = log.Message
	}

	jsonFields, err := json.Marshal(fieldMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log fields: %v\n", err)
		jsonFields = []byte("{}")
	}

	service := ""
	if v, ok := fieldMap["service"]; ok {
		service = fmt.Sprint(v)
	}

	entry := models.Logger{
		Level:     log.Level,
		Service:   service,
		Method:    fmt.Sprint(fieldMap["method"]),
		Path:      fmt.Sprint(fieldMap["path"]),
		UserAuth:  fmt.Sprint(fieldMap["user"]),
		Timestamp: log.Timestamp,
		Fields:    datatypes.JSON(jsonFields),
	}

	// status can arrive as int64 (local zap path, via zapcore.Int64Type) or
	// float64 (any NATS/JSON path — encoding/json always decodes numbers
	// into interface{} as float64, never int64). Missing the float64 case
	// meant every status code published over NATS silently became 0.
	if v, ok := fieldMap["status"]; ok {
		switch sv := v.(type) {
		case int64:
			entry.Status = uint16(sv)
		case float64:
			entry.Status = uint16(sv)
		case int:
			entry.Status = uint16(sv)
		}
	}

	// duration can arrive two ways: float64 milliseconds (from the local zap
	// path, via zapcore.DurationType) or a Go duration string like "12.5ms"
	// (from NATS payloads, which send time.Duration.String()). The old code
	// only handled the numeric cases, so every NATS-sourced duration was
	// silently recorded as 0.
	var durationVal float64
	if v, ok := fieldMap["duration"]; ok {
		switch d := v.(type) {
		case float64:
			durationVal = d
		case int64:
			durationVal = float64(d)
		case string:
			if parsed, err := time.ParseDuration(d); err == nil {
				durationVal = float64(parsed.Milliseconds())
			}
		}
	}
	entry.Duration = durationVal

	if err := database.GDB.WithContext(ctx).Create(&entry).Error; err != nil {
		fmt.Fprintf(os.Stderr, "GORM log insert failed: %v\n", err)
	}
}

func InitLogger(ctx context.Context) error {
	LogQueue = make(chan logEntry, 1000)
	go StartLogWorker(ctx, LogQueue)

	pgxCore := NewAsyncPGXCore(zapcore.InfoLevel, LogQueue)
	ZapLogger = zap.New(pgxCore)

	return nil
}
