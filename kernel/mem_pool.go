package kernel

import (
	"log"
)

type MemPool struct {
	pos uint32
	buf []uint8
}

const MAX_POOL_SIZE = 500 * 1024

func (pool *MemPool) GetMem(size uint32) []byte {

	if size >= MAX_POOL_SIZE {
		pool.buf = make([]uint8, size)
		pool.pos = size

		log.Println("warning: Gem Memory, too large, size=", size)

		return pool.buf
	}

	if MAX_POOL_SIZE-pool.pos < size {
		pool.pos = 0
		pool.buf = make([]uint8, MAX_POOL_SIZE)
	}
	b := pool.buf[pool.pos : pool.pos+size]
	pool.pos += size
	return b
}

func NewMemPool() *MemPool {
	return &MemPool{
		buf: make([]uint8, MAX_POOL_SIZE),
	}
}
