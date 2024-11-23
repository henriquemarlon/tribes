package router

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/rollmelette/rollmelette"
)

type ctxKey string

type AdvanceHandlerFunc func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error
type InspectHandlerFunc func(ctx context.Context, env rollmelette.EnvInspector) error

type AdvanceRequest struct {
	Path    string          `json:"path"`
	Payload json.RawMessage `json:"payload"`
}

type Router struct {
	AdvanceHandlers map[string]AdvanceHandlerFunc
	InspectHandlers map[string]struct {
		Handler InspectHandlerFunc
		Regex   *regexp.Regexp
	}
}

func NewRouter() *Router {
	return &Router{
		AdvanceHandlers: make(map[string]AdvanceHandlerFunc),
		InspectHandlers: make(map[string]struct {
			Handler InspectHandlerFunc
			Regex   *regexp.Regexp
		}),
	}
}

func (r *Router) parseAdvanceRawPayload(rawRequest []byte) (*AdvanceRequest, error) {
	var input AdvanceRequest
	if err := json.Unmarshal(rawRequest, &input); err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}
	return &input, nil
}

func (r *Router) HandleAdvance(path string, handler AdvanceHandlerFunc) {
	r.AdvanceHandlers[path] = handler
}

func (r *Router) Advance(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	input, err := r.parseAdvanceRawPayload(payload)
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}
	handler, ok := r.AdvanceHandlers[input.Path]
	if !ok {
		return fmt.Errorf("path '%s' not found", input.Path)
	}
	return handler(env, metadata, deposit, input.Payload)
}

func (r *Router) HandleInspect(path string, handler InspectHandlerFunc) {
	pattern := transformBracePattern(path)
	regex := regexp.MustCompile(pattern)
	r.InspectHandlers[pattern] = struct {
		Handler InspectHandlerFunc
		Regex   *regexp.Regexp
	}{handler, regex}
}

func transformBracePattern(pattern string) string {
	pattern = "^" + pattern + "$"
	return regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllStringFunc(pattern, func(m string) string {
		paramName := m[1 : len(m)-1]
		return fmt.Sprintf(`(?P<%s>[^/]+)`, paramName)
	})
}

func (r *Router) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	request := string(payload)
	for _, handler := range r.InspectHandlers {
		matches := handler.Regex.FindStringSubmatch(request)
		if matches != nil {
			ctx := context.Background()
			for i, name := range handler.Regex.SubexpNames() {
				if i > 0 && name != "" && i < len(matches) {
					ctx = context.WithValue(ctx, ctxKey(name), matches[i])
				}
			}
			return handler.Handler(ctx, env)
		}
	}
	return fmt.Errorf("no handler found for request: %s", request)
}
