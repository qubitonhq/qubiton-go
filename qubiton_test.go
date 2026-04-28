package qubiton

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ── Helpers ───────────────────────────────────────────────────────────────

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := NewClient("test-key", WithBaseURL(srv.URL))
	return srv, c
}

func mustReadJSON(t *testing.T, r io.Reader, into interface{}) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(into); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}

// ── Single-object response ────────────────────────────────────────────────

func TestValidateAddress_SingleObject(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("apikey"); got != "test-key" {
			t.Errorf("apikey header: got %q want %q", got, "test-key")
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("content-type: got %q want %q", got, "application/json")
		}
		if !strings.HasPrefix(r.Header.Get("User-Agent"), "qubiton-go-sdk/") {
			t.Errorf("user-agent: %q", r.Header.Get("User-Agent"))
		}
		if r.URL.Path != "/api/address/validate" {
			t.Errorf("path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"address1":"1600 Pennsylvania Ave NW","city":"Washington","state":"DC","postalCode":"20500"}`))
	})
	resp, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	if resp.City != "Washington" || resp.State != "DC" {
		t.Errorf("unexpected: %+v", resp)
	}
}

// ── Array-decode (tests the round-2 critical fix) ─────────────────────────

func TestLookupBusinessRegistration_Array(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"businessRegistrations":[{"registrationId":"R1","entityName":"Apple Inc","status":"Active"}],"validationPass":true,"validationDescription":"OK"},
			{"businessRegistrations":[{"registrationId":"R2","entityName":"Apple Inc (DBA)","status":"Active"}]}
		]`))
	})
	resp, err := c.LookupBusinessRegistration(context.Background(), BusinessRegistrationRequest{EntityName: "Apple Inc", Country: "US"})
	if err != nil {
		t.Fatalf("LookupBusinessRegistration: %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("len: got %d want 2", len(resp))
	}
	if len(resp[0].BusinessRegistrations) != 1 || resp[0].BusinessRegistrations[0].EntityName != "Apple Inc" {
		t.Errorf("first wrapper: %+v", resp[0])
	}
	if resp[0].ValidationPass == nil || *resp[0].ValidationPass != true {
		t.Errorf("validationPass: %+v", resp[0].ValidationPass)
	}
	if len(resp[1].BusinessRegistrations) != 1 || resp[1].BusinessRegistrations[0].RegistrationId != "R2" {
		t.Errorf("second wrapper: %+v", resp[1])
	}
}

func TestCheckSanctions_Array(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"isMatch":true,"score":0.85,"description":"OFAC"},{"isMatch":false,"score":0,"description":"EU"}]`))
	})
	resp, err := c.CheckSanctions(context.Background(), SanctionsRequest{CompanyName: "X"})
	if err != nil {
		t.Fatalf("CheckSanctions: %v", err)
	}
	if len(resp) != 2 || !resp[0].IsMatch || resp[1].IsMatch {
		t.Fatalf("unexpected: %+v", resp)
	}
}

