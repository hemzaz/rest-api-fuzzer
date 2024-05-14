package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var (
	probeBaseURL       string
	probeAuthType      string
	probeAuthToken     string
	probeRateLimit     int
	probeCustomHeaders []string
)

func init() {
	probeCmd.Flags().StringVarP(&probeBaseURL, "url", "u", "", "Base URL of the API (required)")
	probeCmd.Flags().StringVarP(&probeAuthType, "auth-type", "a", "", "Authentication type (e.g., 'Bearer', 'Basic')")
	probeCmd.Flags().StringVarP(&probeAuthToken, "auth-token", "t", "", "Authentication token")
	probeCmd.Flags().IntVarP(&probeRateLimit, "rate-limit", "r", 0, "Rate limit in requests per second")
	probeCmd.Flags().StringSliceVar(&probeCustomHeaders, "header", nil, "Custom headers (key:value)")
	probeCmd.MarkFlagRequired("url")
}

var probeCmd = &cobra.Command{
	Use:     "probe",
	Short:   "Probe the REST API structure",
	Long:    `Probes the structure of a REST API by sending requests to common endpoints and recording the responses.`,
	Example: `  # Probe the structure of a REST API\n  fuzzer probe --url https://localhost:8080`,
	Run: func(cmd *cobra.Command, args []string) {
		if probeBaseURL == "" {
			log.Fatal("Base URL is required")
		}
		probeAPI(probeBaseURL)
	},
}

func probeAPI(baseURL string) {
	client := resty.New()

	if probeAuthType != "" && probeAuthToken != "" {
		if probeAuthType == "Bearer" {
			client.SetAuthToken(probeAuthToken)
		} else if probeAuthType == "Basic" {
			client.SetHeader("Authorization", "Basic "+probeAuthToken)
		}
	}

	for _, header := range probeCustomHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			client.SetHeader(parts[0], parts[1])
		}
	}

	url := fmt.Sprintf("%s", baseURL)

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

	for _, method := range methods {
		resp, err := client.R().Execute(method, url)
		if err != nil {
			log.Printf("Error probing %s %s: %v", method, url, err)
			continue
		}

		if resp.StatusCode() != 405 {
			fmt.Printf("Discovered endpoint: %s %s\n", method, url)
			storeProbeResult(url, method, resp)
		}

		if probeRateLimit > 0 {
			time.Sleep(time.Second / time.Duration(probeRateLimit))
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}
