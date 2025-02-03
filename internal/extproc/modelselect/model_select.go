package modelselect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/envoyproxy/ai-gateway/internal/apischema/openai"
)

// SemanticProcessorService handles communication with the semantic analysis service
type SemanticProcessorService struct {
	url          string
	simpleModels []string
	strongModels []string
}

// NewSemanticProcessorService creates a new semantic processor service client
func NewSemanticProcessorService(url string, simpleModels, strongModels []string) *SemanticProcessorService {
	return &SemanticProcessorService{
		url:          url,
		simpleModels: simpleModels,
		strongModels: strongModels,
	}
}

// SelectModel calls the SemanticProcessor service and returns an appropriate model
func (b *SemanticProcessorService) SelectModel(request *openai.ChatCompletionRequest) (string, error) {
	// Extract text from messages
	var chatText string
	for _, msg := range request.Messages {
		if msg.Type == openai.ChatMessageRoleUser {
			if userMsg, ok := msg.Value.(openai.ChatCompletionUserMessageParam); ok {
				if content, ok := userMsg.Content.Value.(string); ok {
					chatText += content + " "
				}
			}
		}
	}
	// Call SemanticProcessor service
	req := struct {
		Text         string   `json:"text"`
		SimpleModels []string `json:"simple_models,omitempty"`
		StrongModels []string `json:"strong_models,omitempty"`
	}{
		Text:         chatText,
		SimpleModels: b.simpleModels,
		StrongModels: b.strongModels,
	}

	// Ensure SimpleModels and StrongModels are not nil
	if req.SimpleModels == nil {
		req.SimpleModels = []string{}
	}
	if req.StrongModels == nil {
		req.StrongModels = []string{}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(b.url, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("error calling semantic processor service: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		SelectedModel string `json:"selected_model"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("error decoding response: %v\n", err)
		return "", err
	}
	return result.SelectedModel, nil
}
