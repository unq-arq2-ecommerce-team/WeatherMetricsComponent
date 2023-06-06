package logger

import (
	"context"
	"github.com/google/uuid"
)

type requestId string

const requestIdKeyStr = "system-request-id"

var requestIdKey = requestId(requestIdKeyStr)

// SetRequestId set given id as "system-request-id". If id is "", set a random uuid value
func SetRequestId(ctx context.Context, id string) context.Context {
	if id == "" {
		id = uuid.New().String()
	}
	return context.WithValue(ctx, requestIdKey, id)
}

// GetRequestId get context value from key "system-request-id"
func GetRequestId(ctx context.Context) string {
	if id, ok := ctx.Value(requestIdKey).(string); ok {
		return id
	}
	return ""
}

// RequestIdHeaderKey return key of send requestId header
func RequestIdHeaderKey() string {
	return requestIdKeyStr
}
