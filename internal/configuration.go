package internal

import "os"

type Configuration struct {
	ClickHouseHost     string
	ClickHousePort     string
	ClickHouseDatabase string
	ClickHouseUsername string
	ClickHousePassword string
	Port               string
}

func GetEnvironment(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetConfiguration() Configuration {
	return Configuration{
		ClickHouseHost:     GetEnvironment("CLICKHOUSE_HOST", "localhost"),
		ClickHousePort:     GetEnvironment("CLICKHOUSE_PORT", "9000"),
		ClickHouseDatabase: GetEnvironment("CLICKHOUSE_DATABASE", "default"),
		ClickHouseUsername: GetEnvironment("CLICKHOUSE_USERNAME", "default"),
		ClickHousePassword: GetEnvironment("CLICKHOUSE_PASSWORD", ""),
		Port:               GetEnvironment("PORT", "3000"),
	}
}
