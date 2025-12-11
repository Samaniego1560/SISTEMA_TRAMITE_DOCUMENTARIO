package sanction

import (
	"dbu-api/internal/models"
	"dbu-api/pkg/residence/residence_configuration"
	"dbu-api/pkg/residence/residences"
	"dbu-api/pkg/residence/rooms"
	"dbu-api/pkg/sanction/fault"

	"github.com/jmoiron/sqlx"
)

type ServerSanction struct {
	SrvResidence              residences.PortsServerResidence
	SrvRoom                   rooms.PortsServerRoom
	SrvResidenceConfiguration residence_configuration.PortsServerResidenceConfiguration
	SrvFault                  fault.PortsServerFault
}

func NewServerSanction(db *sqlx.DB, usr *models.User, txID string) *ServerSanction {
	return &ServerSanction{
		SrvResidence:              residences.NewResidenceService(residences.FactoryStorage(db, txID), usr, txID),
		SrvRoom:                   rooms.NewRoomService(rooms.FactoryStorage(db, usr, txID), usr, txID),
		SrvResidenceConfiguration: residence_configuration.NewResidenceConfigurationService(residence_configuration.FactoryStorage(db, txID), usr, txID),
		SrvFault:                  fault.NewFaultService(fault.FactoryStorage(db, txID), usr, txID),
	}
}
