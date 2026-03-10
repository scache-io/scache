// Package scache provides a high-performance, thread-safe in-memory cache for Go.
// It supports multiple data types (String, List, Hash, Struct), TTL expiration,
// and LRU eviction policy.
//
// Quick Start:
//
//	// Local cache
//	cache := scache.New(config.MediumConfig...)
//	cache.SetString("key", "value", time.Hour)
//
//	// Global cache
//	scache.SetString("key", "value", time.Hour)
package scache

import (
	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/errors"
	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/pkg/api"
	"github.com/scache-io/scache/types"
)

// Public types
type (
	// LocalCache 局部缓存实例
	LocalCache = api.LocalCache

	// EngineConfig 缓存引擎配置
	EngineConfig = config.EngineConfig

	// DataObject 通用数据对象接口
	DataObject = interfaces.DataObject

	// StringObject 字符串对象接口
	StringObject = interfaces.StringObject

	// ListObject 列表对象接口
	ListObject = interfaces.ListObject

	// HashObject 哈希对象接口
	HashObject = interfaces.HashObject

	// StructObject 结构体对象接口
	StructObject = interfaces.StructObject

	// DataType 数据类型
	DataType = interfaces.DataType
)

// Public errors
var (
	ErrKeyEmpty        = errors.ErrKeyEmpty
	ErrInvalidArgument = errors.ErrInvalidArgument
	ErrTypeMismatch    = errors.ErrTypeMismatch
	ErrKeyNotFound     = errors.ErrKeyNotFound
	ErrFieldNotFound   = errors.ErrFieldNotFound
	ErrIndexOutOfRange = errors.ErrIndexOutOfRange
	ErrListEmpty       = errors.ErrListEmpty
)

// Public constants
const (
	DataTypeString = interfaces.DataTypeString
	DataTypeList   = interfaces.DataTypeList
	DataTypeHash   = interfaces.DataTypeHash
	DataTypeStruct = interfaces.DataTypeStruct
)

// Local cache API
var (
	New             = api.New
	GetGlobalCache  = api.GetGlobalCache
	InitGlobalCache = api.InitGlobalCache
	SetString       = api.SetString
	GetString       = api.GetString
	SetList         = api.SetList
	GetList         = api.GetList
	SetHash         = api.SetHash
	GetHash         = api.GetHash
	Store           = api.Store
	Load            = api.Load
	Delete          = api.Delete
	Exists          = api.Exists
	Keys            = api.Keys
	Flush           = api.Flush
	Size            = api.Size
	Expire          = api.Expire
	TTL             = api.TTL
	Stats           = api.Stats
)

// Config helpers
var (
	DefaultEngineConfig = config.DefaultEngineConfig
)

// Type constructors
var (
	NewStringObject = types.NewStringObject
	NewListObject   = types.NewListObject
	NewHashObject   = types.NewHashObject
	NewStructObject = types.NewStructObject
)
