package main

import (
	"Calls-Data-Engine/notionapi" // Matches `go.mod`
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func writeLastUpdated() {
	// Get current timestamp
	timestamp := time.Now().Format("01/02/2006 - 15:04:05")

	// Write timestamp to a second CSV file
	lastUpdatedFile := "last_updated.csv"
	csvFile, err := os.Create(lastUpdatedFile)
	if err != nil {
		fmt.Println("❌ ERROR: Creating Last Updated CSV file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"Last Updated", timestamp})

	fmt.Printf("✅ Timestamp successfully written to %s!\n", lastUpdatedFile)
}

func main() {
	// Load configuration first
	err := notionapi.LoadConfig()
	if err != nil {
		fmt.Println("❌ Error loading config:", err)
		return
	}

	fmt.Println("🚀 Fetching Notion Data...")
	data, err := notionapi.FetchNotionData()
	if err != nil {
		fmt.Println("❌ Error fetching data from Notion:", err)
		return
	}

	totalEntries := len(data)
	if totalEntries == 0 {
		fmt.Println("❌ No data found in Notion.")
		return
	}

	fmt.Printf("✅ Successfully retrieved %d entries from Notion.\n", totalEntries)

	// Save data to CSV
	csvFileName := "notion_data.csv"
	csvFile, err := os.Create(csvFileName)
	if err != nil {
		fmt.Println("❌ ERROR: Creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	headers := []string{
		"Property",
		"Missed Percentage",
		"Period Start",
		"Period End",
	}
	writer.Write(headers)

	// Process each entry
	for i, entry := range data {
		props, ok := entry["properties"].(map[string]interface{})
		if !ok {
			fmt.Printf("⚠️ Skipping entry %d: Invalid properties format.\n", i+1)
			continue
		}

		propertyStr := notionapi.GetRollupText(props, "Property (As Text)")
		missedPercentageStr := notionapi.GetFloatValue(props, "Missed Percentage")
		periodStartStr := notionapi.GetDateValue(props, "Period Start")
		periodEndStr := notionapi.GetDateValue(props, "Period End")

		row := []string{
			propertyStr,
			missedPercentageStr,
			periodStartStr,
			periodEndStr}
		writer.Write(row)

		// Print progress
		fmt.Printf("📊 Progress: %d/%d (%.2f%%)\n", i+1, totalEntries, float64(i+1)/float64(totalEntries)*100)
	}

	fmt.Printf("✅ Data successfully written to %s!\n", csvFileName) // Save data to CSV

	// Write last updated timestamp
	writeLastUpdated()
}
