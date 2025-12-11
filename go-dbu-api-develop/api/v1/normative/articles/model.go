package articles

import "github.com/asaskevich/govalidator"

type ArticleRequest struct {
	ID          string `json:"id" valid:"-"`
	Description string `json:"descripcion" valid:"required"`
	Severity    string `json:"gravedad" valid:"in(leve1|leve2|grave),required"`
}

func (m *ArticleRequest) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
