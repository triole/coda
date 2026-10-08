package main

import (
	"regexp"
	"sync"
)

// regexCache stores compiled regexes for performance
var (
	regexCache     map[string]*regexp.Regexp
	regexCacheMu   sync.RWMutex
	regexCacheSize = 100
)

func init() {
	regexCache = make(map[string]*regexp.Regexp, regexCacheSize)
}

// Find uses cached regex compiles to avoid repeated parsing with double-check locking
func find(rx string, str string) (result string) {
	// Check cache first (read lock)
	regexCacheMu.RLock()
	temp, cached := regexCache[rx]
	regexCacheMu.RUnlock()

	if !cached {
		// Compile regex (outside lock to avoid blocking readers)
		var err error
		temp, err = regexp.Compile(rx)
		if err != nil {
			logger.Warn("invalid regex %q: %v\n", rx, err)
			return ""
		}

		// Insert into cache with write lock and double-check
		regexCacheMu.Lock()
		defer regexCacheMu.Unlock()

		// Double-check: another goroutine may have inserted it while we compiled
		if temp2, ok := regexCache[rx]; ok {
			temp = temp2
		} else {
			// Evict one entry if cache is full (random eviction, not LRU)
			if len(regexCache) >= regexCacheSize {
				for key := range regexCache {
					delete(regexCache, key)
					break
				}
			}
			regexCache[rx] = temp
		}
	}

	result = temp.FindString(str)
	return
}
