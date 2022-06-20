package consistenthash

import "hash/crc32"

// Hash map bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed values
type Map struct {
	hash     Hash
	replicas int
	keys     []int // sorted keys
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return nil
}