func TestLookupExchangeRates_PathParamArrayBody(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/currency/exchange-rates/USD" {
			t.Errorf("path: %q", r.URL.Path)
		}
		var dates []string
		mustReadJSON(t, r.Body, &dates)
		if len(dates) != 2 {
			t.Errorf("expected 2 dates, got %d", len(dates))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"date":"2026-04-01","baseCurrency":"USD","exchangeRates":[{"currency":"EUR","rate":0.91,"dataSource":"ECB"}]}]`))
	})
	when, _ := time.Parse(time.RFC3339, "2026-04-01T00:00:00Z")
	when2, _ := time.Parse(time.RFC3339, "2026-04-02T00:00:00Z")
	resp, err := c.LookupExchangeRates(context.Background(), ExchangeRateRequest{
		BaseCurrency: "USD",
		Dates:        []time.Time{when, when2},
	})
	if err != nil {
		t.Fatalf("LookupExchangeRates: %v", err)
	}
	if len(resp) != 1 || len(resp[0].ExchangeRates) != 1 || resp[0].ExchangeRates[0].Currency != "EUR" || resp[0].ExchangeRates[0].Rate != 0.91 {
		t.Errorf("unexpected: %+v", resp)
	}
}

// ── OAuth flow ────────────────────────────────────────────────────────────

func TestOAuth_FetchesTokenAndCachesIt(t *testing.T) {
	t.Parallel()
	var fetches int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/oauth/token":
			atomic.AddInt32(&fetches, 1)
			var body map[string]string
			mustReadJSON(t, r.Body, &body)
			if body["key"] != "client-id" || body["secret"] != "client-secret" {
				t.Errorf("oauth body: %+v", body)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"jwe-token-xyz","expiresInSeconds":900,"subscriptionName":"Test"}`))
		case "/api/address/validate":
			if got := r.Header.Get("Authorization"); got != "Bearer jwe-token-xyz" {
				t.Errorf("authz: %q", got)
			}
			if r.Header.Get("apikey") != "" {
				t.Errorf("apikey header should be absent under OAuth, got %q", r.Header.Get("apikey"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"address1":"X"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	c := NewClient("", WithBaseURL(srv.URL), WithOAuth("client-id", "client-secret", ""))
	for i := 0; i < 3; i++ {
		_, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
		if err != nil {
			t.Fatalf("ValidateAddress[%d]: %v", i, err)
		}
	}
	if got := atomic.LoadInt32(&fetches); got != 1 {
		t.Errorf("expected 1 oauth fetch (cached), got %d", got)
	}
}

func TestOAuth_OptionOrderingDoesNotMatter(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"t","expiresInSeconds":900}`))
		}
	}))
	t.Cleanup(srv.Close)

	// WithOAuth before WithBaseURL — token URL should still resolve to srv.URL.
	c := NewClient("", WithOAuth("k", "s", ""), WithBaseURL(srv.URL))
	tok, err := c.oauth.getToken(context.Background())
	if err != nil {
		t.Fatalf("getToken: %v", err)
	}
	if tok != "t" {
		t.Errorf("token: %q", tok)
	}
	if c.oauth.tokenURL != srv.URL+"/api/oauth/token" {
		t.Errorf("tokenURL: %q", c.oauth.tokenURL)
	}
}

func TestOAuth_BadResponseWrapsAsApiError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"expiresInSeconds":900}`)) // missing "token"
	}))
	t.Cleanup(srv.Close)
	c := NewClient("", WithBaseURL(srv.URL), WithOAuth("k", "s", ""))
	_, err := c.oauth.getToken(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *ApiError, got %T: %v", err, err)
	}
}

func TestOAuth_DeduplicatesConcurrentFetches(t *testing.T) {
	t.Parallel()
	var fetches int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&fetches, 1)
		// simulate slow auth server
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"t","expiresInSeconds":900}`))
	}))
	t.Cleanup(srv.Close)
	c := NewClient("", WithBaseURL(srv.URL), WithOAuth("k", "s", ""))
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = c.oauth.getToken(context.Background())
		}()
	}
	wg.Wait()
	if got := atomic.LoadInt32(&fetches); got != 1 {
		t.Errorf("expected 1 deduped fetch, got %d", got)
	}
}

// ── Retry behavior ────────────────────────────────────────────────────────

func TestRetry_5xxThenSuccess(t *testing.T) {
	t.Parallel()
	var calls int32
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 2 {
			http.Error(w, `{"error":"transient"}`, http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"address1":"OK"}`))
	})
	// Reduce backoff under test by overriding via reflection-free workaround:
	// the backoff floor is 50ms — fine for tests.
	resp, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	if resp.Address1 != "OK" {
		t.Errorf("address1: %q", resp.Address1)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("calls: got %d want 2", got)
	}
}

