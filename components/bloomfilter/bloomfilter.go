package main

import (
	"fmt"
	"hash/fnv"
)

type BloomFilter struct {
	bitArray         []bool
	numHashFunctions int
}

func NewBloomFilter(size int, numHashFunctions int) *BloomFilter {
	return &BloomFilter{
		bitArray:         make([]bool, size),
		numHashFunctions: numHashFunctions,
	}
}

func (bf *BloomFilter) Add(item string) {
	for i := 0; i < bf.numHashFunctions; i++ {
		hash := hash(item, i) % len(bf.bitArray)
		bf.bitArray[hash] = true
	}
}

func (bf *BloomFilter) Contains(item string) bool {
	for i := 0; i < bf.numHashFunctions; i++ {
		hash := hash(item, i) % len(bf.bitArray)
		if !bf.bitArray[hash] {
			return false
		}
	}
	return true
}

func hash(item string, seed int) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(item))
	hashValue := hasher.Sum32()
	// 1001^1111=0110 same return 0, else return 1
	return int(hashValue) ^ seed
}

func main() {
	filterSize := 1000
	numHashFunctions := 5
	bloomFilter := NewBloomFilter(filterSize, numHashFunctions)

	// 添加一些元素到布隆过滤器中
	items := []string{"apple", "banana", "cherry", "date"}
	for _, item := range items {
		bloomFilter.Add(item)
	}

	// 检查元素是否存在于布隆过滤器中
	fmt.Println("apple exists:", bloomFilter.Contains("apple"))   // 应返回 true
	fmt.Println("grape exists:", bloomFilter.Contains("grape"))   // 应返回 false
	fmt.Println("cherry exists:", bloomFilter.Contains("cherry")) // 应返回 true
}
