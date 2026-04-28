// Package qubiton provides a Go client for the QubitOn API.
//
// 45 APIs for validating, enriching, and assessing business data across 250+ countries.
//
// Usage:
//
//	client := qubiton.NewClient("svm...")
//	resp, err := client.ValidateAddress(ctx, qubiton.AddressRequest{
//	    AddressLine1: "123 Main St", City: "New York", State: "NY",
//	    PostalCode: "10001", Country: "US",
//	})
package qubiton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.qubiton.com"
	maxRetries     = 3
	maxBackoff     = 30 * time.Second
	minBackoff     = 50 * time.Millisecond
)

// version is overridable at link time via -ldflags
// "-X github.com/qubitonhq/qubiton-go.version=0.4.1" so that release builds
// can stamp the published tag without editing source. Default is the current
// in-tree version.
var version = "0.4.0"

// Client is the QubitOn API client. Safe for concurrent use by multiple goroutines.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	oauth      *oauthTokenManager
	userAgent  string

	// OAuth bootstrap fields — captured by WithOAuth and applied in NewClient
	// after all options have run, so that WithHTTPClient/WithBaseURL ordering
	// does not matter.
	oauthClientID     string
	oauthClientSecret string
	oauthTokenURL     string
	oauthEnabled      bool
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL (default: https://api.qubiton.com).
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(url, "/") }
}

// WithTimeout sets the HTTP client timeout.
//
// Ordering with WithHTTPClient matters:
//
//   - WithHTTPClient(hc) then WithTimeout(d): WithTimeout shallow-clones the
//     caller-supplied http.Client and sets the new timeout on the clone, so
//     the caller's hc is left untouched. The Transport pointer is copied,
//     not deep-cloned — the SDK therefore shares Transport state (connection
//     pool, TLS session cache) with the caller's client. That is the
//     intended behaviour for a "give me a client and I'll use its
//     transport" contract; if you do NOT want shared Transport state, build
//     a separate http.Client with its own *http.Transport.
//
//   - WithTimeout(d) then WithHTTPClient(hc): WithHTTPClient overrides the
//     intermediate clone — the final client is the caller-supplied hc, so
//     the earlier WithTimeout call has no effect. This matches the general
//     "later option wins" rule. Use WithTimeout AFTER WithHTTPClient when
//     you need the timeout applied.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		cloned := *c.httpClient
		cloned.Timeout = d
		c.httpClient = &cloned
	}
}

// WithHTTPClient provides a custom HTTP client. The client is shared across
// all requests. The SDK does NOT clone the Transport, so any subsequent
// option that mutates client state (e.g. WithTimeout) operates on a clone of
// the http.Client struct but keeps the same *http.Transport.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithOAuth enables OAuth authentication using a key/secret pair.
//
// The QubitOn OAuth flow exchanges {key, secret} for a JWE access token via
// POST /api/oauth/token. If tokenURL is empty, it defaults to
// "{baseURL}/api/oauth/token".
//
// When OAuth is configured, the access token is sent in the Authorization
// header as a Bearer token; the apikey header is not used.
//
// The token manager is constructed after all options have been applied, so
// option ordering (e.g., WithHTTPClient before/after WithOAuth) does not
// matter — the manager always sees the final client and base URL.
func WithOAuth(clientID, clientSecret, tokenURL string) Option {
	return func(c *Client) {
		c.oauthEnabled = true
		c.oauthClientID = clientID
		c.oauthClientSecret = clientSecret
		c.oauthTokenURL = tokenURL
	}
}

// WithOAuth2 is a deprecated alias for WithOAuth retained for backwards compatibility.
//
// Deprecated: use WithOAuth. The QubitOn OAuth endpoint is not standard OAuth2.
func WithOAuth2(clientID, clientSecret, tokenURL string) Option {
	return WithOAuth(clientID, clientSecret, tokenURL)
}

// defaultTransport returns a tuned http.Transport for service-to-service traffic.
//
// Pool sizing:
//   - MaxIdleConns 100 — total cached idle connections across all hosts.
//   - MaxIdleConnsPerHost 20 — Go's stdlib default of 2 is too low for
//     service-to-service traffic; 20 lets a single host re-use connections
//     under modest concurrency without churning sockets.
//   - MaxConnsPerHost 100 — hard cap on simultaneously open connections per
//     host. Acts as a back-pressure ceiling: requests beyond 100 in-flight
//     to one host are queued at the transport rather than opening
//     unbounded sockets. Tune up if running highly concurrent batch loads
//     against a single QubitOn host; a 5x ratio over MaxIdleConnsPerHost
//     is a deliberate overhead allowance for short bursts.
//
// ResponseHeaderTimeout caps how long the transport waits for the first
// response byte (status line + headers) after fully writing the request.
// 30s matches the default Client.Timeout so a stuck server cannot pin a
// connection past the overall request budget.
func defaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}
}

