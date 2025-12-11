package low_code_medical_area

import (
	"dbu-api/internal/logger"
	"dbu-api/internal/models"
	"dbu-api/pkg/medical_area"
	"dbu-api/pkg/medical_area/consultation_medical_area"
	"dbu-api/pkg/medical_area/nursing_consultation_routine_review"
	"dbu-api/pkg/medical_area/nursing_consultation_vaccine"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sort"
	"strings"
	"time"
)

// Insertar una consulta de enfermería en la base de datos

type NursingConsultationService struct {
	db   *sqlx.DB
	usr  *models.User
	txID string
}

type PortsServerNursingConsultation interface {
	CreateNursingConsultationLowCode(m *models.RequestNursingConsultation, isAssigned bool) (int, error)
	UpdateNursingConsultationLowCode(m *models.RequestNursingConsultation) (int, error)
	DeleteNursingConsultationLowCode(id string, isAssigned bool) (int, error)
	GetNursingConsultationByIdConsultationLowCode(id string) (*models.RequestNursingConsultation, int, error)
	GetNursingConsultationByIdPatientLowCode(id string) ([]*models.RequestNursingConsultation, int, error)
	GetNursingConsultationByDNILowCode(dni string) (*models.GetAllConsultationMedicalArea, int, error)
	GetAllNursingConsultationLowCode() ([]*models.GetAllConsultationMedicalArea, error)
	GetAllTypesVaccines() ([]*models.VaccineType, error)
	GetVaccinesRequired(patientId string) (*models.NextVaccinePatient, int, error)
}

func NewNursingConsultation(db *sqlx.DB, usr *models.User, txID string) PortsServerNursingConsultation {
	return &NursingConsultationService{db: db, usr: usr, txID: txID}
}

const (
	StatusSuccess  = 20
	StatusNotFound = 301
)