func TestRetry_4xxNotRetried(t *testing.T) {
	t.Parallel()
	var calls int32
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
	})
	_, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *ApiError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 ApiError, got %v", err)
	}
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation sentinel")
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("expected 1 call, got %d", got)
	}
}

func TestRetry_OAuth401DoesNotConsumeAttempt(t *testing.T) {
	t.Parallel()
	var calls int32
	var tokenFetches int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/oauth/token":
			atomic.AddInt32(&tokenFetches, 1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"t","expiresInSeconds":900}`))
		case "/api/address/validate":
			n := atomic.AddInt32(&calls, 1)
			if n == 1 {
				http.Error(w, `{"error":"expired"}`, http.StatusUnauthorized)
				return
			}
			// After 401 + token refresh, second attempt should not consume an attempt slot,
			// and a subsequent 5xx should still get the full retry budget.
			if n == 2 {
				http.Error(w, `{"error":"oops"}`, http.StatusBadGateway)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"address1":"OK"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	c := NewClient("", WithBaseURL(srv.URL), WithOAuth("k", "s", ""))
	resp, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	if resp.Address1 != "OK" {
		t.Errorf("address1: %q", resp.Address1)
	}
	// Expect calls = 3: initial 401, retry-after-refresh = 502, retry-5xx = 200.
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Errorf("expected 3 calls, got %d", got)
	}
	if got := atomic.LoadInt32(&tokenFetches); got != 2 {
		t.Errorf("expected 2 token fetches (initial + after invalidate), got %d", got)
	}
}

// ── Error parsing ─────────────────────────────────────────────────────────

func TestErrorParsing_ProblemDetailsErrorsMap(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"title":"Validation failed","status":400,"errors":{"Country":["The Country field is required."]}}`))
	})
	_, err := c.ValidateAddress(context.Background(), AddressRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *ApiError
	if !errors.As(err, &apiErr) {
		t.Fatalf("not ApiError: %v", err)
	}
	if !strings.Contains(apiErr.Message, "Country") {
		t.Errorf("message should include errors map: %q", apiErr.Message)
	}
}

func TestWithTimeout_DoesNotMutateSharedClient(t *testing.T) {
	t.Parallel()
	hc := &http.Client{Timeout: 5 * time.Second}
	_ = NewClient("k", WithHTTPClient(hc), WithTimeout(1*time.Second))
	if hc.Timeout != 5*time.Second {
		t.Errorf("WithTimeout mutated caller-supplied client: got %v want 5s", hc.Timeout)
	}
}

// ── splitNonEmpty has been removed; ensure ExchangeRateRequest serialises empty as [] ──

func TestExchangeRates_EmptyDatesSendsEmptyArray(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if strings.TrimSpace(string(body)) != "[]" {
			t.Errorf("expected empty array body, got %q", string(body))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	})
	resp, err := c.LookupExchangeRates(context.Background(), ExchangeRateRequest{BaseCurrency: "USD"})
	if err != nil {
		t.Fatalf("LookupExchangeRates: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected 0 results, got %d", len(resp))
	}
}

// TestExchangeRates_BaseCurrencyURLEscaped verifies the path parameter is
// escaped — guards against a regression where weird/empty currencies could
// inject path segments.
func TestExchangeRates_BaseCurrencyURLEscaped(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// %20 is space; should arrive as "/api/currency/exchange-rates/US%20D" → decoded to "US D".
		// httptest decodes URL.Path by default; we assert on RawPath when present.
		if !strings.Contains(r.URL.Path, "US D") && !strings.Contains(r.URL.RawPath, "US%20D") {
			t.Errorf("path: %q raw=%q", r.URL.Path, r.URL.RawPath)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	})
	_, err := c.LookupExchangeRates(context.Background(), ExchangeRateRequest{BaseCurrency: "US D"})
	if err != nil {
		t.Fatalf("LookupExchangeRates: %v", err)
	}
}

