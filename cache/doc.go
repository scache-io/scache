/*
Package cache provides high-performance caching implementations with support for
multiple eviction policies (LRU, LFU, FIFO), TTL-based expiration, and concurrent access.

# Features

- Sharded Design: The cache is divided into multiple shards to reduce lock contention and improve concurrent performance
- Pluggable Eviction Policies: Support for LRU (Least Recently Used), LFU (Least Frequently Used), and FIFO (First In First Out) strategies
- TTL Support: Time-to-Live based expiration with lazy cleanup
- Global Cache Management: Singleton pattern for managing multiple named caches
- Comprehensive Statistics: Hit rates, operation counts, and performance metrics
- Thread Safety: All operations are safe for concurrent use

# Basic Usage

Create a cache with default settings:

	c := cache.New()
	c.Set("key", "value")
	value, found := c.Get("key")

Create a cache with specific eviction policy:

	lruCache := cache.NewLRU(1000)        // LRU with max 1000 items
	lfuCache := cache.NewLFU(1000)        // LFU with max 1000 items
	fifoCache := cache.NewFIFO(1000)      // FIFO with max 1000 items

Create a cache with custom configuration:

	c := cache.New(
		cache.WithMaxSize(10000),
		cache.WithShards(32),
		cache.WithEvictionPolicy("lru"),
		cache.WithDefaultTTL(5*time.Minute),
		cache.WithCleanupInterval(1*time.Minute),
		cache.WithStatistics(true),
	)

# Advanced Usage

Global Cache Management:

	// Register named caches
	cache.RegisterLRU("users", 1000)
	cache.RegisterLFU("sessions", 500)

	// Use global caches
	usersCache, _ := cache.Get("users")
	usersCache.Set("user:123", userData)

	// Use default cache
	cache.Set("app:version", "1.0.0")
	version, _ := cache.GetFromDefault("app:version")

# Performance Considerations

- Choose appropriate shard count based on expected concurrency
- Use TTL to automatically expire stale data
- Enable statistics only when needed (small performance overhead)
- Consider memory usage when setting max cache size
- Use batch operations for bulk operations

# Thread Safety

All cache operations are thread-safe. The cache uses sharding to minimize lock contention
and provides atomic operations for statistics updates.

# Configuration Options

The cache supports the following configuration options via functional options:

- WithMaxSize(int): Maximum number of cache items
- WithShards(int): Number of cache shards for concurrency
- WithEvictionPolicy(string): Eviction strategy ("lru", "lfu", "fifo")
- WithDefaultTTL(time.Duration): Default TTL for cache items
- WithCleanupInterval(time.Duration): Interval for cleanup goroutine
- WithStatistics(bool): Enable detailed statistics
- WithSerializer(string): Serialization method for persistence
- WithLazyExpiration(bool): Enable lazy expiration checking
- WithMetrics(bool): Enable detailed performance metrics

# Error Handling

The cache returns errors for the following conditions:

- Invalid configuration parameters
- Serialization/deserialization failures
- Cache closed errors
- Invalid key/value parameters

All methods that can return errors should be checked in production code.

# Statistics

The cache provides comprehensive statistics including:

- Hit/miss counts and rates
- Current cache size
- Creation and last access times
- Eviction counts (when applicable)

Example:

	stats := c.Stats()
	fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
	fmt.Printf("Current size: %d/%d\n", stats.Size, stats.MaxSize)

*/
package cache