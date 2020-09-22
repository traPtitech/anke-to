package model

import (
	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
	"log"
	"os"
)

func GetLogger() (*logging.Logger, error) {

	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Println("no GCP Project ID")
		return nil, nil
	}

	// Sets the name of the log to write to.
	logName := os.Getenv("LOG_NAME")
	if logName == "" {
		logName = "anke-to-log"
	}

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	logger := client.Logger(logName)

	// Logs "hello world", log entry is visible at
	// Stackdriver Logs.
	logger.StandardLogger(logging.Info).Println("hello world")

	return logger, nil
}
