package bloomfilter

import (
	"hash/fnv"
	"sync"
)

//实现一个布隆过滤器

const (
	DefaultSize = 1000000
	Hashes      = 1
)

type BloomFilter struct {
	set    []byte
	hashes [Hashes]func(data []byte) uint32
	size   uint
	mutex  sync.RWMutex
}

func NewBloomFilter() *BloomFilter {
	bf := &BloomFilter{}
	bf.set = make([]byte, DefaultSize)
	bf.hashes = [Hashes]func(data []byte) uint32{
		fnvHash1,
	}
	bf.size = uint(len(bf.set)) * 8
	return bf
}

func fnvHash1(data []byte) uint32 {
	h := fnv.New32a()
	_, err := h.Write(data)
	if err != nil {
		return 0
	}
	return h.Sum32()
}

func (bf *BloomFilter) Add(data []byte) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for _, hash := range bf.hashes {
		index := hash(data) % bf.size
		byteIndex := index / 8
		bitIndex := byte(index % 8)
		bf.set[byteIndex] |= 1 << bitIndex

	}

}

func (bf *BloomFilter) MayExists(data []byte) bool {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	for _, hash := range bf.hashes {
		index := hash(data) % bf.size
		byteIndex := index / 8
		bitIndex := byte(index % 8)
		if bf.set[byteIndex]&(1<<bitIndex) == 0 {
			return false
		}
	}
	return true
}