// TestRequestedByClient_DefaultedFromUserAgent verifies the SDK injects the
// User-Agent into the body's requestedByClient field when callers leave it
// blank — without this every authenticated request returns 400.
func TestRequestedByClient_DefaultedFromUserAgent(t *testing.T) {
	t.Parallel()
	var seen string
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		mustReadJSON(t, r.Body, &body)
		if v, ok := body["requestedByClient"].(string); ok {
			seen = v
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"address1":"OK"}`))
	})
	_, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	if !strings.HasPrefix(seen, "qubiton-go-sdk/") {
		t.Errorf("requestedByClient: got %q, expected SDK user-agent", seen)
	}
}

// TestRequestedByClient_PreservesCallerValue verifies an explicit caller value
// is not overwritten by the default injector.
func TestRequestedByClient_PreservesCallerValue(t *testing.T) {
	t.Parallel()
	var seen string
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		mustReadJSON(t, r.Body, &body)
		if v, ok := body["requestedByClient"].(string); ok {
			seen = v
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"address1":"OK"}`))
	})
	req := AddressRequest{Country: "US"}
	req.RequestedByClient = "my-app/1.0"
	_, err := c.ValidateAddress(context.Background(), req)
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	if seen != "my-app/1.0" {
		t.Errorf("requestedByClient: got %q, expected 'my-app/1.0'", seen)
	}
}

// TestRetryAfter_HonoredOn429 ensures the SDK respects the server's Retry-After
// header on rate-limit responses.
func TestRetryAfter_HonoredOn429(t *testing.T) {
	t.Parallel()
	var calls int32
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 2 {
			w.Header().Set("Retry-After", "0")
			http.Error(w, `{"error":"slow down"}`, http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"address1":"OK"}`))
	})
	resp, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	if resp.Address1 != "OK" {
		t.Errorf("expected OK, got %q", resp.Address1)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Errorf("expected 2 calls, got %d", got)
	}
}

// TestNetworkError_RetriesUntilExhausted verifies network errors are retried
// the full attempt budget then surfaced as the final error.
func TestNetworkError_RetriesUntilExhausted(t *testing.T) {
	t.Parallel()
	// Point the client at a closed server — Dial will fail.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()
	c := NewClient("k", WithBaseURL(srv.URL))
	_, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err == nil {
		t.Fatal("expected network error")
	}
	// Either a raw url.Error or a wrapped *ApiError with StatusCode=0.
	if !strings.Contains(err.Error(), "connect") &&
		!strings.Contains(err.Error(), "refused") &&
		!strings.Contains(err.Error(), "EOF") &&
		!strings.Contains(err.Error(), "request failed after retries") {
		t.Errorf("expected network error, got: %v", err)
	}
}

// TestRetry_Repeated401WithOAuthDoesNotLoop ensures a server that keeps
// returning 401 even after a token refresh terminates with a 401, not an
// infinite loop. authRetried guards exactly one transparent retry.
func TestRetry_Repeated401WithOAuthDoesNotLoop(t *testing.T) {
	t.Parallel()
	var validateCalls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/oauth/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"t","expiresInSeconds":900}`))
		case "/api/address/validate":
			atomic.AddInt32(&validateCalls, 1)
			http.Error(w, `{"error":"still expired"}`, http.StatusUnauthorized)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	c := NewClient("", WithBaseURL(srv.URL), WithOAuth("k", "s", ""))
	_, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
	if err == nil {
		t.Fatal("expected 401")
	}
	var apiErr *ApiError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 ApiError, got: %v", err)
	}
	// Expect exactly 2 calls: initial 401 and one transparent retry after refresh.
	if got := atomic.LoadInt32(&validateCalls); got != 2 {
		t.Errorf("expected 2 validate calls, got %d", got)
	}
}

