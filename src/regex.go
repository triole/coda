package main

import (
	"regexp"
	"sync"
)

// LRUCache implements a simple LRU cache using a doubly-linked list and a map
type LRUCache struct {
	mu      sync.RWMutex
	cap     int
	entries map[string]*cacheEntry
	head    *cacheEntry // most recently used
	tail    *cacheEntry // least recently used
}

// cacheEntry holds a cached regex and linked list pointers
type cacheEntry struct {
	key        string
	value      *regexp.Regexp
	prev, next *cacheEntry
}

// NewLRUCache creates a new LRU cache with the given capacity
func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		cap:     cap,
		entries: make(map[string]*cacheEntry, cap),
	}
}

// Get retrieves a value from the cache, moving it to the front (MRU)
func (lru *LRUCache) Get(key string) (*regexp.Regexp, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	entry, ok := lru.entries[key]
	if !ok {
		return nil, false
	}

	// Move to front (most recently used)
	lru.moveToFront(entry)
	return entry.value, true
}

// Set adds or updates a value in the cache
func (lru *LRUCache) Set(key string, value *regexp.Regexp) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if entry, ok := lru.entries[key]; ok {
		// Update existing entry
		entry.value = value
		lru.moveToFront(entry)
		return
	}

	// Evict if at capacity
	if len(lru.entries) >= lru.cap {
		lru.evict()
	}

	// Create new entry at front
	entry := &cacheEntry{
		key:   key,
		value: value,
	}
	lru.insertBefore(lru.head, entry)
	lru.entries[key] = entry
}

// moveToFront moves an entry to the front of the list
func (lru *LRUCache) moveToFront(entry *cacheEntry) {
	if entry == lru.head {
		return // already at front
	}

	// Remove from current position
	if entry.prev != nil {
		entry.prev.next = entry.next
	}
	if entry.next != nil {
		entry.next.prev = entry.prev
	}
	if entry == lru.tail {
		lru.tail = entry.prev
	}

	// Insert at front
	lru.insertBefore(lru.head, entry)
}

// insertBefore inserts entry before target
func (lru *LRUCache) insertBefore(target, entry *cacheEntry) {
	entry.prev = target.prev
	entry.next = target
	if target.prev != nil {
		target.prev.next = entry
	}
	target.prev = entry
	if lru.head == target {
		lru.head = entry
	}
}

// evict removes the least recently used entry
func (lru *LRUCache) evict() {
	if lru.tail != nil {
		delete(lru.entries, lru.tail.key)
		if lru.tail.prev != nil {
			lru.tail.prev.next = nil
		}
		lru.tail = lru.tail.prev
	}
}

// regexCache stores compiled regexes with LRU eviction
var regexCache = NewLRUCache(100)

// find uses LRU cache to avoid repeated parsing
func find(rx, str string) (result string) {
	// Check cache first
	if temp, ok := regexCache.Get(rx); ok {
		return temp.FindString(str)
	}

	// Compile regex (outside lock)
	temp, err := regexp.Compile(rx)
	if err != nil {
		log.Warn("invalid regex", map[string]any{"regex": rx, "error": err})
		return ""
	}

	// Add to cache
	regexCache.Set(rx, temp)

	return temp.FindString(str)
}
