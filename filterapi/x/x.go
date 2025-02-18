// Copyright Envoy AI Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

// Package x is an experimental package that provides the customizability of the AI Gateway filter.
package x

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	pb "github.com/envoyproxy/ai-gateway/docs/proposals/002-routing-processor-proposal/poc/semantic_cache/client/routing_processor"
	"github.com/envoyproxy/ai-gateway/filterapi"
	"google.golang.org/grpc"
)

// Router is the interface that wraps the Calculate method.
type Router interface {
	Calculate(requestHeaders map[string]string, requestBody []byte) (*filterapi.Backend, error)
	StoreInCache(ctx context.Context, model string, request []byte, response []byte) error
}

// SemanticCacheRouter implements the Router interface with semantic caching
type SemanticCacheRouter struct {
	defaultRouter Router
	cacheClient   pb.SemanticCacheServiceClient
	threshold     float32
	logger        *slog.Logger
	conn          *grpc.ClientConn
}

// SemanticCacheConfig holds configuration for semantic cache
type SemanticCacheConfig struct {
	Address             string  `yaml:"address"`
	SimilarityThreshold float32 `yaml:"similarity_threshold"`
}

// NewSemanticCacheRouter creates a new router with semantic caching capabilities
func NewSemanticCacheRouter(defaultRouter Router, config *filterapi.Config, cacheConfig *SemanticCacheConfig, logger *slog.Logger) (Router, error) {
	conn, err := grpc.Dial(cacheConfig.Address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cache service: %v", err)
	}

	return &SemanticCacheRouter{
		defaultRouter: defaultRouter,
		cacheClient:   pb.NewSemanticCacheServiceClient(conn),
		threshold:     cacheConfig.SimilarityThreshold,
		logger:        logger,
		conn:          conn,
	}, nil
}

// helper function to parse request messages
func parseRequest(messagesStr string) ([]*pb.Message, string, error) {
	// This is an OpenAI chat completion request body.
	var chatRequest struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}

	if err := json.Unmarshal([]byte(messagesStr), &chatRequest); err != nil {
		return nil, "", fmt.Errorf("failed to parse messages: %v", err)
	}

	messages := make([]*pb.Message, len(chatRequest.Messages))
	for i, msg := range chatRequest.Messages {
		messages[i] = &pb.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return messages, chatRequest.Model, nil
}

// helper function to parse response
func parseResponse(responseStr string) ([]*pb.Message, error) {
	var responseData struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal([]byte(responseStr), &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	responseMessages := make([]*pb.Message, len(responseData.Choices))
	for i, choice := range responseData.Choices {
		responseMessages[i] = &pb.Message{
			Role:    choice.Message.Role,
			Content: choice.Message.Content,
		}
	}
	return responseMessages, nil
}

// Calculate implements the Router interface
func (r *SemanticCacheRouter) Calculate(requestHeaders map[string]string, requestBody []byte) (*filterapi.Backend, error) {
	// First try to get from cache
	if len(requestBody) > 0 {
		messages, model, err := parseRequest(string(requestBody))
		r.logger.Info("Parsed messages", "messages", messages, "model", model)
		if err != nil {
			r.logger.Error("Failed to parse messages", "error", err)
			return nil, err
		}

		cacheResp, err := r.cacheClient.SearchCache(context.Background(), &pb.SearchRequest{
			Messages:            messages,
			SimilarityThreshold: r.threshold,
			Model:               model,
		})

		r.logger.Info("Cache response", "cacheResp", cacheResp)
		if err != nil {
			r.logger.Error("Cache search error", "error", err)
		} else if cacheResp.Found {
			r.logger.Info("Cache hit", "similarity_score", cacheResp.SimilarityScore)
			// Create OpenAI format response
			response := map[string]interface{}{
				"choices": []map[string]interface{}{{
					"message": map[string]interface{}{
						"role":    cacheResp.ResponseMessages[0].Role,
						"content": cacheResp.ResponseMessages[0].Content,
					},
					"finish_reason": "stop",
				}},
				"model": model,
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				r.logger.Error("Failed to marshal cache response", "error", err)
				return nil, err
			}
			// Store the response in the context for the processor to use
			r.logger.Info("Storing cache response in context", "response", response)
			requestHeaders["x-cache-response"] = string(responseBytes)
			return &filterapi.Backend{
				Name:   "semantic_cache",
				Schema: filterapi.VersionedAPISchema{Name: filterapi.APISchemaOpenAI},
				Weight: 100,
			}, nil
		}
	}

	// If no cache hit, delegate to default router
	return r.defaultRouter.Calculate(requestHeaders, requestBody)
}

// StoreInCache stores the response in the semantic cache
func (r *SemanticCacheRouter) StoreInCache(ctx context.Context, model string, request []byte, response []byte) error {
	r.logger.Info("Storing in cache", "model", model)
	requestMessages, _, err := parseRequest(string(request))
	if err != nil {
		r.logger.Error("Failed to parse request", "error", err)
		return err
	}
	// Parse the response message
	responseMessages, err := parseResponse(string(response))
	if err != nil {
		r.logger.Error("Failed to parse response", "error", err)
		return err
	}
	// Store in cache
	cacheRequest := &pb.StoreChatRequest{
		RequestMessages:  requestMessages,
		ResponseMessages: responseMessages,
		Model:            model,
	}
	output, err := r.cacheClient.StoreChat(ctx, cacheRequest)
	if err != nil {
		r.logger.Error("Failed to store in cache", "error", err)
		return err
	}
	r.logger.Info("Stored in cache", "output", output)
	return nil
}

// Close closes the gRPC connection
func (r *SemanticCacheRouter) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// NewCustomRouter is the function to create a custom router with semantic caching
var NewCustomRouter NewCustomRouterFn = func(defaultRouter Router, config *filterapi.Config) Router {
	cacheConfig := &SemanticCacheConfig{
		Address:             config.SemanticCache.Address,
		SimilarityThreshold: config.SemanticCache.SimilarityThreshold,
	}

	router, err := NewSemanticCacheRouter(defaultRouter, config, cacheConfig, slog.Default())
	if err != nil {
		slog.Error("Failed to create semantic cache router", "error", err)
		return defaultRouter
	}
	return router
}

// NewCustomRouterFn is the function signature for [NewCustomRouter].
//
// It accepts the exptproc config passed to the AI Gateway filter and returns a [Router].
// This is called when the new configuration is loaded.
//
// The defaultRouter can be used to delegate the calculation to the default router implementation.
type NewCustomRouterFn func(defaultRouter Router, config *filterapi.Config) Router
