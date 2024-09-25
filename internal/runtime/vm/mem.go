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
		free  *memEntry
		size  int
		spans int
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
		res.buckets[i].size = i + 1
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

func (s *memBucket) alloc() data.Vector {
	s.Lock()
	if next := s.free; next != nil {
		s.free = next.next
		s.Unlock()
		return next.values
	}
	total := spanSize / s.size
	values := make(data.Vector, spanSize)
	var next *memEntry
	for i := s.size; i < total; i += s.size {
		next = &memEntry{
			values: values[i : i+s.size],
			next:   next,
		}
	}
	s.free = next
	s.spans++
	s.Unlock()
	return values[0:s.size]
}

func (s *memBucket) dealloc(v data.Vector) {
	copy(v, zeros)
	s.Lock()
	s.free = &memEntry{
		values: v,
		next:   s.free,
	}
	s.Unlock()
}
