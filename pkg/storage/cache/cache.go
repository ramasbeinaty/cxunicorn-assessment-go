package cache

import (
	"clinicapp/pkg/storage/postgres"
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

func (c cacheMem) GetDoctor(id int) (Doctor, error) {

	_cachedDoctor, found := c.cache.Get(cachedDoctor)

	if !found {
		return Doctor{}, errors.New("no doctor found in cache memory")
	}

	if _cachedDoctor.(Doctor).ID != id {
		return Doctor{}, errors.New("requested doctor not found in cache memory")
	}

	return _cachedDoctor.(Doctor), nil
}

func (c cacheMem) SetDoctor(doctor postgres.Doctor) {
	var d Doctor

	d.ID = doctor.ID
	d.FirstName = doctor.FirstName
	d.LastName = doctor.LastName
	d.DOB = doctor.DOB
	d.PhoneNumber = doctor.PhoneNumber
	d.Email = doctor.Email
	d.Role = doctor.Role
	d.WorkDays = doctor.WorkDays
	d.WorkTime = doctor.WorkTime
	d.BreakTime = doctor.BreakTime
	d.Specialization = doctor.Specialization

	c.cache.Set(cachedDoctor, d, cache.DefaultExpiration)

}