func (s *NursingConsultationService) CreateNursingConsultationLowCode(m *models.RequestNursingConsultation, isAssigned bool) (int, error) {

	srv := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	dataNursingConsultation, code, err := srv.SrvConsultationMedicalArea.CreateConsultationMedicalArea(m.ConsultaAreaMedica.ID, m.ConsultaAreaMedica.IDPaciente, m.ConsultaAreaMedica.FechaConsulta, "enfermería")
	if err != nil {
		logger.Error.Printf("couldn't create nursing consultation, error: %v", dataNursingConsultation, err)
		return code, err
	}

	if !isAssigned {
		dataConsultationAssignment, code, err := srv.SrvConsultationAssignment.CreateConsultationAssignment("", m.ConsultaAreaMedica.ID, "enfermería")
		if err != nil {
			logger.Error.Printf("couldn't create consultation assignment, error: %v", dataConsultationAssignment, err)
			return code, err
		}
	}
	dataRoutineReview, code, err := srv.SrvNursingConsultationRoutineReview.CreateRoutineReview(m.RevisionRutina.ID, m.ConsultaAreaMedica.ID, m.RevisionRutina.FiebreUltimoQuinceDias, m.RevisionRutina.TosMasQuinceDias, m.RevisionRutina.SecrecionLesionGenitales, m.RevisionRutina.FechaUltimaRegla, m.RevisionRutina.Comentarios)
	if err != nil {
		logger.Error.Printf("couldn't create routine review, error: %v", dataRoutineReview, err)
		return code, err
	}

	if m.DatosAcompanante != nil {
		dataAccompanyingData, code, err := srv.SrvNursingConsultationAccompanyingData.CreateAccompanyingData(m.DatosAcompanante.ID, m.ConsultaAreaMedica.ID, m.DatosAcompanante.DNI, m.DatosAcompanante.NombresApellidos, m.DatosAcompanante.Edad)
		if err != nil {
			logger.Error.Printf("couldn't create accompanying data, error: %v", dataAccompanyingData, err)
			return code, err
		}
	}

	if m.Examenes.Vacunas != nil {
		for _, vacuna := range m.Examenes.Vacunas {
			dataVaccine, code, err := srv.SrvNursingConsultationVaccine.CreateVaccine(vacuna.ID, m.ConsultaAreaMedica.ID, vacuna.TipoVacuna, vacuna.FechaDosis, vacuna.Minsa, vacuna.TipoAtencion, vacuna.Indicaciones, vacuna.Observaciones)
			if err != nil {
				logger.Error.Printf("couldn't create vaccine, error: %v", dataVaccine, err)
				return code, err
			}
		}
	}

	if m.Examenes.ExamenFisico != nil {
		dataPhysicalTest, code, err := srv.SrvNursingConsultationPhysicalTest.CreatePhysicalTest(m.Examenes.ExamenFisico.ID, m.ConsultaAreaMedica.ID, m.Examenes.ExamenFisico.TallaPesos, m.Examenes.ExamenFisico.PerimetroCintura, m.Examenes.ExamenFisico.IndiceMasaCorporalImg, m.Examenes.ExamenFisico.PresionArterial, m.Examenes.ExamenFisico.Comentarios)
		if err != nil {
			logger.Error.Printf("couldn't create physical test, error: %v", dataPhysicalTest, err)
			return code, err

		}
	}

	if m.Examenes.ExamenLaboratorio != nil {
		dataLaboratoryTest, code, err := srv.SrvNursingConsultationLaboratoryTest.CreateLaboratoryTest(m.Examenes.ExamenLaboratorio.ID, m.ConsultaAreaMedica.ID, m.Examenes.ExamenLaboratorio.Serologia, m.Examenes.ExamenLaboratorio.Bk, m.Examenes.ExamenLaboratorio.Hemograma, m.Examenes.ExamenLaboratorio.ExamenOrina, m.Examenes.ExamenLaboratorio.Colesterol, m.Examenes.ExamenLaboratorio.Glucosa, m.Examenes.ExamenLaboratorio.Comentarios)
		if err != nil {
			logger.Error.Printf("couldn't create laboratory test, error: %v", dataLaboratoryTest, err)
			return code, err
		}
	}

	if m.Examenes.ExamenPreferencial != nil {
		dataPreferentialTest, code, err := srv.SrvNursingConsultationPreferentialTest.CreatePreferentialTest(m.Examenes.ExamenPreferencial.ID, m.ConsultaAreaMedica.ID, m.Examenes.ExamenPreferencial.AparatoRespiratorio, m.Examenes.ExamenPreferencial.AparatoCardiovascular, m.Examenes.ExamenPreferencial.AparatoDigestivo, m.Examenes.ExamenPreferencial.AparatoGenitourinario, m.Examenes.ExamenPreferencial.Papanicolau, m.Examenes.ExamenPreferencial.ExamenMama, m.Examenes.ExamenPreferencial.Comentarios)
		if err != nil {
			logger.Error.Printf("couldn't create preferential test, error: %v", dataPreferentialTest, err)
			return code, err
		}
	}

	if m.Examenes.ExamenSexualidad != nil {
		dataSexualityTest, code, err := srv.SrvNursingConsultationSexualityTest.CreateSexualityTest(m.Examenes.ExamenSexualidad.ID, m.ConsultaAreaMedica.ID, m.Examenes.ExamenSexualidad.ActividadSexual, m.Examenes.ExamenSexualidad.PlanificacionFamiliar, m.Examenes.ExamenSexualidad.Comentarios)
		if err != nil {
			logger.Error.Printf("couldn't create sexuality test, error: %v", dataSexualityTest, err)
			return code, err
		}
	}

	if m.Examenes.ExamenVisual != nil {
		dataVisualTest, code, err := srv.SrvNursingConsultationVisualTest.CreateVisualTest(m.Examenes.ExamenVisual.ID, m.ConsultaAreaMedica.ID, m.Examenes.ExamenVisual.OjoDerecho, m.Examenes.ExamenVisual.OjoIzquierdo, m.Examenes.ExamenVisual.PresionOcular, m.Examenes.ExamenVisual.Comentarios)
		if err != nil {
			logger.Error.Printf("couldn't create visual test, error: %v", dataVisualTest, err)
			return code, err
		}
	}

	if m.Examenes.ProcedimientoRealizado != nil {
		dataPerformedProcedures, code, err := srv.SrvNursingConsultationPerformedProcedures.CreatePerformedProcedures(
			m.Examenes.ProcedimientoRealizado.ID, m.ConsultaAreaMedica.ID, m.Examenes.ProcedimientoRealizado.Procedimiento,
			m.Examenes.ProcedimientoRealizado.NumeroRecibo, m.Examenes.ProcedimientoRealizado.Comentarios,
			m.Examenes.ProcedimientoRealizado.Costo, m.Examenes.ProcedimientoRealizado.FechaPago,
			m.Examenes.ProcedimientoRealizado.AreaSolicitante, m.Examenes.ProcedimientoRealizado.EspecialistaSolicitante,
		)
		if err != nil {
			logger.Error.Printf("couldn't create performed procedures, error: %v", dataPerformedProcedures, err)
			return code, err
		}
	}

	if m.Examenes.TratamientoMedicamentoso != nil {
		dataMedicationTreatment, code, err := srv.SrvNursingConsultationMedicationTreatment.CreateMedicationTreatment(
			m.Examenes.TratamientoMedicamentoso.ID, m.ConsultaAreaMedica.ID, m.Examenes.TratamientoMedicamentoso.NombreGenericoMedicamento,
			m.Examenes.TratamientoMedicamentoso.ViaAdministracion, m.Examenes.TratamientoMedicamentoso.HoraAplicacion,
			m.Examenes.TratamientoMedicamentoso.ResponsableAtencion, m.Examenes.TratamientoMedicamentoso.Observaciones,
			m.Examenes.TratamientoMedicamentoso.AreaSolicitante, m.Examenes.TratamientoMedicamentoso.EspecialistaSolicitante,
		)
		if err != nil {
			logger.Error.Printf("couldn't create medication treatment, error: %v", dataMedicationTreatment, err)
			return code, err
		}
	}

	if m.Examenes.AtencionIntegralOtros != nil {
		dataIntegralAttentionOther, code, err := srv.SrvConsultationIntegralAttention.CreateConsultationIntegralAttention(m.Examenes.AtencionIntegralOtros.ID, m.ConsultaAreaMedica.ID, m.Examenes.AtencionIntegralOtros.Fecha, m.Examenes.AtencionIntegralOtros.Hora, m.Examenes.AtencionIntegralOtros.Edad, m.Examenes.AtencionIntegralOtros.MotivoConsulta, m.Examenes.AtencionIntegralOtros.TiempoEnfermedad, m.Examenes.AtencionIntegralOtros.Apetito, m.Examenes.AtencionIntegralOtros.Sed, m.Examenes.AtencionIntegralOtros.Suenio, m.Examenes.AtencionIntegralOtros.EstadoAnimo, m.Examenes.AtencionIntegralOtros.Orina, m.Examenes.AtencionIntegralOtros.Deposiciones, m.Examenes.AtencionIntegralOtros.Temperatura, m.Examenes.AtencionIntegralOtros.PresionArterial, m.Examenes.AtencionIntegralOtros.FrecuenciaCardiaca, m.Examenes.AtencionIntegralOtros.FrecuenciaRespiratoria, m.Examenes.AtencionIntegralOtros.Peso, m.Examenes.AtencionIntegralOtros.Talla, m.Examenes.AtencionIntegralOtros.IndiceMasaCorporal, m.Examenes.AtencionIntegralOtros.Diagnostico, m.Examenes.AtencionIntegralOtros.Tratamiento, m.Examenes.AtencionIntegralOtros.ExamenesAxuliares, m.Examenes.AtencionIntegralOtros.Referencia, m.Examenes.AtencionIntegralOtros.Observacion, m.Examenes.AtencionIntegralOtros.NumeroRecibo, m.Examenes.AtencionIntegralOtros.Costo, m.Examenes.AtencionIntegralOtros.FechaPago)
		if err != nil {
			logger.Error.Printf("couldn't create consultation integral attention, error: %v", dataIntegralAttentionOther, err)
			return code, err
		}

		s.addAnnouncement(srv, m.ConsultaAreaMedica.IDPaciente, m.ConsultaAreaMedica.ID, "enfermería")
	}

	if m.Examenes.ConsultaBucal != nil {
		dataBuccalConsultation, code, err := srv.SrvDentistryConsultationBuccalConsultation.CreateBuccalConsultation(m.Examenes.ConsultaBucal.ID, m.ConsultaAreaMedica.ID, m.Examenes.ConsultaBucal.Relato, m.Examenes.ConsultaBucal.Diagnostico, m.Examenes.ConsultaBucal.ExamenAuxiliar, m.Examenes.ConsultaBucal.ExamenClinico, m.Examenes.ConsultaBucal.Tratamiento, m.Examenes.ConsultaBucal.Indicaciones, m.Examenes.ConsultaBucal.Comentarios)
		if err != nil {
			logger.Error.Printf("couldn't create buccal consultation, error: %v", dataBuccalConsultation, err)
			return code, err
		}
	}

	if m.Examenes.ConsultaMedicinaGeneral != nil {
		dataGeneralMedicineConsultation, code, err := srv.SrvMedicalGeneralMedicineConsultation.CreateGeneralMedicineConsultation(m.Examenes.ConsultaMedicinaGeneral.ID, m.ConsultaAreaMedica.ID, m.Examenes.ConsultaMedicinaGeneral.FechaHora, m.Examenes.ConsultaMedicinaGeneral.Anamnesis, m.Examenes.ConsultaMedicinaGeneral.ExamenClinico, m.Examenes.ConsultaMedicinaGeneral.Indicaciones)
		if err != nil {
			logger.Error.Printf("couldn't create general medicine consultation , error: %v", dataGeneralMedicineConsultation, err)
			return code, err
		}
		s.addAnnouncement(srv, m.ConsultaAreaMedica.IDPaciente, m.ConsultaAreaMedica.ID, "medicina")
	}

	return StatusSuccess, nil
}

