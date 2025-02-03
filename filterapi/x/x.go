// Package x is an experimental package that provides the customizability of the AI Gateway filter.
package x

import "github.com/envoyproxy/ai-gateway/filterapi"

// NewCustomRouter is the function to create a custom router over the default router.
// This is nil by default and can be set by the custom build of external processor.
var NewCustomRouter NewCustomRouterFn

// NewCustomRouterFn is the function signature for [NewCustomRouter].
//
// It accepts the exptproc config passed to the AI Gateway filter and returns a [Router].
// This is called when the new configuration is loaded.
//
// The defaultRouter can be used to delegate the calculation to the default router implementation.
type NewCustomRouterFn func(defaultRouter Router, config *filterapi.Config) Router

// Router is the interface for the router.
//
// Router must be goroutine-safe as it is shared across multiple requests.
type Router interface {
	// Calculate determines the backend to route to based on the request headers.
	//
	// The request headers include the populated [filterapi.Config.ModelNameHeaderKey]
	// with the parsed model name based on the [filterapi.Config] given to the NewCustomRouterFn.
	//
	// The request body is passed when there is a semantic processor service in place to handle semanatic related analysis.
	//
	// Returns the backend.
	Calculate(requestHeaders map[string]string, requestBody any) (backend *filterapi.Backend, model string, err error)
}
