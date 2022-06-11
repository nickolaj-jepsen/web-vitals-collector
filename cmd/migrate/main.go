package main

import (
	"errors"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"web-vitals-collector/internal"
)

func main() {
	configuration := internal.GetConfiguration()
	m, err := migrate.New(
		"file://./migrations",
		fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
			configuration.ClickHouseUsername,
			configuration.ClickHousePassword,
			configuration.ClickHouseHost,
			configuration.ClickHousePort,
			configuration.ClickHouseDatabase))
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil && err != errors.New("no change") {
		log.Fatal(err)
	}
}