// buildUserAgent returns "qubiton-go-sdk/<version> (<go-version>; <os>/<arch>)".
func buildUserAgent() string {
	return fmt.Sprintf("qubiton-go-sdk/%s (%s; %s/%s)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// NewClient creates a new QubitOn API client. apiKey may be empty when WithOAuth is used.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: defaultTransport(),
		},
		userAgent: buildUserAgent(),
	}
	for _, opt := range opts {
		opt(c)
	}
	// Apply OAuth deferred: needs final httpClient + baseURL.
	if c.oauthEnabled {
		tokenURL := c.oauthTokenURL
		if tokenURL == "" {
			tokenURL = c.baseURL + "/api/oauth/token"
		}
		c.oauth = newOAuthTokenManager(c.oauthClientID, c.oauthClientSecret, tokenURL, c.httpClient)
	}
	return c
}

// doRequest sends an HTTP request, applying authentication, retries, and error parsing.
// Returns the raw JSON response body. Callers decode into the appropriate type
// (object or array) themselves — this lets each public method decode against
// the correct shape (e.g., []FooResponse for endpoints that return arrays).
//
// Before marshalling, the SDK populates BaseRequest.RequestedByClient with the
// User-Agent if the caller has left it empty — the .NET server marks the field
// as [Required] [StringLength(350)] and rejects the request with a 400
// validation error otherwise.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyBytes = injectRequestedByClient(bodyBytes, c.userAgent)
	}

	// Local variable name `endpoint` avoids shadowing the imported net/url package.
	endpoint := c.baseURL + path
	var lastErr error
	authRetried := false

	attempt := 0
	for attempt < maxRetries {
		// Per-attempt context check; bail early if caller cancelled.
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		req, err := c.newRequest(ctx, method, endpoint, bodyBytes)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			// Network errors are retryable.
			if attempt < maxRetries-1 {
				sleep(ctx, c.backoff(attempt, 0))
				attempt++
				continue
			}
			break
		}

		respBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		// Auto-retry once on 401 when using OAuth: token may have expired
		// out-of-band. Invalidate the cached token and retry with a fresh one.
		// Does NOT consume an attempt slot — the client should not be punished
		// for an out-of-band token expiration.
		if resp.StatusCode == http.StatusUnauthorized && c.oauth != nil && !authRetried {
			c.oauth.invalidate()
			authRetried = true
			continue
		}

		// Rate-limited: respect Retry-After if present, otherwise backoff.
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
			if attempt < maxRetries-1 {
				delay := c.backoff(attempt, retryAfter)
				sleep(ctx, delay)
				attempt++
				continue
			}
			return nil, buildApiError(resp.StatusCode, respBody, retryAfter)
		}

		// 5xx: retryable for any HTTP method (server-side issue, no client side-effect).
		// Exception: 501 Not Implemented is terminal — retrying does not help.
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
			apiErr := buildApiError(resp.StatusCode, respBody, retryAfter)
			if resp.StatusCode == http.StatusNotImplemented {
				return nil, apiErr
			}
			lastErr = apiErr
			if attempt < maxRetries-1 {
				sleep(ctx, c.backoff(attempt, retryAfter))
				attempt++
				continue
			}
			return nil, lastErr
		}

		// 408 Request Timeout: retryable.
		if resp.StatusCode == http.StatusRequestTimeout {
			lastErr = buildApiError(resp.StatusCode, respBody, 0)
			if attempt < maxRetries-1 {
				sleep(ctx, c.backoff(attempt, 0))
				attempt++
				continue
			}
			return nil, lastErr
		}

		// All other 4xx are terminal — do not retry.
		if resp.StatusCode >= 400 {
			return nil, buildApiError(resp.StatusCode, respBody, 0)
		}

		// 2xx: return raw bytes.
		return respBody, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &ApiError{StatusCode: 0, Message: "request failed after retries"}
}

// decodeJSON unmarshals body into out, treating empty body as a no-op (out is
// left at its zero value).
//
// Slice/array contract: when out points to a slice (e.g. *[]Foo) and the
// server returns an empty body, out is left as a nil slice. Callers must
// treat nil and empty slices identically (range works on both, len returns
// 0). decodeJSON does NOT replace a nil slice with an empty one because
// doing so would silently mask a server contract regression where an
// endpoint that should always return an array starts returning empty
// bodies. Use len(slice) == 0 rather than slice == nil for portability.
func decodeJSON(body []byte, out interface{}) error {
	if len(body) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// newRequest builds an authenticated HTTP request. `endpoint` is the fully
// qualified URL (named to avoid shadowing the net/url package).
func (c *Client) newRequest(ctx context.Context, method, endpoint string, bodyBytes []byte) (*http.Request, error) {
	var reqBody io.Reader
	if bodyBytes != nil {
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Idempotency key for safe retries on POST/PUT mutations.
	if v := ctx.Value(idempotencyKeyCtxKey{}); v != nil {
		if k, ok := v.(string); ok && k != "" {
			req.Header.Set("X-Idempotency-Key", k)
		}
	}

	// Authentication: prefer OAuth Bearer when configured, otherwise apikey header.
	// Note: the .NET server reads the lowercase "apikey" header. "X-Api-Key" is ignored.
	if c.oauth != nil {
		token, err := c.oauth.getToken(ctx)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	} else if c.apiKey != "" {
		req.Header.Set("apikey", c.apiKey)
	}

	return req, nil
}

// idempotencyKeyCtxKey is the context-key type used to attach an idempotency
// key to a single request. Use WithIdempotencyKey to opt in.
type idempotencyKeyCtxKey struct{}

// WithIdempotencyKey returns a derived context that, when passed to a Client
// method, adds an X-Idempotency-Key header to the outgoing request. Used to
// make retries of POST/PUT mutations safe (the server stores the result under
// the key for 24 h and returns the same response on duplicate submissions).
//
// Generate a fresh, cryptographically random key per logical operation.
func WithIdempotencyKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKeyCtxKey{}, key)
}

