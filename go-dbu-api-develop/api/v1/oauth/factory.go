package oauth

import (
	oauth "dbu-api/pkg/auth/oauth/providers"
	"dbu-api/pkg/auth/oauth/providers/microsoft"
)

// createMicrosoftProvider crea una instancia del proveedor OAuth de Microsoft
func createMicrosoftProvider() oauth.OAuthProvider {
	return microsoft.NewMicrosoftOAuthProvider()
}

// Esta función permite agregar más proveedores fácilmente en el futuro:
// func createGoogleProvider() oauth.OAuthProvider {
//     return google.NewGoogleOAuthProvider()
// }
