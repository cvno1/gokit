package log

import "go.uber.org/zap"

type Field interface {
	Key() string
	Value() any
}

type field struct {
	key   string
	value any
}

var _ Field = (*field)(nil)

func (f *field) Key() string {
	return f.key
}

func (f *field) Value() any {
	return f.value
}

func NewField(key string, value any) Field {
	return &field{
		key:   key,
		value: value,
	}
}

func Error(err error) Field {
	return &field{
		key:   "error",
		value: err,
	}
}

func toZapField(fields ...Field) []zap.Field {
	zfields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zfields = append(zfields, zap.Any(f.Key(), f.Value()))
	}
	return zfields
}
