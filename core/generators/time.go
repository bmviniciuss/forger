package generators

import (
	"context"
	"errors"
	"time"
)

const (
	utcLayout = "2006-01-02T15:04:05.000Z"
)

type TimeType string

var (
	TimeTypeIso8601 TimeType = "iso8601"
	TimeTypeRfc3339 TimeType = "rfc3339"
)

func (t TimeType) Format() string {
	switch t {
	case TimeTypeIso8601:
		return utcLayout
	case TimeTypeRfc3339:
		return time.RFC3339
	default:
		return utcLayout
	}
}

func NewTimeType(timeType string) TimeType {
	switch timeType {
	case "iso8601":
		return TimeTypeIso8601
	case "rfc3339":
		return TimeTypeRfc3339
	default:
		return TimeTypeIso8601
	}
}

// ctx type now
func Time(ctx context.Context, options ...interface{}) (string, error) {
	// Signature Time(ctx)
	if len(options) <= 0 {
		return time.Now().Format(utcLayout), nil
	}

	// Signature UUID(ctx, type)
	tRaw, ok := options[0].(string)
	if !ok {
		return "", errors.New("invalid type for time function")
	}
	t := NewTimeType(tRaw)
	return time.Now().Format(t.Format()), nil
}
