package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vicsobdev/LittleBruteClient/brute"
	"go.uber.org/zap"
)

// init configures the standard logger to include file name and line number for easier debugging.
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Initialize structured logging with zap.
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize zap logger: ", err)
	}
	defer logger.Sync() // Flushes buffer, if any, to disk

	// Retrieve configurations and paths from user inputs and environment variables.
	wordListPath := getWordlistPath()
	proxyPath := getProxyPath()
	outPath := getOutPath()
	cfg := getRabbitConfig()

	// Initialize the brute force attack client with provided configurations.
	verbose := true // Enables verbose logging
	bruteClient, err := brute.NewBrute(wordListPath, proxyPath, outPath, verbose)
	if err != nil {
		log.Fatalf("Failed to create a new Brute instance: %v", err)
	}

	bruteClient.SetDebug(true)

	// Set up the logger within the brute client.
	err = bruteClient.SetLogger(*logger)
	if err != nil {
		log.Fatalf("Failed to set logger for Brute: %v", err)
	}

	// Initialize message queues for task distribution.
	err = bruteClient.InitializeQueues(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize queues: %v", err)
	}

	// Start the brute force attack process.
	err = bruteClient.Start()
	if err != nil {
		log.Fatalf("Failed to start the brute force attack: %v", err)
	}

	// Prevent the main function from exiting immediately.
	select {}
}

// getWordlistPath prompts the user for the path to the wordlist file.
func getWordlistPath() string {
	fmt.Print("Wordlist Path: ")
	var wordListPath string
	fmt.Scanln(&wordListPath)
	return wordListPath
}

// getProxyPath prompts the user for the path to the proxy file.
func getProxyPath() string {
	fmt.Print("Proxy Path: ")
	var proxyPath string
	fmt.Scanln(&proxyPath)
	return proxyPath
}

// getOutPath prompts the user for the path to the output file.
func getOutPath() string {
	fmt.Print("Output Path: ")
	var outPath string
	fmt.Scanln(&outPath)
	return outPath
}

// getRabbitConfig retrieves RabbitMQ configuration from environment variables.
func getRabbitConfig() brute.RabbitConfig {
	return brute.RabbitConfig{
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
	}
}
