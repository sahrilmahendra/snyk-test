package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DeepseekRequest struct {
	Model    string     `json:"model"`
	Messages []Messages `json:"messages"`
	Stream   bool       `json:"stream"`
}

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepseekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func CallDeepseekAI(apiKey, prompt string) (string, error) {
	url := "https://api.deepseek.com/chat/completions"

	reqBody := DeepseekRequest{
		Model: "deepseek-chat",
		Messages: []Messages{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
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

	fmt.Println("Logged response: ", string(respBytes))

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %s", string(respBytes))
	}

	var deepSeekResponse DeepseekResponse
	err = json.Unmarshal(respBytes, &deepSeekResponse)
	if err != nil {
		return "", err
	}

	if len(deepSeekResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in deepseek response")
	}

	return deepSeekResponse.Choices[0].Message.Content, nil
}
