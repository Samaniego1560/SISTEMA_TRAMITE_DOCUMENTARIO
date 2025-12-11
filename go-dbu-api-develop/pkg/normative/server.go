package normative

import (
	"dbu-api/internal/models"
	articles "dbu-api/pkg/normative/articulos"
	chapters "dbu-api/pkg/normative/capitulos"
	incisos "dbu-api/pkg/normative/incisos"
	resolutions "dbu-api/pkg/normative/resoluciones"

	"github.com/jmoiron/sqlx"
)

type ServerResolution struct {
	SrvResolution resolutions.PortsServerResolution
	SrvArticle    articles.PortsServerArticle
	SrvChapter    chapters.PortsServerChapter
	SrvInciso     incisos.PortsServerInciso
}

func NewServerResolution(db *sqlx.DB, usr *models.User, txID string) *ServerResolution {
	return &ServerResolution{
		SrvResolution: resolutions.NewResolutionService(resolutions.FactoryStorage(db, txID), usr, txID),
		SrvArticle:    articles.NewArticleService(articles.FactoryStorage(db, txID), usr, txID),
		SrvChapter:    chapters.NewChapterService(chapters.FactoryStorage(db, txID), usr, txID),
		SrvInciso:     incisos.NewIncisoService(incisos.FactoryStorage(db, txID), usr, txID),
	}
}
