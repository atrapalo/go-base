package mysql

import (
	"fmt"
	"strconv"
	"time"
)

type Field struct {
	value  interface{}
	dbType string
}

func NewField(v interface{}, dbType string) Field {
	return Field{value: v, dbType: dbType}
}

func (f Field) RawVal() interface{} {
	return f.value
}

func (f Field) IntVal() int {
	switch f.value.(type) {
	case int:
		return f.value.(int)
	case uint8:
		return int(f.value.(uint8))
	case uint16:
		return int(f.value.(uint16))
	case uint32:
		return int(f.value.(uint32))
	case uint64:
		return int(f.value.(uint64))
	case int8:
		return int(f.value.(int8))
	case int16:
		return int(f.value.(int16))
	case int32:
		return int(f.value.(int32))
	case int64:
		return int(f.value.(int64))
	case nil:
		return 0
	default:
		val, err := strconv.Atoi(f.StringVal())
		if err != nil {
			panic("unable to convert to int value")
		}

		return val
	}
}

func (f Field) Float32Val() float32 {
	switch f.value.(type) {
	case float64:
		return float32(f.value.(float64))
	case float32:
		return f.value.(float32)
	case nil:
		return .0
	}

	panic("unable to convert to float32 value")
}

func (f Field) Float64Val() float64 {
	switch f.value.(type) {
	case float64:
		return f.value.(float64)
	case float32:
		return float64(f.value.(float32))
	case nil:
		return .0
	}

	panic("unable to convert to float64 value")
}

func (f Field) StringVal() string {
	switch f.value.(type) {
	case nil:
		return ""
	default:
		return fmt.Sprintf("%s", f.value)
	}
}

func (f Field) BoolVal() bool {
	switch f.IntVal() {
	case 0:
		return false
	case 1:
		return true
	default:
		panic("unable to convert to bool value")
	}
}

func (f Field) TimeValue() *time.Time {
	var format string
	switch f.dbType {
	case "TIMESTAMP", "DATETIME":
		format = dateTimeFormat
		break
	case "DATE":
		format = dateFormat
	default:
		panic("unknown database type time format")
	}

	stringVal := f.StringVal()
	if stringVal == "" {
		return nil
	}

	value, err := time.Parse(format, stringVal)
	if err != nil {
		panic("unable to convert to time value")
	}

	return &value
}
