package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func postPRComment(comment string) {
	repo := os.Getenv("GITHUB_REPOSITORY")
	pr := os.Getenv("PR_NUMBER")
	token := os.Getenv("GITHUB_TOKEN")

	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments", repo, pr)

	data := map[string]string{"body": comment}
	body, _ := json.Marshal(data)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to post comment:", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("Posted suggestion to PR.")
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <snyk-result.json>", os.Args[0])
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("API_KEY environment variable is required")
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

		fmt.Printf("Issue detected in %s lines %d-%d: %s\n", file, startLine, endLine, msg)
		fmt.Println("Requesting fix suggestion from OpenAI...")

		suggestion, err := ai.CallGeminiAI(apiKey, prompt)
		if err != nil {
			fmt.Printf("Failed to get suggestion: %v\n\n", err)
		} else {
			fmt.Println("Fix suggestion:")
			comment := fmt.Sprintf("🔍 **AI Fix Suggestion** for `%s` at line %d-%d:\n\n%s", file, startLine, endLine, suggestion)
			postPRComment(comment)
			fmt.Println(strings.Repeat("-", 80))
		}

		time.Sleep(62 * time.Second)
	}
}
