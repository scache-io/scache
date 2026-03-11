# Performance Benchmarks

Performance test results for scache on Apple M4 (ARM64).

> **GC Optimization (v0.1.4+)**: Object pooling with sync.Pool dramatically reduces memory allocations in high-throughput scenarios. See [GC Optimization Impact](#gc-optimization-impact) for details.

## Basic Operations

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Store String | 38,013 ns/op | 355 B/op | 3 allocs/op |
| Load String | 124.5 ns/op | 15 B/op | 1 allocs/op |
| Delete String | 37,524 ns/op | **27 B/op** | **1 alloc/op** |
| Store Struct | 39,660 ns/op | 493 B/op | 6 allocs/op |
| Load Struct | 510.6 ns/op | 326 B/op | 8 allocs/op |

## Concurrent Operations

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Concurrent Store | 39,717 ns/op | 199 B/op | 2 allocs/op |
| Concurrent Load | 148.8 ns/op | 16 B/op | 1 allocs/op |
| Concurrent Read/Write | 20,866 ns/op | 112 B/op | 2 allocs/op |

## LRU Eviction

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| LRU Eviction | 184.5 ns/op | 256 B/op | 6 allocs/op |
| LRU Eviction with Load | 203.8 ns/op | 257 B/op | 6 allocs/op |

## TTL Expiration

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| TTL Expiration | 38,315 ns/op | 352 B/op | 3 allocs/op |
| TTL Check | 124.1 ns/op | 15 B/op | 1 allocs/op |

## Data Types

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Store List | 38,608 ns/op | 616 B/op | 4 allocs/op |
| Load List | 149.7 ns/op | 103 B/op | 2 allocs/op |
| Store Hash | 39,972 ns/op | 678 B/op | 5 allocs/op |
| Load Hash | 265.3 ns/op | 359 B/op | 3 allocs/op |

## Large Data

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Store Large Struct (1KB) | 41,141 ns/op | 3,503 B/op | 6 allocs/op |
| Load Large Struct (1KB) | 5,424 ns/op | 3,079 B/op | 12 allocs/op |

## Mixed Workload

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Mixed Workload (50% read, 20% write, 10% delete, 10% exists, 10% struct) | 11,945 ns/op | 126 B/op | 2 allocs/op |

## Capacity Tests

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Unlimited Capacity | 361.8 ns/op | 352 B/op | 3 allocs/op |
| Large Capacity (100K) | 214.1 ns/op | 274 B/op | 5 allocs/op |

## Stress Tests

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| High Concurrency (100 workers, 30 ops each) | 43,773,772 ns/op | **25,589 B/op** | **1,222 allocs/op** |

## Code Generation

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Generic Generator | 13,483,180 ns/op | 47,452 B/op | 570 allocs/op |
| Classic Generator | 13,568,362 ns/op | 66,977 B/op | 724 allocs/op |

---

## GC Optimization Benchmarks

Benchmarks measuring GC optimization benefits with sync.Pool object reuse.

### High-Throughput Scenarios

| Test | Speed | Memory | Allocations |
|------|-------|--------|-------------|
| String Operations | 38.1 µs/op | 193 B/op | 4 allocs/op |
| List Operations | 38.6 µs/op | 396 B/op | 3 allocs/op |
| Hash Operations | 39.4 µs/op | 454 B/op | 4 allocs/op |
| High Delete Rate | 38.6 µs/op | **96 B/op** | **1 alloc/op** |
| Mixed Data Types | 39.3 µs/op | 335 B/op | 4 allocs/op |

### Memory Efficiency

| Test | Speed | Memory | Allocations |
|------|-------|--------|-------------|
| Stress Test (1000 ops) | 414.6 µs/op | **2,032 B/op** | **38 allocs/op** |
| Memory Usage | 38.7 µs/op | 247 B/op | 4 allocs/op |
| Long Running | 38.2 µs/op | 209 B/op | 3 allocs/op |
| Concurrent Read/Write | 19.6 µs/op | 107 B/op | 2 allocs/op |

### Object Pool Efficiency

| Test | Speed | Memory | Allocations |
|------|-------|--------|-------------|
| Concurrent Delete | **57.6 ns/op** | 13 B/op | **1 alloc/op** |
| Preallocated Map | **159.0 ns/op** | 167 B/op | 2 allocs/op |
| Flush (1000 items) | 18.9 µs/op | 19,126 B/op | 315 allocs/op |

---

## GC Optimization Impact

Comparison before and after sync.Pool implementation:

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Delete operation memory | 187 B/op | 27 B/op | **86% reduction** |
| Delete operation allocs | 3 allocs/op | 1 alloc/op | **67% reduction** |
| High concurrency memory | 187,815 B/op | 25,589 B/op | **86% reduction** |
| High concurrency allocs | 3,230 allocs/op | 1,222 allocs/op | **62% reduction** |
| Store operation allocs | 4 allocs/op | 3 allocs/op | **25% reduction** |
| Concurrent store allocs | 3 allocs/op | 2 allocs/op | **33% reduction** |

### Key Benefits

1. **Reduced GC Pressure**: Object reuse eliminates 60-86% of allocations
2. **Lower Memory Footprint**: Up to 86% reduction in memory usage
3. **Better Performance**: Fewer allocations = faster operations
4. **Stable Under Load**: High concurrency scenarios benefit most

---

## Running Benchmarks

```bash
# Run all benchmarks
go test ./tests/... -bench=. -benchmem

# Run specific benchmark
go test ./tests/... -bench=BenchmarkDeleteString -benchmem

# Run GC-specific benchmarks
go test ./tests/... -bench=BenchmarkGC -benchmem

# Run with custom duration
go test ./tests/... -bench=. -benchmem -benchtime=5s
```

## Key Insights

1. **Fast Reads**: Load operations are extremely fast (124-510 ns/op)
2. **Efficient Memory**: Low memory footprint for most operations
3. **Concurrent Performance**: Good parallel performance with minimal contention
4. **LRU Efficiency**: LRU eviction adds minimal overhead (184 ns/op)
5. **Scalability**: Handles large capacities efficiently
6. **GC Optimized**: Object pooling dramatically reduces allocation overhead
