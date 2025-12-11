package room_assignment

import "github.com/asaskevich/govalidator"

type ResponseResidence struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
	Code  int    `json:"code"`
	Type  string `json:"type"`
}

type RequestRoomAssignment struct {
	ResidenceId  string `json:"residence_id" valid:"required,uuid"`
	SubmissionId int64  `json:"submission_id,omitempty" valid:"optional"`
}

type ResponseAutomaticAssignationResidence struct {
	RegisterID string `json:"register_id"`
}

func (m *RequestRoomAssignment) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
