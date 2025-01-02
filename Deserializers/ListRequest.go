package Deserializers

import (
	"encoding/json"
	"net/http"
)

// Deserialize ListTasks request body or query parameters to a map
func ListRequest(r *http.Request) (map[string]interface{}, error) {
	params := make(map[string]interface{})

	// Deserialize query parameters for page, limit, status, and sort
	page := r.URL.Query().Get("page")
	if page != "" {
		params["page"] = page
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		params["limit"] = limit
	}

	status := r.URL.Query().Get("status")
	if status != "" {
		params["status"] = status
	}

	sort := r.URL.Query().Get("sort")
	if sort != "" {
		params["sort"] = sort
	}

	// If JSON body exists (e.g., for POST), parse the body
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		var bodyParams map[string]interface{}
		if err := decoder.Decode(&bodyParams); err != nil {
			return nil, err
		}

		// Overwrite query params with body params if they exist
		for key, value := range bodyParams {
			params[key] = value
		}
	}

	return params, nil
}
