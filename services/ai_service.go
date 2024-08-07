package services

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateAIResponse(prompt string) (string, error) {
	var apiKey = os.Getenv("GEMINI_API_KEY")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	modelAI := client.GenerativeModel("gemini-1.5-flash")
	resp, err := modelAI.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}

	textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected content type")
	}

	return string(textPart), nil
}

func CreatePrompt(username string, data interface{}, readme string) string {
	prompt := fmt.Sprintf("Buatkan pujian dalam bahasa gaul (gunakan lu-gue) yang berlebihan dan sangat melebih-lebihkan seakan dia adalah developer yang sangat hebat serta jenius dan project-projectnya sangat masterpiece untuk profil GitHub %v. Berikut detailnya: %s", username, data)

	if readme != "" {
		prompt += fmt.Sprintf(", Profile Markdown: ```%s```", readme)
	} else {
		prompt += ", Profile Markdown: Not Found"
	}

	prompt += ". (Berikan response dalam bahasa indonesia, dan gunakan emoticon supaya lebih menarik)"

	return prompt
}
