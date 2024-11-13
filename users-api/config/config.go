package config

import "time"

const (
	MySQLHost     = "mysql"
	MySQLPort     = "3306"
	MySQLDatabase = "users-api"
	MySQLUsername = "root"
	MySQLPassword = "ladrillo753"

	CacheDuration = 30 * time.Second //tiempo que datos se mantienen en cache 30 segundos

	MemcachedHost = "localhost"
	MemcachedPort = "11211"

	JWTKey      = "ThisIsAnExampleJWTKey!"
	JWTDuration = 24 * time.Hour //tiempo de validez de token: 24 horas
)
