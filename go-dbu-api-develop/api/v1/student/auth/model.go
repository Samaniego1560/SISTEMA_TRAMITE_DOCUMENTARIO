package auth

// RequestOTPRequest request para solicitar código OTP
type RequestOTPRequest struct {
	DNI string `json:"dni" validate:"required"`
}

// RequestOTPResponse response al solicitar OTP
type RequestOTPResponse struct {
	CorreoParcial string `json:"correo_parcial"`
	ExpiraEnSeg   int    `json:"expira_en_segundos"`
}

// VerifyOTPRequest request para verificar código OTP
type VerifyOTPRequest struct {
	DNI       string `json:"dni" validate:"required"`
	CodigoOTP string `json:"codigo_otp" validate:"required,len=6"`
}

// AlumnoInfo información básica del alumno
type AlumnoInfo struct {
	ID                  int64  `json:"id"`
	CodigoEstudiante    string `json:"codigo_estudiante"`
	DNI                 string `json:"DNI"`
	Nombres             string `json:"nombres"`
	ApellidoPaterno     string `json:"apellido_paterno"`
	ApellidoMaterno     string `json:"apellido_materno"`
	CorreoInstitucional string `json:"correo_institucional"`
}

// VerifyOTPResponse response al verificar OTP exitosamente
type VerifyOTPResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}
