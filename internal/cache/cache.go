package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var GlobalCache *cache.Cache

func InitCache() {
	GlobalCache = cache.New(20*time.Minute, 30*time.Minute)
}
