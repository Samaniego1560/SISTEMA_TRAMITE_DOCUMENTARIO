package reportes

import (
	"dbu-api/internal/models"

	"github.com/jmoiron/sqlx"
)

const (
	Postgresql = "postgres"
	SqlServer  = "sqlserver"
	Oracle     = "oci8"
)

type ServiceReportesRepository interface {
	GetReportAttentionsDataTeachers(fecha_inicio, fecha_fin string) ([]*models.SRQMonthlySummary, error)
	GetReportAttentionsDataStudents(fecha_inicio, fecha_fin string) ([]*models.ConsultationAttentionExcel, error)
	GetReportPatientsByDateRange(startDate, endDate string) ([]*models.PatientReportExcel, error)
}

func FactoryStorage(db *sqlx.DB) ServiceReportesRepository {
	return NewReportesRepository(db)
}
