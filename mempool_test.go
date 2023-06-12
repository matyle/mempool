package mempool

import (
	"testing"
)

func TestPool(t *testing.T) {
	poolSize := 5
	bufferCap := 1024
	pool := NewPool(poolSize, bufferCap)

	// Test Get and Put
	for i := 0; i < poolSize*2; i++ {
		buf := pool.Get()
		if buf.Cap() != bufferCap {
			t.Errorf("Expected buffer capacity to be %d, but got %d", bufferCap, buf.Cap())
		}
		pool.Put(buf)
	}

	// Test Resize
	newSize := 10
	pool.Resize(newSize)
	if pool.GetPoolSize() != newSize {
		t.Errorf("Expected pool size to be %d, but got %d", newSize, pool.GetPoolSize())
	}
}
