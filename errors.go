package qubiton

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors that may be wrapped by *ApiError. Use errors.Is to test for them.
var (
	// ErrAuth indicates an authentication or authorization failure (HTTP 401 or 403).
	ErrAuth = errors.New("authentication failed")
	// ErrRateLimit indicates the request was rate-limited (HTTP 429).
	ErrRateLimit = errors.New("rate limit exceeded")
	// ErrServerError indicates the API returned a 5xx response.
	ErrServerError = errors.New("server error")
	// ErrNotFound indicates the requested resource was not found (HTTP 404).
	ErrNotFound = errors.New("not found")
	// ErrValidation indicates the request was rejected as invalid (HTTP 400 or 422).
	ErrValidation = errors.New("validation failed")
)

// ApiError represents a non-2xx response from the QubitOn API.
// The wrapped sentinel error (accessible via errors.Is or Unwrap) classifies the failure.
type ApiError struct {
	// StatusCode is the HTTP status code returned by the server.
	StatusCode int
	// Message is the human-readable error message extracted from the response body.
	Message string
	// Raw is the parsed JSON body, if the response was JSON.
	Raw map[string]interface{}
	// RetryAfter is the duration suggested by the server's Retry-After header (only set on 429/503).
	RetryAfter int
	// sentinel is the wrapped sentinel error for use with errors.Is.
	sentinel error
}

// Error implements the error interface.
//
// Format:
//
//   - StatusCode == 0 (transport / decode failure, no HTTP response):
//     "qubiton api error: <message>". The status segment is omitted because
//     "status 0" is meaningless — no HTTP exchange completed.
//   - StatusCode > 0: "qubiton api error: status <code>: <message>".
//   - RetryAfter > 0: appends " (retry after <N>s)" so the hint surfaces
//     in logs without callers having to type-assert the *ApiError to read
//     RetryAfter manually.
func (e *ApiError) Error() string {
	var head string
	switch {
	case e.StatusCode == 0 && e.Message != "":
		head = fmt.Sprintf("qubiton api error: %s", e.Message)
	case e.StatusCode == 0:
		head = "qubiton api error"
	case e.Message == "":
		head = fmt.Sprintf("qubiton api error: status %d", e.StatusCode)
	default:
		head = fmt.Sprintf("qubiton api error: status %d: %s", e.StatusCode, e.Message)
	}
	if e.RetryAfter > 0 {
		head += fmt.Sprintf(" (retry after %ds)", e.RetryAfter)
	}
	return head
}

// Unwrap returns the wrapped sentinel error, enabling errors.Is matching.
func (e *ApiError) Unwrap() error {
	return e.sentinel
}

// IsAuthError returns true for 401/403 responses.
func (e *ApiError) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

// IsRateLimit returns true for 429 responses.
func (e *ApiError) IsRateLimit() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServerError returns true for 5xx responses.
//
// HTTP status codes are 3-digit integers in [100, 599]; the upper bound
// (< 600) is a no-op given the wire protocol, so we test only the lower
// boundary. Documenting here explicitly rather than perpetuating the
// dead check.
func (e *ApiError) IsServerError() bool {
	return e.StatusCode >= http.StatusInternalServerError
}

// classifyStatus returns the sentinel error appropriate for an HTTP status code.
func classifyStatus(status int) error {
	switch {
	case status == http.StatusUnauthorized, status == http.StatusForbidden:
		return ErrAuth
	case status == http.StatusNotFound:
		return ErrNotFound
	case status == http.StatusTooManyRequests:
		return ErrRateLimit
	case status == http.StatusBadRequest, status == http.StatusUnprocessableEntity:
		return ErrValidation
	case status >= http.StatusInternalServerError:
		return ErrServerError
	default:
		return nil
	}
}
