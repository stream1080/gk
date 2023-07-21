package lru

import "container/list"

type Cache struct {
	mBytes    int64 // max cache
	nBytes    int64 // use cache
	list      *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		mBytes:    maxBytes,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if e, ok := c.cache[key]; ok {
		c.list.MoveToFront(e)
		kv := e.Value.(*entry)
		return kv.value, ok
	}

	return
}

func (c *Cache) Del() {
	e := c.list.Back()
	if e == nil {
		return
	}

	c.list.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())

	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

func (c *Cache) Add(key string, value Value) {
	if e, ok := c.cache[key]; ok {
		c.list.MoveToFront(e)
		kv := e.Value.(*entry)
		c.nBytes += int64(len(key)) - int64(kv.value.Len())
		kv.value = value
	} else {
		e := c.list.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = e
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	//
	for c.mBytes != 0 && c.mBytes < c.nBytes {
		c.Del()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.list.Len()
}
