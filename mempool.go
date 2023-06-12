package mempool

import (
	"bytes"
	"sync"
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

// reuse the memory slice
// and never gc

// Pool represents a buffer pool that can be used to reduce allocations and
// improve performance when dealing with byte buffers.
type Pool struct {
	// pool is a channel that holds reusable byte buffers.
	// Using a channel ensures that the pool is concurrency-safe.
	pool chan *bytes.Buffer

	// poolSize represents the number of goroutines that can be used to
	// access the pool.
	poolSize int

	// bufferInitCap represents the capacity of the one buffer  that will be
	// added to the pool for future use.
	bufferInitCap int

	// maxBufferSize represents the maximum size of the buffer
	// if the buffer is larger than this size, it will not be added to the pool.
	maxBufferSize int

	// lock is a read-write mutex that can be used to safely modify the
	// pool attributes.
	lock *sync.RWMutex

	// len is the number of buffers currently in the pool.
	len int
}

// NewPool creates a new Pool object with a given routineSize and cap.
func NewPool(routineSize, capacity int) *Pool {
	// Create a channel that can hold *bytes.Buffer objects with a capacity of routineSize.
	pool := make(chan *bytes.Buffer, routineSize)
	// Create a Mutex object to manage access to the channel.
	lock := &sync.RWMutex{}
	// Create a new Pool object with the given parameters.
	return &Pool{
		pool:          pool,
		lock:          lock,
		bufferInitCap: capacity,
		poolSize:      routineSize,
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
	if p.len < p.poolSize {
		return bytes.NewBuffer(make([]byte, 0, p.bufferInitCap))
	}

	// Otherwise, retrieve a buffer from the pool and return it.
	buf := <-p.pool
	p.len--
	return buf
}

// Put adds a buffer to the pool.
//
// If the buffer's capacity is at least twice its length, a new buffer
// is created with the same length as the input buffer and added to the pool.
// Otherwise, the input buffer is reset and added to the pool.
//
// Finally, the length of the pool is incremented and the mutex is unlocked
// to allow other goroutines to access the pool.
func (p *Pool) Put(b *bytes.Buffer) {
	cap := b.Cap()
	len := b.Len()

	if len != 0 && cap >= 2*len {
		// Create a new buffer with the same length as the input buffer
		newBuf := bytes.NewBuffer(make([]byte, 0, len))
		p.pool <- newBuf
	} else {
		// Reset the input buffer before adding it to the pool
		b.Reset()
		p.pool <- b
	}

	// Increment the length of the pool and unlock the mutex
	p.lock.Lock()
	p.len++
	p.lock.Unlock()
}

// Resize resizes the pool to a new maxSize. If the new size is smaller than the current size,
// it does nothing. If the new size is larger than the current size, then it creates a new pool
// with the desired size and adds elements from the current pool.
func (p *Pool) Resize(poolSize int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// do nothing if new size is smaller than the current size
	if poolSize < len(p.pool) {
		return
	}

	// create a new pool with the desired size
	newPool := make(chan *bytes.Buffer, poolSize)

	// add elements from current pool to new pool
	for i := 0; i < len(p.pool); i++ {
		newPool <- <-p.pool
	}

	// replace the current pool with the new pool
	p.pool = newPool
	p.poolSize = poolSize
}

func (p *Pool) GetLen() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.len
}

func (p *Pool) GetPoolSize() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.poolSize
}
