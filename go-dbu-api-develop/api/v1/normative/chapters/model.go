package chapters

import "github.com/asaskevich/govalidator"

type ChapterRequest struct {
	ID           string `json:"id" valid:"uuid,required"`
	name         string `json:"nombre" valid:"required"`
	Description  string `json:"descripcion" valid: "required"`
	ResolucionId string `json:"resolucion_id" valid: "required"`
}

func (m *ChapterRequest) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
