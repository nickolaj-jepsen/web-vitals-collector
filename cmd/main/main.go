package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"log"
	"net/http"
	"time"
	"web-vitals-collector/internal"
)

type Server struct {
	batch chan *internal.VitalRow
}

type VitalDTO struct {
	Url  string   `json:"url"`
	Cls  *float64 `json:"cls"`
	Fcp  *float64 `json:"fcp"`
	Fid  *float64 `json:"fid"`
	Lcp  *float64 `json:"lcp"`
	Ttfb *float64 `json:"ttfb"`
}

func GenerateIdentifier(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	hasher := sha1.New()
	hasher.Write([]byte(IPAddress))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))

}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("only POST is supported"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var data VitalDTO
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse body"))
		return
	}

	s.batch <- &internal.VitalRow{
		Timestamp:  time.Now(),
		Url:        data.Url,
		Identifier: GenerateIdentifier(r),
		Cls:        data.Cls,
		Fcp:        data.Fcp,
		Fid:        data.Fid,
		Lcp:        data.Lcp,
		Ttfb:       data.Ttfb,
	}

	w.Write([]byte("ok"))
}

func main() {
	batchingChannel := make(chan *internal.VitalRow, 100)
	configuration := internal.GetConfiguration()

	var conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", configuration.ClickHouseHost, configuration.ClickHousePort)},
		Auth: clickhouse.Auth{
			Database: configuration.ClickHouseDatabase,
			Username: configuration.ClickHouseUsername,
			Password: configuration.ClickHousePassword,
		},
		DialTimeout:     time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})

	if err != nil {
		log.Fatal(err)
	}

	var s = Server{batch: batchingChannel}
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop() // release resources

		data := make([]*internal.VitalRow, 0)

		for {
			select {
			case i := <-batchingChannel:
				data = append(data, i)
			case <-ticker.C:
				err = internal.Insert(conn, data)
				if err != nil {
					fmt.Println("ERROR: Lost connection to clickhouse")
					log.Fatal(err)
				}
				data = make([]*internal.VitalRow, 0)
			}
		}
	}()

	http.HandleFunc("/", s.Handler)
	log.Print(fmt.Sprintf("Accepting requests on port %s", configuration.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configuration.Port), nil))
}