func (s *NursingConsultationService) UpdateNursingConsultationLowCode(m *models.RequestNursingConsultation) (int, error) {

	code, err := s.DeleteNursingConsultationLowCode(m.ConsultaAreaMedica.ID, true)

	if err != nil {
		logger.Error.Printf("couldn't delete nursing consultation, error: %v", err)
		return code, err
	}

	code, err = s.CreateNursingConsultationLowCode(m, true)

	if err != nil {
		logger.Error.Printf("couldn't create nursing consultation, error: %v", err)
		return code, err
	}

	return StatusSuccess, nil
}

func (s *NursingConsultationService) DeleteNursingConsultationLowCode(id string, isAssigned bool) (int, error) {

	srv := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	codeRoutineReview, err := srv.SrvNursingConsultationRoutineReview.DeleteRoutineReviewByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete routine review, error: %v", err)
		return codeRoutineReview, err
	}

	codeAccompanyingData, err := srv.SrvNursingConsultationAccompanyingData.DeleteAccompanyingDataByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete accompanying data, error: %v", err)
		return codeAccompanyingData, err
	}

	codeVaccine, err := srv.SrvNursingConsultationVaccine.DeleteVaccineByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete vaccine, error: %v", err)
		return codeVaccine, err
	}

	codePhysicalTest, err := srv.SrvNursingConsultationPhysicalTest.DeletePhysicalTestByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete physical test, error: %v", err)
		return codePhysicalTest, err
	}

	codeLaboratoryTest, err := srv.SrvNursingConsultationLaboratoryTest.DeleteLaboratoryTestByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete laboratory test, error: %v", err)
		return codeLaboratoryTest, err
	}

	codePreferentialTest, err := srv.SrvNursingConsultationPreferentialTest.DeletePreferentialTestByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete preferential test, error: %v", err)
		return codePreferentialTest, err
	}

	codeSexualityTest, err := srv.SrvNursingConsultationSexualityTest.DeleteSexualityTestByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete sexuality test, error: %v", err)
		return codeSexualityTest, err
	}

	codeVisualTest, err := srv.SrvNursingConsultationVisualTest.DeleteVisualTestByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete visual test, error: %v", err)
		return codeVisualTest, err
	}

	codeMedicationTreatment, err := srv.SrvNursingConsultationMedicationTreatment.DeleteMedicationTreatmentByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't deletem medication treatment, error: %v", err)
		return codeMedicationTreatment, err
	}

	codePerformedProcedures, err := srv.SrvNursingConsultationPerformedProcedures.DeletePerformedProceduresByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete performed procedures, error: %v", err)
		return codePerformedProcedures, err
	}

	if !isAssigned {
		codeConsultationAssignment, err := srv.SrvConsultationAssignment.DeleteConsultationAssignmentByIDConsultation(id)
		if err != nil {
			logger.Error.Printf("couldn't delete consultation assignment, error: %v", err)
			return codeConsultationAssignment, err
		}
	}

	codeConsultationIntegralAttention, err := srv.SrvConsultationIntegralAttention.DeleteConsultationIntegralAttentionByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't delete integral attention, error: %v", err)
		return codeConsultationIntegralAttention, err
	}

	codeBuccalConsultation, err := srv.SrvDentistryConsultationBuccalConsultation.DeleteBuccalConsultationByIDConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't create buccal consultation, error: %v", err)
		return codeBuccalConsultation, err
	}

	codeGeneralMedicineConsultation, err := srv.SrvMedicalGeneralMedicineConsultation.DeleteGeneralMedicineConsultation(id)
	if err != nil {
		logger.Error.Printf("couldn't create general medicine consultation , error: %v", err)
		return codeGeneralMedicineConsultation, err
	}

	codeNursingConsultation, err := srv.SrvConsultationMedicalArea.DeleteConsultationMedicalArea(id)
	if err != nil {
		logger.Error.Printf("couldn't delete nursing consultation, error: %v", err)
		return codeNursingConsultation, err
	}

	return StatusSuccess, nil
}