// injectRequestedByClient ensures the marshalled JSON body has a
// "requestedByClient" field. The .NET server's BaseRequest marks the field as
// [Required] [StringLength(350)], so without it every authenticated request
// returns a 400 validation error.
//
// Strategy: detect whether the field is already present and non-empty using a
// json.Decoder with UseNumber() (so we don't lose precision on large int64
// values during the inspection phase), and if absent, inject the user-agent
// at the byte level by splicing a new key-value pair before the closing brace
// of the top-level object. We never fully re-marshal the body, so int64
// values >2^53, original key ordering, and any custom marshallers in the
// caller's struct are all preserved on the wire.
//
// Duplicate-key tolerance: if the body already contains an empty
// requestedByClient (key present but value is the empty string), this
// function appends a fresh key-value pair rather than rewriting in place,
// which produces a JSON object with two "requestedByClient" entries. This
// is safe because System.Text.Json (the .NET server's deserializer) and
// Go's encoding/json both apply "last value wins" semantics on duplicate
// keys, so the injected user-agent is the value seen by the server. If a
// future deserializer tightens this and rejects duplicate keys, the
// inspection-then-rewrite path here will need to switch to a structured
// rewrite rather than a splice.
//
// Arrays (e.g. the LookupExchangeRates dates body) are left alone — they
// don't carry a BaseRequest. Non-object JSON is returned unchanged.
//
// The user-agent is truncated to 350 characters to satisfy [StringLength(350)].
func injectRequestedByClient(bodyBytes []byte, userAgent string) []byte {
	if len(bodyBytes) == 0 || userAgent == "" {
		return bodyBytes
	}
	// Cheap check: only objects need the field. Arrays / scalars return as-is.
	trim := bytes.TrimLeft(bodyBytes, " \t\r\n")
	if len(trim) == 0 || trim[0] != '{' {
		return bodyBytes
	}

	// Inspect (don't decode-and-re-encode) to see whether a non-empty
	// requestedByClient is already present. UseNumber preserves int64
	// precision in case anything else matters during inspection.
	dec := json.NewDecoder(bytes.NewReader(bodyBytes))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil || tok != json.Delim('{') {
		return bodyBytes
	}
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return bodyBytes
		}
		key, ok := keyTok.(string)
		if !ok {
			return bodyBytes
		}
		valTok, err := dec.Token()
		if err != nil {
			return bodyBytes
		}
		// If the value is the start of a nested object/array, skip past it.
		if d, ok := valTok.(json.Delim); ok && (d == '{' || d == '[') {
			depth := 1
			for depth > 0 {
				t, err := dec.Token()
				if err != nil {
					return bodyBytes
				}
				if dd, ok := t.(json.Delim); ok {
					switch dd {
					case '{', '[':
						depth++
					case '}', ']':
						depth--
					}
				}
			}
			continue
		}
		if key == "requestedByClient" {
			if s, ok := valTok.(string); ok && s != "" {
				return bodyBytes
			}
			// Field present but empty — fall through to splicing in a fresh value.
			// We still rewrite the full body via byte manipulation below;
			// duplicate keys are tolerated by encoding/json (last value wins
			// on the .NET side and on Go), so we can simply append.
			break
		}
	}

	ua := userAgent
	if len(ua) > 350 {
		ua = ua[:350]
	}

	// Splice "requestedByClient":"<ua>" before the closing brace of the
	// top-level object. We find the rightmost '}' that is at brace-depth 1
	// in the original byte stream, ignoring braces that appear inside
	// string literals.
	insertAt := findTopLevelClosingBraceOptimistic(bodyBytes)
	if insertAt < 0 {
		return bodyBytes
	}
	uaJSON, err := json.Marshal(ua)
	if err != nil {
		return bodyBytes
	}
	// Detect whether the object already has any keys (so we know whether to
	// prepend a comma). Walk left from insertAt skipping whitespace.
	prev := insertAt - 1
	for prev >= 0 {
		ch := bodyBytes[prev]
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			prev--
			continue
		}
		break
	}
	var prefix []byte
	if prev >= 0 && bodyBytes[prev] != '{' {
		prefix = []byte{','}
	}
	out := make([]byte, 0, len(bodyBytes)+len(uaJSON)+24)
	out = append(out, bodyBytes[:insertAt]...)
	out = append(out, prefix...)
	out = append(out, []byte(`"requestedByClient":`)...)
	out = append(out, uaJSON...)
	out = append(out, bodyBytes[insertAt:]...)
	return out
}

