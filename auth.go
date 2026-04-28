package qubiton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
	"sync"
	"time"
)

// oauthFetchTimeout caps how long a background OAuth token fetch can run.
// Set to twice maxBackoff so the maxRetries backoff schedule fits with
// headroom for the actual request RTT.
const oauthFetchTimeout = 2 * maxBackoff

// oauthBackoff mirrors Client.backoff but lives here to avoid pulling the data
// plane's *Client receiver into the auth manager. Same shape: exponential
// (1s * 2^attempt) with ±25% jitter, clamped to [minBackoff, maxBackoff].
// retryAfterSeconds is used as a floor when the server provided a hint.
func oauthBackoff(attempt, retryAfterSeconds int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 30 {
		attempt = 30
	}
	base := time.Second * time.Duration(1<<attempt)
	if retryAfterSeconds > 0 {
		hinted := time.Duration(retryAfterSeconds) * time.Second
		if hinted > base {
			base = hinted
		}
	}
	jitter := (rand.Float64() - 0.5) * 0.5
	d := time.Duration(float64(base) * (1 + jitter))
	if d > maxBackoff {
		d = maxBackoff
	}
	if d < minBackoff {
		d = minBackoff
	}
	return d
}

// oauthSleep waits for d or until ctx is cancelled, whichever comes first.
func oauthSleep(ctx context.Context, d time.Duration) {
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

// oauthTokenManager caches access tokens issued by POST /api/oauth/token.
//
// The QubitOn OAuth endpoint is NOT a standard OAuth2 endpoint: it accepts a
// JSON body of the form {"key": "...", "secret": "..."} and returns a JWE
// access token along with an integer "expiresInSeconds" lifetime.
//
// Concurrent callers requesting a fresh token are deduped through a small
// in-process singleflight implementation: only one HTTP fetch is in flight
// at any time per token URL, and all waiters receive the same token/error.
// Cache hits never block on the network.
type oauthTokenManager struct {
	clientID     string
	clientSecret string
	tokenURL     string
	httpClient   *http.Client

	mu    sync.Mutex
	token string
	expAt time.Time

	// inFlight is non-nil while a fetch is in progress; concurrent callers
	// wait on it instead of issuing duplicate HTTP requests.
	inFlight *fetchCall
}

// fetchCall tracks a single in-flight token fetch.
type fetchCall struct {
	done  chan struct{}
	token string
	err   error
}

// newOAuthTokenManager constructs a token manager for the given credentials.
func newOAuthTokenManager(clientID, clientSecret, tokenURL string, httpClient *http.Client) *oauthTokenManager {
	return &oauthTokenManager{
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     tokenURL,
		httpClient:   httpClient,
	}
}

// getToken returns a cached access token if still valid, otherwise fetches a fresh one.
// Concurrent callers do not serialize on the network — only the cache check is
// guarded; the fetch itself is shared via singleflight.
//
// Leader cancellation isolation: the goroutine that performs the network
// round-trip runs against a context derived from context.Background() with
// the same write/read deadlines as the data plane (oauthFetchTimeout). This
// decouples the shared fetch from any single caller's cancellation — if the
// leader cancels its own ctx, the in-flight fetch continues so followers
// still receive a usable token. The leader still observes its own ctx via
// the early select below, returning ctx.Err() to the caller while the
// background fetch keeps running.
func (m *oauthTokenManager) getToken(ctx context.Context) (string, error) {
	// Fast path: cache hit under brief lock.
	m.mu.Lock()
	if m.token != "" && time.Now().UTC().Before(m.expAt.Add(-30*time.Second)) {
		t := m.token
		m.mu.Unlock()
		return t, nil
	}

	// No usable cached token. If a fetch is already in flight, wait for it.
	call := m.inFlight
	if call == nil {
		call = &fetchCall{done: make(chan struct{})}
		m.inFlight = call
		m.mu.Unlock()

		// Run the fetch on a background context so the leader's ctx
		// cancellation does not abort it for the followers. The fetch has
		// its own bounded deadline (oauthFetchTimeout) to prevent runaway
		// goroutine leaks if the auth server is unreachable.
		go func() {
			fetchCtx, cancel := context.WithTimeout(context.Background(), oauthFetchTimeout)
			defer cancel()
			token, expAt, err := m.fetchToken(fetchCtx)
			call.token, call.err = token, err
			close(call.done)

			// Publish token AND expiry atomically under the same lock so
			// that followers (or a concurrent fast-path read) never observe
			// a token without its expiry, or vice versa.
			m.mu.Lock()
			if err == nil {
				m.token = token
				m.expAt = expAt
			}
			m.inFlight = nil
			m.mu.Unlock()
		}()

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-call.done:
			return call.token, call.err
		}
	}
	m.mu.Unlock()

	// A fetch is already happening; wait for its outcome.
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-call.done:
		return call.token, call.err
	}
}