// TestExtractErrorMessage_DeterministicErrorOrder verifies that the
// fields-map → message rendering produces a stable order even when Go map
// iteration is randomised.
func TestExtractErrorMessage_DeterministicErrorOrder(t *testing.T) {
	t.Parallel()
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"title":"Validation","errors":{"Z_Field":["zz"],"A_Field":["aa"],"M_Field":["mm"]}}`))
	})
	const N = 5
	var msgs []string
	for i := 0; i < N; i++ {
		_, err := c.ValidateAddress(context.Background(), AddressRequest{})
		var apiErr *ApiError
		if !errors.As(err, &apiErr) {
			t.Fatalf("not ApiError: %v", err)
		}
		msgs = append(msgs, apiErr.Message)
	}
	for i := 1; i < N; i++ {
		if msgs[i] != msgs[0] {
			t.Errorf("non-deterministic message order:\n[0]=%q\n[%d]=%q", msgs[0], i, msgs[i])
		}
	}
	if !strings.Contains(msgs[0], "A_Field") || !strings.Contains(msgs[0], "Z_Field") {
		t.Errorf("missing fields: %q", msgs[0])
	}
	// A_Field should appear before Z_Field (alphabetical).
	if a, z := strings.Index(msgs[0], "A_Field"), strings.Index(msgs[0], "Z_Field"); a > z || a < 0 || z < 0 {
		t.Errorf("alphabetical order broken: %q", msgs[0])
	}
}

// TestScreenContinuous_NotImplementedNotRetried verifies 501 is terminal —
// the SDK does not waste retries against an unimplemented endpoint.
func TestScreenContinuous_NotImplementedNotRetried(t *testing.T) {
	t.Parallel()
	var calls int32
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		http.Error(w, `{"error":"not implemented"}`, http.StatusNotImplemented)
	})
	_, err := c.ScreenContinuous(context.Background(), ContinuousScreeningRequest{EntityName: "X"})
	if err == nil {
		t.Fatal("expected 501")
	}
	var apiErr *ApiError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotImplemented {
		t.Fatalf("expected 501 ApiError, got: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("expected 1 call (no retries), got %d", got)
	}
}

// TestIdempotencyKey_AddedToHeader verifies WithIdempotencyKey adds the header.
func TestIdempotencyKey_AddedToHeader(t *testing.T) {
	t.Parallel()
	var seen string
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("X-Idempotency-Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"address1":"OK"}`))
	})
	ctx := WithIdempotencyKey(context.Background(), "test-key-abc")
	_, err := c.ValidateAddress(ctx, AddressRequest{Country: "US"})
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	if seen != "test-key-abc" {
		t.Errorf("X-Idempotency-Key: got %q want %q", seen, "test-key-abc")
	}
}

