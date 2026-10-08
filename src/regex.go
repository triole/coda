package main

import (
	"regexp"
	"sync"
)

// regexCache stores compiled regexes for performance
var (
	regexCache     map[string]*regexp.Regexp
	regexCacheMu   sync.RWMutex
	regexCacheSize = 100 // Limit cache size to prevent memory bloat
)

func init() {
	regexCache = make(map[string]*regexp.Regexp, regexCacheSize)
}

// find uses cached regex compiles to avoid repeated parsing
func find(rx string, str string) (result string) {
	// Check cache first
	regexCacheMu.RLock()
	temp, cached := regexCache[rx]
	regexCacheMu.RUnlock()

	if !cached {
		// Compile and cache
		var err error
		temp, err = regexp.Compile(rx)
		if err != nil {
			logger.Warn("invalid regex %q: %v\n", rx, err)
			return ""
		}
		// Add to cache with size limit
		regexCacheMu.Lock()
		if len(regexCache) >= regexCacheSize {
			// Simple LRU eviction: remove first key (not optimal but functional)
			for key := range regexCache {
				delete(regexCache, key)
				break
			}
		}
		regexCache[rx] = temp
		regexCacheMu.Unlock()
	}

	result = temp.FindString(str)
	return
}
