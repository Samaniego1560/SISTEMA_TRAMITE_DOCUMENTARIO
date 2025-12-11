package sanctioned

import "github.com/asaskevich/govalidator"

type SanctionedRequest struct {
	ID       string `json:"id" valid:"uuid,required"`
	Capacity int    `json:"capacity" valid:"numeric,required"`
	Status   string `json:"status" valid:"in(mantenimiento|deshabilitado|habilitado),required"`
}

func (m *SanctionedRequest) Valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
