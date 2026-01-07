package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// ImageCache provides in-memory caching for images during the refine flow
type ImageCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
}

// CacheEntry represents a cached image with metadata
type CacheEntry struct {
	OriginalImageBase64 string
	MimeType            string
	CreatedAt           time.Time
}

// Global image cache instance
var imageCache = NewImageCache()

// NewImageCache creates a new image cache
func NewImageCache() *ImageCache {
	cache := &ImageCache{
		entries: make(map[string]*CacheEntry),
	}
	// Start background cleanup goroutine
	go cache.cleanupLoop()
	return cache
}

// GetImageCache returns the global image cache instance
func GetImageCache() *ImageCache {
	return imageCache
}

// generateCacheID creates a random cache ID
func generateCacheID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Store saves an image to the cache and returns a cache ID
func (c *ImageCache) Store(imageBase64, mimeType string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheID := generateCacheID()
	c.entries[cacheID] = &CacheEntry{
		OriginalImageBase64: imageBase64,
		MimeType:            mimeType,
		CreatedAt:           time.Now(),
	}

	return cacheID
}

// Get retrieves an image from the cache
func (c *ImageCache) Get(cacheID string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[cacheID]
	return entry, exists
}

// Delete removes an image from the cache
func (c *ImageCache) Delete(cacheID string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.entries[cacheID]; exists {
		delete(c.entries, cacheID)
		return true
	}
	return false
}

// cleanupLoop removes expired cache entries (older than 30 minutes)
func (c *ImageCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup removes entries older than 30 minutes
func (c *ImageCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiry := 30 * time.Minute
	now := time.Now()

	for id, entry := range c.entries {
		if now.Sub(entry.CreatedAt) > expiry {
			delete(c.entries, id)
		}
	}
}

// Count returns the number of cached entries (for debugging)
func (c *ImageCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
