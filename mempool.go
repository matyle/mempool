package mempool

import "bytes"

// reuse the memory slice
type Pool struct {
	pool        chan *bytes.Buffer
	routineSize int
	off         int
	// slice cp
	bufferCap int
}

func NewPool(routineSize, cap int) *Pool {
	p := &Pool{
		pool: make(chan *bytes.Buffer, routineSize),
	}
	p.bufferCap = cap
	p.routineSize = routineSize
	return p
}

func (p *Pool) Get() *bytes.Buffer {
	// if pool chan size less than , create a new one
	if len(p.pool) < p.routineSize {
		return bytes.NewBuffer(make([]byte, 0, p.bufferCap))
	}
	p.off++
	return <-p.pool
}

func (p *Pool) Put(b *bytes.Buffer) {
	b.Reset()
	select {
	case p.pool <- b:
	}
	p.off--
}