func (s *NursingConsultationService) GetNursingConsultationByIdConsultationLowCode(id string) (*models.RequestNursingConsultation, int, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	nursingConsultation, code, err := srvService.SrvConsultationMedicalArea.GetConsultationMedicalAreaByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get Nursing Consultation by ID:", err)
		return nil, code, fmt.Errorf("failed to get Nursing Consultation: %w", err)
	}
	if nursingConsultation == nil {
		return nil, StatusNotFound, nil
	}

	routineReview, err := fetchRoutineReview(srvService, nursingConsultation.ID, s.txID)
	if err != nil {
		return nil, code, err
	}

	accompanyingData, err := fetchAccompanyingData(srvService, nursingConsultation.ID, s.txID)
	if err != nil {
		return nil, code, err
	}

	exams, code, err := s.fetchConsultationData(srvService, id)
	if err != nil {
		return nil, code, err
	}

	data := &models.RequestNursingConsultation{
		ConsultaAreaMedica: mapNursingConsultation(nursingConsultation),
		RevisionRutina:     routineReview,
		DatosAcompanante:   accompanyingData,
		Examenes:           exams,
	}

	return data, StatusSuccess, nil
}

func (s *NursingConsultationService) GetNursingConsultationByIdPatientLowCode(id string) ([]*models.RequestNursingConsultation, int, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	nursingConsultations, code, err := srvService.SrvConsultationMedicalArea.GetConsultationMedicalAreasByPatientID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get Nursing Consultations by Patient ID:", err)
		return nil, code, fmt.Errorf("failed to get Nursing Consultations: %w", err)
	}
	if nursingConsultations == nil {
		return nil, StatusNotFound, nil
	}

	var result []*models.RequestNursingConsultation

	for _, nursingConsultation := range nursingConsultations {

		routineReview, err := fetchRoutineReview(srvService, nursingConsultation.ID, s.txID)
		if err != nil {
			return nil, code, err
		}

		accompanyingData, err := fetchAccompanyingData(srvService, nursingConsultation.ID, s.txID)
		if err != nil {
			return nil, code, err
		}

		exams, code, err := s.fetchConsultationData(srvService, nursingConsultation.ID)
		if err != nil {
			return nil, code, err
		}

		data := &models.RequestNursingConsultation{
			ConsultaAreaMedica: mapNursingConsultation(nursingConsultation),
			RevisionRutina:     routineReview,
			DatosAcompanante:   accompanyingData,
			Examenes:           exams,
		}

		result = append(result, data)
	}

	return result, StatusSuccess, nil
}

func (s *NursingConsultationService) GetAllNursingConsultationLowCode() ([]*models.GetAllConsultationMedicalArea, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	consultations, err := srvService.SrvConsultationMedicalArea.GetAllConsultationMedicalArea()

	if err != nil {
		logger.Error.Println(s.txID, " - failed to get all Nursing Consultations:", err)
		return nil, fmt.Errorf("failed to get all Nursing Consultations: %w", err)
	}
	if consultations == nil {
		return nil, nil
	}

	//nursingConsultations, _, err := srvService.SrvConsultationAssignment.GetConsultationAssignmentByArea("enfermería")
	//
	//if err != nil {
	//	logger.Error.Println(s.txID, " - failed to get all Nursing Consultations:", err)
	//	return nil, fmt.Errorf("failed to get all Nursing Consultations: %w", err)
	//}
	//if nursingConsultations == nil {
	//	return nil, nil
	//}

	var result []*models.GetAllConsultationMedicalArea

	for _, consultation := range consultations {

		patient, err := fetchResposePatientInfo(srvService, consultation.IDPaciente, s.txID)
		if err != nil {
			return nil, err
		}

		ConsultationAssignment, err := fetchConsultationAssignment(srvService, consultation.ID, s.txID)
		if err != nil {
			return nil, err
		}

		services, _, err := s.fetchServicesNursing(srvService, consultation.ID)
		if err != nil {
			return nil, err
		}

		consultationResponse := models.ResponseConsultationMedicalArea{
			ID:            consultation.ID,
			FechaConsulta: consultation.FechaConsulta,
			AreaAsignada:  ConsultationAssignment.AreaAsignada,
			AreaOrigen:    consultation.AreaMedica,
			Servicios:     services,
		}
		data := &models.GetAllConsultationMedicalArea{
			Paciente:  patient,
			Consultas: []models.ResponseConsultationMedicalArea{consultationResponse},
		}

		result = append(result, data)
	}

	return result, nil
}

func (s *NursingConsultationService) GetNursingConsultationByDNILowCode(dni string) (*models.GetAllConsultationMedicalArea, int, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	consultations, code, err := srvService.SrvConsultationMedicalArea.GetConsultationMedicalAreaByPatientDNI(dni)

	if err != nil {
		logger.Error.Println(s.txID, " - failed to get all Nursing Consultations:", err)
		return nil, code, fmt.Errorf("failed to get all Nursing Consultations: %w", err)
	}
	if consultations == nil {
		return nil, code, nil
	}

	var result []models.ResponseConsultationMedicalArea

	patient, err := fetchResposePatientInfo(srvService, dni, s.txID)
	if err != nil {
		return nil, code, err
	}

	for _, consultation := range consultations {

		ConsultationAssignment, err := fetchConsultationAssignment(srvService, consultation.ID, s.txID)
		if err != nil {
			return nil, code, err
		}

		services, _, err := s.fetchServicesNursing(srvService, consultation.ID)
		if err != nil {
			return nil, code, err
		}

		consultationResponse := models.ResponseConsultationMedicalArea{
			ID:            consultation.ID,
			FechaConsulta: consultation.FechaConsulta,
			AreaAsignada:  ConsultationAssignment.AreaAsignada,
			AreaOrigen:    consultation.AreaMedica,
			Servicios:     services,
		}

		result = append(result, consultationResponse)

	}

	data := &models.GetAllConsultationMedicalArea{
		Paciente:  patient,
		Consultas: result,
	}

	return data, code, nil
}

