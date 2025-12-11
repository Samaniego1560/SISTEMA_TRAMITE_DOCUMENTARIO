// pkg/notificacion/mapper.go
package notificacion

import (
	"dbu-api/internal/models"
	"dbu-api/pkg/notificacion/smtp"
	"dbu-api/pkg/sanction/fault"
	"fmt"
)

// Nueva versión: recibe convocatorias y sanción asignada
func MapeaFaltaAgrupadaANotificationData(
	df fault.DetalleFaltaAgrupado,
	convocatorias []models.Convocatoria,
	sancion *models.SancionAsignadaDetalle,
) smtp.NotificationData {
	notificationData := smtp.NotificationData{
		FaltaID:           df.FaltaID,
		FuenteInformacion: df.FuenteInformacion,
		Alumno: smtp.Alumno{
			ID:                  int(df.Alumno.ID),
			CodigoEstudiante:    df.Alumno.CodigoEstudiante,
			DNI:                 df.Alumno.DNI,
			Nombres:             df.Alumno.Nombres,
			ApellidoPaterno:     df.Alumno.ApellidoPaterno,
			ApellidoMaterno:     df.Alumno.ApellidoMaterno,
			Sexo:                df.Alumno.Sexo,
			Facultad:            df.Alumno.Facultad,
			EscuelaProfesional:  df.Alumno.EscuelaProfesional,
			Edad:                df.Alumno.Edad,
			CorreoInstitucional: df.Alumno.CorreoInstitucional,
			Direccion:           df.Alumno.Direccion,
			LugarProcedencia:    df.Alumno.LugarProcedencia,
			CelularEstudiante:   df.Alumno.CelularEstudiante,
		},
		Resolucion: smtp.Resolucion{
			ResolucionID:     df.Resolucion.ResolucionID,
			ResolucionNombre: df.Resolucion.ResolucionNombre,
			Capitulos:        MapeaCapitulos(df.Resolucion.Capitulos),
		},
		Documentos: df.Documentos,
		FechaFalta: df.FechaFalta,
		Servicio:   df.Servicio,
	}

	// Buscar el semestre académico
	semestre := ""
	for _, conv := range convocatorias {
		if conv.ID == df.ConvocatoriaId {
			semestre = conv.Nombre // O el campo que contenga el semestre
			break
		}
	}

	// Poblar obsLabel y fechaHtml según la sanción asignada
	obsLabel := ""
	fechaHtml := ""
	if sancion != nil {
		obsLabel = sancion.Observaciones
		if obsLabel == "separacion temporal" || obsLabel == "separacion definitiva" {
			if sancion.FechaInicio != nil && sancion.FechaFin != nil {
				fechaHtml = fmt.Sprintf(" (del %s al %s)", sancion.FechaInicio.Format("02/01/2006"), sancion.FechaFin.Format("02/01/2006"))
			}
		}
	}

	notificationData.SemestreAcademico = semestre
	notificationData.ObsLabel = obsLabel
	notificationData.FechaHtml = fechaHtml
	return notificationData
}

// Mapear capítulos
func MapeaCapitulos(caps []fault.CapituloDetalle) []smtp.Capitulo {
	var out []smtp.Capitulo
	for _, c := range caps {
		out = append(out, smtp.Capitulo{
			CapituloID:     c.CapituloID,
			CapituloNombre: c.CapituloNombre,
			Articulos:      MapeaArticulos(c.Articulos),
		})
	}
	return out
}

// Mapear artículos
func MapeaArticulos(arts []fault.ArticuloDetalle) []smtp.Articulo {
	var out []smtp.Articulo
	for _, a := range arts {
		out = append(out, smtp.Articulo{
			ArticuloID:          a.ArticuloID,
			ArticuloDescripcion: a.ArticuloDescripcion,
			ArticuloGravedad:    a.ArticuloGravedad,
			Incisos:             MapeaIncisos(a.Incisos),
		})
	}
	return out
}

// Mapear incisos
func MapeaIncisos(incs []fault.IncisoDetalle) []smtp.Inciso {
	var out []smtp.Inciso
	for _, i := range incs {
		out = append(out, smtp.Inciso{
			IncisoID:          i.IncisoID,
			IncisoNombre:      i.IncisoNombre,
			IncisoDescripcion: i.IncisoDescripcion,
		})
	}
	return out
}
