package gpdfs

import (
	"dbu-api/internal/logger"
	"errors"
	"time"
)

type PortsServerPDFs interface {
	GeneratePDF_SRQ(participantID int) (SRQPDFData, error)
}

type service struct {
	repository ServicesGpdfsRepository
	txID       string
}

func NewPDFService(repository ServicesGpdfsRepository, txID string) PortsServerPDFs {
	return &service{repository: repository, txID: txID}
}

var (
	ErrInvalidParticipantID    = errors.New("invalid participant ID")
	ErrFetchingParticipantData = errors.New("error fetching participant data")
	ErrFetchingSurveyResponses = errors.New("error fetching survey responses")
	ErrGeneratingPDF           = errors.New("error generating PDF")
)

func (s *service) GeneratePDF_SRQ(participantID int) (SRQPDFData, error) {
	if participantID == 0 {
		return SRQPDFData{}, errors.New("invalid participant ID")
	}

	participant, err := s.repository.GetParticipantByID(participantID)
	if err != nil {
		logger.Error.Printf("error fetching participant data: %v", err)
		return SRQPDFData{}, errors.New("error fetching participant data")
	}

	responses, err := s.repository.GetResponsesPerParticipant(participantID, 1)
	if err != nil {
		logger.Error.Printf("error fetching survey responses: %v", err)
		return SRQPDFData{}, errors.New("error fetching survey responses")
	}

	var questions []any
	for i, r := range responses {
		questions = append(questions, map[string]interface{}{
			"Number": i + 1,
			"Text":   r.TextoPregunta,
			"Answer": r.Respuesta,
		})
	}

	data := SRQPDFData{
		Participant: *participant,
		Questions:   questions,
		GeneratedAt: time.Now(),
	}
	return data, nil
}
