package qqwry

import (
	"container/list"
	"sync"
)

// cacheEntry 缓存条目
type cacheEntry struct {
	key   string
	value *Location
}

// Cache 简单的LRU缓存实现
// 字段: capacity 容量, list 双向链表, items哈希表, lock 互斥锁
type Cache struct {
	capacity int
	list     *list.List
	items    map[string]*list.Element
	lock     sync.Mutex
}

// NewCache 创建新的LRU缓存
// 入参: capacity 缓存容量
// 返回: *Cache 缓存实例
func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		list:     list.New(),
		items:    make(map[string]*list.Element),
	}
}

// Get 获取缓存
// 入参: key 键
// 返回: value 值, ok 是否存在
func (c *Cache) Get(key string) (value *Location, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if ent, ok := c.items[key]; ok {
		c.list.MoveToFront(ent)
		return ent.Value.(*cacheEntry).value, true
	}
	return nil, false
}

// Add 添加缓存
// 入参: key 键, value 值
func (c *Cache) Add(key string, value *Location) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if ent, ok := c.items[key]; ok {
		c.list.MoveToFront(ent)
		ent.Value.(*cacheEntry).value = value
		return
	}
	ent := c.list.PushFront(&cacheEntry{key: key, value: value})
	c.items[key] = ent
	if c.list.Len() > c.capacity {
		back := c.list.Back()
		if back != nil {
			c.list.Remove(back)
			delete(c.items, back.Value.(*cacheEntry).key)
		}
	}
}
