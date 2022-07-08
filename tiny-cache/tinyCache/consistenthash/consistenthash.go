package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32 // 将 byte 数组转换为 uint32 的函数

// Map 存储 hash 所有的 key
type Map struct {
	hash     Hash           // 允许用户替换成自定义的 hash 函数
	replicas int            // 代表虚拟节点的倍数
	keys     []int          // sorted keys
	hashMap  map[int]string // 虚拟节点和 真实节点的映射表，键是虚拟节点的 hashValue key 是真实节点的名称
}

// New 函数允许自定义虚拟节点的倍数和 hash 函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE // 默认使用这个 hash 算法
	}
	return m
}

// Add 添加一个真实节点
func (m *Map) Add(keys ...string) { // 允许传入 0 或多个 真实节点的名称
	for _, key := range keys { // 对于每一个真实节点 key 创建 m.replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			hashValue := int(m.hash([]byte(strconv.Itoa(i) + key))) // 虚拟节点的名称，使用 m.hash 计算出具体的 hashValue
			m.keys = append(m.keys, hashValue)
			m.hashMap[hashValue] = key
		}
	}
	sort.Ints(m.keys) // 将环上面的 hashValue 进行排序
}

// Get 根据给定的 key 得到具体使用的 node
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key))) // 获取具体的 hashValue
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]] // 在数组上面取余数就可以成环
}
