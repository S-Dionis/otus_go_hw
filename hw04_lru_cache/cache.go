package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type KeyValue struct {
	Key   Key
	Value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	kv := KeyValue{Key: key, Value: value}
	element, exists := c.items[key]
	if exists {
		element.Value = kv
		c.queue.MoveToFront(element)
	} else {
		if c.queue.Len()+1 > c.capacity {
			back := c.queue.Back()
			c.queue.Remove(back)
			delete(c.items, back.Value.(KeyValue).Key)
		}
		c.items[key] = c.queue.PushFront(kv)
	}
	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	element, exists := c.items[key]
	if !exists {
		return nil, false
	}
	c.queue.MoveToFront(element)
	return element.Value.(KeyValue).Value, true
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
