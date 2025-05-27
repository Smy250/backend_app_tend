package apis

import (
	"context"

	"google.golang.org/genai"
)

func ConsultGemini(gemini_key string, consult string) (*genai.GenerateContentResponse, error) {

	var ctx = context.Background()
	var ai_Models = []string{
		"gemini-2.0-flash",
		"gemini-2.5-flash-preview-04-17",
		"gemini-2.5-pro-preview-05-06",
	}

	client, err1 := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  gemini_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err1 != nil {
		return nil, err1
	}

	result, err2 := client.Models.GenerateContent(
		ctx,
		ai_Models[0],
		genai.Text(consult),
		nil,
	)
	if err2 != nil {
		return nil, err2
	}

	return result, nil
}
