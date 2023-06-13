package main

import (
	"sync"

	"github.com/matyle/mempool"
)

func main() {
	poolSize := 100
	bufferCap := 1024
	pool := mempool.NewPool(poolSize, bufferCap)

	wg := &sync.WaitGroup{}
	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		go func(i int) {
			buffer := pool.Get()

			buffer.Write([]byte("Hello World!"))
			// println(i, buffer.String())
			_ = buffer.String()
			pool.Put(buffer)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
