package mempool

import "testing"

// func TestPool(t *testing.T) {
// 	poolSize := 5
// 	bufferCap := 1024
// 	pool := NewPool(poolSize, bufferCap)
//
// 	// Test Get and Put
// 	for i := 0; i < 100000; i++ {
// 		go func() {
// 			buffer := pool.Get()
// 			if buffer == nil {
// 				t.Errorf("Expected buffer to be not nil")
// 			}
// 			if buffer.Cap() != bufferCap {
// 				t.Errorf("Expected buffer capacity to be %d, but got %d", bufferCap, buffer.Cap())
// 			}
//
// 			pool.Put(buffer)
// 		}()
// 	}
//
// 	time.Sleep(5 * time.Second)
//
// 	// Test Resize
// 	newSize := 10
// 	pool.Resize(newSize)
// 	if pool.GetPoolSize() != newSize {
// 		t.Errorf("Expected pool size to be %d, but got %d", newSize, pool.GetPoolSize())
// 	}
// }

// benchmark pool and sync.Pool concurrency performance
// func BenchmarkPool(b *testing.B) {
// 	// mempool
// 	poolSize := 1000
// 	bufferCap := 1024
// 	pool := NewPool(poolSize, bufferCap)
//
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			buffer := pool.Get()
// 			defer func() {
// 				pool.Put(buffer)
// 			}()
// 			if buffer == nil {
// 				b.Errorf("Expected buffer to be not nil")
// 			}
// 			if buffer.Cap() != bufferCap {
// 				b.Errorf("Expected buffer capacity to be %d, but got %d", bufferCap, buffer.Cap())
// 			}
// 		}
// 	})
//
// }

// benchmark pool and sync.Pool concurrency performance
func BenchmarkPool(b *testing.B) {
	// mempool
	poolSize := 1000
	bufferCap := 1024
	pool := NewPool(poolSize, bufferCap)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			func() {
				buffer := pool.Get()
				defer func() {
					pool.Put(buffer)
				}()
				if buffer == nil {
					b.Errorf("Expected buffer to be not nil")
				}
				if buffer.Cap() != bufferCap {
					b.Errorf("Expected buffer capacity to be %d, but got %d", bufferCap, buffer.Cap())
				}
			}()
		}
	})

}

// func BenchmarkSyncPool(b *testing.B) {
// 	bufferCap := 1024
// 	// sync.Pool
// 	syncPool := &sync.Pool{
// 		New: func() interface{} {
// 			return make([]byte, bufferCap)
// 		},
// 	}
//
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			buffer := syncPool.Get().([]byte)
// 			if buffer == nil {
// 				b.Errorf("Expected buffer to be not nil")
// 			}
// 			if cap(buffer) != bufferCap {
// 				b.Errorf("Expected buffer capacity to be %d, but got %d", bufferCap, cap(buffer))
// 			}
// 			syncPool.Put(buffer)
// 		}
// 	})
//
// }
