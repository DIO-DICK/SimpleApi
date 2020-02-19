package simplecache


import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 定义缓存存放的值类型和失效时间，Object为空接口类型
type Item struct {
	Object interface{}
	Expiration int64
}

// 定义无失效时间和默认失效时间
const (
	NoExpiration time.Duration = -1
	DefaultExpiration time.Duration = 0
)


type cache struct {
	defaultexpiration time.Duration
	items map[string]Item
	mu sync.RWMutex
	timechannel *timechannel
}

type Cache struct {
	*cache
}

// 当Expiration大于0时设置此Item的过期时间--非线程安全
func (c *cache) set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultexpiration
	}

	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object: v,
		Expiration: e,
	}
}

// 当Expiration大于0时设置此Item的过期时间--线程安全
func (c *cache) Set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultexpiration
	}

	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	c.mu.Lock()
	c.items[k] = Item{
		Object: v,
		Expiration: e,
	}
	c.mu.Unlock()
}

// 替换已存在的数据，如果不存在相关key，输出错误
func (c *cache) replace (k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		fmt.Println("此Item不存在")
	}
	c.set(k, v, d)
	c.mu.Unlock()
	return nil
}

// 获取相关键值对，当key不存在或时间大于超时时间时，返回空值 -- 非线程安全
func (c *cache) get (k string) (interface{}, bool) {
	key, value := c.items[k]
	if !value {
		return nil, false
	}
	if key.Expiration > 0 {
		if time.Now().UnixNano() > key.Expiration {
			return nil, false
		}
	}
	return key.Object, true
}

// 获取相关键值对，当key不存在或时间大于超时时间时，返回空值 -- 线程安全
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	key, value := c.items[k]
	if !value {
		c.mu.RUnlock()
		return nil, false
	}
	if key.Expiration > 0 {
		if time.Now().UnixNano() > key.Expiration {
			c.mu.RUnlock()
			return nil,false
		}
	}
	c.mu.RUnlock()
	return key.Object, true
}

// 删除缓存中的数据，判断没有相关数据则返回空和false
func (c *cache) delete (k string) (interface{}, bool) {
	if v, found := c.items[k]; found {
		delete(c.items, k)
		return v.Object, true
	}
	return nil, false
}

func (c *cache) Delete (k string) (interface{}, bool) {
	c.mu.Lock()
	ob, bo := c.delete(k)
	c.mu.Unlock()
	return ob, bo
}

// 判断map中数据的过期时间，删除过期时间小于当前时间的数据
func (c *cache) DeleteExpiration() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	for key, value := range c.items {
		if value.Expiration > 0 && value.Expiration < now {
			delete(c.items, key)
		}
	}
}

func (c *cache) Flush() {
	c.mu.Lock()
	c.items = map[string]Item{}
	c.mu.Unlock()
}

/*确定一个时间结构体，用于定期删除缓存内的数据*/
type timechannel struct {
	Interval time.Duration
	stop chan bool
}

func (t *timechannel) Run(c *cache) {
	ticker := time.NewTicker(t.Interval)
	for {
		select {
		case <- ticker.C:
			c.DeleteExpiration()
		case <- t.stop:
			ticker.Stop()
			return
		}
	}
}

func StopTimechannel (c *Cache) {
	c.timechannel.stop <- true
}

func RunTimechannel(c *cache, d time.Duration) {
	t := &timechannel{
		Interval: d,
		stop: make(chan bool),
	}
	c.timechannel = t
	go t.Run(c)
}

func InitCache(d time.Duration, m map[string]Item) *cache {
	c := &cache{
		defaultexpiration: d,
		items: m,
	}
	return c
}

func InitCacheWithTimechannel(de time.Duration, ci time.Duration, m map[string]Item) *Cache {
	c := InitCache(de, m)
	C := &Cache{c}
	if ci > 0 {
		RunTimechannel(c, ci)
		runtime.SetFinalizer(C,StopTimechannel)
	}
	return C
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	item := make(map[string]Item)
	return InitCacheWithTimechannel(defaultExpiration, cleanupInterval, item)
}