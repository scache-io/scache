/*
Package constants defines all configuration constants, error messages, and default values
used throughout the SCache framework.

# Constants Categories

## Configuration Constants

Default values for cache configuration:

	DefaultMaxSize         = 10000    // Default maximum cache size
	DefaultShards          = 16       // Default number of shards
	DefaultCleanupInterval = 10m      // Default cleanup interval
	DefaultTTL             = 0        // Default TTL (0 = no expiration)

## Strategy Constants

Identifiers for different eviction strategies:

	LRUStrategy  = "lru"  // Least Recently Used
	LFUStrategy  = "lfu"  // Least Frequently Used
	FIFOStrategy = "fifo" // First In First Out

## Global Cache Constants

Constants for global cache management:

	DefaultCacheName = "default" // Name of the default cache
	ManagerTimeout   = 30s      // Timeout for manager operations

## Performance Constants

Performance-related limits and thresholds:

	MaxKeyLength  = 256  // Maximum key length in characters
	MinKeyLength  = 1    // Minimum key length in characters
	MaxValueSize  = 10MB // Maximum value size in bytes

## Error Message Constants

Standardized error messages for consistent error reporting:

	ErrCacheNotFound      = "cache not found"
	ErrInvalidCacheName   = "invalid cache name"
	ErrCacheAlreadyExists = "cache already exists"
	ErrInvalidStrategy    = "invalid eviction strategy"
	ErrKeyNotFound        = "key not found"
	ErrKeyTooLong         = "key too long"
	ErrKeyEmpty           = "key empty"
	ErrValueTooLarge      = "value too large"
	ErrCacheClosed        = "cache is closed"

## Logging Constants

Log message prefixes for different components:

	LogPrefixCache   = "[SCache]"
	LogPrefixManager = "[SCache-Manager]"
	LogPrefixGlobal  = "[SCache-Global]"

## Statistics Constants

Constants for statistics and monitoring:

	StatsUpdateInterval = 1s   // Interval for updating statistics
	HitRateThreshold   = 0.8  // Threshold for good hit rate

## Serialization Constants

Supported serialization formats:

	JSONEncoding = "json" // JSON serialization
	GobEncoding  = "gob"  // Go gob encoding

# Usage

These constants are intended for use throughout the SCache framework to ensure
consistency in configuration, error messages, and logging. They should not be
modified at runtime.

Example:

	// Use strategy constants
	policy := constants.LRUStrategy

	// Use error constants for consistent error messages
	return errors.New(constants.ErrKeyNotFound)

	// Use default configuration values
	config.MaxSize = constants.DefaultMaxSize

*/
package constants