# Mempool

Mempool is a memory pool library for reducing allocations and improving performance when dealing with byte buffers.

Mempool is Limited GC

## Installation

```sh
go get github.com/matyle/mempool
```

## Usage

```go
package main

import (
"fmt"

"github.com/matyle/mempool"
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

## License

MIT
