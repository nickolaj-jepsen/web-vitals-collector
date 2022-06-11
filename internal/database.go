package internal

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"time"
)

type VitalRow struct {
	Timestamp  time.Time `ch:"timestamp"`
	Url        string    `ch:"url"`
	Identifier string    `ch:"identifier"`
	Cls        *float64  `ch:"cls"`
	Fcp        *float64  `ch:"fcp"`
	Fid        *float64  `ch:"fid"`
	Lcp        *float64  `ch:"lcp"`
	Ttfb       *float64  `ch:"ttfb"`
}

func Insert(conn clickhouse.Conn, rows []*VitalRow) error {
	var batch, err = conn.PrepareBatch(context.Background(), "INSERT INTO vitals")
	if err != nil {
		return err
	}
	for _, row := range rows {
		if err := batch.AppendStruct(row); err != nil {
			return err
		}
	}
	return batch.Send()
}
