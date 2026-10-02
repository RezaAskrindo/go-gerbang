package types

import "time"

type ConsumerMetrics struct {
	Name         string `json:"name"`
	Pending      uint64 `json:"pending"`
	AckPending   int    `json:"ack_pending"`
	Redelivered  int    `json:"redelivered"`
	Waiting      int    `json:"waiting"`
	DeliveredSeq uint64 `json:"delivered_seq"`
	AckFloorSeq  uint64 `json:"ack_floor_seq"`
}

type StreamMetrics struct {
	Name      string            `json:"name"`
	Subjects  []string          `json:"subjects"`
	Messages  uint64            `json:"messages"`
	Bytes     uint64            `json:"bytes"`
	FirstSeq  uint64            `json:"first_seq"`
	LastSeq   uint64            `json:"last_seq"`
	Created   time.Time         `json:"created"`
	Consumers []ConsumerMetrics `json:"consumers"`
	Schedules *ScheduleMetrics  `json:"schedules"`
}

type DetailedMetrics struct {
	MemoryUsed     uint64          `json:"memory_used"`
	StoreUsed      uint64          `json:"store_used"`
	MaxMemory      int64           `json:"max_memory"`
	MaxStore       int64           `json:"max_store"`
	TotalStreams   int             `json:"total_streams"`
	TotalConsumers int             `json:"total_consumers"`
	Streams        []StreamMetrics `json:"streams"`
}

type ScheduledMessage struct {
	Subject  string    `json:"subject"`
	Seq      uint64    `json:"seq"`
	Schedule string    `json:"schedule"`
	Target   string    `json:"target,omitempty"`
	Source   string    `json:"source,omitempty"`
	TTL      string    `json:"ttl,omitempty"`
	TimeZone string    `json:"time_zone,omitempty"`
	Stored   time.Time `json:"stored"`
}

type ScheduleMetrics struct {
	Count     int                `json:"count"`
	Truncated bool               `json:"truncated"`
	Items     []ScheduledMessage `json:"items"`
}
