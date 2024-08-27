package cache

// Cache is a cache of string. When it is "full" the oldest entry is dropped.
type Cache struct {
	cache   map[string]any
	entries []*string
}

// Exists add the key to the cache and returns true when the doesn't exist.
// Otherwise, the key is added to the cache and it returns false.
func (c *Cache) Exists(key string) bool {
	if _, exists := c.cache[key]; exists {
		return true
	} else {
		c.add(key)
		return false
	}
}

// add adds the key to the cache and drop the oldest entry when the cache is "full".
func (c *Cache) add(key string) {
	if len(c.entries) == cap(c.entries) {
		c.removeOldest()
	}
	c.cache[key] = nil
	c.entries = append(c.entries, &key)
}

// removeOldest drops the oldest entry.
func (c *Cache) removeOldest() {
	delete(c.cache, *c.entries[0])
	c.entries = c.entries[1:]
}

// New returns a new cache instance.
func New(size uint) *Cache {
	return &Cache{
		cache:   make(map[string]any, size),
		entries: make([]*string, 0, size),
	}
}
