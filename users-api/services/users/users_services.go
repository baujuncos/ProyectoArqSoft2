package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	dao "users-api/dao"
	dto "users-api/dto"
)

type Repository interface {
	//Metodos de users
	GetUserByID(id int64) (dao.Users, error)
	GetUserByEmail(email string) (dao.Users, error)
	CreateUser(user dao.Users) (int64, error)
}

type Tokenizer interface {
	GenerateToken(username string, userID int64) (string, error)
}

type Service struct {
	mainRepository      Repository
	cacheRepository     Repository
	memcachedRepository Repository
	tokenizer           Tokenizer
}

func NewService(mainRepository, cacheRepository, memcachedRepository Repository, tokenizer Tokenizer) Service {
	return Service{
		mainRepository:      mainRepository,
		cacheRepository:     cacheRepository,
		memcachedRepository: memcachedRepository,
		tokenizer:           tokenizer,
	}
}

func (service Service) GetUserByID(id int64) (dto.UserDto, error) {
	// Primero chequeamos en cache si esta el user
	user, err := service.cacheRepository.GetUserByID(id)
	if err == nil {
		return service.buildUserDto(user), nil
	}

	// Si no esta en cache buscamos en memcached
	user, err = service.memcachedRepository.GetUserByID(id)
	if err == nil {
		if _, err := service.cacheRepository.CreateUser(user); err != nil {
			return dto.UserDto{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.buildUserDto(user), nil
	}

	// Como no esta en ambas cache y memcached buscamos en la bd principal
	user, err = service.mainRepository.GetUserByID(id)
	if err != nil {
		return dto.UserDto{}, fmt.Errorf("error getting user by ID: %w", err)
	}

	// Guardamos la informacion encontrada en cache y memcached
	if _, err := service.cacheRepository.CreateUser(user); err != nil {
		return dto.UserDto{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.CreateUser(user); err != nil {
		return dto.UserDto{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.buildUserDto(user), nil
}

func (service Service) GetUserByEmail(email string) (dto.UserDto, error) {
	// Check in cache first
	user, err := service.cacheRepository.GetUserByEmail(email)
	if err == nil {
		return service.buildUserDto(user), nil
	}

	// Check memcached
	user, err = service.memcachedRepository.GetUserByEmail(email)
	if err == nil {
		if _, err := service.cacheRepository.CreateUser(user); err != nil {
			return dto.UserDto{}, fmt.Errorf("error caching user after memcached retrieval: %w", err)
		}
		return service.buildUserDto(user), nil
	}

	// Check main repository
	user, err = service.mainRepository.GetUserByEmail(email)
	if err != nil {
		return dto.UserDto{}, fmt.Errorf("error getting user by username: %w", err)
	}

	// Save in cache and memcached
	if _, err := service.cacheRepository.CreateUser(user); err != nil {
		return dto.UserDto{}, fmt.Errorf("error caching user after main retrieval: %w", err)
	}
	if _, err := service.memcachedRepository.CreateUser(user); err != nil {
		return dto.UserDto{}, fmt.Errorf("error saving user in memcached: %w", err)
	}

	return service.buildUserDto(user), nil
}

func (service Service) CreateUser(user dto.UserDto) (int64, error) {
	// Hasheo de la contraseña
	passwordHash := Hash(user.Password)

	newUser := dao.Users{
		Email:    user.Email,
		Password: passwordHash,
		Nombre:   user.Nombre,
		Apellido: user.Apellido,
		Admin:    user.Admin,
	}

	fmt.Printf("Mapped user: %+v\n", newUser) // Depuración

	// Creamos en main repository
	id, err := service.mainRepository.CreateUser(newUser)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	// Add to cache and memcached
	newUser.User_id = id
	if _, err := service.cacheRepository.CreateUser(newUser); err != nil {
		fmt.Printf("warning: failed to cache user: %v\n", err)
	}
	if _, err := service.memcachedRepository.CreateUser(newUser); err != nil {
		fmt.Printf("warning: failed to save user in memcached: %v\n", err)
	}

	return id, nil
}

func (service Service) Login(email string, password string) (dto.LoginDtoResponse, error) {
	// Hash the password
	passwordHash := Hash(password)

	// Try to get user from cache repository first
	user, err := service.cacheRepository.GetUserByEmail(email)
	if err != nil {
		// If not found in cache, try to get user from memcached repository
		user, err = service.memcachedRepository.GetUserByEmail(email)
		if err != nil {
			// If not found in memcached, try to get user from the main repository (database)
			user, err = service.mainRepository.GetUserByEmail(email)
			if err != nil {
				return dto.LoginDtoResponse{}, fmt.Errorf("error getting user by username from main repository: %w", err)
			}

			// Save the found user in both cache and memcached repositories
			if _, err := service.cacheRepository.CreateUser(user); err != nil {
				return dto.LoginDtoResponse{}, fmt.Errorf("error caching user in cache repository: %w", err)
			}
			if _, err := service.memcachedRepository.CreateUser(user); err != nil {
				return dto.LoginDtoResponse{}, fmt.Errorf("error caching user in memcached repository: %w", err)
			}
		} else {
			// Save the found user in the cache repository for future access
			if _, err := service.cacheRepository.CreateUser(user); err != nil {
				return dto.LoginDtoResponse{}, fmt.Errorf("error caching user in cache repository: %w", err)
			}
		}
	}

	// Comparar contraseñas
	if user.Password != passwordHash {
		return dto.LoginDtoResponse{}, fmt.Errorf("invalid credentials")
	}

	// Generate token
	token, err := service.tokenizer.GenerateToken(user.Email, int64(user.User_id))
	if err != nil {
		return dto.LoginDtoResponse{}, fmt.Errorf("error generating token: %w", err)
	}

	// Send the login response
	return dto.LoginDtoResponse{
		User_id: int(user.User_id),
		Token:   token,
		Admin:   user.Admin,
	}, nil
}

func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

func (service Service) buildUserDto(user dao.Users) dto.UserDto {
	return dto.UserDto{
		User_id:  user.User_id,
		Email:    user.Email,
		Password: user.Password,
		Nombre:   user.Nombre,
		Apellido: user.Apellido,
		Admin:    user.Admin,
	}
}
