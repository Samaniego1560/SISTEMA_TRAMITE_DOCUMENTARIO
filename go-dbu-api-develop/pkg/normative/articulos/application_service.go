package articles

import (
	"database/sql"
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"fmt"

	"github.com/google/uuid"
)

type PortsServerArticle interface {
	CreateArticle(id string, descripcion string, gravedad string, capituloId string) (*Article, int, error)
	UpdateArticle(id string, descripcion string, gravedad string, capituloId string) (*Article, int, error)
	DeleteArticle(id string) (int, error)
	GetArticleByID(id string) (*Article, int, error)
	GetAllArticle() ([]*Article, error)
	GetArticlesByChapter(CapituloId int64) ([]*Article, int, error)
}

type service struct {
	repository ServiceArticleRepository
	user       *models.User
	txID       string
}

func NewArticleService(repository ServiceArticleRepository, user *models.User, TxID string) PortsServerArticle {
	return &service{repository: repository, user: user, txID: TxID}
}
func (s *service) CreateArticle(id string, descripcion string, gravedad string, capituloId string) (*Article, int, error) {
	m := NewArticle(id, descripcion, gravedad, capituloId)
	valid, err := m.valid()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't meet validations:", err)
		return nil, 15, err
	}
	if !valid {
		logger.Error.Println(s.txID, " - don't meet validations:")
		return nil, 15, err
	}
	if err := s.repository.create(m); err != nil {
		if err.Error() == "rows affected error" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Article :", err)
		return nil, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateArticle(id string, gravedad string, descripcion string, CapituloId string) (*Article, int, error) {
	m := NewArticle(id, descripcion, gravedad, CapituloId)
	valid, err := m.valid()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn't meet validations:", err)
		return nil, 15, err
	}
	if !valid {
		logger.Error.Println(s.txID, " - don't meet validations:")
		return nil, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Article :", err)
		return nil, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteArticle(id string) (int, error) {
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "rows affected error" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetArticleByID(id string) (*Article, int, error) {
	if err := uuid.Validate(id); err != nil {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllArticle() ([]*Article, error) {
	return s.repository.getAll()
}

func (s *service) GetArticlesByChapter(ChapterId int64) ([]*Article, int, error) {
	articles, err := s.repository.GetByChapterID(ChapterId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			logger.Error.Println(s.txID, " - no rows found for Chapter: ", ChapterId)
			return nil, 14, err
		default:
			logger.Error.Println(s.txID, " - couldn't get articles: ", err)
			return nil, 13, err
		}
	}

	return articles, 29, nil
}

/*func (s *service) UpdateOnlyCharacteristicsArticle(id string, inciso string, descripcion string, graedad string) (*Article, int, error) {
	m := NewArticle(id, inciso, descripcion, graedad, s.user.ID)
	if err := s.repository.updateOnlyCharacteristics(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update article :", err)
		return m, 18, err
	}
	return m, 29, nil
}*/
