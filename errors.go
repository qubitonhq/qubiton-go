package qubiton

import "fmt"

// ApiError represents an API error response.
type ApiError struct {
	StatusCode int
	Message    string
	Raw        map[string]interface{}
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

// IsAuthError returns true for 401/403 responses.
func (e *ApiError) IsAuthError() bool {
	return e.StatusCode == 401 || e.StatusCode == 403
}

// IsRateLimit returns true for 429 responses.
func (e *ApiError) IsRateLimit() bool {
	return e.StatusCode == 429
}
