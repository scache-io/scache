# GC Optimization Implementation

## Overview

This document describes the GC optimization implementation for scache, which uses object pooling and pre-allocated data structures to reduce memory allocations and GC pressure.

## Changes Made

### 1. Object Pooling (types/data_objects.go)

Added sync.Pool support for all data object types:

- **stringObjectPool**: Reuses StringObject instances
- **listObjectPool**: Reuses ListObject instances
- **hashObjectPool**: Reuses HashObject instances

Each pool has a New function that creates pre-initialized objects with sensible defaults:
- ListObject: Pre-allocates slice capacity of 16 elements
- HashObject: Pre-allocates empty map

### 2. Reset/Clear Methods

Added methods to all data objects for reuse:

- `Reset()`: Resets object state for reuse (internal)
- `Clear()`: Public method that calls Reset()

These methods:
- Clear data values (strings, slices, maps)
- Reset BaseObject fields (dataType, expiresAt, created, accessed)
- Don't deallocate underlying memory, allowing reuse

### 3. Storage Engine Enhancements (storage/engine.go)

#### Pre-allocated Map Capacity

```go
initialCapacity := 64
if engineConfig.MaxSize > 0 && engineConfig.MaxSize < 10000 {
    initialCapacity = engineConfig.MaxSize
}
engine := &StorageEngine{
    data: make(map[string]interfaces.DataObject, initialCapacity),
    // ...
}
```

This reduces map reallocations during initial population.

#### Object Pool Integration

Added `returnObjectToPool()` method that:
- Returns deleted/evicted objects to appropriate pool
- Records pool hit statistics
- Handles unsupported object types gracefully

Modified methods to use object pooling:
- `Delete()`: Returns objects to pool before deletion
- `deleteExpired()`: Returns expired objects to pool
- `evictOne()`: Returns evicted objects to pool
- `cleanupExpired()`: Returns expired objects to pool
- `Flush()`: Returns all objects to pool before clearing

### 4. Enhanced Statistics

Added new fields to EngineStats:
- `gcCycles`: Tracks GC cycle count
- `poolHits`: Tracks objects returned to pool
- `poolAllocs`: Tracks new object allocations
- `lastGCTime`: Timestamp of last GC

Enhanced Stats() output with:
- GC metrics (gc_cycles, pool_hits, pool_allocs)
- Runtime memory stats (heap_alloc, heap_sys, num_gc, gc_cpu_frac)

## Benefits

### Memory Allocation Reduction

Object pooling significantly reduces allocations by reusing objects instead of creating new ones. This is especially effective for:
- High-churn workloads (frequent deletes/evictions)
- TTL-based expiration scenarios
- Large-scale cache operations

### GC Pressure Reduction

By reusing objects:
- Fewer objects for GC to scan
- Shorter GC pause times
- Lower CPU overhead from GC

### Performance Improvements

Pre-allocated map capacity:
- Reduces map reallocations
- Improves cache locality
- Faster initial population

## Benchmark Results

Example benchmark results (Apple M4):

```
BenchmarkGCStringOperations-10     31570    38652 ns/op    209 B/op    5 allocs/op
BenchmarkGCListOperations-10       31486    38442 ns/op    197 B/op    3 allocs/op
BenchmarkGCHashOperations-10       30998    38690 ns/op    181 B/op    3 allocs/op
BenchmarkGCConcurrentDelete-10  20602346      56.83 ns/op     13 B/op    1 allocs/op
BenchmarkGCPreallocatedMap-10    8087797     148.5 ns/op    167 B/op    3 allocs/op
```

Key observations:
- Very low allocation counts (1-5 allocs/op)
- Efficient memory usage (13-209 B/op)
- Excellent concurrent performance (56.83 ns/op)

## Backward Compatibility

All changes maintain backward compatibility:
- Public API unchanged
- Existing code works without modifications
- Pool usage is internal to the engine
- New stats are additive, not breaking

## Usage Recommendations

### Best Workloads

Object pooling is most effective for:
1. High-frequency delete operations
2. TTL-based expiration with frequent renewal
3. LRU eviction scenarios
4. Mixed data type operations

### Monitoring

Use the Stats() output to monitor pool effectiveness:
```go
stats := cache.Stats()
poolHits := stats["pool_hits"].(int64)
poolAllocs := stats["pool_allocs"].(int64)
hitRate := float64(poolHits) / float64(poolHits + poolAllocs)
```

A high pool hit rate (>80%) indicates effective pool usage.

## Testing

Run the new GC-specific benchmarks:
```bash
go test ./tests/... -bench=BenchmarkGC -benchmem
```

Run all tests to ensure compatibility:
```bash
go test ./tests/... -run "Test(String|List|Hash|Struct|TTL|Expire|Concurrent|MaxSize|Global|Stats)" -v
```

## Future Improvements

Potential enhancements:
1. Pool size limits to prevent unbounded growth
2. Pool statistics per object type
3. Configurable pool behavior
4. Pool warmup strategies
5. Adaptive pool sizing based on workload patterns
