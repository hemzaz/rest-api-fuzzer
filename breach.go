package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
)

var (
	breachBaseURL       string
	threads             int
	breachAuthType      string
	breachAuthToken     string
	breachCustomHeaders []string
)

func init() {
	breachCmd.Flags().StringVarP(&breachBaseURL, "url", "u", "", "Base URL of the API (required)")
	breachCmd.Flags().IntVarP(&threads, "threads", "t", 4, "Number of concurrent threads (default is 4)")
	breachCmd.Flags().StringVarP(&breachAuthType, "auth-type", "a", "", "Authentication type (e.g., 'Bearer', 'Basic')")
	breachCmd.Flags().StringVarP(&breachAuthToken, "auth-token", "k", "", "Authentication token")
	breachCmd.Flags().StringSliceVar(&breachCustomHeaders, "header", nil, "Custom headers (key:value)")
	breachCmd.MarkFlagRequired("url")
}

var breachCmd = &cobra.Command{
	Use:     "breach",
	Short:   "Breach the REST API for security testing",
	Long:    `Attempts to breach the security of a REST API by sending malicious payloads to discovered endpoints.`,
	Example: `  # Breach the security of a REST API with 4 concurrent threads\n  fuzzer breach --url https://localhost:8080\n\n  # Breach the security of a REST API with 8 concurrent threads\n  fuzzer breach --url https://localhost:8080 --threads 8`,
	Run: func(cmd *cobra.Command, args []string) {
		if breachBaseURL == "" {
			log.Fatal("Base URL is required")
		}
		breachAPI(breachBaseURL, threads)
	},
}

func breachAPI(baseURL string, threads int) {
	client := resty.New()

	if breachAuthType != "" && breachAuthToken != "" {
		if breachAuthType == "Bearer" {
			client.SetAuthToken(breachAuthToken)
		} else if breachAuthType == "Basic" {
			client.SetHeader("Authorization", "Basic "+breachAuthToken)
		}
	}

	for _, header := range breachCustomHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			client.SetHeader(parts[0], parts[1])
		}
	}

	rand.Seed(time.Now().UnixNano())

	endpoints := getProbedEndpoints()

	payloads := getAIPayloads()

	payloadChan := make(chan map[string]interface{}, len(payloads)*len(endpoints))
	responseChan := make(chan *resty.Response, len(payloads)*len(endpoints))
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(client, payloadChan, responseChan, &wg)
	}

	go func() {
		for _, endpoint := range endpoints {
			for _, payload := range payloads {
				payloadChan <- map[string]interface{}{
					"url":     fmt.Sprintf("%s/%s", baseURL, endpoint["url"]),
					"method":  endpoint["method"],
					"payload": payload,
				}
			}
		}
		close(payloadChan)
	}()

	go func() {
		wg.Wait()
		close(responseChan)
	}()

	for resp := range responseChan {
		fmt.Printf("Payload: %v\nResponse: %s\n", resp.Request.Body, resp)
		storeResponse(resp.Request.Body, resp)
	}
}

func worker(client *resty.Client, payloadChan <-chan map[string]interface{}, responseChan chan<- *resty.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	for payload := range payloadChan {
		url := payload["url"].(string)
		method := payload["method"].(string)
		body := payload["payload"]

		resp, err := client.R().
			SetBody(body).
			Execute(method, url)
		if err != nil {
			log.Printf("Error sending request to %s: %v", url, err)
			continue
		}
		responseChan <- resp
		time.Sleep(500 * time.Millisecond)
	}
}

func getProbedEndpoints() []map[string]string {
	db, err := openDB()
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var endpoints []map[string]string

	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(probeBucketName))
		if b == nil {
			return nil
		}

		b.ForEach(func(k, v []byte) error {
			var endpoint map[string]string
			if err := json.Unmarshal(v, &endpoint); err != nil {
				return err
			}
			endpoints = append(endpoints, endpoint)
			return nil
		})
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to retrieve probe data: %v", err)
	}

	return endpoints
}
