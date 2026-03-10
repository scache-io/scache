# Performance Benchmarks

Performance test results for scache on Apple M4 (ARM64).

## Basic Operations

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Store String | 37,895 ns/op | 355 B/op | 4 allocs/op |
| Load String | 126.0 ns/op | 15 B/op | 1 allocs/op |
| Delete String | 38,674 ns/op | 187 B/op | 3 allocs/op |
| Store Struct | 39,953 ns/op | 487 B/op | 7 allocs/op |
| Load Struct | 523.2 ns/op | 326 B/op | 8 allocs/op |

## Concurrent Operations

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Concurrent Store | 38,807 ns/op | 198 B/op | 3 allocs/op |
| Concurrent Load | 154.7 ns/op | 16 B/op | 1 allocs/op |
| Concurrent Read/Write | 21,051 ns/op | 111 B/op | 2 allocs/op |

## LRU Eviction

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| LRU Eviction | 172.3 ns/op | 256 B/op | 7 allocs/op |
| LRU Eviction with Load | 189.7 ns/op | 257 B/op | 7 allocs/op |

## TTL Expiration

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| TTL Expiration | 38,637 ns/op | 357 B/op | 4 allocs/op |
| TTL Check | 128.8 ns/op | 15 B/op | 1 allocs/op |

## Data Types

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Store List | 38,656 ns/op | 376 B/op | 4 allocs/op |
| Load List | 160.8 ns/op | 103 B/op | 2 allocs/op |
| Store Hash | 38,569 ns/op | 359 B/op | 4 allocs/op |
| Load Hash | 289.5 ns/op | 359 B/op | 3 allocs/op |

## Large Data

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Store Large Struct (1KB) | 40,903 ns/op | 3,497 B/op | 7 allocs/op |
| Load Large Struct (1KB) | 5,453 ns/op | 3,079 B/op | 12 allocs/op |

## Mixed Workload

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Mixed Workload (50% read, 20% write, 10% delete, 10% exists, 10% struct) | 11,939 ns/op | 126 B/op | 2 allocs/op |

## Capacity Tests

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Unlimited Capacity | 349.2 ns/op | 348 B/op | 4 allocs/op |
| Large Capacity (100K) | 203.1 ns/op | 274 B/op | 6 allocs/op |

## Stress Tests

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| High Concurrency (100 workers, 30 ops each) | 43,759,411 ns/op | 185,910 B/op | 3,223 allocs/op |

## Code Generation

| Operation | Speed | Memory | Allocations |
|-----------|-------|--------|-------------|
| Generic Generator | 13,508,133 ns/op | 47,440 B/op | 570 allocs/op |
| Classic Generator | 13,704,580 ns/op | 66,980 B/op | 724 allocs/op |

## Running Benchmarks

```bash
# Run all benchmarks
go test ./tests/... -bench=. -benchmem

# Run specific benchmark
go test ./tests/... -bench=BenchmarkStoreString -benchmem

# Run with custom duration
go test ./tests/... -bench=. -benchmem -benchtime=5s
```

## Key Insights

1. **Fast Reads**: Load operations are extremely fast (126-523 ns/op)
2. **Efficient Memory**: Low memory footprint for most operations
3. **Concurrent Performance**: Good parallel performance with minimal contention
4. **LRU Efficiency**: LRU eviction adds minimal overhead (172 ns/op)
5. **Scalability**: Handles large capacities efficiently
