# Mempool

Mempool 是一个用于减少内存分配和提高字节缓冲区处理性能的内存池库。
实现了一个 Limited GC 内存池，用于重用 bytes.Buffer 对象，以减少内存分配和垃圾回收的开销。

## 安装

```sh
go get github.com/matyle/mempool
```

## 使用方法

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

- 内存池可以减少内存分配和垃圾回收的开销，提程序的性能。
- 内存池可以重用.Buffer 对象，减少内存碎片的产生，提高内存的利用率。
- 内存池可以避免频繁的内存分配和垃圾回收减少内存泄漏的风险，提高程序的稳定性。
