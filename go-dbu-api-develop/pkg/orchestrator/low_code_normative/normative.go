package low_code_normative

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"dbu-api/pkg/normative"
	"errors"

	"github.com/jmoiron/sqlx"
)

// Errores y Códigos
// 101 // Datos requeridos faltantes
// 201 // Error al insertar en la base de datos
// 301 // Elemento no encontrado
// 20  // Operación exitosa

type ResolutionService struct {
	db   *sqlx.DB
	usr  *models.User
	txID string
}

type PortsServerResolution interface {
	CreateResolutionLowCode(m *models.Normative) (int, error)
	GetAllResolutionsLowCode() ([]*models.Normative, int, error)
	CreateChapterLowCode(m *models.Chapter) (int, error)
	GetChaptersByResolutionLowCode(resolutionID string) ([]*models.Chapter, int, error)
	GetArticlesByChapterLowCode(chapterID int64) ([]*models.Article, int, error)
	CreateArticleLowCode(m *models.Article) (int, error)
}

func NewResolution(db *sqlx.DB, usr *models.User, txID string) PortsServerResolution {
	return &ResolutionService{db: db, usr: usr, txID: txID}
}

// Crear una nueva Resolución
func (s *ResolutionService) CreateResolutionLowCode(m *models.Normative) (int, error) {
	if m.Nombre == "" || m.Descripcion == "" || m.ServicioId == "" || m.RutaArchivo == "" {
		logger.Error.Println(s.txID, " - Missing required fields in resolution")
		return 101, errors.New("missing required fields: name, description, service, or document path")
	}

	valid, err := m.Valid()
	if err != nil || !valid {
		logger.Error.Println(s.txID, " - Invalid resolution data:", err)
		return 102, errors.New("invalid resolution data")
	}

	srvService := normative.NewServerResolution(s.db, s.usr, s.txID)
	dataResolution, _, err := srvService.SrvResolution.CreateResolution(
		m.ID,
		m.Nombre,
		m.Descripcion,
		m.Estado,
		m.ServicioId,
		m.RutaArchivo,
	)

	if err != nil {
		logger.Error.Println(s.txID, " - Couldn't create resolution:", err)
		return 201, err
	}

	if dataResolution == nil {
		logger.Error.Println(s.txID, " - Resolution not found")
		return 301, errors.New("resolution not found")
	}

	return 20, nil
}

// Crear un nuevo Capítulo
func (s *ResolutionService) CreateChapterLowCode(m *models.Chapter) (int, error) {
	if m.Nombre == "" || m.NormativeId == "" || m.Descripcion == "" {
		logger.Error.Println(s.txID, " - Missing required fields in chapter")
		return 101, errors.New("missing required fields: title or resolution ID")
	}

	srvService := normative.NewServerResolution(s.db, s.usr, s.txID)
	dataChapter, _, err := srvService.SrvChapter.CreateChapter(
		m.ID,
		m.Nombre,
		m.Descripcion,
		m.NormativeId,
	)

	if err != nil {
		logger.Error.Println(s.txID, " - Couldn't create chapter:", err)
		return 201, err
	}

	if dataChapter == nil {
		logger.Error.Println(s.txID, " - Chapter not found")
		return 301, errors.New("chapter not found")
	}

	return 20, nil
}

// Crear un nuevo Artículo
func (s *ResolutionService) CreateArticleLowCode(m *models.Article) (int, error) {
	if m.Gravedad == "" || m.CapituloId == "" || m.Descripcion == "" {
		logger.Error.Println(s.txID, " - Missing required fields in article")
		return 101, errors.New("missing required fields: title, chapter ID or content")
	}

	srvService := normative.NewServerResolution(s.db, s.usr, s.txID)
	dataArticle, _, err := srvService.SrvArticle.CreateArticle(
		m.ID,
		m.Descripcion,
		m.Gravedad,
		m.CapituloId,
	)
	if err != nil {
		logger.Error.Println(s.txID, " - Couldn't create article:", err)
		return 201, err
	}

	if dataArticle == nil {
		logger.Error.Println(s.txID, " - Article not found")
		return 301, errors.New("article not found")
	}

	return 20, nil
}

// Obtener todas las resoluciones
func (s *ResolutionService) GetAllResolutionsLowCode() ([]*models.Normative, int, error) {
	srvService := normative.NewServerResolution(s.db, s.usr, s.txID)
	resolutions, err := srvService.SrvResolution.GetAllResolution()
	if err != nil {
		logger.Error.Println(s.txID, " - Couldn't fetch resolutions:", err)
		return nil, 201, err
	}

	if resolutions == nil || len(resolutions) == 0 {
		logger.Error.Println(s.txID, " - No resolutions found")
		return []*models.Normative{}, 301, nil
	}

	// Convertir de []*resolutions.Resolution a []*models.Normative
	normatives := make([]*models.Normative, len(resolutions))
	for i, res := range resolutions {
		normatives[i] = &models.Normative{
			ID:          res.ID,
			Nombre:      res.Nombre,
			Descripcion: res.Descripcion,
			ServicioId:  res.Servicio_id,
			RutaArchivo: res.Ruta_archivo,
			Estado:      res.Estado,
		}
	}

	return normatives, 20, nil
}

// Obtener capítulos por resolución
func (s *ResolutionService) GetChaptersByResolutionLowCode(resolutionID string) ([]*models.Chapter, int, error) {
	srvService := normative.NewServerResolution(s.db, s.usr, s.txID)
	chaptersData, _, err := srvService.SrvChapter.GetChaptersByResolution(resolutionID)
	if err != nil {
		logger.Error.Println(s.txID, " - Couldn't fetch chapters:", err)
		return nil, 201, err
	}

	if chaptersData == nil || len(chaptersData) == 0 {
		logger.Error.Println(s.txID, " - No chapters found")
		return []*models.Chapter{}, 301, nil
	}

	// Convertir de []*chapters.Chapter a []*models.Chapter
	modelChapters := make([]*models.Chapter, len(chaptersData))
	for i, ch := range chaptersData {
		modelChapters[i] = &models.Chapter{
			ID:          ch.ID,
			Nombre:      ch.Nombre,
			Descripcion: ch.Descripcion,
			NormativeId: ch.Resolucion_id,
		}
	}

	return modelChapters, 20, nil
}

// Obtener artículos por capítulo
func (s *ResolutionService) GetArticlesByChapterLowCode(chapterID int64) ([]*models.Article, int, error) {
	srvService := normative.NewServerResolution(s.db, s.usr, s.txID)
	articlesData, _, err := srvService.SrvArticle.GetArticlesByChapter(chapterID)
	if err != nil {
		logger.Error.Println(s.txID, " - Couldn't fetch articles:", err)
		return nil, 201, err
	}

	if articlesData == nil || len(articlesData) == 0 {
		logger.Error.Println(s.txID, " - No articles found")
		return []*models.Article{}, 301, nil
	}

	// Convertir de []*articles.Article a []*models.Article
	modelArticles := make([]*models.Article, len(articlesData))
	for i, art := range articlesData {
		modelArticles[i] = &models.Article{
			ID:          art.ID,
			Descripcion: art.Descripcion,
			CapituloId:  art.Capitulo_id,
		}
	}

	return modelArticles, 20, nil
}
