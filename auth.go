package qubiton

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type oauth2TokenManager struct {
	clientID     string
	clientSecret string
	tokenURL     string
	httpClient   *http.Client

	mu    sync.Mutex
	token string
	expAt time.Time
}

func newOAuth2TokenManager(clientID, clientSecret, tokenURL string, httpClient *http.Client) *oauth2TokenManager {
	return &oauth2TokenManager{
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     tokenURL,
		httpClient:   httpClient,
	}
}

func (m *oauth2TokenManager) getToken() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.token != "" && time.Now().UTC().Before(m.expAt.Add(-30*time.Second)) {
		return m.token, nil
	}

	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {m.clientID},
		"client_secret": {m.clientSecret},
	}

	resp, err := m.httpClient.Post(m.tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("oauth2 token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oauth2 token request failed: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}

	expiresIn := result.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}

	m.token = result.AccessToken
	m.expAt = time.Now().UTC().Add(time.Duration(expiresIn) * time.Second)

	return m.token, nil
}
