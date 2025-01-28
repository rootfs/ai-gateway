package router

import (
	"encoding/json"
	"fmt"
	"regexp"

	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

	"github.com/envoyproxy/ai-gateway/filterconfig"
	"github.com/envoyproxy/ai-gateway/internal/apischema/openai"
)

// RequestBodyParser is a function that parses the body of the request.
type RequestBodyParser func(path string, body *extprocv3.HttpBody) (modelName string, rb RequestBody, err error)

// NewRequestBodyParser creates a new RequestBodyParser based on the schema.
func NewRequestBodyParser(schema filterconfig.VersionedAPISchema) (RequestBodyParser, error) {
	if schema.Name == filterconfig.APISchemaOpenAI {
		return openAIParseBody, nil
	}
	return nil, fmt.Errorf("unsupported API schema: %s", schema)
}

// RequestBody is the union of all request body types.
type RequestBody any

// UpdateModel updates the model field in the request body and returns it as bytes
func UpdateModel(model string, body []byte) ([]byte, error) {
	// This regex looks for "model": "any-characters-here"
	//TODO this should really go to the openai json marshalling code
	re := regexp.MustCompile(`"model"\s*:\s*"[^"]*"`)
	newModelField := fmt.Sprintf(`"model": "%s"`, model)
	return re.ReplaceAll(body, []byte(newModelField)), nil
}

// openAIParseBody parses the body of the request for OpenAI.
func openAIParseBody(path string, body *extprocv3.HttpBody) (modelName string, rb RequestBody, err error) {
	if path == "/v1/chat/completions" {
		var openAIReq openai.ChatCompletionRequest
		if err := json.Unmarshal(body.Body, &openAIReq); err != nil {
			return "", nil, fmt.Errorf("failed to unmarshal body: %w", err)
		}
		return openAIReq.Model, &openAIReq, nil
	} else {
		return "", nil, fmt.Errorf("unsupported path: %s", path)
	}
}
