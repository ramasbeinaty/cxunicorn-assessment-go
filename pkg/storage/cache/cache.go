package cache

import (
	"clinicapp/pkg/listing"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

func NewCacheMem() *cacheMem {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	_cache := cache.New(5*time.Minute, 10*time.Minute)

	return &cacheMem{
		cache: _cache,
	}
}

type cacheMem struct {
	cache *cache.Cache
}

func (c cacheMem) GetDoctor(id int) (listing.Doctor, error) {

	_cachedDoctor, found := c.cache.Get(cachedDoctor)

	if !found {
		return listing.Doctor{}, errors.New("no doctor found in cache memory")
	}

	cachedDoctor := _cachedDoctor.(listing.Doctor)

	if cachedDoctor.ID != id {
		return listing.Doctor{}, errors.New("requested doctor not found in cache memory")
	}

	return cachedDoctor, nil
}

func (c cacheMem) SetDoctor(doctor listing.Doctor) error {
	c.cache.Set(cachedDoctor, doctor, cache.DefaultExpiration)

	return nil
}
