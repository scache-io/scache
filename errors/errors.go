package errors

import "errors"

// Error定义
var (
	// ErrKeyEmpty 键为空Error
	ErrKeyEmpty = errors.New("key cannot be empty")

	// ErrInvalidArgument 无效ParameterError
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrTypeMismatch Type不匹配Error
	ErrTypeMismatch = errors.New("type mismatch")

	// ErrKeyNotFound 键不存在Error
	ErrKeyNotFound = errors.New("key not found")

	// ErrFieldNotFound 字段不存在Error
	ErrFieldNotFound = errors.New("field not found")

	// ErrIndexOutOfRange 索引超出范围Error
	ErrIndexOutOfRange = errors.New("index out of range")

	// ErrListEmpty 列表为空Error
	ErrListEmpty = errors.New("list is empty")
)
