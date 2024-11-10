package config

import "time"

const (
	MySQLHost     = "MySQL"
	MySQLPort     = "3306"
	MySQLDatabase = "users-api"
	MySQLUsername = "root"
	MySQLPassword = "ladrillo753"

	CacheDuration = 30 * time.Second //tiempo que datos se mantienen en cache 30 segundos

	MemcachedHost = "memCached"
	MemcachedPort = "11211"

	JWTKey      = "ThisIsAnExampleJWTKey!"
	JWTDuration = 24 * time.Hour //tiempo de validez de token: 24 horas
)
