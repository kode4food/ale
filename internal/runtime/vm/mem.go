package vm

import (
	"sync"
	"unsafe"

	"github.com/kode4food/ale/internal/debug"
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

	Allocator struct {
		buckets [bucketCount]memBucket
	}
)

var mem = NewAllocator()

func malloc(size int) data.Vector {
	return mem.Malloc(size)
}

func free(vals data.Vector) {
	mem.Free(vals)
}

func NewAllocator() *Allocator {
	res := &Allocator{}
	for i := range len(res.buckets) {
		b := &res.buckets[i]
		b.size = i + 1
	}
	return res
}

func (a *Allocator) Malloc(size int) data.Vector {
	switch {
	case size > 0 && size <= bucketCount:
		return a.getBucket(size).alloc()
	case size == 0:
		return data.EmptyVector
	case size > bucketCount:
		return make(data.Vector, size)
	default:
		panic(debug.ProgrammerError("invalid malloc: %d", size))
	}
}

func (a *Allocator) Free(vals data.Vector) {
	if size := len(vals); size > 0 && size <= bucketCount {
		a.getBucket(size).dealloc(vals)
	}
}

func (a *Allocator) getBucket(size int) *memBucket {
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
	panic(debug.ProgrammerError("no free memory entries"))
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
	entries := make([]memEntry, total)
	var next *memEntry
	for i, off := 1, b.size; i < total; i, off = i+1, off+b.size {
		e := &entries[i]
		e.values = values[off : off+b.size]
		e.next = next
		next = e
	}
	b.putEntry(&entries[0])
	b.free = next
	b.spans++
	b.Unlock()
	return values[0:b.size]
}

func (b *memBucket) dealloc(v data.Vector) {
	clear(v)
	b.Lock()
	e := b.getEntry()
	e.values = v
	b.pushFree(e)
	b.Unlock()
}
