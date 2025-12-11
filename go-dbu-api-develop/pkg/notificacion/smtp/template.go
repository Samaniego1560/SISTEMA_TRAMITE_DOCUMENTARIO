package smtp

import (
	"bytes"
	"dbu-api/internal/models"
	"fmt"
	"html/template"
	"os"
	"strings"
)

const emailTemplateHTML = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Notificación UNAS</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;
            line-height: 1.6;
            color: #1f2937;
            background: #f3f4f6;
            padding: 20px 10px;
        }
        .email-wrapper {
            max-width: 680px;
            margin: 0 auto;
            background: #ffffff;
        }
        .header {
            background: linear-gradient(135deg, #1e5128 0%, #2c5530 100%);
            padding: 30px 40px;
            text-align: center;
            color: #ffffff;
        }
        .header-logo {
            font-size: 24px;
            font-weight: 700;
            letter-spacing: 0.5px;
            margin-bottom: 8px;
        }
        .header-subtitle {
            font-size: 14px;
            opacity: 0.95;
            font-weight: 500;
        }
        .notification-badge {
            background: #dc2626;
            color: #ffffff;
            padding: 12px 24px;
            text-align: center;
            font-weight: 700;
            font-size: 16px;
            letter-spacing: 0.5px;
        }
        .content-wrapper {
            padding: 40px;
        }
        .notification-number {
            display: inline-block;
            background: #f3f4f6;
            padding: 8px 16px;
            border-radius: 6px;
            font-weight: 600;
            color: #1e5128;
            margin-bottom: 24px;
            font-size: 15px;
        }
        .student-card {
            background: #f9fafb;
            border-left: 4px solid #1e5128;
            padding: 20px;
            margin: 24px 0;
            border-radius: 8px;
        }
        .student-card-row {
            display: flex;
            padding: 8px 0;
            border-bottom: 1px solid #e5e7eb;
        }
        .student-card-row:last-child {
            border-bottom: none;
        }
        .student-card-label {
            font-weight: 600;
            color: #1e5128;
            min-width: 120px;
            font-size: 14px;
        }
        .student-card-value {
            color: #374151;
            font-size: 14px;
        }
        .main-text {
            font-size: 15px;
            line-height: 1.8;
            color: #374151;
            margin: 24px 0;
            text-align: justify;
        }
        .highlight-red {
            color: #dc2626;
            font-weight: 700;
        }
        .highlight-blue {
            color: #2563eb;
            font-weight: 600;
        }
        .violations-box {
            background: #fef3c7;
            border: 2px solid #fbbf24;
            border-radius: 8px;
            padding: 20px;
            margin: 24px 0;
        }
        .violations-title {
            font-weight: 700;
            color: #92400e;
            margin-bottom: 12px;
            font-size: 15px;
        }
        .violations-box ul {
            margin-left: 20px;
        }
        .violations-box li {
            margin: 8px 0;
            color: #78350f;
            font-size: 14px;
        }
        .sanction-box {
            background: #fee2e2;
            border: 2px solid #ef4444;
            border-radius: 8px;
            padding: 20px;
            margin: 24px 0;
        }
        .sanction-title {
            font-weight: 700;
            color: #991b1b;
            margin-bottom: 8px;
            font-size: 15px;
        }
        .sanction-text {
            color: #7f1d1d;
            font-size: 14px;
            line-height: 1.6;
        }
        .note-box {
            background: #dbeafe;
            border: 2px solid #3b82f6;
            border-radius: 8px;
            padding: 20px;
            margin: 24px 0;
        }
        .note-title {
            font-weight: 700;
            color: #1e40af;
            margin-bottom: 8px;
            font-size: 15px;
        }
        .note-text {
            color: #1e3a8a;
            font-size: 14px;
            line-height: 1.6;
        }
        .footer {
            background: #1f2937;
            padding: 30px 40px;
            display: flex;
            align-items: center;
            gap: 20px;
        }
        .footer-logo {
            height: 70px;
            width: 70px;
            background: #374151;
            border-radius: 8px;
            flex-shrink: 0;
        }
        .footer-content {
            color: #e5e7eb;
            font-size: 13px;
            line-height: 1.8;
        }
        .footer-title {
            font-weight: 700;
            color: #ffffff;
            margin-bottom: 4px;
            font-size: 14px;
        }
        .footer-link {
            color: #60a5fa;
            text-decoration: none;
        }
        .footer-link:hover {
            text-decoration: underline;
        }
        .disclaimer {
            text-align: center;
            padding: 20px;
            background: #f9fafb;
            color: #6b7280;
            font-size: 12px;
            line-height: 1.6;
        }
        @media only screen and (max-width: 640px) {
            .content-wrapper {
                padding: 24px 20px;
            }
            .footer {
                flex-direction: column;
                text-align: center;
            }
            .student-card-row {
                flex-direction: column;
            }
            .student-card-label {
                margin-bottom: 4px;
            }
        }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <!-- Header -->
        <div class="header">
            <div class="header-logo">UNIVERSIDAD NACIONAL AGRARIA DE LA SELVA</div>
            <div class="header-subtitle">Dirección de Bienestar Universitario - Tingo María</div>
			<div class="header-subtitlee"> SEMESTRE ACADEMICO {{.SemestreAcademico}}</div>
        </div>

        <!-- Notification Badge -->
        <div class="notification-badge">
            ⚠️ NOTIFICACIÓN OFICIAL DE FALTA
        </div>
        <!-- Main Content -->
        <div class="content-wrapper">
            <div class="notification-number">
                📄 Notificación N° {{.NumeroNotificacion}}
            </div>

            <!-- Student Information Card -->
            <div class="student-card">
                <div class="student-card-row">
                    <span class="student-card-label">Estudiante:</span>
                    <span class="student-card-value">{{.Estudiante}}</span>
                </div>
                <div class="student-card-row">
                    <span class="student-card-label">Facultad:</span>
                    <span class="student-card-value">{{.Facultad}}</span>
                </div>
                <div class="student-card-row">
                    <span class="student-card-label">Fecha:</span>
                    <span class="student-card-value">{{.Fecha}}</span>
                </div>
                <div class="student-card-row">
                    <span class="student-card-label">Dirección:</span>
                    <span class="student-card-value">{{.Direccion}}</span>
                </div>
            </div>

            <!-- Main Message -->
            <div class="main-text">
                Por la presente se le notifica que según el <span class="highlight-blue">{{or .FuenteInformacion .fuente_informacion}}</span>, usted ha incurrido en <span class="highlight-red">{{.ContadorFaltas}} FALTA {{.GravedadFalta}}</span>, de acuerdo con la Directiva que regula el uso del servicio de <span class="highlight-blue">{{.Servicio}}</span> de la UNAS, aprobada mediante <span class="highlight-blue">{{.Resolucion.ResolucionNombre}}</span>.
            </div>

            <!-- Violations Box -->
            <div class="violations-box">
                <div class="violations-title">📋 La falta está tipificada en los siguientes casos:</div>
                <ul>
                    {{range .FaltasDetalle}}
                    <li>{{.}}</li>
                    {{end}}
                </ul>
            </div>

            <!-- Sanction Box -->
            <div class="sanction-box">
                <div class="sanction-title">⚖️ Sanción Aplicada</div>
                <div class="sanction-text">
                    La sanción correspondiente se encuentra en el <span class="highlight-blue">{{.CapituloSancion}}</span>, <span class="highlight-blue">{{.ArticuloSancion}}</span>, Inc. <span class="highlight-blue">{{.IncisoSancion}}</span> de la Directiva, que establece: <strong>"{{.SancionDetalle}}"</strong>
                    <br><br>
                    Esta sanción se determina como: <strong>{{.ObsLabel}}</strong>{{.FechaHtml}} y ha sido registrada en el sistema de la Dirección de Bienestar Universitario (DBU) y en su expediente personal.
                </div>
            </div>

            <!-- Important Note -->
            <div class="note-box">
                <div class="note-title">⚠️ IMPORTANTE</div>
                <div class="note-text">
                    En caso de omitir la presente notificación, se procederá a tomar otras medidas que defina la Comisión de Disciplina de acuerdo con el reglamento vigente.
                </div>
            </div>
        </div>

        <!-- Footer -->
        <div class="footer">
            <div class="footer-logo">
                <img src="https://i.imgur.com/ChHqC1A.png" alt="Logo DBU" style="width: 100%; height: 100%; object-fit: contain;">
            </div>
            <div class="footer-content">
                <div class="footer-title">Comisión de Sanciones - DBU UNAS</div>
                <div>Universidad Nacional Agraria de la Selva</div>
                <div>
                    <a href="https://bienestar.unas.edu.pe" class="footer-link">www.bienestar.unas.edu.pe</a> | 
                    <a href="mailto:obu@unas.edu.pe" class="footer-link">obu@unas.edu.pe</a>
                </div>
                <div>📞 Teléfono: +51 922 012 218</div>
            </div>
        </div>

        <!-- Disclaimer -->
        <div class="disclaimer">
            Este es un documento oficial generado por el Sistema de Gestión de la Dirección de Bienestar Universitario de la Universidad Nacional Agraria de la Selva - UNAS
        </div>
    </div>
</body>
</html>`

// GenerateNotificationHTML genera el HTML de la notificación usando el template mejorado
func GenerateNotificationHTML(data EmailTemplate) (string, error) {
	// Crear funciones personalizadas para el template
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}

	tmpl, err := template.New("notification").Funcs(funcMap).Parse(emailTemplateHTML)
	if err != nil {
		return "", fmt.Errorf("error al parsear template: %v", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("error al ejecutar template: %v", err)
	}

	return buffer.String(), nil
}

// GenerateEmailHTML es un alias para mantener compatibilidad
func GenerateEmailHTML(data EmailTemplate) (string, error) {
	return GenerateNotificationHTML(data)
}

// MapNotificationDataToTemplate convierte tu estructura JSON al template de email
func MapNotificationDataToTemplate(data NotificationData, numeroNotificacion string, sancion *models.SancionAsignadaDetalle) EmailTemplate {
	nombreCompleto := strings.TrimSpace(data.Alumno.Nombres + " " + data.Alumno.ApellidoPaterno + " " + data.Alumno.ApellidoMaterno)

	// Contar faltas basado en los incisos
	contadorFaltas := 0
	gravedadMasAlta := "LEVE"

	for _, capitulo := range data.Resolucion.Capitulos {
		for _, articulo := range capitulo.Articulos {
			contadorFaltas += len(articulo.Incisos)
			if strings.ToUpper(articulo.ArticuloGravedad) == "GRAVE" {
				gravedadMasAlta = "GRAVE"
			}
		}
	}

	faltasDetalle := []string{}
	faltasDetalleText := BuildFaltasDetalleText(data)
	for _, line := range strings.Split(strings.TrimSpace(faltasDetalleText), "\n") {
		if line != "" {
			faltasDetalle = append(faltasDetalle, line)
		}
	}

	sancionDetalle := ""
	capituloSancion := ""
	articuloSancion := ""
	incisoSancion := ""
	descripcionSancion := ""

	if sancion != nil {
		sancionDetalle = sancion.DetalleSancion
		capituloSancion = sancion.CapituloSancion
		articuloSancion = sancion.ArticuloSancion
		incisoSancion = sancion.IncisoSancion
		descripcionSancion = sancion.DetalleSancion
	}

	return EmailTemplate{
		NumeroNotificacion: numeroNotificacion,
		Estudiante:         nombreCompleto,
		Facultad:           data.Alumno.Facultad,
		Fecha:              data.FechaFalta,
		Direccion:          data.Alumno.Direccion,
		ContadorFaltas:     contadorFaltas,
		GravedadFalta:      gravedadMasAlta,
		Servicio:           data.Servicio,
		NombreNotificador:  "Comité de disciplina DBU-UNAS",
		NombreNotificado:   nombreCompleto,
		FaltasDetalle:      faltasDetalle,
		SancionDetalle:     sancionDetalle,
		FuenteInformacion:  data.FuenteInformacion,
		CapituloSancion:    capituloSancion,
		ArticuloSancion:    articuloSancion,
		IncisoSancion:      incisoSancion,
		DescripcionSancion: descripcionSancion,
		Resolucion:         data.Resolucion,
		SemestreAcademico:  data.SemestreAcademico,
		ObsLabel:           data.ObsLabel,
		FechaHtml:          data.FechaHtml,
	}
}

// OTPEmailData estructura para datos del template OTP
type OTPEmailData struct {
	NombreEstudiante  string
	CodigoOTP         string
	MinutosExpiracion int
}

// GenerateOTPHTML genera el HTML del email OTP usando el template
func GenerateOTPHTML(templatePath, nombreEstudiante, codigoOTP string, minutosExpiracion int) (string, error) {
	// Leer el archivo del template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error al leer template OTP: %v", err)
	}

	// Parsear el template
	tmpl, err := template.New("otp").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("error al parsear template OTP: %v", err)
	}

	// Preparar datos
	data := OTPEmailData{
		NombreEstudiante:  nombreEstudiante,
		CodigoOTP:         codigoOTP,
		MinutosExpiracion: minutosExpiracion,
	}

	// Ejecutar template
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("error al ejecutar template OTP: %v", err)
	}

	return buffer.String(), nil
}