func (s *NursingConsultationService) GetAllTypesVaccines() ([]*models.VaccineType, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	typesVaccines, err := srvService.SrvNursingConsultationVaccine.GetAllTypesVaccines()
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get all Types Vaccines:", err)
		return nil, fmt.Errorf("failed to get all Types Vaccines: %w", err)
	}
	if typesVaccines == nil {
		return nil, nil
	}

	var result []*models.VaccineType

	for _, typeVaccine := range typesVaccines {

		data := &models.VaccineType{
			ID:                  typeVaccine.ID,
			Nombre:              typeVaccine.Nombre,
			DosisMinimas:        typeVaccine.DosisMinimas,
			DosisMaximas:        typeVaccine.DosisMaximas,
			IntervaloEntreDosis: typeVaccine.IntervaloEntreDosis,
			Estado:              typeVaccine.Estado,
		}

		result = append(result, data)
	}

	return result, nil
}

func (s *NursingConsultationService) GetVaccinesRequired(patientId string) (*models.NextVaccinePatient, int, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)
	result := &models.NextVaccinePatient{
		Required:    nil,
		NotRequired: nil,
	}

	vaccines, code, err := srvService.SrvNursingConsultationVaccine.GetVaccineByIDPatient(patientId)
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get Types Vaccines Required:", err)
		return nil, code, fmt.Errorf("failed to get Types Vaccines Required: %w", err)
	}
	if vaccines == nil {
		requiredAllVaccines, err := s.getVaccinesRequired()
		if err != nil {
			logger.Error.Println(s.txID, " - failed to get all required Vaccines:", err)
			return nil, 0, fmt.Errorf("failed to get all required Vaccines: %w", err)
		}
		result.Required = requiredAllVaccines
		return result, code, nil
	}
	nextVaccinesRequired, code, err := s.getNextVaccinesRequired(vaccines)
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get next Vaccines Required:", err)
		return nil, code, fmt.Errorf("failed to get next Vaccines Required: %w", err)
	}

	result.Required = nextVaccinesRequired

	nextVaccinesNotRequired, code, err := s.getNextVaccinesNotRequired(vaccines)
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get next Vaccines Required:", err)
		return nil, code, fmt.Errorf("failed to get next Vaccines Required: %w", err)
	}

	result.NotRequired = nextVaccinesNotRequired

	return result, StatusSuccess, nil
}

func (s *NursingConsultationService) getNextVaccinesRequired(vaccines []*nursing_consultation_vaccine.Vaccine) ([]*models.TypesVaccineRequired, int, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	requiredVaccines, err := srvService.SrvNursingConsultationVaccine.GetAllRequiredVaccineTypes()
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get all Types Vaccines:", err)
		return nil, 0, fmt.Errorf("failed to get all Types Vaccines: %w", err)
	}
	if requiredVaccines == nil {
		return nil, 0, nil
	}

	var result []*models.TypesVaccineRequired

	for _, requiredVaccine := range requiredVaccines {

		vaccinesPatient := s.getPatientVaccinesByName(requiredVaccine.Nombre, vaccines)
		if vaccinesPatient == nil {
			result = append(result, &models.TypesVaccineRequired{
				Nombre:         requiredVaccine.Nombre,
				FechaRequerida: nil,
				Requerido:      true,
			})
			continue
		}

		totalVaccinesPatient := len(vaccinesPatient)
		if requiredVaccine.DosisMinimas == totalVaccinesPatient {
			continue
		}

		vaccinesSorted := s.sortVaccinesByDosisDate(vaccinesPatient)
		lastVaccinePatient := vaccinesSorted[totalVaccinesPatient-1]

		interval, code, err := srvService.SrvNursingConsultationVaccine.GetNextVaccineIntervalsByVaccineType(requiredVaccine.ID, totalVaccinesPatient)
		if err != nil {
			logger.Error.Println(s.txID, " - failed to get interval Types Vaccines:", err)
			return nil, code, fmt.Errorf("failed to get interval Types Vaccines: %w", err)
		}

		isRequiredNextVaccine, dateNextVaccine, err := s.isRequiredNextVaccine(lastVaccinePatient.FechaDosis, interval.MesesEspera)
		if err != nil {
			logger.Error.Println(s.txID, " - failed validate next vaccine:", err)
			return nil, 0, fmt.Errorf("failed validate next vaccine: %w", err)
		}

		result = append(result, &models.TypesVaccineRequired{
			Nombre:         requiredVaccine.Nombre,
			FechaRequerida: dateNextVaccine,
			Requerido:      isRequiredNextVaccine,
		})

	}

	return result, StatusSuccess, nil
}

func (s *NursingConsultationService) getVaccinesRequired() ([]*models.TypesVaccineRequired, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)
	requiredVaccines, err := srvService.SrvNursingConsultationVaccine.GetAllRequiredVaccineTypes()
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get all Types Vaccines:", err)
		return nil, fmt.Errorf("failed to get all Types Vaccines: %w", err)
	}
	if requiredVaccines == nil {
		return nil, nil
	}

	var result []*models.TypesVaccineRequired

	for _, requiredVaccine := range requiredVaccines {

		result = append(result, &models.TypesVaccineRequired{
			Nombre:         requiredVaccine.Nombre,
			FechaRequerida: nil,
			Requerido:      true,
		})
	}

	return result, nil
}

