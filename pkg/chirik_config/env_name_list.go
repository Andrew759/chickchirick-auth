package chirik_config

// server
const (
	ServerUrl = "SERVER_URL"
)

// environment
const (
	Enviroment = "ENVIRONMENT"
	Dev        = "DEV"
	Test       = "TEST"
	Prod       = "PROD"
)

// database
const (
	DbHost           = "DB_HOST"
	DbPort           = "DB_PORT"
	DbName           = "DB_NAME"
	DbUser           = "DB_USER"
	DbPass           = "DB_PASS"
	DbTimezone       = "DB_TIMEZONE"
	DbMigrationPatch = "DB_MIGRATION_PATCH"
)

// redis
const (
	RedisHost     = "REDIS_HOST"
	RedisPort     = "REDIS_PORT"
	RedisUser     = "REDIS_USER"
	RedisPassword = "REDIS_PASSWORD"
)
