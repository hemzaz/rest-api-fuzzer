package main

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
	"go.etcd.io/bbolt"
)

const (
	dbFile          = "responses.db"
	bucketName      = "Responses"
	probeBucketName = "Probes"
)

func openDB() (*bbolt.DB, error) {
	return bbolt.Open(dbFile, 0600, nil)
}

func storeResponse(payload interface{}, resp *resty.Response) {
	db, err := openDB()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		id, _ := b.NextSequence()
		responseData := map[string]interface{}{
			"payload":  payload,
			"response": resp.String(),
			"status":   resp.StatusCode(),
		}
		responseBytes, err := json.Marshal(responseData)
		if err != nil {
			return err
		}

		return b.Put(itob(int(id)), responseBytes)
	})

	if err != nil {
		log.Fatalf("Failed to store response: %v", err)
	}
}

func storeProbeResult(url, method string, resp *resty.Response) {
	db, err := openDB()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(probeBucketName))
		if err != nil {
			return err
		}

		id, _ := b.NextSequence()
		probeData := map[string]interface{}{
			"url":      url,
			"method":   method,
			"response": resp.String(),
			"status":   resp.StatusCode(),
		}
		probeBytes, err := json.Marshal(probeData)
		if err != nil {
			return err
		}

		return b.Put(itob(int(id)), probeBytes)
	})

	if err != nil {
		log.Fatalf("Failed to store probe result: %v", err)
	}
}

func itob(v int) []byte {
	b := make([]byte, 8)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	return b
}
