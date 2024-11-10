package users

import (
	"fmt"
	"github.com/karlseguin/ccache"
	"time"
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

func (repository Cache) GetAll() ([]users.User, error) {
	// Since it's not typical to cache all users in one request, you might skip caching here
	// Alternatively, you can cache a summary list if needed
	return nil, fmt.Errorf("GetAll not implemented in cache")
}

func (repository Cache) GetUserByID(id int64) (users.User, error) {
	// Convert ID to string for cache key
	idKey := fmt.Sprintf("user:id:%d", id)

	// Try to get from cache
	item := repository.client.Get(idKey)
	if item != nil && !item.Expired() {
		// Return cached value
		user, ok := item.Value().(users.User)
		if !ok {
			return users.User{}, fmt.Errorf("failed to cast cached value to user")
		}
		return user, nil
	}

	// If not found, return cache miss error
	return users.User{}, fmt.Errorf("cache miss for user ID %d", id)
}

func (repository Cache) GetUserByEmail(email string) (users.User, error) {
	userKey := fmt.Sprintf("user:email:%s", email)

	// Try to get from cache
	item := repository.client.Get(userKey)
	if item != nil && !item.Expired() {
		// Return cached value
		user, ok := item.Value().(users.User)
		if !ok {
			return users.User{}, fmt.Errorf("failed to cast cached value to user")
		}
		return user, nil
	}

	// If not found, return cache miss error
	return users.User{}, fmt.Errorf("cache miss for email %s", email)
}

func (repository Cache) CreateUser(user users.User) (int64, error) {
	// Cache user by ID and by username after creation
	idKey := fmt.Sprintf("user:id:%d", user.User_id)
	userKey := fmt.Sprintf("user:email:%s", user.Email)

	// Set user in cache
	repository.client.Set(idKey, user, repository.ttl)
	repository.client.Set(userKey, user, repository.ttl)

	// Return the user ID as if it was created successfully
	return user.User_id, nil
}