// findTopLevelClosingBraceOptimistic returns the index of the closing '}' of
// the top-level JSON object. Returns -1 if the input is malformed or not an
// object. Brace and quote tracking is just enough to skip string literals
// (with backslash escapes) — we DO NOT fully validate the JSON. Callers must
// have already validated structure (e.g. via the upstream json.Marshal that
// produced the body); this helper is intentionally optimistic so the hot
// path stays a single byte scan.
func findTopLevelClosingBraceOptimistic(b []byte) int {
	depth := 0
	inStr := false
	for i := 0; i < len(b); i++ {
		ch := b[i]
		if inStr {
			if ch == '\\' {
				i++ // skip escaped char
				continue
			}
			if ch == '"' {
				inStr = false
			}
			continue
		}
		switch ch {
		case '"':
			inStr = true
		case '{', '[':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		case ']':
			depth--
		}
	}
	return -1
}

// buildApiError parses a non-2xx response body and returns a typed *ApiError.
//
// The QubitOn API returns errors in one of these shapes:
//
//	{"error": "...", "statusCode": N}                  // most endpoints (BaseErrorResponse)
//	{"errorResponse": {"statusCode": N, "message": ...}} // legacy/internal cache responses
//	{"title": "...", "detail": "...", "errors": {...}}   // ASP.NET ProblemDetails / validation
//
// We try the modern "error" field first, fall back to "message", then
// "errorResponse.message", then ProblemDetails fields, and finally the raw body.
// Validation "errors" maps are folded into the message ("Field: msg, Other: msg").
func buildApiError(status int, body []byte, retryAfter int) *ApiError {
	apiErr := &ApiError{
		StatusCode: status,
		RetryAfter: retryAfter,
		sentinel:   classifyStatus(status),
	}

	var raw map[string]interface{}
	if len(body) > 0 && json.Unmarshal(body, &raw) == nil {
		apiErr.Raw = raw
		apiErr.Message = extractErrorMessage(raw)
	}
	if apiErr.Message == "" && len(body) > 0 {
		// Cap raw body for human-readable fallback (avoid HTML pages flooding logs).
		txt := string(body)
		if len(txt) > 256 {
			txt = txt[:256]
		}
		apiErr.Message = strings.TrimSpace(txt)
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(status)
	}
	return apiErr
}

// extractErrorMessage looks at the common error envelope shapes used by the
// QubitOn API and returns a human-readable message, or "" if none can be found.
func extractErrorMessage(raw map[string]interface{}) string {
	// Modern envelope: {"error": "...", "statusCode": N}
	if v, ok := raw["error"].(string); ok && v != "" {
		return v
	}
	// Some endpoints use {"message": "..."} directly.
	if v, ok := raw["message"].(string); ok && v != "" {
		return v
	}
	// Legacy/internal cache: {"errorResponse": {"message": "...", ...}}
	if er, ok := raw["errorResponse"].(map[string]interface{}); ok {
		if v, ok := er["message"].(string); ok && v != "" {
			return v
		}
		if v, ok := er["status"].(string); ok && v != "" {
			return v
		}
	}
	// ASP.NET ProblemDetails / validation errors.
	// {"title": "...", "detail": "...", "errors": {"Field": ["msg1", "msg2"]}}
	base := ""
	if v, ok := raw["detail"].(string); ok && v != "" {
		base = v
	} else if v, ok := raw["title"].(string); ok && v != "" {
		base = v
	}
	if errs, ok := raw["errors"].(map[string]interface{}); ok && len(errs) > 0 {
		// Sort fields for deterministic output — Go's map iteration order is
		// randomised so consumers logging or asserting on the message would
		// otherwise see flaky text.
		fields := make([]string, 0, len(errs))
		for field := range errs {
			fields = append(fields, field)
		}
		sort.Strings(fields)
		parts := make([]string, 0, len(errs))
		for _, field := range fields {
			switch v := errs[field].(type) {
			case string:
				parts = append(parts, fmt.Sprintf("%s: %s", field, v))
			case []interface{}:
				msgs := make([]string, 0, len(v))
				for _, item := range v {
					if s, ok := item.(string); ok {
						msgs = append(msgs, s)
					}
				}
				if len(msgs) > 0 {
					parts = append(parts, fmt.Sprintf("%s: %s", field, strings.Join(msgs, "; ")))
				}
			}
		}
		if len(parts) > 0 {
			joined := strings.Join(parts, "; ")
			if base == "" {
				return joined
			}
			return base + " (" + joined + ")"
		}
	}
	return base
}