// TestIdempotencyKey_PersistsAcrossRetries verifies the idempotency key is sent
// on every retry attempt — without persistence, a 5xx-then-success path would
// generate two distinct logical operations on the server side.
func TestIdempotencyKey_PersistsAcrossRetries(t *testing.T) {
	t.Parallel()
	var seenKeys []string
	var mu sync.Mutex
	var calls int32
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seenKeys = append(seenKeys, r.Header.Get("X-Idempotency-Key"))
		mu.Unlock()
		n := atomic.AddInt32(&calls, 1)
		if n < 2 {
			http.Error(w, `{"error":"transient"}`, http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"address1":"OK"}`))
	})
	ctx := WithIdempotencyKey(context.Background(), "persist-key-xyz")
	_, err := c.ValidateAddress(ctx, AddressRequest{Country: "US"})
	if err != nil {
		t.Fatalf("ValidateAddress: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(seenKeys) < 2 {
		t.Fatalf("expected ≥2 attempts, got %d", len(seenKeys))
	}
	for i, k := range seenKeys {
		if k != "persist-key-xyz" {
			t.Errorf("attempt %d: idempotency key %q (want %q)", i, k, "persist-key-xyz")
		}
	}
}

// TestSentinelErrors verifies errors.Is matches the documented sentinels for
// each HTTP status class.
func TestSentinelErrors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		status   int
		sentinel error
	}{
		{"Unauthorized→ErrAuth", http.StatusUnauthorized, ErrAuth},
		{"Forbidden→ErrAuth", http.StatusForbidden, ErrAuth},
		{"NotFound→ErrNotFound", http.StatusNotFound, ErrNotFound},
		{"TooManyRequests→ErrRateLimit", http.StatusTooManyRequests, ErrRateLimit},
		{"BadRequest→ErrValidation", http.StatusBadRequest, ErrValidation},
		{"UnprocessableEntity→ErrValidation", http.StatusUnprocessableEntity, ErrValidation},
		{"InternalServerError→ErrServerError", http.StatusInternalServerError, ErrServerError},
		{"BadGateway→ErrServerError", http.StatusBadGateway, ErrServerError},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, `{"error":"x"}`, tc.status)
			})
			_, err := c.ValidateAddress(context.Background(), AddressRequest{Country: "US"})
			if err == nil {
				t.Fatalf("expected error for status %d", tc.status)
			}
			if !errors.Is(err, tc.sentinel) {
				t.Errorf("status %d: errors.Is did not match sentinel %v: %v", tc.status, tc.sentinel, err)
			}
		})
	}
}

// TestArrayEndpoints_AllDecodeAsArrays exercises every documented array
// endpoint — guards against silent regression when adding new endpoints.
func TestArrayEndpoints_AllDecodeAsArrays(t *testing.T) {
	t.Parallel()
	bodies := map[string]string{
		"/api/businessregistration/lookup":     `[{"businessRegistrations":[{"entityName":"X"}]}]`,
		"/api/prohibited/lookup":               `[{"isMatch":true}]`,
		"/api/pep/lookup":                      `[{"score":0.5}]`,
		"/api/criminalprosecution/validate":    `[{"name":"X"}]`,
		"/api/criminalprosecution/lookup":      `[{"name":"Y"}]`,
		"/api/providerexclusion/validate":      `[{"firstName":"X","isProviderExcluded":true}]`,
		"/api/risk/lookup":                     `[{"totalArticles":1}]`,
		"/api/creditanalysis/lookup":           `[{"applicationId":"A1"}]`,
		"/api/esg/Scores":                      `[{"name":"X"}]`,
		"/api/beneficialownership/lookup":      `[{"responseCodeField":"OK"}]`,
		"/api/corporatehierarchy/lookup":       `[{"sourceResultCode":"OK"}]`,
		"/api/duns-number-lookup":              `[{"dunsNumber":"123"}]`,
		"/api/dot/fmcsa/lookup":                `[{"dotNumber":"99"}]`,
		"/api/businessclassification/lookup":   `[{"naicsCode":"541"}]`,
		"/api/aribasupplierprofile/lookup":     `[{"aNID":"AN1"}]`,
		"/api/certification/lookup":            `[{"certificationType":"MBE"}]`,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := bodies[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	c := NewClient("k", WithBaseURL(srv.URL))
	ctx := context.Background()

	if r, err := c.LookupBusinessRegistration(ctx, BusinessRegistrationRequest{EntityName: "X", Country: "US"}); err != nil || len(r) != 1 {
		t.Errorf("LookupBusinessRegistration: %v len=%d", err, len(r))
	}
	if r, err := c.CheckSanctions(ctx, SanctionsRequest{CompanyName: "X"}); err != nil || len(r) != 1 {
		t.Errorf("CheckSanctions: %v len=%d", err, len(r))
	}
	if r, err := c.ScreenPEP(ctx, PEPRequest{Name: "X", Country: "US"}); err != nil || len(r) != 1 {
		t.Errorf("ScreenPEP: %v len=%d", err, len(r))
	}
	if r, err := c.CheckEPAProsecution(ctx, EPARequest{Name: "X"}); err != nil || len(r) != 1 {
		t.Errorf("CheckEPAProsecution: %v len=%d", err, len(r))
	}
	if r, err := c.LookupEPAProsecution(ctx, EPARequest{Name: "X"}); err != nil || len(r) != 1 {
		t.Errorf("LookupEPAProsecution: %v len=%d", err, len(r))
	}
	if r, err := c.CheckHealthcareExclusion(ctx, HealthcareExclusionRequest{HealthCareType: "HCP"}); err != nil || len(r) != 1 || !r[0].IsProviderExcluded {
		t.Errorf("CheckHealthcareExclusion: %v len=%d", err, len(r))
	}
	if r, err := c.LookupRisk(ctx, RiskLookupRequest{EntityName: "X", Category: "Social"}); err != nil || len(r) != 1 {
		t.Errorf("LookupRisk: %v len=%d", err, len(r))
	}
	if r, err := c.LookupCreditAnalysis(ctx, CreditAnalysisRequest{CompanyName: "X", AddressLine1: "Y", City: "C", State: "S", Country: "US"}); err != nil || len(r) != 1 {
		t.Errorf("LookupCreditAnalysis: %v len=%d", err, len(r))
	}
	if r, err := c.LookupESGScore(ctx, ESGRequest{CompanyName: "X"}); err != nil || len(r) != 1 {
		t.Errorf("LookupESGScore: %v len=%d", err, len(r))
	}
	if r, err := c.LookupBeneficialOwnership(ctx, BeneficialOwnershipRequest{CompanyName: "X", CountryISO2: "US"}); err != nil || len(r) != 1 || r[0].ResponseCodeField != "OK" {
		t.Errorf("LookupBeneficialOwnership: %v len=%d", err, len(r))
	}
	if r, err := c.LookupCorporateHierarchy(ctx, CorporateHierarchyRequest{CompanyName: "X"}); err != nil || len(r) != 1 {
		t.Errorf("LookupCorporateHierarchy: %v len=%d", err, len(r))
	}
	if r, err := c.LookupDUNS(ctx, DUNSRequest{DunsNumber: "123"}); err != nil || len(r) != 1 {
		t.Errorf("LookupDUNS: %v len=%d", err, len(r))
	}
	if r, err := c.LookupDOTCarrier(ctx, DOTCarrierRequest{DotNumber: "99"}); err != nil || len(r) != 1 {
		t.Errorf("LookupDOTCarrier: %v len=%d", err, len(r))
	}
	if r, err := c.LookupBusinessClassification(ctx, BusinessClassificationRequest{CompanyName: "X", City: "C", State: "S", Country: "US"}); err != nil || len(r) != 1 {
		t.Errorf("LookupBusinessClassification: %v len=%d", err, len(r))
	}
	if r, err := c.LookupAribaSupplier(ctx, AribaSupplierRequest{ANID: "AN1"}); err != nil || len(r) != 1 {
		t.Errorf("LookupAribaSupplier: %v len=%d", err, len(r))
	}
	if r, err := c.LookupCertification(ctx, CertificationRequest{CompanyName: "X", Country: "US"}); err != nil || len(r) != 1 {
		t.Errorf("LookupCertification: %v len=%d", err, len(r))
	}
}

// TestOptionOrdering_TimeoutThenHTTPClient verifies the documented behaviour
// of WithTimeout vs WithHTTPClient regardless of ordering.
//
// Order 1: WithHTTPClient(hc) then WithTimeout(d) → cloned client gets the new
// timeout, caller's `hc` is left untouched (already covered by
// TestWithTimeout_DoesNotMutateSharedClient).
//
// Order 2: WithTimeout(d) first then WithHTTPClient(hc) → the caller-supplied
// client wins outright; the earlier WithTimeout call is overridden because
// the option that runs last sets the final httpClient. Documenting this here.
func TestOptionOrdering_TimeoutBeforeHTTPClient_HTTPClientWins(t *testing.T) {
	t.Parallel()
	hc := &http.Client{Timeout: 7 * time.Second}
	c := NewClient("k", WithTimeout(1*time.Second), WithHTTPClient(hc))
	if c.httpClient != hc {
		t.Errorf("expected client to be the caller-supplied hc")
	}
	if hc.Timeout != 7*time.Second {
		t.Errorf("WithTimeout (earlier) should not have mutated caller's hc; got %v", hc.Timeout)
	}
}

// TestOAuth_FollowerContextCancellation verifies that a follower passing a
// cancelled context returns ctx.Err() immediately and does not corrupt the
// in-flight slot for the leader.
func TestOAuth_FollowerContextCancellation(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow auth server so the follower can race the leader.
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"t","expiresInSeconds":900}`))
	}))
	t.Cleanup(srv.Close)
	c := NewClient("", WithBaseURL(srv.URL), WithOAuth("k", "s", ""))

	leaderDone := make(chan error, 1)
	go func() {
		_, err := c.oauth.getToken(context.Background())
		leaderDone <- err
	}()

	// Give the leader time to start the in-flight fetch.
	time.Sleep(20 * time.Millisecond)

	followerCtx, followerCancel := context.WithCancel(context.Background())
	followerCancel() // cancel before calling
	if _, err := c.oauth.getToken(followerCtx); !errors.Is(err, context.Canceled) {
		t.Errorf("follower expected context.Canceled, got %v", err)
	}

	if err := <-leaderDone; err != nil {
		t.Errorf("leader getToken: %v", err)
	}
}

