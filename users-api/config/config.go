import "time"

const (
MySQLHost     = "localhost"
MySQLHost     = "MySQL"
MySQLPort     = "3306"
MySQLDatabase = "users-api"
MySQLUsername = "root"
MySQLPassword = "root"
MySQLPassword = "ladrillo753"

CacheDuration = 30 * time.Second
CacheDuration = 30 * time.Second //tiempo que datos se mantienen en cache 30 segundos

MemcachedHost = "localhost"
MemcachedHost = "memCached"
MemcachedPort = "11211"

JWTKey      = "ThisIsAnExampleJWTKey!"
JWTDuration = 24 * time.Hour
JWTDuration = 24 * time.Hour //tiempo de validez de token: 24 horas
)