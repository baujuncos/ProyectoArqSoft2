package users

import (
	"fmt"
	"github.com/karlseguin/ccache"
	"time"
	dao "users-api/dao"
)

type CacheConfig struct {
	TTL time.Duration // Cache expiration time
}

type Cache struct {
	client *ccache.Cache
	ttl    time.Duration
}

func NewCache(config CacheConfig) Cache {
	// Initialize ccache with default settings
	cache := ccache.New(ccache.Configure())
	return Cache{
		client: cache,
		ttl:    config.TTL,
	}
}

func (repository Cache) GetUserByID(id int64) (dao.Users, error) {
	// Convert ID to string for cache key
	idKey := fmt.Sprintf("user:id:%d", id)

	// Try to get from cache
	item := repository.client.Get(idKey)
	if item != nil && !item.Expired() {
		// Return cached value
		user, ok := item.Value().(dao.Users)
		if !ok {
			return dao.Users{}, fmt.Errorf("failed to cast cached value to user")
		}
		return user, nil
	}

	// If not found, return cache miss error
	return dao.Users{}, fmt.Errorf("cache miss for user ID %d", id)
}

func (repository Cache) GetUserByEmail(email string) (dao.Users, error) {
	userKey := fmt.Sprintf("user:email:%s", email)

	// Try to get from cache
	item := repository.client.Get(userKey)
	if item != nil && !item.Expired() {
		// Return cached value
		user, ok := item.Value().(dao.Users)
		if !ok {
			return dao.Users{}, fmt.Errorf("failed to cast cached value to user")
		}
		return user, nil
	}

	// If not found, return cache miss error
	return dao.Users{}, fmt.Errorf("cache miss for email %s", email)
}

func (repository Cache) CreateUser(user dao.Users) (int64, error) {
	// Cache user by ID and by username after creation
	idKey := fmt.Sprintf("user:id:%d", user.User_id)
	userKey := fmt.Sprintf("user:email:%s", user.Email)

	// Set user in cache
	repository.client.Set(idKey, user, repository.ttl)
	repository.client.Set(userKey, user, repository.ttl)

	// Return the user ID as if it was created successfully
	return user.User_id, nil
}
