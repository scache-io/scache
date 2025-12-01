package internal

import (
	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/types"
)

// ExtractStringValue 从数据对象中提取字符串值
func ExtractStringValue(obj interfaces.DataObject) (string, bool) {
	if obj.Type() != interfaces.DataTypeString {
		return "", false
	}

	if strObj, ok := obj.(*types.StringObject); ok {
		return strObj.Value(), true
	}
	return "", false
}

// ExtractListValue 从数据对象中提取列表值
func ExtractListValue(obj interfaces.DataObject) ([]interface{}, bool) {
	if obj.Type() != interfaces.DataTypeList {
		return nil, false
	}

	if listObj, ok := obj.(*types.ListObject); ok {
		return listObj.Values(), true
	}
	return nil, false
}

// ExtractHashValue 从数据对象中提取哈希值
func ExtractHashValue(obj interfaces.DataObject) (map[string]interface{}, bool) {
	if obj.Type() != interfaces.DataTypeHash {
		return nil, false
	}

	if hashObj, ok := obj.(*types.HashObject); ok {
		return hashObj.Fields(), true
	}
	return nil, false
}

// ExtractStructValue 从数据对象中提取结构体值（JSON字符串）
func ExtractStructValue(obj interfaces.DataObject) (string, bool) {
	// 结构体对象底层是StringObject，所以检查字符串类型
	if obj.Type() != interfaces.DataTypeString {
		return "", false
	}

	if strObj, ok := obj.(*types.StringObject); ok {
		return strObj.Value(), true
	}
	return "", false
}

// IsDataTypeCompatible 检查数据类型是否兼容
func IsDataTypeCompatible(obj interfaces.DataObject, expectedType interfaces.DataType) bool {
	return obj.Type() == expectedType
}