// invalidate clears any cached token, forcing the next getToken call to refetch.
func (m *oauthTokenManager) invalidate() {
	m.mu.Lock()
	m.token = ""
	m.expAt = time.Time{}
	m.mu.Unlock()
}

// fetchToken issues a fresh OAuth token request. Returns the (token, expAt)
// pair on success so getToken can publish them atomically under a single
// lock acquisition. Caller must NOT hold m.mu while invoking fetchToken —
// the network round-trip would otherwise block the cache.
//
// The fetch retries up to maxRetries times on transient failures (5xx
// responses or network errors) with the same exponential-backoff schedule
// the data plane uses, so a momentary auth-server hiccup does not bubble up
// as a hard failure to the first caller after token expiry. 4xx responses
// (invalid credentials, malformed body) are terminal — retrying does not
// help.
func (m *oauthTokenManager) fetchToken(ctx context.Context) (string, time.Time, error) {
	body, err := json.Marshal(map[string]string{
		"key":    m.clientID,
		"secret": m.clientSecret,
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("marshal oauth request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return "", time.Time{}, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.tokenURL, bytes.NewReader(body))
		if err != nil {
			return "", time.Time{}, fmt.Errorf("create oauth request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", buildUserAgent())

		resp, err := m.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("oauth token request: %w", err)
			if attempt < maxRetries-1 {
				oauthSleep(ctx, oauthBackoff(attempt, 0))
				continue
			}
			return "", time.Time{}, lastErr
		}

		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		_ = resp.Body.Close()

		// 5xx is retryable; 429 is retryable (with Retry-After honoured).
		// 4xx other than 429 and 2xx are terminal.
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			lastErr = &ApiError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("oauth token request failed: %s", strings.TrimSpace(string(respBody))),
				RetryAfter: retryAfter,
				sentinel:   classifyStatus(resp.StatusCode),
			}
			if attempt < maxRetries-1 {
				oauthSleep(ctx, oauthBackoff(attempt, retryAfter))
				continue
			}
			return "", time.Time{}, lastErr
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = &ApiError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("oauth token request rate-limited: %s", strings.TrimSpace(string(respBody))),
				RetryAfter: retryAfter,
				sentinel:   classifyStatus(resp.StatusCode),
			}
			if attempt < maxRetries-1 {
				oauthSleep(ctx, oauthBackoff(attempt, retryAfter))
				continue
			}
			return "", time.Time{}, lastErr
		}

		if resp.StatusCode != http.StatusOK {
			return "", time.Time{}, &ApiError{
				StatusCode: resp.StatusCode,
				Message:    fmt.Sprintf("oauth token request failed: %s", strings.TrimSpace(string(respBody))),
				RetryAfter: retryAfter,
				sentinel:   classifyStatus(resp.StatusCode),
			}
		}

		var result struct {
			Token            string `json:"token"`
			ExpiresInSeconds int    `json:"expiresInSeconds"`
			SubscriptionName string `json:"subscriptionName"`
		}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return "", time.Time{}, &ApiError{
				StatusCode: 0,
				Message:    fmt.Sprintf("decode oauth response: %s", err),
			}
		}

		if result.Token == "" {
			return "", time.Time{}, &ApiError{
				StatusCode: 0,
				Message:    "oauth response missing token",
			}
		}

		expiresIn := result.ExpiresInSeconds
		if expiresIn <= 0 {
			// Conservative fallback: 15 minutes.
			expiresIn = 900
		}
		expAt := time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)
		return result.Token, expAt, nil
	}

	if lastErr == nil {
		lastErr = &ApiError{StatusCode: 0, Message: "oauth token fetch failed after retries"}
	}
	return "", time.Time{}, lastErr
}
