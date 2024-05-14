package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fuzzer",
	Short: "REST API Fuzzer",
	Long:  `A CLI tool to fuzz and test REST APIs for structure discovery and security vulnerabilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use either the 'probe' or 'breach' command to start.")
	},
	Example: `  # Probe the structure of a REST API
  fuzzer probe --url https://localhost:8080

  # Breach the security of a REST API with 8 concurrent threads
  fuzzer breach --url https://localhost:8080 --threads 8`,
}

func main() {
	logFile, err := os.OpenFile("fuzzer.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	rootCmd.AddCommand(probeCmd)
	rootCmd.AddCommand(breachCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
