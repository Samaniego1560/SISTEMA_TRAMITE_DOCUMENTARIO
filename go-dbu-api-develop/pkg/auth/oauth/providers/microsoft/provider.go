package microsoft

import (
	"context"
	"dbu-api/internal/env"
	oauth "dbu-api/pkg/auth/oauth/providers"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

// MicrosoftUserInfo representa la respuesta de Microsoft Graph API
type MicrosoftUserInfo struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	UserPrincipalName string `json:"userPrincipalName"`
	Mail              string `json:"mail"`
}

// MicrosoftOAuthProvider implementa OAuthProvider para Microsoft
type MicrosoftOAuthProvider struct {
	config *oauth2.Config
}

// NewMicrosoftOAuthProvider crea una nueva instancia del proveedor de Microsoft
func NewMicrosoftOAuthProvider() *MicrosoftOAuthProvider {
	cfg := env.NewConfiguration()
	msConfig := cfg.OAuth.Microsoft

	return &MicrosoftOAuthProvider{
		config: &oauth2.Config{
			ClientID:     msConfig.ClientID,
			ClientSecret: msConfig.ClientSecret,
			RedirectURL:  msConfig.RedirectURI,
			Scopes: []string{
				"openid",
				"profile",
				"email",
				"User.Read",
			},
			Endpoint: microsoft.AzureADEndpoint(msConfig.TenantID),
		},
	}
}

// GetAuthURL genera la URL de autorización de Microsoft
func (m *MicrosoftOAuthProvider) GetAuthURL(state string) string {
	// prompt=select_account: Fuerza a mostrar el selector de cuentas, permitiendo al usuario
	// cambiar de cuenta o agregar una nueva en cada login
	return m.config.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)
}

// ExchangeCode intercambia el código por un token y obtiene la información del usuario
func (m *MicrosoftOAuthProvider) ExchangeCode(ctx context.Context, code string) (*oauth.OAuthUserInfo, error) {
	// Intercambiar código por token
	token, err := m.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Crear cliente HTTP con el token
	client := m.config.Client(ctx, token)

	// Obtener información del usuario desde Microsoft Graph API
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	// Decodificar respuesta
	var msUserInfo MicrosoftUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&msUserInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Convertir a estructura genérica
	// Email: correo personal (Mail)
	// InstitutionalEmail: correo institucional (UserPrincipalName)
	return &oauth.OAuthUserInfo{
		ID:                 msUserInfo.ID,
		Email:              msUserInfo.Mail,
		InstitutionalEmail: msUserInfo.UserPrincipalName,
		Name:               msUserInfo.DisplayName,
		Picture:            "", // Microsoft requiere llamada adicional a /me/photo/$value
		Provider:           "microsoft",
	}, nil
}

// GetProviderName retorna el nombre del proveedor
func (m *MicrosoftOAuthProvider) GetProviderName() string {
	return "microsoft"
}
