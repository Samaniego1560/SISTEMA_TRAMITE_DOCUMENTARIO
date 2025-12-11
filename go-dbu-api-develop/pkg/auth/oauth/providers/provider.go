package oauth

import "context"

// OAuthUserInfo representa la información básica del usuario obtenida del proveedor OAuth
type OAuthUserInfo struct {
	ID                 string
	Email              string // Correo personal del usuario
	InstitutionalEmail string // Correo institucional (ej: UserPrincipalName en Microsoft)
	Name               string
	Picture            string
	Provider           string // "microsoft", "google", etc.
}

// OAuthProvider define la interfaz para el patrón Strategy
// Permite implementar diferentes proveedores OAuth (Microsoft, Google, GitHub, etc.)
type OAuthProvider interface {
	// GetAuthURL genera la URL de autorización para redirigir al usuario
	GetAuthURL(state string) string

	// ExchangeCode intercambia el código de autorización por un token de acceso
	// y obtiene la información del usuario
	ExchangeCode(ctx context.Context, code string) (*OAuthUserInfo, error)

	// GetProviderName retorna el nombre del proveedor (ej: "microsoft", "google")
	GetProviderName() string
}
