package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"snyk/ai"
	"strings"
	"time"
)

// SarifResult JSON result from snyk code
type SarifResult struct {
	Runs []struct {
		Results []struct {
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
			Locations []struct {
				PhysicalLocation struct {
					ArtifactLocation struct {
						URI string `json:"uri"`
					} `json:"artifactLocation"`
					Region struct {
						StartLine   int `json:"startLine"`
						EndLine     int `json:"endLine"`
						StartColumn int `json:"startColumn"`
						EndColumn   int `json:"endColumn"`
					} `json:"region"`
				} `json:"physicalLocation"`
			} `json:"locations"`
		} `json:"results"`
	} `json:"runs"`
}

func extractLines(filename string, start, end int) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read file %s: %v\n", filename, err)
		return ""
	}

	lines := strings.Split(string(content), "\n")
	if start < 1 || start > len(lines) {
		return ""
	}
	if end > len(lines) {
		end = len(lines)
	}

	snippet := lines[start-1 : end]
	return strings.Join(snippet, "\n")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <snyk-result.json>", os.Args[0])
	}

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var result SarifResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(result.Runs) == 0 || len(result.Runs[0].Results) == 0 {
		fmt.Println("No issues found by Snyk Code.")
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("API_KEY environment variable is required")
	}

	for _, issue := range result.Runs[0].Results {
		msg := issue.Message.Text
		if len(issue.Locations) == 0 {
			fmt.Printf("Issue: %s\nNo location info available.\n\n", msg)
			continue
		}

		loc := issue.Locations[0].PhysicalLocation
		file := loc.ArtifactLocation.URI
		startLine := loc.Region.StartLine
		endLine := loc.Region.EndLine

		codeSnippet := extractLines(file, startLine, endLine)
		if codeSnippet == "" {
			codeSnippet = "<code snippet not available>"
		}

		prompt := fmt.Sprintf(
			"You are a senior software engineer. Review this vulnerable Golang code snippet and provide detailed explanation and fix suggestion:\n\nFile: %s\n\n%s\n\nIssue: %s",
			file, codeSnippet, msg,
		)

		fmt.Printf("Issue detected in %s lines %d-%d, column %d-%d: %s\n", file, startLine, endLine, loc.Region.StartColumn, loc.Region.EndColumn, msg)
		fmt.Println("Requesting fix suggestion from AI...")

		suggestion, err := ai.CallGeminiAI(apiKey, prompt)
		if err != nil {
			fmt.Printf("Failed to get suggestion: %v\n\n", err)
		} else {
			fmt.Println("Fix suggestion:")
			fmt.Println(suggestion)
			fmt.Println(strings.Repeat("-", 90))
		}

		time.Sleep(62 * time.Second)
	}
}
