package logs

import (
	"context"
	"sync"
)

type KeyVal struct {
	key interface{}
	val interface{}
}

type LogField struct {
	kvs       []KeyVal
	fieldLock sync.RWMutex
}

type kvsIdKey struct{}

func getFields(ctx context.Context) *LogField {
	field, ok := ctx.Value(kvsIdKey{}).(*LogField)
	if !ok {
		return nil
	}

	return field
}
