# Mempool

Mempool is a memory pool library for reducing allocations and improving performance when dealing with byte buffers.

Mempool is never GC

## Installation

```sh
go get github.com/yourusername/mempool
```

## Usage

```go
package main

import (
"fmt"

"github.com/yourusername/mempool"
)

func main() {
poolSize := 5
bufferCap := 1024
pool := mempool.NewPool(poolSize, bufferCap)

buf := pool.Get()
buf.WriteString("Hello, Mempool!")
fmt.Println(buf.String())

pool.Put(buf)
}
```

## API

- NewPool(routineSize int, bufferCap int) \*Pool
  Creates a new Pool object, where routineSize represents the number of goroutines that can access the pool concurrently, and bufferCap represents the estimated size of the buffers that will be stored in the pool.

- Get() \*bytes.Buffer
  Retrieves a buffer from the pool. If the pool does not have any available buffers, it will create a new buffer and return it.

- Put(b \*bytes.Buffer)
  Adds a buffer to the pool. It resets the buffer before adding it to the pool.

- Resize(maxSize int)
  Resizes the pool to a new size. If the new size is smaller than the current size, it does nothing. If the new size is larger than the current size, then it creates a new pool with the desired size and adds elements from the current pool.

## License

MIT
