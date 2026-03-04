// Package scache provides a high-performance, thread-safe in-memory cache for Go.
// This file re-exports the public API from pkg/scache for backward compatibility.
package scache

import (
	"github.com/scache-io/scache/config"
	"github.com/scache-io/scache/errors"
	"github.com/scache-io/scache/pkg/api"
)

// Re-export types
type LocalCache = api.LocalCache
type EngineConfig = config.EngineConfig

// Re-export errors
var (
	ErrKeyEmpty        = errors.ErrKeyEmpty
	ErrInvalidArgument = errors.ErrInvalidArgument
	ErrTypeMismatch    = errors.ErrTypeMismatch
	ErrKeyNotFound     = errors.ErrKeyNotFound
	ErrFieldNotFound   = errors.ErrFieldNotFound
	ErrIndexOutOfRange = errors.ErrIndexOutOfRange
	ErrListEmpty       = errors.ErrListEmpty
)

// Re-export local cache API
var New = api.New

// Re-export global cache API
var (
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
