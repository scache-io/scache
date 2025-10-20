package main

import "errors"

// 错误定义
var (
	// ErrKeyEmpty 键为空错误
	ErrKeyEmpty = errors.New("key cannot be empty")

	// ErrInvalidArgument 无效参数错误
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrTypeMismatch 类型不匹配错误
	ErrTypeMismatch = errors.New("type mismatch")

	// ErrKeyNotFound 键不存在错误
	ErrKeyNotFound = errors.New("key not found")

	// ErrFieldNotFound 字段不存在错误
	ErrFieldNotFound = errors.New("field not found")

	// ErrIndexOutOfRange 索引超出范围错误
	ErrIndexOutOfRange = errors.New("index out of range")

	// ErrListEmpty 列表为空错误
	ErrListEmpty = errors.New("list is empty")
)
