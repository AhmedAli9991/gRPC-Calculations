package utils

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func timeFormat(dateString string) (string, error) {
	// Splitting the date string to handle time zone separately
	parts := strings.Fields(dateString)
	if len(parts) < 5 {
		log.Fatal("Invalid date string format")
	}

	dateStr := fmt.Sprintf("%s %s %s %s %s", parts[0], parts[1], parts[2], parts[3], parts[4])

	// Parse the date string into a time.Time object
	date, err := time.Parse("Mon Jan 02 2006 15:04:05", dateStr)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
	}

	// Format the date as ISO 8601 (YYYY-MM-DD HH:MM:SS)
	formattedDate := date.Format("2006-01-02 15:04:05")

	fmt.Println("Formatted Date:", formattedDate)
	return formattedDate, nil
}
