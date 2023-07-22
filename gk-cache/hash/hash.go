package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int // Sorted
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}

	return &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= int(m.hash([]byte(key)))
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}

// Remove use to remove a key and its virtual keys on the ring and map
func (m *Map) Del(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
}
