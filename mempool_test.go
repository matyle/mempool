package mempool

import (
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	poolSize := 5
	bufferCap := 1024
	pool := NewPool(poolSize, bufferCap)

	// Test Get and Put
	for i := 0; i < poolSize*2; i++ {
		go func() {
			buffer := pool.Get()
			if buffer == nil {
				t.Errorf("Expected buffer to be not nil")
			}
			if buffer.Cap() != bufferCap {
				t.Errorf("Expected buffer capacity to be %d, but got %d", bufferCap, buffer.Cap())
			}

			pool.Put(buffer)
		}()
	}

	time.Sleep(5 * time.Second)

	// Test Resize
	newSize := 10
	pool.Resize(newSize)
	if pool.GetPoolSize() != newSize {
		t.Errorf("Expected pool size to be %d, but got %d", newSize, pool.GetPoolSize())
	}
}
