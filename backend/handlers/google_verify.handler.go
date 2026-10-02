package handlers

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type GoogleLoginRequest struct {
	Token string `json:"token" validate:"required"`
}

var (
	googleOIDCProvider *oidc.Provider
)

func init() {
	var err error
	ctx := context.Background()

	// Create provider
	googleOIDCProvider, err = oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		fmt.Printf("Failed to create OIDC provider: %v\n", err)
		return
	}

	fmt.Println("✓ Google OIDC provider initialized successfully")
}

func VerifyIdTokenGoogle(ctx context.Context, tokenString string, clientID string) (*GoogleClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}

	if clientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}

	// Create a verifier with the specific client ID
	verifier := googleOIDCProvider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	// Verify the token
	idToken, err := verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	// Extract claims
	var claims GoogleClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	// Validate email
	if !claims.EmailVerified {
		return nil, fmt.Errorf("email not verified by Google")
	}

	if claims.Email == "" {
		return nil, fmt.Errorf("no email in token")
	}

	// fmt.Printf("✓ Token verified for: %s\n", claims.Email)
	return &claims, nil
}
