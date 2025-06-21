package auth

import (
	"golang.org/x/oauth2"
)

// Replace these values with your Civic app details or load from env/config in production.
var CivicOauthConfig = &oauth2.Config{
	ClientID:     "YOUR_CIVIC_CLIENT_ID",
	ClientSecret: "YOUR_CIVIC_CLIENT_SECRET",                        // May not be needed if using PKCE
	RedirectURL:  "https://yourbackend.com/api/auth/civic/callback", // Must match your Civic app config
	Scopes:       []string{"openid"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://auth.civic.com/oauth",
		TokenURL: "https://auth.civic.com/oauth/token",
	},
}
