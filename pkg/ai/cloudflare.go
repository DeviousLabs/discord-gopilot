package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const cloudflareAPIURLTemplate = "https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s"

type CloudflareRequest struct {
	Prompt string `json:"prompt"`
}

type CloudflareResponse struct {
	Result struct {
		Response string `json:"response"`
	} `json:"result"`
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

func RunCloudflareAI(accountID, apiKey, message string) (*CloudflareResponse, error) {
	modelMap := map[int]string{
		1: "@cf/meta/llama-3.1-70b-instruct",
		2: "@hf/thebloke/deepseek-coder-6.7b-instruct-awq",
		3: "@hf/google/gemma-7b-it",
		4: "@hf/mistral/mistral-7b-instruct-v0.2",
		5: "@cf/qwen/qwen1.5-14b-chat-awq",
		6: "@cf/microsoft/phi-2",
		7: "@cf/stabilityai/stable-diffusion-xl-base-1.0",
	}

	modelIDStr := os.Getenv("MODEL")
	if modelIDStr == "" {
		return nil, fmt.Errorf("MODEL environment variable is not set")
	}

	modelID, err := strconv.Atoi(modelIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MODEL environment variable: %s", modelIDStr)
	}

	model, exists := modelMap[modelID]
	if !exists {
		return nil, fmt.Errorf("Invalid Model ID: %d", modelID)
	}

	url := fmt.Sprintf(cloudflareAPIURLTemplate, accountID, model)

	requestBody, err := json.Marshal(CloudflareRequest{Prompt: message})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var response CloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
