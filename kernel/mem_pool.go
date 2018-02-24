package kernel

import (
	"log"
)

// MemPool memory pool
type MemPool struct {
	pos uint32
	buf []uint8
}

const maxPoolSize = 500 * 1024

// GetMem get memory space from pool
func (pool *MemPool) GetMem(size uint32) []byte {

	if size >= maxPoolSize {
		pool.buf = make([]uint8, size)
		pool.pos = size

		log.Println("warning: Gem Memory, too large, size=", size)

		return pool.buf
	}

	if maxPoolSize-pool.pos < size {
		pool.pos = 0
		pool.buf = make([]uint8, maxPoolSize)
	}
	b := pool.buf[pool.pos : pool.pos+size]
	pool.pos += size
	return b
}

// NewMemPool create a new memory pool
func NewMemPool() *MemPool {
	return &MemPool{
		buf: make([]uint8, maxPoolSize),
	}
}
