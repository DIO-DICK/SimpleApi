package cache

import (
	"time"
	cache "simpleapi/app/module/simplecache"
)

var LabelCache *cache.Cache
func init() {
	LabelCache = cache.New(time.Duration(2)*time.Minute,time.Duration(2)*time.Minute)
}
