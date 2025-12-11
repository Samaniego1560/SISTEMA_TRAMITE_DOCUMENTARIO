package middleware

import (
	"dbu-api/internal/env"
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"errors"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/jmoiron/sqlx"
)

type StudentJWTClaims struct {
	AlumnoID  int64  `json:"alumno_id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"`
	DNI       string `json:"dni"`
	jwt.StandardClaims
}

// StudentJWTProtected middleware personalizado para proteger rutas de estudiantes
func StudentJWTProtected() fiber.Handler {
	e := env.NewConfiguration()
	config := jwtware.Config{
		ErrorHandler:  studentJWTError,
		SigningKey:    e.Key.PublicKey,
		SigningMethod: "RS256",
	}
	return jwtware.New(config)
}

func studentJWTError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token no proporcionado o mal formado", Type: "error"})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token inválido o expirado", Type: "error"})
}

func CreateStudentJWT(alumnoID int64, dni, sessionID string) (string, error) {
	e := env.NewConfiguration()

	claims := StudentJWTClaims{
		AlumnoID:  alumnoID,
		SessionID: sessionID,
		Role:      "student",
		DNI:       dni,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    "dbu-api-student",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(e.Key.PrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateStudentSession middleware que valida la sesión del estudiante y almacena datos en c.Locals
func ValidateStudentSession(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el token del header
		bearer := c.Get("Authorization")
		if bearer == "" {
			logger.Warning.Printf("Token vacío para estudiante")
			return c.Status(fiber.StatusUnauthorized).
				JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token no proporcionado", Type: "error"})
		}

		if len(bearer) < 7 {
			logger.Warning.Printf("Token inválido para estudiante")
			return c.Status(fiber.StatusUnauthorized).
				JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token inválido", Type: "error"})
		}

		tkn := bearer[7:] // Remover "Bearer "

		// Parsear y validar el token
		e := env.NewConfiguration()
		verifyFunction := func(tkn *jwt.Token) (interface{}, error) {
			return e.Key.PublicKey, nil
		}

		token, err := jwt.ParseWithClaims(tkn, &StudentJWTClaims{}, verifyFunction)
		if err != nil {
			var validationError *jwt.ValidationError
			if !errors.As(err, &validationError) {
				logger.Warning.Printf("Error al procesar el token de estudiante: %v", err)
				return c.Status(fiber.StatusUnauthorized).
					JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token inválido", Type: "error"})
			}

			if validationError.Errors == jwt.ValidationErrorExpired {
				logger.Warning.Printf("Token de estudiante expirado: %v", err)
				return c.Status(fiber.StatusUnauthorized).
					JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token expirado", Type: "error"})
			}

			logger.Warning.Printf("Error de validación del token de estudiante: %v", err)
			return c.Status(fiber.StatusUnauthorized).
				JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token inválido", Type: "error"})
		}

		if !token.Valid {
			logger.Warning.Printf("Token de estudiante no válido")
			return c.Status(fiber.StatusUnauthorized).
				JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Token inválido", Type: "error"})
		}

		claims := token.Claims.(*StudentJWTClaims)

		// Verificar que el role sea "student"
		if claims.Role != "student" {
			logger.Warning.Printf("Token no es de estudiante, role: %s", claims.Role)
			return c.Status(fiber.StatusForbidden).
				JSON(models.Response{Error: true, Data: nil, Code: 0, Msg: "Acceso denegado", Type: "error"})
		}

		// Guardar datos del estudiante en c.Locals para uso en handlers
		c.Locals("alumno_id", claims.AlumnoID)
		c.Locals("session_id", claims.SessionID)
		c.Locals("dni", claims.DNI)

		logger.Info.Printf("Estudiante autenticado - alumno_id: %d, session_id: %s", claims.AlumnoID, claims.SessionID)

		return c.Next()
	}
}
