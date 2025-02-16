package ratelimit

import (
	"sync"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
)

type RateLimiter interface {
	Allow(key string) bool
	RemainingTokens(key string) int
	Reset(key string) time.Time
}

type tokenBucket struct {
	tokens     int
	capacity   int
	lastRefill time.Time
	refillRate float64
}

type rateLimiter struct {
	buckets   sync.Map
	config    *config.RateLimitConfig
	cleanupCh chan struct{}
	stopOnce  sync.Once
}

func NewRateLimiter(cfg *config.RateLimitConfig) RateLimiter {
	rl := &rateLimiter{
		config:    cfg,
		cleanupCh: make(chan struct{}),
	}

	go rl.cleanupLoop()
	return rl
}

func (rl *rateLimiter) Allow(key string) bool {
	if !rl.config.Enabled {
		return true
	}

	bucket := rl.getBucket(key)
	return bucket.tryConsume()
}

func (rl *rateLimiter) RemainingTokens(key string) int {
	if !rl.config.Enabled {
		return rl.config.RequestsPerMinute
	}

	bucket := rl.getBucket(key)
	return bucket.tokens
}

func (rl *rateLimiter) Reset(key string) time.Time {
	bucket := rl.getBucket(key)
	return bucket.lastRefill.Add(time.Minute)
}

func (rl *rateLimiter) getBucket(key string) *tokenBucket {
	bucketI, _ := rl.buckets.LoadOrStore(key, &tokenBucket{
		tokens:     rl.config.RequestsPerMinute,
		capacity:   rl.config.RequestsPerMinute,
		lastRefill: time.Now(),
		refillRate: float64(rl.config.RequestsPerMinute) / 60.0, // tokens per second
	})
	return bucketI.(*tokenBucket)
}

func (b *tokenBucket) tryConsume() bool {
	b.refill()
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

func (b *tokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	tokensToAdd := int(elapsed * b.refillRate)

	if tokensToAdd > 0 {
		b.tokens = min(b.capacity, b.tokens+tokensToAdd)
		b.lastRefill = now
	}
}

func (rl *rateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.cleanupCh:
			return
		}
	}
}

func (rl *rateLimiter) cleanup() {
	now := time.Now()
	rl.buckets.Range(func(key, value interface{}) bool {
		bucket := value.(*tokenBucket)
		if now.Sub(bucket.lastRefill) > rl.config.TTL {
			rl.buckets.Delete(key)
		}
		return true
	})
}

func (rl *rateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.cleanupCh)
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
