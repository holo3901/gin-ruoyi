package models

type SysCache struct {
	CacheName  string `json:"cacheName,omitempty"`
	CacheKey   string `json:"cacheKey,omitempty"`
	CacheValue string `json:"cacheValue,omitempty"`
}