func (s *NursingConsultationService) getNextVaccinesNotRequired(vaccines []*nursing_consultation_vaccine.Vaccine) ([]*models.TypesVaccineRequired, int, error) {
	srvService := medical_area.NewServerMedicalArea(s.db, s.usr, s.txID)

	notRequiredVaccines, err := srvService.SrvNursingConsultationVaccine.GetAllNotRequiredVaccineTypes()
	if err != nil {
		logger.Error.Println(s.txID, " - failed to get all Types Vaccines:", err)
		return nil, 0, fmt.Errorf("failed to get all Types Vaccines: %w", err)
	}
	if notRequiredVaccines == nil {
		return nil, 0, nil
	}

	var result []*models.TypesVaccineRequired

	for _, notRequiredVaccine := range notRequiredVaccines {

		vaccinesPatient := s.getPatientVaccinesByName(notRequiredVaccine.Nombre, vaccines)
		if vaccinesPatient == nil {
			continue
		}

		totalVaccinesPatient := len(vaccinesPatient)
		if notRequiredVaccine.DosisMaximas == nil {
			continue
		}

		if *notRequiredVaccine.DosisMaximas == totalVaccinesPatient {
			continue
		}

		vaccinesSorted := s.sortVaccinesByDosisDate(vaccinesPatient)
		lastVaccinePatient := vaccinesSorted[totalVaccinesPatient-1]

		interval, code, err := srvService.SrvNursingConsultationVaccine.GetNextVaccineIntervalsByVaccineType(notRequiredVaccine.ID, totalVaccinesPatient)
		if err != nil {
			logger.Error.Println(s.txID, " - failed to get interval Types Vaccines:", err)
			return nil, code, fmt.Errorf("failed to get interval Types Vaccines: %w", err)
		}

		if interval == nil {
			continue
		}

		isRequiredNextVaccine, dateNextVaccine, err := s.isRequiredNextVaccine(lastVaccinePatient.FechaDosis, interval.MesesEspera)
		if err != nil {
			logger.Error.Println(s.txID, " - failed validate next vaccine:", err)
			return nil, 0, fmt.Errorf("failed validate next vaccine: %w", err)
		}

		result = append(result, &models.TypesVaccineRequired{
			Nombre:         notRequiredVaccine.Nombre,
			FechaRequerida: dateNextVaccine,
			Requerido:      isRequiredNextVaccine,
		})

	}

	return result, StatusSuccess, nil
}

func (s *NursingConsultationService) getPatientVaccinesByName(vaccineName string, vaccines []*nursing_consultation_vaccine.Vaccine) []*nursing_consultation_vaccine.Vaccine {
	var vaccinesPatient []*nursing_consultation_vaccine.Vaccine
	for _, vaccine := range vaccines {
		if vaccine.TipoVacuna == vaccineName {
			vaccinesPatient = append(vaccinesPatient, vaccine)
		}
	}

	return vaccinesPatient
}

func (s *NursingConsultationService) sortVaccinesByDosisDate(vaccines []*nursing_consultation_vaccine.Vaccine) []*nursing_consultation_vaccine.Vaccine {
	sort.Slice(vaccines, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", vaccines[i].FechaDosis)
		dateJ, errJ := time.Parse("2006-01-02", vaccines[j].FechaDosis)

		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}

		return dateI.Before(dateJ)
	})

	return vaccines
}

func (s *NursingConsultationService) isRequiredNextVaccine(dateStr string, monthsToAdd int) (bool, *string, error) {
	baseDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false, nil, fmt.Errorf("error parseando fecha: %v", err)
	}
	baseDate = baseDate.UTC()

	futureDate := baseDate.AddDate(0, monthsToAdd, 0)
	futureDateStr := futureDate.Format("2006-01-02")

	today := time.Now().UTC()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	return today.After(futureDate), &futureDateStr, nil
}

