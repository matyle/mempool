# mempool

实现了一个 No GC 内存池，用于重用 bytes.Buffer 对象，以减少内存分配和垃圾回收的开销。
使用方法

```go
import "github.com/mempool"

// 创建一个内存池，池的大小为 10，每个 bytes.Buffer 对象的容量为 1024
pool := mempool.NewPool(10, 1024)

// 从内存池中获取一个 bytes.Buffer 对象
buf := pool.Get()

// 使用 bytes.Buffer 对象
buf.WriteString("hello, world")

// 将 bytes.Buffer 对象放回内存池中
pool.Put(buf)
```

- 内存池可以减少内存分配和垃圾回收的开销，提程序的性能。
- 内存池可以重用.Buffer 对象，减少内存碎片的产生，提高内存的利用率。
- 内存池可以避免频繁的内存分配和垃圾回收减少内存泄漏的风险，提高程序的稳定性。
