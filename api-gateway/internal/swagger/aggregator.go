package swagger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var httpClient = &http.Client{}

type ServiceSpec struct {
	Name string
	URL  string
}

var Services = []ServiceSpec{
	{Name: "user-service", URL: "http://user-service:8081/swagger/doc.json"},
	{Name: "product-service", URL: "http://product-service:8082/swagger/doc.json"},
	{Name: "warehouse-service", URL: "http://warehouse-service:8083/swagger/doc.json"},
	{Name: "merchant-service", URL: "http://merchant-service:8084/swagger/doc.json"},
	{Name: "transaction-service", URL: "http://transaction-service:8085/swagger/doc.json"},
	{Name: "notification-service", URL: "http://notification-service:8086/swagger/doc.json"},
}

func FetchSpec(url string) (map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("[Swagger] FetchSpec - 1: failed to create request for %s: %w", url, err)
	}

	req.Header.Set("X-Gateway", "warehouse-api-gateway")
	req.Header.Set("X-Internal-Request", "true")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[Swagger] FetchSpec - 2: failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Swagger] FetchSpec - 3: failed to read body from %s: %w", url, err)
	}

	var spec map[string]any
	if err := json.Unmarshal(body, &spec); err != nil {
		return nil, fmt.Errorf("[Swagger] FetchSpec - 4: failed to unmarshal spec from %s: %w", url, err)
	}

	return spec, nil
}

func AggregateSpecs(gatewayURL string) (map[string]any, error) {
	merged := map[string]any{
		"swagger":     "2.0",
		"host":        "localhost:8080",
		"basePath":    "/",
		"definitions": map[string]any{},
		"info": map[string]any{
			"title":   "Warehouse Microservices API",
			"version": "1.0.0",
			"servers": []any{
				map[string]any{"url": gatewayURL, "description": "API Gateway"},
			},
			"paths": map[string]any{},
		},
		"components": map[string]any{
			"schemas": map[string]any{},
		},
	}

	type result struct {
		ss   ServiceSpec
		spec map[string]any
		err  error
	}

	results := make(chan result, len(Services))

	for _, ss := range Services {
		go func(ss ServiceSpec) {
			spec, err := FetchSpec(ss.URL)
			results <- result{ss, spec, err}
		}(ss)
	}

	mergedPaths := merged["paths"].(map[string]any)
	mergedSchemas := merged["components"].(map[string]any)["schemas"].(map[string]any)

	var mu sync.Mutex
	var errs []error

	for range Services {
		r := <-results
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}

		// build rename map — support definitions (2.0) dan components.schemas (3.0)
		renameMap := map[string]string{}

		// Swagger 2.0
		if defs, ok := r.spec["definitions"].(map[string]any); ok {
			for name := range defs {
				renameMap[name] = r.ss.Name + "_" + name
			}
		}
		// OpenAPI 3.0
		if components, ok := r.spec["components"].(map[string]any); ok {
			if schemas, ok := components["schemas"].(map[string]any); ok {
				for name := range schemas {
					renameMap[name] = r.ss.Name + "_" + name
				}
			}
		}

		mu.Lock()
		// merge paths
		if paths, ok := r.spec["paths"].(map[string]any); ok {
			for path, pathItem := range paths {
				mergedPaths[path] = rewriteRefs(pathItem, renameMap)
			}
		}

		// merge schemas dari definitions (2.0) → components/schemas (3.0)
		if defs, ok := r.spec["definitions"].(map[string]any); ok {
			for name, schema := range defs {
				mergedSchemas[r.ss.Name+"_"+name] = rewriteRefs(schema, renameMap)
			}
		}
		// merge schemas dari components.schemas (3.0)
		if components, ok := r.spec["components"].(map[string]any); ok {
			if schemas, ok := components["schemas"].(map[string]any); ok {
				for name, schema := range schemas {
					mergedSchemas[r.ss.Name+"_"+name] = rewriteRefs(schema, renameMap)
				}
			}
		}
		mu.Unlock()
	}

	if len(errs) > 0 {
		return merged, fmt.Errorf("[Swagger] AggregateSpecs - 1: some services unavailable: %v", errs[0])
	}
	return merged, nil
}

// rewriteRefs mengganti semua $ref lama dengan nama schema baru (prefixed)
func rewriteRefs(v any, renameMap map[string]string) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, v2 := range val {
			if k == "$ref" {
				if ref, ok := v2.(string); ok {
					// handle Swagger 2.0: #/definitions/Foo → #/components/schemas/service_Foo
					if strings.HasPrefix(ref, "#/definitions/") {
						oldName := strings.TrimPrefix(ref, "#/definitions/")
						if newName, ok := renameMap[oldName]; ok {
							out[k] = "#/components/schemas/" + newName
							continue
						}
					}
					// handle OpenAPI 3.0: #/components/schemas/Foo → #/components/schemas/service_Foo
					if strings.HasPrefix(ref, "#/components/schemas/") {
						oldName := strings.TrimPrefix(ref, "#/components/schemas/")
						if newName, ok := renameMap[oldName]; ok {
							out[k] = "#/components/schemas/" + newName
							continue
						}
					}
				}
			}
			out[k] = rewriteRefs(v2, renameMap)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = rewriteRefs(item, renameMap)
		}
		return out
	}
	return v
}
