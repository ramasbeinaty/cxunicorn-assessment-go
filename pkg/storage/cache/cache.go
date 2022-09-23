package cache

import (
	"clinicapp/pkg/storage/postgres"
	"errors"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
)

func NewCacheMem() (*cacheMem, error) {
	c := new(cacheMem)

	// parse and define the expiration duration of the cache
	_expiryDuration, err := time.ParseDuration(os.Getenv("CACHE_EXPIRY_DURATION"))
	if err != nil {
		return c, errors.New("NewCacheMem - Failed to parse expiry duration of type string to type duration -" + err.Error())
	}

	// parse and define the purging duration of the cache after which expired data will be deleted
	_purgeDuration, err := time.ParseDuration(os.Getenv("CACHE_PURGE_DURATION"))
	if err != nil {
		return c, errors.New("NewCacheMem - Failed to parse purge duration of type string to type duration -" + err.Error())
	}

	// create a new cache memory
	_cache := cache.New(_expiryDuration, _purgeDuration)

	c.cache = _cache
	return c, nil
}

type cacheMem struct {
	cache *cache.Cache
}

func (c cacheMem) GetDoctor(id int) (Doctor, error) {

	// attempt to get a doctor from cache memory
	_cachedDoctor, found := c.cache.Get(cachedDoctor)

	// check if any doctor is found
	if !found {
		return Doctor{}, errors.New("no doctor found in cache memory")
	}

	// then check if the cached doctor found has the same id of the requested doctor
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
