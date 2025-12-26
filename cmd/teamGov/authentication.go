package teamGov

import (
	"fmt"

	"github.com/saricon83-sudo/Tyr365AdminCli/internal/auth"
)

// TokenCache is kept for backward compatibility.
// Deprecated: Use auth.GetGovernanceToken() instead.
var TokenCache string

// AuthGovernanceApi gets an authentication token for the Teams Governance API.
// Deprecated: Use auth.GetGovernanceToken() directly instead.
func AuthGovernanceApi() (string, error) {
	return auth.GetGovernanceToken()
}

// AuthGraphApi gets an authentication token for the Graph API.
// Deprecated: Use auth.GetGraphToken() directly instead.
func AuthGraphApi() (string, error) {
	return auth.GetGraphToken()
}

// RetrieveAuthToken returns the cached token.
// Deprecated: Use auth.GetGovernanceToken() instead.
func RetrieveAuthToken() (string, error) {
	return TokenCache, nil
}

// PrintToken prints the cached token.
// Deprecated: Consider using the auth package directly.
func PrintToken() {
	fmt.Println(TokenCache)
}
