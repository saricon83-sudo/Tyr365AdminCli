// Package auth provides shared authentication and token caching functionality.
package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/saricon83-sudo/Tyr365AdminCli/internal/config"
	"github.com/saricon83-sudo/Tyr365AdminCli/internal/httputil"
)

// Token represents a cached API token with expiration
type Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

// IsValid checks if the token is still valid (not expired)
func (t *Token) IsValid() bool {
	return t != nil && time.Now().Before(t.ExpiresAt)
}

// TokenCache provides thread-safe token caching
type TokenCache struct {
	tokens map[string]*Token
	mutex  sync.RWMutex
}

// Global token cache instance
var cache = &TokenCache{
	tokens: make(map[string]*Token),
}

// GetCache returns the global token cache instance
func GetCache() *TokenCache {
	return cache
}

// Get retrieves a token from the cache
func (c *TokenCache) Get(resource string) (*Token, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	token, found := c.tokens[resource]
	if !found || !token.IsValid() {
		return nil, false
	}
	return token, true
}

// Set stores a token in the cache
func (c *TokenCache) Set(resource string, token *Token) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.tokens[resource] = token
}

// Clear removes a token from the cache
func (c *TokenCache) Clear(resource string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.tokens, resource)
}

// ClearAll removes all tokens from the cache
func (c *TokenCache) ClearAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.tokens = make(map[string]*Token)
}

// OAuthConfig holds OAuth2 client credentials configuration
type OAuthConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	Resource     string
	GrantType    string
}

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   interface{} `json:"expires_in"` // Can be string or number
	Resource    string      `json:"resource"`
}

// GetExpiresInSeconds parses the expires_in field which can be string or number
func (tr *TokenResponse) GetExpiresInSeconds() int {
	switch v := tr.ExpiresIn.(type) {
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 3600 // Default to 1 hour
}

// AcquireToken gets a token using OAuth2 client credentials flow.
// It uses the cache and only fetches a new token if necessary.
func AcquireToken(cfg OAuthConfig) (string, error) {
	// Check cache first
	if token, found := cache.Get(cfg.Resource); found {
		return token.AccessToken, nil
	}

	// Build token endpoint URL
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", cfg.TenantID)

	// Build request body
	body := url.Values{}
	body.Set("grant_type", cfg.GrantType)
	body.Set("client_id", cfg.ClientID)
	body.Set("client_secret", cfg.ClientSecret)
	body.Set("resource", cfg.Resource)

	// Make the request
	resp, err := httputil.MakePOSTRequest(tokenURL, []byte(body.Encode()))
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := httputil.ReadResponseBody(resp)
		return "", fmt.Errorf("token request returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	// Cache the token with 5-minute buffer before expiry
	expiresIn := tokenResp.GetExpiresInSeconds()
	cache.Set(cfg.Resource, &Token{
		AccessToken: tokenResp.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn-300) * time.Second),
	})

	return tokenResp.AccessToken, nil
}

// ============================================================================
// Pre-configured auth functions for specific APIs
// ============================================================================

// GetGovernanceToken gets a token for the Teams Governance API
func GetGovernanceToken() (string, error) {
	cfg := config.Get()
	return AcquireToken(OAuthConfig{
		TenantID:     cfg.GetString("O365TenantName"),
		ClientID:     cfg.GetString("client_id"),
		ClientSecret: cfg.GetString("client_secret"),
		Resource:     cfg.GetString("resource"),
		GrantType:    cfg.GetString("grant_type"),
	})
}

// GetGraphToken gets a token for Microsoft Graph API
func GetGraphToken() (string, error) {
	cfg := config.Get()
	return AcquireToken(OAuthConfig{
		TenantID:     cfg.GetString("O365TenantName"),
		ClientID:     cfg.GetString("client_id"),
		ClientSecret: cfg.GetString("client_secret"),
		Resource:     "https://graph.microsoft.com",
		GrantType:    "client_credentials",
	})
}

// GetArchiverToken gets a token for the M365 Archiver API
func GetArchiverToken() (string, error) {
	cfg := config.Get()
	resource := cfg.GetString("archiverResource")
	if resource == "" {
		return "", fmt.Errorf("archiverResource not found in configuration")
	}

	return AcquireToken(OAuthConfig{
		TenantID:     cfg.GetString("O365TenantName"),
		ClientID:     cfg.GetString("archiverApp"),
		ClientSecret: cfg.GetString("archiverSecret"),
		Resource:     resource,
		GrantType:    "client_credentials",
	})
}

// GetTeamToolboxToken gets a token for the Team Toolbox API
func GetTeamToolboxToken() (string, error) {
	cfg := config.Get()
	return AcquireToken(OAuthConfig{
		TenantID:     cfg.GetString("O365TenantName"),
		ClientID:     cfg.GetString("365ManagementAppId"),
		ClientSecret: cfg.GetString("365ManagementAppSecret"),
		Resource:     fmt.Sprintf("api://%s", cfg.GetString("teamToolboxAppId")),
		GrantType:    "client_credentials",
	})
}
