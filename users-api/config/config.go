package config

import "time"

const (
	MySQLHost     = "mysql"
	MySQLPort     = "3306"
	MySQLDatabase = "users_api"
	MySQLUsername = "root"
	MySQLPassword = "belusql1"

	CacheDuration = 30 * time.Second //tiempo que datos se mantienen en cache 30 segundos

	MemcachedHost = "memcached"
	MemcachedPort = "11211"

	JWTKey      = "ThisIsAnExampleJWTKey!"
	JWTDuration = 24 * time.Hour //tiempo de validez de token: 24 horas
)