// parseRetryAfter parses an HTTP Retry-After header value. Supports
// delta-seconds (integer or fractional, per RFC 9110 errata for some
// servers) and HTTP-date formats. Returns 0 on parse failure.
//
// Fractional values are rounded up (ceiling) so we never under-wait — a
// 0.4s hint becomes 1s of delay rather than 0s.
func parseRetryAfter(h string) int {
	h = strings.TrimSpace(h)
	if h == "" {
		return 0
	}
	if n, err := strconv.Atoi(h); err == nil {
		if n < 0 {
			return 0
		}
		return n
	}
	if f, err := strconv.ParseFloat(h, 64); err == nil {
		if f <= 0 {
			return 0
		}
		// Round up so a fractional hint is never under-honoured.
		n := int(f)
		if float64(n) < f {
			n++
		}
		return n
	}
	if t, err := http.ParseTime(h); err == nil {
		d := time.Until(t)
		if d <= 0 {
			return 0
		}
		return int(d.Seconds())
	}
	return 0
}

// backoff returns the delay before the next retry attempt.
//
// Uses exponential backoff (1s * 2^attempt, computed via integer shift) with
// ±25% jitter to avoid thundering-herd. If the server provided a Retry-After
// hint, that is used as the floor. Capped at maxBackoff (30s); floored at
// minBackoff (50ms).
//
// Uses the process-wide concurrency-safe rand.Float64() from math/rand/v2.
func (c *Client) backoff(attempt int, retryAfterSeconds int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 30 {
		attempt = 30 // saturate; (1<<31) seconds is meaningless and overflows int32 platforms.
	}
	base := time.Second * time.Duration(1<<attempt)
	if retryAfterSeconds > 0 {
		hinted := time.Duration(retryAfterSeconds) * time.Second
		if hinted > base {
			base = hinted
		}
	}
	// ±25% jitter.
	jitter := (rand.Float64() - 0.5) * 0.5
	delay := time.Duration(float64(base) * (1 + jitter))
	if delay > maxBackoff {
		delay = maxBackoff
	}
	if delay < minBackoff {
		delay = minBackoff
	}
	return delay
}

