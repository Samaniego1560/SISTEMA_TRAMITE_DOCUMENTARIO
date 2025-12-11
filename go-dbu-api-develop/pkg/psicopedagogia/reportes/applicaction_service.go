package reportes

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
)

type PortsServerReportes interface {
	GetReportAttentionsDataTeachers(year, month string) ([]*models.SRQMonthlySummary, error)
	GetReportAttentionsDataStudents(year, month string) ([]*models.ConsultationAttentionExcel, error)
	GetReportPatientsByDateRange(startDate, endDate string) ([]*models.PatientReportExcel, error)
}

type service struct {
	repository ServiceReportesRepository
	txtID      string
}

func NewReportesService(repository ServiceReportesRepository, txID string) PortsServerReportes {
	return &service{repository: repository, txtID: txID}

}

func (s *service) GetReportAttentionsDataStudents(year, month string) ([]*models.ConsultationAttentionExcel, error) {
	m, err := s.repository.GetReportAttentionsDataStudents(year, month)
	if err != nil {
		logger.Error.Println(" - couldn't get consultation integran attention excel:", err)
		return nil, err
	}
	return m, nil
}

func (s *service) GetReportAttentionsDataTeachers(year, month string) ([]*models.SRQMonthlySummary, error) {

	m, err := s.repository.GetReportAttentionsDataTeachers(year, month)
	if err != nil {
		logger.Error.Println(" - couldn't get consultation integran attention excel:", err)
		return nil, err
	}
	return m, nil
}

func (s *service) GetReportPatientsByDateRange(startDate, endDate string) ([]*models.PatientReportExcel, error) {
	// Call repository to get the data
	data, err := s.repository.GetReportPatientsByDateRange(startDate, endDate)
	if err != nil {
		logger.Error.Println(" - couldn't get patient report by date range:", err)
		return nil, err
	}
	return data, nil
}
