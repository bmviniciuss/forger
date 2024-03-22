package ctx

import (
	"context"
)

type key string

var (
	reqIdKey = key("request_id")
)

func WithRequestID(c context.Context, id string) context.Context {
	return context.WithValue(c, reqIdKey, id)
}

func GetRequestID(c context.Context) (string, bool) {
	val, ok := c.Value(reqIdKey).(string)
	return val, ok
}