func (s *NursingConsultationService) fetchServicesNursing(srvService *medical_area.ServerMedicalArea, id string) (string, int, error) {
	var services []string

	vaccines, err := fetchVaccines(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if vaccines != nil {
		services = append(services, "Vacunas")
	}

	physicalTest, err := fetchPhysicalTest(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if physicalTest != nil {
		services = append(services, "Examen Físico")
	}

	laboratoryTest, err := fetchLaboratoryTest(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if laboratoryTest != nil {
		services = append(services, "Examen Laboratorio")
	}

	preferentialTest, err := fetchPreferentialTest(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if preferentialTest != nil {
		services = append(services, "Examen Preferencial")
	}

	sexualityTest, err := fetchSexualityTest(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if sexualityTest != nil {
		services = append(services, "Examen Sexualidad")
	}

	visualTest, err := fetchVisualTest(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if visualTest != nil {
		services = append(services, "Examen Visual")
	}

	medicationTreatment, err := fetchMedicationTreatment(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if medicationTreatment != nil {
		services = append(services, "Tratamiento Medicamentoso")
	}

	performedProcedures, err := fetchPerformedProcedures(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if performedProcedures != nil {
		services = append(services, "Procedimientos Realizados")
	}

	consultationIntegralAttention, err := fetchConsultationIntegralAttention(srvService, id, s.txID)
	if err != nil {
		return "", StatusNotFound, err
	}
	if consultationIntegralAttention != nil {
		services = append(services, "Atención Integral")
	}

	result := strings.Join(services, ", ")

	return result, StatusSuccess, nil
}

func (s *NursingConsultationService) fetchConsultationData(srvService *medical_area.ServerMedicalArea, id string) (models.Exams, int, error) {

	exams := models.Exams{}

	vaccines, err := fetchVaccines(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	physicalTest, err := fetchPhysicalTest(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	laboratoryTest, err := fetchLaboratoryTest(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	preferentialTest, err := fetchPreferentialTest(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	sexualityTest, err := fetchSexualityTest(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	visualTest, err := fetchVisualTest(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	medicationTreatment, err := fetchMedicationTreatment(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	performedProcedures, err := fetchPerformedProcedures(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	consultationIntegralAttention, err := fetchConsultationIntegralAttention(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	consultationBuccal, err := fetchBuccalConsultation(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	consultationGeneralMedicine, err := fetchConsultationGeneralMedicine(srvService, id, s.txID)
	if err != nil {
		return exams, StatusNotFound, err
	}

	//buccalConsultation, err := fetchBuccalProcedure(srvService, id, s.txID)
	//if err != nil {
	//	return exams, StatusNotFound, err
	//}
	//
	//generalMedicineConsultation, err := fetchConsultationIntegralAttention(srvService, id, s.txID)
	//if err != nil {
	//	return exams, StatusNotFound, err
	//}

	exams.Vacunas = vaccines
	exams.ExamenFisico = physicalTest
	exams.ExamenLaboratorio = laboratoryTest
	exams.ExamenPreferencial = preferentialTest
	exams.ExamenSexualidad = sexualityTest
	exams.ExamenVisual = visualTest
	exams.TratamientoMedicamentoso = medicationTreatment
	exams.ProcedimientoRealizado = performedProcedures
	exams.AtencionIntegralOtros = consultationIntegralAttention
	exams.ConsultaBucal = consultationBuccal
	exams.ConsultaMedicinaGeneral = consultationGeneralMedicine

	return exams, StatusSuccess, nil
}

func fetchRoutineReview(srvService *medical_area.ServerMedicalArea, id, txID string) (models.RoutineReview, error) {
	routineReview, _, err := srvService.SrvNursingConsultationRoutineReview.GetRoutineReviewByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Routine Review:", err)
		return models.RoutineReview{}, fmt.Errorf("failed to get Routine Review: %w", err)
	}
	if routineReview == nil {
		return models.RoutineReview{}, fmt.Errorf("Routine Review not found")
	}
	return mapRoutineReview(routineReview), nil
}

func fetchAccompanyingData(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.AccompanyingData, error) {
	accompanyingData, _, err := srvService.SrvNursingConsultationAccompanyingData.GetAccompanyingDataByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Accompanying Data:", err)
		return nil, fmt.Errorf("failed to get Accompanying Data: %w", err)
	}
	if accompanyingData == nil {
		return nil, nil
	}
	return &models.AccompanyingData{
		ID:               accompanyingData.ID,
		DNI:              accompanyingData.DNI,
		NombresApellidos: accompanyingData.NombresApellidos,
		Edad:             accompanyingData.Edad,
	}, nil
}

func fetchVaccines(srvService *medical_area.ServerMedicalArea, id, txID string) ([]*models.Vaccine, error) {
	vaccines, _, err := srvService.SrvNursingConsultationVaccine.GetVaccineByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Vaccines:", err)
		return nil, fmt.Errorf("failed to get Vaccines: %w", err)
	}
	if vaccines == nil {
		return nil, nil
	}
	result := make([]*models.Vaccine, len(vaccines))
	for i, vaccine := range vaccines {
		result[i] = &models.Vaccine{
			ID:            vaccine.ID,
			TipoVacuna:    vaccine.TipoVacuna,
			FechaDosis:    vaccine.FechaDosis,
			Minsa:         vaccine.Minsa,
			Observaciones: vaccine.Observaciones,
		}
	}
	return result, nil
}

func fetchPhysicalTest(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.PhysicalTest, error) {
	physicalTest, _, err := srvService.SrvNursingConsultationPhysicalTest.GetPhysicalTestByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Physical Test:", err)
		return nil, fmt.Errorf("failed to get Physical Test: %w", err)
	}
	if physicalTest == nil {
		return nil, nil
	}
	return &models.PhysicalTest{
		ID:                    physicalTest.ID,
		TallaPesos:            physicalTest.TallaPesos,
		PerimetroCintura:      physicalTest.PerimetroCintura,
		IndiceMasaCorporalImg: physicalTest.IndiceMasaCorporalImg,
		PresionArterial:       physicalTest.PresionArterial,
		Comentarios:           physicalTest.Comentarios,
	}, nil

}

func fetchLaboratoryTest(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.LaboratoryTest, error) {
	laboratoryTest, _, err := srvService.SrvNursingConsultationLaboratoryTest.GetLaboratoryTestByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Laboratory Test:", err)
		return nil, fmt.Errorf("failed to get Laboratory Test: %w", err)
	}
	if laboratoryTest == nil {
		return nil, nil
	}
	return &models.LaboratoryTest{
		ID:          laboratoryTest.ID,
		Serologia:   laboratoryTest.Serologia,
		Bk:          laboratoryTest.Bk,
		Hemograma:   laboratoryTest.Hemograma,
		ExamenOrina: laboratoryTest.ExamenOrina,
		Colesterol:  laboratoryTest.Colesterol,
		Glucosa:     laboratoryTest.Glucosa,
		Comentarios: laboratoryTest.Comentarios,
	}, nil
}

func fetchPreferentialTest(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.PreferentialTest, error) {
	preferentialTest, _, err := srvService.SrvNursingConsultationPreferentialTest.GetPreferentialTestByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Preferential Test:", err)
		return nil, fmt.Errorf("failed to get Preferential Test: %w", err)
	}
	if preferentialTest == nil {
		return nil, nil
	}
	return &models.PreferentialTest{
		ID:                    preferentialTest.ID,
		AparatoRespiratorio:   preferentialTest.AparatoRespiratorio,
		AparatoCardiovascular: preferentialTest.AparatoCardiovascular,
		AparatoDigestivo:      preferentialTest.AparatoDigestivo,
		AparatoGenitourinario: preferentialTest.AparatoGenitourinario,
		Papanicolau:           preferentialTest.Papanicolau,
		ExamenMama:            preferentialTest.ExamenMama,
		Comentarios:           preferentialTest.Comentarios,
	}, nil

}

func fetchSexualityTest(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.SexualityTest, error) {
	sexualityTest, _, err := srvService.SrvNursingConsultationSexualityTest.GetSexualityTestByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Sexuality Test:", err)
		return nil, fmt.Errorf("failed to get Sexuality Test: %w", err)
	}
	if sexualityTest == nil {
		return nil, nil
	}
	return &models.SexualityTest{
		ID:                    sexualityTest.ID,
		ActividadSexual:       sexualityTest.ActividadSexual,
		PlanificacionFamiliar: sexualityTest.PlanificacionFamiliar,
		Comentarios:           sexualityTest.Comentarios,
	}, nil
}

func fetchVisualTest(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.VisualTest, error) {
	visualTest, _, err := srvService.SrvNursingConsultationVisualTest.GetVisualTestByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Visual Test:", err)
		return nil, fmt.Errorf("failed to get Visual Test: %w", err)
	}
	if visualTest == nil {
		return nil, nil
	}
	return &models.VisualTest{
		ID:            visualTest.ID,
		OjoDerecho:    visualTest.OjoDerecho,
		OjoIzquierdo:  visualTest.OjoIzquierdo,
		PresionOcular: visualTest.PresionOcular,
		Comentarios:   visualTest.Comentarios,
	}, nil
}

func fetchMedicationTreatment(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.MedicationTreatment, error) {
	medicationTreatment, _, err := srvService.SrvNursingConsultationMedicationTreatment.GetMedicationTreatmentByIDConsultation(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Medication Treatment:", err)
		return nil, fmt.Errorf("failed to get Medication Treatment: %w", err)
	}
	if medicationTreatment == nil {
		return nil, nil
	}

	return &models.MedicationTreatment{
		ID:                        medicationTreatment.ID,
		NombreGenericoMedicamento: medicationTreatment.NombreGenericoMedicamento,
		ViaAdministracion:         medicationTreatment.ViaAdministracion,
		HoraAplicacion:            medicationTreatment.HoraAplicacion,
		ResponsableAtencion:       medicationTreatment.ResponsableAtencion,
		Observaciones:             medicationTreatment.Observaciones,
		AreaSolicitante:           medicationTreatment.AreaSolicitante,
		EspecialistaSolicitante:   medicationTreatment.EspecialistaSolicitante,
	}, nil
}

func fetchPerformedProcedures(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.PerformedProcedure, error) {
	performedProcedure, _, err := srvService.SrvNursingConsultationPerformedProcedures.GetPerformedProcedureByIDConsultation(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Performed Procedures:", err)
		return nil, fmt.Errorf("failed to get Performed Procedures: %w", err)
	}
	if performedProcedure == nil {
		return nil, nil
	}
	return &models.PerformedProcedure{
		ID:                      performedProcedure.ID,
		Procedimiento:           performedProcedure.Procedimiento,
		NumeroRecibo:            performedProcedure.NumeroRecibo,
		Comentarios:             performedProcedure.Comentarios,
		Costo:                   performedProcedure.Costo,
		FechaPago:               performedProcedure.FechaPago,
		AreaSolicitante:         performedProcedure.AreaSolicitante,
		EspecialistaSolicitante: performedProcedure.EspecialistaSolicitante,
	}, nil
}

func fetchConsultationIntegralAttention(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.IntegralAttentionOther, error) {
	integralAttention, _, err := srvService.SrvConsultationIntegralAttention.GetConsultationIntegralAttentionByIDConsultation(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get Consultation Integral Attention:", err)
		return nil, fmt.Errorf("failed to get Visual Test: %w", err)
	}
	if integralAttention == nil {
		return nil, nil
	}
	return &models.IntegralAttentionOther{
		ID:                     integralAttention.ID,
		Fecha:                  integralAttention.Fecha,
		Hora:                   integralAttention.Hora,
		Edad:                   integralAttention.Edad,
		MotivoConsulta:         integralAttention.MotivoConsulta,
		TiempoEnfermedad:       integralAttention.TiempoEnfermedad,
		Apetito:                integralAttention.Apetito,
		Sed:                    integralAttention.Sed,
		Suenio:                 integralAttention.Suenio,
		EstadoAnimo:            integralAttention.EstadoAnimo,
		Orina:                  integralAttention.Orina,
		Deposiciones:           integralAttention.Deposiciones,
		Temperatura:            integralAttention.Temperatura,
		PresionArterial:        integralAttention.PresionArterial,
		FrecuenciaCardiaca:     integralAttention.FrecuenciaCardiaca,
		FrecuenciaRespiratoria: integralAttention.FrecuenciaRespiratoria,
		Peso:                   integralAttention.Peso,
		Talla:                  integralAttention.Talla,
		IndiceMasaCorporal:     integralAttention.IndiceMasaCorporal,
		Diagnostico:            integralAttention.Diagnostico,
		Tratamiento:            integralAttention.Tratamiento,
		ExamenesAxuliares:      integralAttention.ExamenesAxuliares,
		Referencia:             integralAttention.Referencia,
		Observacion:            integralAttention.Observacion,
		NumeroRecibo:           integralAttention.NumeroRecibo,
		Costo:                  integralAttention.Costo,
		FechaPago:              integralAttention.FechaPago,
	}, nil
}

func fetchConsultationGeneralMedicine(srvService *medical_area.ServerMedicalArea, id, txID string) (*models.GeneralMedicineConsultation, error) {
	generalMedicineConsultation, _, err := srvService.SrvMedicalGeneralMedicineConsultation.GetGeneralMedicineConsultationByID(id)
	if err != nil {
		logger.Error.Println(txID, " - failed to get General Medicine Consultation:", err)
		return nil, fmt.Errorf("failed to get General Medicine Consultation: %w", err)
	}
	if generalMedicineConsultation == nil {
		return nil, nil
	}
	return &models.GeneralMedicineConsultation{
		ID:            generalMedicineConsultation.ID,
		FechaHora:     generalMedicineConsultation.FechaHora,
		Anamnesis:     generalMedicineConsultation.Anamnesis,
		ExamenClinico: generalMedicineConsultation.ExamenClinico,
		Indicaciones:  generalMedicineConsultation.Indicaciones,
	}, nil

}

func mapNursingConsultation(nursingConsultation *consultation_medical_area.ConsultationMedicalArea) models.ConsultationMedicalArea {
	return models.ConsultationMedicalArea{
		ID:            nursingConsultation.ID,
		IDPaciente:    nursingConsultation.IDPaciente,
		FechaConsulta: nursingConsultation.FechaConsulta,
		AreaMedica:    &nursingConsultation.AreaMedica,
	}
}

func mapRoutineReview(routineReview *nursing_consultation_routine_review.RoutineReview) models.RoutineReview {
	return models.RoutineReview{
		ID:                       routineReview.ID,
		FiebreUltimoQuinceDias:   routineReview.FiebreUltimoQuinceDias,
		TosMasQuinceDias:         routineReview.TosMasQuinceDias,
		SecrecionLesionGenitales: routineReview.SecrecionLesionGenitales,
		FechaUltimaRegla:         routineReview.FechaUltimaRegla,
		Comentarios:              routineReview.Comentarios,
	}
}
