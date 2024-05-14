package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"

	"go.etcd.io/bbolt"
)

func generateNewPayload(prevPayload map[string]interface{}) map[string]interface{} {
	newPayload := make(map[string]interface{})
	for k, v := range prevPayload {
		switch v := v.(type) {
		case string:
			newPayload[k] = mutateString(v)
		case int:
			newPayload[k] = mutateInt(v)
		case bool:
			newPayload[k] = !v
		case []interface{}:
			newPayload[k] = mutateArray(v)
		case map[string]interface{}:
			newPayload[k] = generateNewPayload(v)
		default:
			newPayload[k] = v
		}
	}
	applyPatterns(newPayload)
	return newPayload
}

func mutateString(s string) string {
	mutations := []string{
		s + " OR 1=1",
		s + " AND 1=0",
		s + "'; DROP TABLE users; --",
		"<script>alert('XSS')</script>",
		strings.Repeat(s, 2),
	}
	return mutations[rand.Intn(len(mutations))]
}

func mutateInt(i int) int {
	return i + rand.Intn(100) - 50
}

func mutateArray(a []interface{}) []interface{} {
	if len(a) == 0 {
		return a
	}
	index := rand.Intn(len(a))
	a[index] = mutateValue(a[index])
	return a
}

func mutateValue(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		return mutateString(v)
	case int:
		return mutateInt(v)
	case bool:
		return !v
	default:
		return v
	}
}

func applyPatterns(payload map[string]interface{}) {
	patterns := []map[string]interface{}{
		{"username": "admin", "password": "password"},
		{"username": "admin", "password": "admin"},
		{"username": "' OR '1'='1", "password": "' OR '1'='1"},
		{"username": "root", "password": "toor"},
	}
	for k, v := range patterns[rand.Intn(len(patterns))] {
		payload[k] = v
	}
}

func getAIPayloads() []map[string]interface{} {
	db, err := openDB()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var payloads []map[string]interface{}

	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			var responseData map[string]interface{}
			if err := json.Unmarshal(v, &responseData); err != nil {
				return err
			}

			if status, ok := responseData["status"].(float64); ok && status >= 400 {
				payloads = append(payloads, generateNewPayload(responseData["payload"].(map[string]interface{})))
			}
			return nil
		})
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to retrieve responses: %v", err)
	}

	if len(payloads) == 0 {
		payloads = []map[string]interface{}{
			{"test": "<script>alert(1)</script>"},
			{"username": "' OR 1=1 --", "password": "password"},
			{"long_string": randomString(5000)},
		}
	}
	return payloads
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
