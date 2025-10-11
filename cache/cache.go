package cache

// New 创建一个新的缓存实例
func New(opts ...Option) Cache {
	return NewMemoryCache(opts...)
}

// NewLRU 创建一个使用 LRU 策略的缓存
func NewLRU(maxSize int, opts ...Option) Cache {
	allOpts := append([]Option{WithMaxSize(maxSize), WithEvictionPolicy("lru")}, opts...)
	return NewMemoryCache(allOpts...)
}

// NewLFU 创建一个使用 LFU 策略的缓存
func NewLFU(maxSize int, opts ...Option) Cache {
	allOpts := append([]Option{WithMaxSize(maxSize), WithEvictionPolicy("lfu")}, opts...)
	return NewMemoryCache(allOpts...)
}

// NewFIFO 创建一个使用 FIFO 策略的缓存
func NewFIFO(maxSize int, opts ...Option) Cache {
	allOpts := append([]Option{WithMaxSize(maxSize), WithEvictionPolicy("fifo")}, opts...)
	return NewMemoryCache(allOpts...)
}
