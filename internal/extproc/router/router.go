package router

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/rand"

	"github.com/envoyproxy/ai-gateway/filterapi"
	"github.com/envoyproxy/ai-gateway/filterapi/x"

	"github.com/envoyproxy/ai-gateway/internal/apischema/openai"
	"github.com/envoyproxy/ai-gateway/internal/extproc/modelselect"
)

// router implements [filterapi.Router].
type router struct {
	rules  []filterapi.RouteRule
	rng    *rand.Rand
	config *filterapi.Config
}

// NewRouter creates a new [extprocapi.Router] implementation for the given config.
func NewRouter(config *filterapi.Config, newCustomFn x.NewCustomRouterFn) (x.Router, error) {
	r := &router{
		rules:  config.Rules,
		rng:    rand.New(rand.NewSource(uint64(time.Now().UnixNano()))), //nolint:gosec
		config: config,
	}
	if newCustomFn != nil {
		customRouter := newCustomFn(r, config)
		return customRouter, nil
	}
	return r, nil
}

// Calculate implements [extprocapi.Router.Calculate].
func (r *router) Calculate(headers map[string]string, requestBody any) (backend *filterapi.Backend, model string, err error) {
	modelName, ok := headers[r.config.ModelNameHeaderKey]
	if !ok {
		return nil, "", errors.New("model name not found in headers")
	}
	var rule *filterapi.RouteRule

	// Handle auto model selection for OpenAI backend
	// Handle auto model selection for OpenAI backend
	if modelName == "auto" {
		var simpleModels []string
		var strongModels []string
		var semanticURL string
		// Currently, we only support OpenAI backend for auto model selection.
		for i := range r.rules {
			_rule := &r.rules[i]
			for _, backend := range _rule.Backends {
				if backend.Schema.Name == filterapi.APISchemaOpenAI && backend.AutoModelConfig != nil {
					// Append models from the backend to the lists
					simpleModels = append(simpleModels, backend.AutoModelConfig.SimpleModels...)
					strongModels = append(strongModels, backend.AutoModelConfig.StrongModels...)
					semanticURL = backend.AutoModelConfig.SemanticProcessorServiceURL
				}
			}

			if chatReq, ok := requestBody.(*openai.ChatCompletionRequest); ok {
				semanticService := modelselect.NewSemanticProcessorService(
					semanticURL,
					simpleModels,
					strongModels,
				)
				selectedModel, err := semanticService.SelectModel(chatReq)
				if err != nil || selectedModel == "" {
					return nil, "", errors.New("failed to select model")
				}
				fmt.Printf("Selected model: %s\n", selectedModel)
				headers[r.config.ModelNameHeaderKey] = selectedModel
				modelName = selectedModel
			}
		}
	}

	for i := range r.rules {
		_rule := &r.rules[i]
		for _, hdr := range _rule.Headers {
			v, ok := headers[string(hdr.Name)]
			// Currently, we only do the exact matching.
			if ok && v == hdr.Value {
				rule = _rule
				break
			}
		}
	}
	if rule == nil {
		return nil, "", errors.New("no matching rule found")
	}

	return r.selectBackendFromRule(rule, modelName), modelName, nil
}

// selectBackendFromRule selects a backend based on the model name and weight.
func (r *router) selectBackendFromRule(rule *filterapi.RouteRule, modelName string) (backend *filterapi.Backend) {
	// Filter backends that contain the selected model
	var candidateBackends []filterapi.Backend
	for i := range rule.Backends {
		b := &rule.Backends[i]
		for _, model := range b.AutoModelConfig.SimpleModels {
			if model == modelName {
				candidateBackends = append(candidateBackends, *b)
				break
			}
		}
		for _, model := range b.AutoModelConfig.StrongModels {
			if model == modelName {
				candidateBackends = append(candidateBackends, *b)
				break
			}
		}
	}

	// If no backends contain the selected model, fall back to weight-based selection
	if len(candidateBackends) == 0 {
		return r.selectBackendByWeight(rule.Backends)
	}

	// Select a backend from the candidates based on weight
	return r.selectBackendByWeight(candidateBackends)
}

// selectBackendByWeight selects a backend based on weight.
func (r *router) selectBackendByWeight(backends []filterapi.Backend) *filterapi.Backend {
	totalWeight := 0
	for _, b := range backends {
		totalWeight += b.Weight
	}
	if totalWeight == 0 {
		return &backends[0]
	}
	selected := r.rng.Intn(totalWeight)
	for _, b := range backends {
		if selected < b.Weight {
			return &b
		}
		selected -= b.Weight
	}
	return &backends[0]
}
