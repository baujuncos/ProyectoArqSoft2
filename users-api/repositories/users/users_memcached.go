package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	dao "users-api/dao"
)

type MemcachedConfig struct {
	Host string
	Port string
}

type Memcached struct {
	client *memcache.Client
}

func idKey(id int64) string {
	return fmt.Sprintf("id:%d", id)
}

func emailKey(email string) string {
	return fmt.Sprintf("email:%s", email)
}

func NewMemcached(config MemcachedConfig) Memcached {
	// Connect to Memcached
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	client := memcache.New(address)

	return Memcached{client: client}
}

func (repository Memcached) GetUserByID(id int64) (dao.Users, error) {
	// Retrieve the user from Memcached
	key := idKey(id)
	item, err := repository.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return dao.Users{}, fmt.Errorf("user not found")
		}
		return dao.Users{}, fmt.Errorf("error fetching user from memcached: %w", err)
	}

	// Deserialize the data
	var user dao.Users
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return dao.Users{}, fmt.Errorf("error unmarshaling user: %w", err)
	}
	return user, nil
}

func (repository Memcached) GetUserByEmail(email string) (dao.Users, error) {
	// Assume we store users with "email:<email>" as key
	key := emailKey(email)
	item, err := repository.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return dao.Users{}, fmt.Errorf("user not found")
		}
		return dao.Users{}, fmt.Errorf("error fetching user by email from memcached: %w", err)
	}

	// Deserialize the data
	var user dao.Users
	if err := json.Unmarshal(item.Value, &user); err != nil {
		return dao.Users{}, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return user, nil
}

func (repository Memcached) CreateUser(user dao.Users) (int64, error) {
	// Serialize user data
	data, err := json.Marshal(user)
	if err != nil {
		return 0, fmt.Errorf("error marshaling user: %w", err)
	}

	// Store user with ID as key and email as an alternate key
	idKey := idKey(user.User_id)
	if err := repository.client.Set(&memcache.Item{Key: idKey, Value: data}); err != nil {
		return 0, fmt.Errorf("error storing user in memcached: %w", err)
	}

	// Set key for email as well for easier lookup by email
	emailKey := emailKey(user.Email)
	if err := repository.client.Set(&memcache.Item{Key: emailKey, Value: data}); err != nil {
		return 0, fmt.Errorf("error storing email in memcached: %w", err)
	}

	return user.User_id, nil
}
