package logs

import (
	"context"
	"fmt"
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

func (f *LogField) AddField(key string, val interface{}) {
	f.fieldLock.Lock()
	f.kvs = append(f.kvs, KeyVal{key: key, val: val})
	f.fieldLock.Unlock()
}

type kvsIdKey struct{}

// 初始化context中key为kvsIdKey{}的字段（*LogField）
func WithFieldContext(ctx context.Context) context.Context {
	field := getFields(ctx)
	if field != nil {
		return ctx
	}
	return context.WithValue(ctx, kvsIdKey{}, &LogField{})
}

func getFields(ctx context.Context) *LogField {
	field, ok := ctx.Value(kvsIdKey{}).(*LogField)
	if !ok {
		return nil
	}

	return field
}

func AddField(ctx context.Context, key string, val interface{}) {
	field := getFields(ctx)
	if field == nil {
		fmt.Printf("field is null")
		return
	}
	field.AddField(key, val)
}
