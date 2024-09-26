package vm

import (
	"sync"
	"unsafe"

	"github.com/kode4food/ale/pkg/data"
)

const (
	bucketCount = 64
	mb          = 1024 * 1024
	spanSize    = mb / int(unsafe.Sizeof(data.Value(nil)))
)

type (
	memBucket struct {
		free    *memEntry
		entries *memEntry
		size    int
		spans   int
		sync.Mutex
	}

	memEntry struct {
		next   *memEntry
		values data.Vector
	}

	allocator struct {
		buckets [bucketCount]memBucket
	}
)

var (
	zeros = make(data.Vector, bucketCount)
	mem   = newAllocator()
)

func newAllocator() *allocator {
	res := &allocator{}
	for i := 0; i < len(res.buckets); i++ {
		b := &res.buckets[i]
		b.size = i + 1
	}
	return res
}

func malloc(size int) data.Vector {
	if size == 0 {
		return data.EmptyVector
	}
	if size > bucketCount {
		return make(data.Vector, size)
	}
	return mem.getBucket(size).alloc()
}

func free(vals data.Vector) {
	if size := len(vals); size > 0 && size <= bucketCount {
		mem.getBucket(size).dealloc(vals)
	}
}

func (a *allocator) getBucket(size int) *memBucket {
	return &a.buckets[size-1]
}

func (b *memBucket) putEntry(e *memEntry) {
	e.next = b.entries
	b.entries = e
}

func (b *memBucket) getEntry() *memEntry {
	if e := b.entries; e != nil {
		b.entries = e.next
		return e
	}
	return &memEntry{}
}

func (b *memBucket) pushFree(e *memEntry) {
	e.next = b.free
	b.free = e
}

func (b *memBucket) popFree() *memEntry {
	if e := b.free; e != nil {
		b.free = e.next
		return e
	}
	return nil
}

func (b *memBucket) alloc() data.Vector {
	b.Lock()
	if next := b.popFree(); next != nil {
		res := next.values
		b.putEntry(next)
		b.Unlock()
		return res
	}
	total := spanSize / b.size
	values := make(data.Vector, spanSize)
	var next *memEntry
	for i := b.size; i < total; i += b.size {
		next = &memEntry{
			values: values[i : i+b.size],
			next:   next,
		}
	}
	b.free = next
	b.spans++
	b.Unlock()
	return values[0:b.size]
}

func (b *memBucket) dealloc(v data.Vector) {
	copy(v, zeros)
	b.Lock()
	e := b.getEntry()
	e.values = v
	b.pushFree(e)
	b.Unlock()
}
