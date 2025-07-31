package ai

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"log"
)

type GeminiRequest struct {
	Contents []Contents `json:"contents"`
}

type Contents struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

func CallGeminiAI(apiKey, prompt string) (string, error) {
	ctx := context.Background()

	cc := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}
	client, err := genai.NewClient(ctx, cc)
	if err != nil {
		log.Fatal(err)
	}

	var thinkingBudget int32 = 0

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &thinkingBudget, // Disable thinking
			},
		},
	)

	if err != nil {
		fmt.Println("error call gemini ai: ", err)
		return "", err
	}

	var resp string

	if len(result.Text()) > 0 {
		resp = result.Text()
	}

	return resp, nil
}
