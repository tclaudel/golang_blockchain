package values

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Timestamp struct {
	timestamp time.Time
}

func TimestampNow() Timestamp {
	return Timestamp{timestamp: time.Now()}
}

func TimestampFromTime(timestamp time.Time) Timestamp {
	return Timestamp{timestamp: timestamp}
}

func (t Timestamp) String() string {
	return t.timestamp.Format(time.RFC3339)
}

func (t Timestamp) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("timestamp", t.String())
	return nil
}

func (t Timestamp) Time() time.Time {
	return t.timestamp
}

func (t Timestamp) Equal(other Timestamp) bool {
	return t.timestamp.Equal(other.timestamp)
}
