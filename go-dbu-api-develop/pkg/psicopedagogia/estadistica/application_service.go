package estadistica

import (
	"errors"
	"fmt"
)

type PortsServerEstadistica interface {
	GetDataChartArea(inicio, fin string) ([]DailyStat, error)
	GetDataChartPie(inicio, fin string) ([]DataPie, error)
	GetDataChartBarra(inicio, fin string) ([]DataBarras, error)
}
type service struct {
	repository ServiceEstadisticaRepository
	txID       string
}

func NewStadisticaService(repository ServiceEstadisticaRepository, txID string) PortsServerEstadistica {
	return &service{repository: repository, txID: txID}
}

func (s service) GetDataChartArea(inicio, fin string) ([]DailyStat, error) {

	if inicio == "" || fin == "" {
		return nil, errors.New("las fechas estan vacias")
	}

	var responses []DailyStat

	responses, err := s.repository.GetDataChartArea(inicio, fin)
	if err != nil {
		return nil, fmt.Errorf("error al extraer datos estadísticos: %w", err)
	}

	return responses, nil
}

func (s *service) GetDataChartPie(inicio, fin string) ([]DataPie, error) {
	if inicio == "" || fin == "" {
		return nil, errors.New("las fechas estan vacias")
	}

	var responses []DataPie

	responses, err := s.repository.GetDataChartPie(inicio, fin)

	if err != nil {
		return nil, fmt.Errorf("error al extraer datos para el chart pie: %w", err)
	}

	return responses, nil
}

func (s *service) GetDataChartBarra(inicio, fin string) ([]DataBarras, error) {

	if inicio == "" || fin == "" {
		return nil, errors.New("las fechas estan vacias")
	}

	var responses []DataBarras

	responses, err := s.repository.GetDataChartBarra(inicio, fin)

	if err != nil {
		return nil, fmt.Errorf("error al extraer datos para el chart barra: %w", err)
	}

	return responses, nil
}
