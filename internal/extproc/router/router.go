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
	if modelName == "auto" {
		// Currently, we only support OpenAI backend for auto model selection.
		for i := range r.rules {
			_rule := &r.rules[i]
			for _, backend := range _rule.Backends {
				if backend.Schema.Name == filterapi.APISchemaOpenAI && backend.AutoModelConfig != nil {
					if chatReq, ok := requestBody.(*openai.ChatCompletionRequest); ok {
						semanticService := modelselect.NewSemanticProcessorService(
							backend.AutoModelConfig.SemanticProcessorServiceURL,
							backend.AutoModelConfig.SimpleModels,
							backend.AutoModelConfig.StrongModels,
						)
						selectedModel, err := semanticService.SelectModel(chatReq)
						if err != nil {
							return nil, "", errors.New("failed to select model")
						}
						fmt.Printf("Selected model: %s\n", selectedModel)
						headers[r.config.ModelNameHeaderKey] = selectedModel
						modelName = selectedModel
					}
					break
				}
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

	return r.selectBackendFromRule(rule), modelName, nil
}

func (r *router) selectBackendFromRule(rule *filterapi.RouteRule) (backend *filterapi.Backend) {
	// Each backend has a weight, so we randomly select depending on the weight.
	// This is a pretty naive implementation and can be buggy, so fix it later.
	totalWeight := 0
	for _, b := range rule.Backends {
		totalWeight += b.Weight
	}
	if totalWeight == 0 {
		return &rule.Backends[0]
	}
	selected := r.rng.Intn(totalWeight)
	for i := range rule.Backends {
		b := &rule.Backends[i]
		if selected < b.Weight {
			return b
		}
		selected -= b.Weight
	}
	return &rule.Backends[0]
}
