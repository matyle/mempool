package mempool

import (
	"bytes"
	"fmt"
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

	for i := 0; i < routineSize; i++ {
		// Add a new buffer to the pool.
		pool <- bytes.NewBuffer(make([]byte, 0, capacity))
	}
	fmt.Println("pool size:", len(pool))
	// Create a new Pool object with the given parameters.
	return &Pool{
		pool:          pool,
		lock:          lock,
		bufferInitCap: capacity,
		poolSize:      routineSize,
	}
}

//Get and Put should be used in pairs

// Get retrieves a buffer from the pool. If the pool does not have
// any available buffers, it will create a new buffer and return it.
func (p *Pool) Get() *bytes.Buffer {
	// Acquire lock to ensure thread safety
	// If the number of available buffers is less than the pool size,
	// create and return a new buffer.
	// Otherwise, retrieve a buffer from the pool and return it.
	buf := <-p.pool
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
	// // cap := b.Cap()
	// // len := b.Len()
	//
	// if len != 0 && cap >= 2*len {
	// 	// Create a new buffer with the same length as the input buffer
	// 	newBuf := bytes.NewBuffer(make([]byte, 0, len))
	// 	// p.pool <- newBuf
	// 	select {
	// 	case p.pool <- newBuf:
	// 	default:
	// 		// Discard the buffer if the pool is full
	// 		log.Println("new pool is full")
	// 	}
	// } else {
	// 	// Reset the input buffer before adding it to the pool
	// 	b.Reset()
	// 	p.pool <- b
	// 	select {
	// 	case p.pool <- b:
	// 	default:
	// 		log.Println("b pool is full")
	// 	}
	// }
	//
	b.Reset()
	p.pool <- b
	// Increment the length of the pool and unlock the mutex
	// p.lock.Lock()
	// p.len++
	// p.lock.Unlock()
}

func (p *Pool) Len() int {
	return len(p.pool)
}

func (p *Pool) Cap() int {
	return cap(p.pool)
}
