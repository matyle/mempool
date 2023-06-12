package mempool

import (
	"bytes"
	"sync"
)

// reuse the memory slice
// and never gc

// Pool represents a buffer pool that can be used to reduce allocations and
// improve performance when dealing with byte buffers.
type Pool struct {
	// pool is a channel that holds reusable byte buffers.
	// Using a channel ensures that the pool is concurrency-safe.
	pool chan *bytes.Buffer

	// routineSize represents the number of goroutines that can be used to
	// access the pool.
	routineSize int

	// bufferCap is the estimated size of the buffers that will be stored
	// in the pool. If a buffer exceeds this size, it will be resized and
	// added to the pool for future use.
	bufferCap int

	// lock is a read-write mutex that can be used to safely modify the
	// pool attributes.
	lock *sync.Mutex

	// len is the number of buffers currently in the pool.
	len int
}

// NewPool creates a new Pool object with a given routineSize and cap.
func NewPool(routineSize, cap int) *Pool {
	// Create a channel that can hold *bytes.Buffer objects with a capacity of routineSize.
	pool := make(chan *bytes.Buffer, routineSize)
	// Create a Mutex object to manage access to the channel.
	lock := &sync.Mutex{}
	// Create a new Pool object with the given parameters.
	return &Pool{
		pool:        pool,
		lock:        lock,
		bufferCap:   cap,
		routineSize: routineSize,
	}
}

// Get retrieves a buffer from the pool. If the pool does not have
// any available buffers, it will create a new buffer and return it.
func (p *Pool) Get() *bytes.Buffer {
	// Acquire lock to ensure thread safety
	p.lock.Lock()
	defer p.lock.Unlock()

	// If the number of available buffers is less than the pool size,
	// create and return a new buffer.
	if p.len < p.routineSize {
		return bytes.NewBuffer(make([]byte, 0, p.bufferCap))
	}

	// Otherwise, retrieve a buffer from the pool and return it.
	buf := <-p.pool
	p.len--
	return buf
}

// Put adds a buffer to the pool.
// It resets the buffer before adding it to the pool.
func (p *Pool) Put(b *bytes.Buffer) {
	b.Reset()
	p.pool <- b

	// Increment the length of the pool and unlock
	// the mutex to allow other goroutines to access
	// the pool.
	//lock here ,cannot block the p.pool chan
	p.lock.Lock()
	p.len++
	p.lock.Unlock()
}

// Resize resizes the pool to a new maxSize. If the new size is smaller than the current size,
// it does nothing. If the new size is larger than the current size, then it creates a new pool
// with the desired size and adds elements from the current pool.
func (p *Pool) Resize(maxSize int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// do nothing if new size is smaller than the current size
	if maxSize < len(p.pool) {
		return
	}

	// create a new pool with the desired size
	newPool := make(chan *bytes.Buffer, maxSize)

	// add elements from current pool to new pool
	for i := 0; i < len(p.pool); i++ {
		newPool <- <-p.pool
	}

	// replace the current pool with the new pool
	p.pool = newPool
	p.routineSize = maxSize
}