// doGet performs a GET request with no body.
func (c *Client) doGet(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// ── Address Validation ────────────────────────────────────────────────────

// ValidateAddress validates and standardizes a postal address across 249 countries.
func (c *Client) ValidateAddress(ctx context.Context, req AddressRequest) (*AddressResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/address/validate", req)
	if err != nil {
		return nil, err
	}
	var resp AddressResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Tax Validation ────────────────────────────────────────────────────────

// ValidateTax validates a tax identification number across 60+ countries with live authority checks.
func (c *Client) ValidateTax(ctx context.Context, req TaxRequest) (*TaxResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/tax/validate", req)
	if err != nil {
		return nil, err
	}
	var resp TaxResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ValidateTaxFormat validates tax ID format using regex and checksum algorithms
// for 193 countries and 242 tax types.
func (c *Client) ValidateTaxFormat(ctx context.Context, req TaxFormatRequest) (*TaxFormatResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/tax/format-validate", req)
	if err != nil {
		return nil, err
	}
	var resp TaxFormatResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Bank Account Validation ───────────────────────────────────────────────

// ValidateBankAccount validates bank accounts across 180+ countries (IBAN, SWIFT, routing numbers).
//
// The /api/bank/validate route is served by BankValidationController and
// returns a BankResponse — i.e. BankAccountNumberResponse extended with a
// BankProValidations analytics block. The SDK's BankAccountResponse already
// carries BankProValidations and AdditionalInfo / IframeDetails inherited
// from the .NET hierarchy, so this method covers both the Pro analytics
// extension and the non-Pro fields with one type.
func (c *Client) ValidateBankAccount(ctx context.Context, req BankAccountRequest) (*BankAccountResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/bank/validate", req)
	if err != nil {
		return nil, err
	}
	var resp BankAccountResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ValidateBankPro performs premium bank analytics with ownership verification and confidence scoring.
func (c *Client) ValidateBankPro(ctx context.Context, req BankProRequest) (*BankProResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/bankaccount/pro/validate", req)
	if err != nil {
		return nil, err
	}
	var resp BankProResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Email & Phone Validation ──────────────────────────────────────────────

// ValidateEmail validates an email address for deliverability and risk.
func (c *Client) ValidateEmail(ctx context.Context, req EmailRequest) (*EmailResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/email/validate", req)
	if err != nil {
		return nil, err
	}
	var resp EmailResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ValidatePhone validates a phone number against carrier databases.
func (c *Client) ValidatePhone(ctx context.Context, req PhoneRequest) (*PhoneResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/phone/validate", req)
	if err != nil {
		return nil, err
	}
	var resp PhoneResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Business Registration ─────────────────────────────────────────────────

// LookupBusinessRegistration looks up official business registration records.
// The server returns an array; results are ordered most-relevant first.
func (c *Client) LookupBusinessRegistration(ctx context.Context, req BusinessRegistrationRequest) ([]BusinessRegistrationResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/businessregistration/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []BusinessRegistrationResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ── Peppol ────────────────────────────────────────────────────────────────

// ValidatePeppol validates a Peppol participant ID against 70+ ISO 6523 ICD schemes.
func (c *Client) ValidatePeppol(ctx context.Context, req PeppolRequest) (*PeppolResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/peppol/validate", req)
	if err != nil {
		return nil, err
	}
	var resp PeppolResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Sanctions & Compliance ────────────────────────────────────────────────

// CheckSanctions screens an entity against 100+ global sanctions lists (OFAC, EU, UN, UK HMT).
// The server returns an array (one entry per matching list/source).
func (c *Client) CheckSanctions(ctx context.Context, req SanctionsRequest) ([]SanctionsResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/prohibited/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []SanctionsResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ScreenPEP screens against Politically Exposed Person databases.
// The server returns an array (one entry per matching source/dataset).
func (c *Client) ScreenPEP(ctx context.Context, req PEPRequest) ([]PEPResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/pep/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []PEPResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CheckDirectors checks for disqualified directors.
func (c *Client) CheckDirectors(ctx context.Context, req DirectorsRequest) (*DirectorsResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/disqualifieddirectors/validate", req)
	if err != nil {
		return nil, err
	}
	var resp DirectorsResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── EPA Prosecution ───────────────────────────────────────────────────────

// CheckEPAProsecution screens against EPA criminal prosecution records.
// The server returns an array.
func (c *Client) CheckEPAProsecution(ctx context.Context, req EPARequest) ([]EPAResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/criminalprosecution/validate", req)
	if err != nil {
		return nil, err
	}
	var resp []EPAResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// LookupEPAProsecution looks up EPA criminal prosecution details.
// The server returns an array.
func (c *Client) LookupEPAProsecution(ctx context.Context, req EPARequest) ([]EPAResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/criminalprosecution/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []EPAResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ── Healthcare Exclusion ──────────────────────────────────────────────────

// CheckHealthcareExclusion screens against healthcare provider exclusion lists.
// The server returns an array (one entry per matching source).
func (c *Client) CheckHealthcareExclusion(ctx context.Context, req HealthcareExclusionRequest) ([]HealthcareExclusionResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/providerexclusion/validate", req)
	if err != nil {
		return nil, err
	}
	var resp []HealthcareExclusionResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// LookupHealthcareExclusion looks up healthcare provider exclusion details.
// The server returns a single aggregated lookup response.
func (c *Client) LookupHealthcareExclusion(ctx context.Context, req HealthcareExclusionRequest) (*HealthcareExclusionLookupResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/providerexclusion/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp HealthcareExclusionLookupResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Risk & Financial ──────────────────────────────────────────────────────
//
// Note on /api/risk/lookup categories:
// The server's RiskController.Lookup endpoint only accepts categories
// "Social", "Governance", and "Environmental" — those return ESG-style
// adverse media analytics. The categories "Bankruptcy", "Credit Score",
// and "Fail Rate" are routed through /api/risk/riskcontrol
// (RiskControlController) which returns a different shape (RiskControlResponse
// with case lists and a recommendation). The methods below send those
// categories to the correct endpoint.

// CheckBankruptcyRisk checks if a company has filed for or is in bankruptcy proceedings.
func (c *Client) CheckBankruptcyRisk(ctx context.Context, req BankruptcyRequest) (*BankruptcyResponse, error) {
	apiReq := riskLookupRequest{EntityName: req.CompanyName, Country: req.Country, Category: "Bankruptcy"}
	body, err := c.doRequest(ctx, http.MethodPost, "/api/risk/riskcontrol", apiReq)
	if err != nil {
		return nil, err
	}
	var resp BankruptcyResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupCreditScore looks up commercial credit score and financial stability.
func (c *Client) LookupCreditScore(ctx context.Context, req CreditScoreRequest) (*CreditScoreResponse, error) {
	apiReq := riskLookupRequest{EntityName: req.CompanyName, Country: req.Country, Category: "Credit Score"}
	body, err := c.doRequest(ctx, http.MethodPost, "/api/risk/riskcontrol", apiReq)
	if err != nil {
		return nil, err
	}
	var resp CreditScoreResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupFailRate looks up payment failure rate and risk classification.
func (c *Client) LookupFailRate(ctx context.Context, req FailRateRequest) (*FailRateResponse, error) {
	apiReq := riskLookupRequest{EntityName: req.CompanyName, Country: req.Country, Category: "Fail Rate"}
	body, err := c.doRequest(ctx, http.MethodPost, "/api/risk/riskcontrol", apiReq)
	if err != nil {
		return nil, err
	}
	var resp FailRateResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupRisk performs a Risk lookup against /api/risk/lookup.
// The server only accepts categories "Social", "Governance", or "Environmental"
// (adverse media analytics) and returns an array of category-scoped results.
func (c *Client) LookupRisk(ctx context.Context, req RiskLookupRequest) ([]RiskResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/risk/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []RiskResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AssessEntityRisk assesses entity fraud risk and adverse media.
func (c *Client) AssessEntityRisk(ctx context.Context, req EntityRiskRequest) (*EntityRiskResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/entity/fraud/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp EntityRiskResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupCreditAnalysis performs comprehensive credit analysis on a business entity.
// The server returns an array (one entry per provider/source).
func (c *Client) LookupCreditAnalysis(ctx context.Context, req CreditAnalysisRequest) ([]CreditAnalysisResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/creditanalysis/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []CreditAnalysisResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ── ESG & Cybersecurity ───────────────────────────────────────────────────

// LookupESGScore looks up ESG (Environmental, Social, Governance) scores.
// The server returns an array (one entry per scoring provider).
func (c *Client) LookupESGScore(ctx context.Context, req ESGRequest) ([]ESGResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/esg/Scores", req)
	if err != nil {
		return nil, err
	}
	var resp []ESGResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DomainSecurityReport assesses domain cybersecurity and threat intelligence.
func (c *Client) DomainSecurityReport(ctx context.Context, req DomainSecurityRequest) (*DomainSecurityResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/itsecurity/domainreport", req)
	if err != nil {
		return nil, err
	}
	var resp DomainSecurityResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CheckIPQuality checks IP address quality and fraud risk.
func (c *Client) CheckIPQuality(ctx context.Context, req IPQualityRequest) (*IPQualityResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/ipquality/validate", req)
	if err != nil {
		return nil, err
	}
	var resp IPQualityResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Corporate Structure ───────────────────────────────────────────────────

// LookupBeneficialOwnership looks up beneficial ownership for corporate transparency.
// The server returns an array (one entry per provider / hierarchy).
func (c *Client) LookupBeneficialOwnership(ctx context.Context, req BeneficialOwnershipRequest) ([]BeneficialOwnershipResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/beneficialownership/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []BeneficialOwnershipResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// LookupCorporateHierarchy looks up corporate hierarchy and ownership structure (US only).
// The server returns an array.
func (c *Client) LookupCorporateHierarchy(ctx context.Context, req CorporateHierarchyRequest) ([]CorporateHierarchyResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/corporatehierarchy/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []CorporateHierarchyResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// LookupDUNS looks up a DUNS number for company identification.
// The server returns an array (one entry per matching record).
func (c *Client) LookupDUNS(ctx context.Context, req DUNSRequest) ([]DUNSResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/duns-number-lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []DUNSResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// LookupHierarchy looks up company parent-child hierarchy.
func (c *Client) LookupHierarchy(ctx context.Context, req HierarchyRequest) (*HierarchyResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/company/hierarchy/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp HierarchyResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Industry Specific ─────────────────────────────────────────────────────

// ValidateNPI validates a US National Provider Identifier.
func (c *Client) ValidateNPI(ctx context.Context, req NPIRequest) (*NPIResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/nationalprovideridentifier/validate", req)
	if err != nil {
		return nil, err
	}
	var resp NPIResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ValidateMedpass validates a healthcare supplier via Medpass.
func (c *Client) ValidateMedpass(ctx context.Context, req MedpassRequest) (*MedpassResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/medpass/validate", req)
	if err != nil {
		return nil, err
	}
	var resp MedpassResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupDOTCarrier looks up USDOT/FMCSA motor carrier safety data.
// The server returns an array (one entry per matching carrier).
func (c *Client) LookupDOTCarrier(ctx context.Context, req DOTCarrierRequest) ([]DOTCarrierResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/dot/fmcsa/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []DOTCarrierResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ValidateIndiaIdentity validates Indian identity documents (Driver License, Voter Registration).
func (c *Client) ValidateIndiaIdentity(ctx context.Context, req IndiaIdentityRequest) (*IndiaIdentityResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/inidentity/validate", req)
	if err != nil {
		return nil, err
	}
	var resp IndiaIdentityResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Certification ─────────────────────────────────────────────────────────

// ValidateCertification validates a business certification (MBE, WBE, DBE, etc.).
func (c *Client) ValidateCertification(ctx context.Context, req CertificationRequest) (*CertificationResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/certification/validate", req)
	if err != nil {
		return nil, err
	}
	var resp CertificationResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupCertification looks up business certifications (diversity, small business).
// The server returns an array (one entry per matching certification).
func (c *Client) LookupCertification(ctx context.Context, req CertificationRequest) ([]CertificationResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/certification/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []CertificationResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ── Business Classification ───────────────────────────────────────────────

// LookupBusinessClassification looks up NAICS/SIC business classification codes.
// The server returns an array (one entry per code/source).
func (c *Client) LookupBusinessClassification(ctx context.Context, req BusinessClassificationRequest) ([]BusinessClassificationResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/businessclassification/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []BusinessClassificationResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ── Financial Operations ──────────────────────────────────────────────────

// AnalyzePaymentTerms analyzes payment terms for optimization and early-pay discounts.
func (c *Client) AnalyzePaymentTerms(ctx context.Context, req PaymentTermsRequest) (*PaymentTermsResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/paymentterms/validate", req)
	if err != nil {
		return nil, err
	}
	var resp PaymentTermsResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// LookupExchangeRates looks up currency exchange rates for specific dates.
// The API expects baseCurrency as a URL-escaped path parameter and dates as a
// JSON array body (RFC3339 timestamps). The server returns an array of per-day
// rates.
//
// Zero-valued time.Time entries in Dates are skipped; an empty/nil Dates slice
// produces an empty JSON array on the wire (rather than `null`).
//
// BaseCurrency must be non-empty after trimming whitespace. An empty
// BaseCurrency would otherwise produce a malformed URL with a trailing
// slash (and no path parameter), which the server rejects as 404 — fail
// fast on the client.
func (c *Client) LookupExchangeRates(ctx context.Context, req ExchangeRateRequest) ([]ExchangeRateResponse, error) {
	if strings.TrimSpace(req.BaseCurrency) == "" {
		return nil, &ApiError{StatusCode: 0, Message: "LookupExchangeRates: BaseCurrency is required", sentinel: ErrValidation}
	}
	path := "/api/currency/exchange-rates/" + url.PathEscape(req.BaseCurrency)
	// Marshal each date as RFC3339; ensure non-nil even when empty so the
	// server gets [] rather than null.
	dates := make([]string, 0, len(req.Dates))
	for _, t := range req.Dates {
		if t.IsZero() {
			continue
		}
		dates = append(dates, t.UTC().Format(time.RFC3339))
	}
	body, err := c.doRequest(ctx, http.MethodPost, path, dates)
	if err != nil {
		return nil, err
	}
	var resp []ExchangeRateResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ── Supplier Profile (SAP Ariba) ──────────────────────────────────────────

// LookupAribaSupplier looks up a SAP Ariba supplier profile by ANID.
// The server returns an array.
func (c *Client) LookupAribaSupplier(ctx context.Context, req AribaSupplierRequest) ([]AribaSupplierResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/aribasupplierprofile/lookup", req)
	if err != nil {
		return nil, err
	}
	var resp []AribaSupplierResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ValidateAribaSupplier validates a SAP Ariba supplier profile by ANID.
func (c *Client) ValidateAribaSupplier(ctx context.Context, req AribaSupplierRequest) (*AribaSupplierResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/aribasupplierprofile/validate", req)
	if err != nil {
		return nil, err
	}
	var resp AribaSupplierResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Gender Identification ─────────────────────────────────────────────────

// IdentifyGender predicts gender from a person's name.
func (c *Client) IdentifyGender(ctx context.Context, req GenderRequest) (*GenderResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/genderize/identifygender", req)
	if err != nil {
		return nil, err
	}
	var resp GenderResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Continuous Screening / Monitoring ─────────────────────────────────────

// ScreenContinuous submits an entity for continuous screening (re-screen on schedule).
//
// Deprecated: the /api/continuous-screening/screen endpoint currently returns
// HTTP 501 Not Implemented on the server. This method is preserved for forward
// compatibility but should not be relied on in production until the server
// implementation ships.
func (c *Client) ScreenContinuous(ctx context.Context, req ContinuousScreeningRequest) (*ContinuousScreeningResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/continuous-screening/screen", req)
	if err != nil {
		return nil, err
	}
	var resp ContinuousScreeningResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Check Status (bulk callback) ──────────────────────────────────────────

// CheckCallbackStatus retrieves the status of a previously submitted bulk
// validation job, identified by its CallBackID. Mirrors the canonical
// CallBackRequest → CheckStatusResponse pair from
// smartvm.BusinessEntities.Client.API.{Status,API.CheckStatus}.
func (c *Client) CheckCallbackStatus(ctx context.Context, req CallBackRequest) (*CheckStatusResponse, error) {
	body, err := c.doRequest(ctx, http.MethodPost, "/api/bulkstatus/check", req)
	if err != nil {
		return nil, err
	}
	var resp CheckStatusResponse
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ── Reference Endpoints ───────────────────────────────────────────────────

// GetSupportedTaxFormats lists all supported country + tax type combinations.
// Returns 193 countries and 242 tax types with format descriptions.
// The server returns a JSON array.
func (c *Client) GetSupportedTaxFormats(ctx context.Context) ([]SupportedTaxFormat, error) {
	body, err := c.doGet(ctx, "/api/tax/format-validate/countries")
	if err != nil {
		return nil, err
	}
	var resp []SupportedTaxFormat
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPeppolSchemes lists all supported Peppol ICD schemes.
// The server returns a JSON array.
func (c *Client) GetPeppolSchemes(ctx context.Context) ([]SupportedPeppolScheme, error) {
	body, err := c.doGet(ctx, "/api/peppol/schemes")
	if err != nil {
		return nil, err
	}
	var resp []SupportedPeppolScheme
	if err := decodeJSON(body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────

// riskLookupRequest is the internal request struct matching the API's unified
// risk endpoint. Embeds BaseRequest for structural consistency with all other
// concrete *Request types — the SDK's User-Agent injector populates
// requestedByClient automatically when the caller leaves it empty.
type riskLookupRequest struct {
	BaseRequest
	EntityName string `json:"entityName"`
	Category   string `json:"category"`
	Country    string `json:"country,omitempty"`
}

// sleep waits for d or until ctx is cancelled, whichever comes first.
func sleep(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

// Compile-time assertion that *ApiError satisfies the error interface.
var _ error = (*ApiError)(nil)