// TestOAuth_RetriesOn5xx verifies the OAuth fetcher retries transient 5xx
// responses with the data-plane backoff schedule.
func TestOAuth_RetriesOn5xx(t *testing.T) {
	t.Parallel()
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 2 {
			http.Error(w, `{"error":"down"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"t","expiresInSeconds":900}`))
	}))
	t.Cleanup(srv.Close)
	c := NewClient("", WithBaseURL(srv.URL), WithOAuth("k", "s", ""))
	tok, err := c.oauth.getToken(context.Background())
	if err != nil {
		t.Fatalf("getToken: %v", err)
	}
	if tok != "t" {
		t.Errorf("token: %q", tok)
	}
	if got := atomic.LoadInt32(&attempts); got < 2 {
		t.Errorf("expected ≥2 attempts (retry on 5xx), got %d", got)
	}
}

// TestInjectRequestedByClient_PreservesLargeInt64 verifies that the body
// injector does not corrupt int64 values larger than 2^53 (which would
// otherwise lose precision through map[string]interface{} round-tripping).
func TestInjectRequestedByClient_PreservesLargeInt64(t *testing.T) {
	t.Parallel()
	// 9_007_199_254_740_993 = 2^53 + 1; exceeds float64 mantissa.
	const bigID = int64(9_007_199_254_740_993)
	var raw []byte
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})
	id := bigID
	_, _ = c.ValidateBankAccount(context.Background(), BankAccountRequest{
		BankNumberType:          "ROUTING",
		Country:                 "US",
		BankAccountValidationId: &id,
	})
	if !strings.Contains(string(raw), "9007199254740993") {
		t.Errorf("expected exact int64 in body, got %s", string(raw))
	}
	// And requestedByClient should still have been injected.
	if !strings.Contains(string(raw), `"requestedByClient":"qubiton-go-sdk/`) {
		t.Errorf("expected requestedByClient in body, got %s", string(raw))
	}
}

// TestParseRetryAfter_FractionalSeconds verifies the Retry-After parser
// accepts fractional values (rounding up).
func TestParseRetryAfter_FractionalSeconds(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want int
	}{
		{"0.4", 1},
		{"1.0", 1},
		{"1.4", 2},
		{"  2.5  ", 3},
		{"-1.5", 0},
	}
	for _, tc := range cases {
		got := parseRetryAfter(tc.in)
		if got != tc.want {
			t.Errorf("parseRetryAfter(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
