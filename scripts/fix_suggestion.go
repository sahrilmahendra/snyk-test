package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type SnykIssue struct {
	Title       string `json:"title"`
	CodeSnippet string `json:"codeSnippet"`
}

type SnykResult struct {
	Issues []SnykIssue `json:"issues"`
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func callOpenAI(apiKey, prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	reqBody := OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   500,
		Temperature: 0.2,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openAIResp OpenAIResponse
	err = json.Unmarshal(respBytes, &openAIResp)
	if err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <snyk-result.json>", os.Args[0])
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	filePath := os.Args[1]
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var snykResult SnykResult
	err = json.Unmarshal(data, &snykResult)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(snykResult.Issues) == 0 {
		fmt.Println("No issues found by Snyk Code.")
		return
	}

	for _, issue := range snykResult.Issues {
		snippet := issue.CodeSnippet
		if snippet == "" {
			snippet = "Code snippet not available."
		}

		prompt := fmt.Sprintf(
			"You are a senior software security engineer. Review the following vulnerable Go code snippet and provide a detailed explanation and fix suggestion:\n\n%s",
			snippet,
		)

		fmt.Printf("Issue: %s\n", issue.Title)

		suggestion, err := callOpenAI(apiKey, prompt)
		if err != nil {
			fmt.Printf("Failed to get suggestion from OpenAI: %v\n", err)
			continue
		}

		fmt.Printf("ChatGPT Suggestion:\n%s\n\n", suggestion)
	}
}
