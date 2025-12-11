package fault

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"dbu-api/internal/models"

	"github.com/jmoiron/sqlx"
)

// sqlserver estructura de conexión a la BD de mssql
type sqlserver struct {
	DB   *sqlx.DB
	user *models.User
	TxID string
}

func newFaultSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// ----------------- MÉTODOS PARA LA TABLA FALTAS -----------------

// Create registra en la BD
func (s *sqlserver) create(m *Fault) error {
	date := time.Now()
	m.UpdatedAt = date
	m.CreatedAt = date

	const sqlInsert = `
INSERT INTO faltas (
    id, alumno_id, convocatoria_id, servicio_id, observacion, fuente_informacion, fecha_falta, 
    estado, apelable, apelacion_documento, motivo_resolucion, created_at, updated_at
) VALUES (
    :id, :alumno_id, :convocatoria_id, :servicio_id, :observacion, :fuente_informacion, :fecha_falta, 
    :estado, :apelable, :apelacion_documento, :motivo_resolucion, :created_at, :updated_at
)`
	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Update actualiza un registro en la BD
func (s *sqlserver) update(m *Fault) error {
	date := time.Now()
	m.UpdatedAt = date

	const sqlUpdate = `
UPDATE faltas SET
    alumno_id = :alumno_id,
    convocatoria_id = :convocatoria_id,
    servicio_id = :servicio_id,
    observacion = :observacion,
    fuente_informacion = :fuente_informacion,
    fecha_falta = :fecha_falta,
    estado = :estado,
    apelable = :apelable,
    apelacion_documento = :apelacion_documento,
    motivo_resolucion = :motivo_resolucion,
    updated_at = :updated_at
WHERE id = :id
`
	rs, err := s.DB.NamedExec(sqlUpdate, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetAllIncisosByAlumnoID obtiene todos los incisos cometidos por un alumno, con su gravedad y fecha
func (s *sqlserver) GetAllIncisosByAlumnoID(alumnoID int64) ([]*FaultIncisoDetalle, error) {
	const query = `
		SELECT 
			f.id AS fault_id,
			i.id AS inciso_id,
			ar.gravedad AS gravedad,
			f.fecha_falta
		FROM faltas f
		LEFT JOIN faltas_incisos fi ON f.id = fi.falta_id
		LEFT JOIN incisos i ON fi.inciso_id = i.id
		LEFT JOIN faltas_articulos fa ON f.id = fa.falta_id
		LEFT JOIN articulos ar ON fa.articulo_id = ar.id
		WHERE f.alumno_id = ?
	`
	var detalles []*FaultIncisoDetalle
	err := s.DB.Select(&detalles, query, alumnoID)
	if err != nil {
		return nil, err
	}
	return detalles, nil
}

// Delete elimina un registro de la BD
func (s *sqlserver) delete(id string) error {
	const sqlDelete = `DELETE FROM faltas WHERE id = :id`
	m := Fault{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// GetByID consulta un registro por su ID
func (s *sqlserver) getByID(id string) (*Fault, error) {
	const sqlGetByID = `
SELECT 
    id, alumno_id, convocatoria_id, servicio_id, observacion, fuente_informacion, fecha_falta, 
    estado, apelable, apelacion_documento, motivo_resolucion, created_at, updated_at
FROM faltas
WHERE id = ?
`
	fmt.Printf("[DEBUG] getByID: buscando falta con id = %s\n", id)
	mdl := Fault{}
	err := s.DB.Get(&mdl, sqlGetByID, id)
	if err != nil {
		fmt.Printf("[DEBUG] getByID: error al buscar falta: %v\n", err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return &mdl, err
	}
	fmt.Printf("[DEBUG] getByID: falta encontrada: %+v\n", mdl)
	return &mdl, nil
}

// GetAll obtiene faltas con datos del alumno
func (s *sqlserver) getAll() ([]*FaultWithStudent, error) {
	const sqlGetAll = `
	SELECT 
	  f.id, f.alumno_id, f.convocatoria_id, f.servicio_id, f.fuente_informacion, f.fecha_falta, 
	  f.estado, f.apelable, f.apelacion_documento, f.motivo_resolucion, f.observacion, 
	  f.created_at, f.updated_at, 
	  a.dni, a.nombres, a.apellido_paterno, a.apellido_materno, a.escuela_profesional,
	  IFNULL((SELECT CAST(c.numero AS CHAR) FROM asignacion_cuartos ac 
		INNER JOIN cuartos c ON c.id = ac.cuarto_id 
		WHERE ac.alumno_id = a.id AND ac.estado = 'activo' ORDER BY ac.fecha_asignacion DESC LIMIT 1), '') AS room_number,
	  IFNULL((SELECT r.nombre FROM asignacion_cuartos ac 
		INNER JOIN cuartos c ON c.id = ac.cuarto_id 
		INNER JOIN residencias r ON r.id = c.residencia_id 
		WHERE ac.alumno_id = a.id AND ac.estado = 'activo' ORDER BY ac.fecha_asignacion DESC LIMIT 1), '') AS residence_name,
	  IFNULL((SELECT ac.fecha_asignacion FROM asignacion_cuartos ac 
		WHERE ac.alumno_id = a.id AND ac.estado = 'activo' ORDER BY ac.fecha_asignacion DESC LIMIT 1), 
		IFNULL((SELECT srv_sol.updated_at FROM solicitudes sol 
		  INNER JOIN servicio_solicitado srv_sol ON srv_sol.solicitud_id = sol.id 
		  WHERE sol.alumno_id = a.id ORDER BY srv_sol.updated_at DESC LIMIT 1), '')) AS admission_date,
	  IFNULL((SELECT ds1.respuesta_formulario FROM detalle_solicitudes ds1 
		  INNER JOIN requisitos req1 ON req1.id = ds1.requisito_id 
		  WHERE ds1.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req1.nombre =
		   'celular de estudiante' LIMIT 1), '') AS celular_estudiante,
	  IFNULL((SELECT ds2.respuesta_formulario FROM detalle_solicitudes ds2 
		  INNER JOIN requisitos req2 ON req2.id = ds2.requisito_id 
		  WHERE ds2.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req2.nombre = 
		  
		  'Celular padre' LIMIT 1), '') AS celular_padre,
	  IFNULL((SELECT ds3.respuesta_formulario FROM detalle_solicitudes ds3 
		  INNER JOIN requisitos req3 ON req3.id = ds3.requisito_id 
		  WHERE ds3.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req3.nombre = 
		  'Departamento de procedencia' LIMIT 1), '') AS departamento_procedencia,
	  IFNULL((SELECT ds4.respuesta_formulario FROM detalle_solicitudes ds4 
		  INNER JOIN requisitos req4 ON req4.id = ds4.requisito_id 
		  WHERE ds4.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req4.nombre = 
		  'Provincia de procedencia' LIMIT 1), '') AS provincia_procedencia,
	  IFNULL((SELECT ds5.respuesta_formulario FROM detalle_solicitudes ds5 
		  INNER JOIN requisitos req5 ON req5.id = ds5.requisito_id 
		  WHERE ds5.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req5.nombre = 
		  'Distrito de procedencia' LIMIT 1), '') AS distrito_procedencia,
	  CONCAT(
		IFNULL((SELECT ds3.respuesta_formulario FROM detalle_solicitudes ds3 
			INNER JOIN requisitos req3 ON req3.id = ds3.requisito_id 
			WHERE ds3.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req3.nombre = 
			'Departamento de procedencia' LIMIT 1), ''), '/',
		IFNULL((SELECT ds4.respuesta_formulario FROM detalle_solicitudes ds4 
			INNER JOIN requisitos req4 ON req4.id = ds4.requisito_id 
			WHERE ds4.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req4.nombre = 
			'Provincia de procedencia' LIMIT 1), ''), '/',
		IFNULL((SELECT ds5.respuesta_formulario FROM detalle_solicitudes ds5 
			INNER JOIN requisitos req5 ON req5.id = ds5.requisito_id 
			WHERE ds5.solicitud_id = (SELECT id FROM solicitudes WHERE alumno_id = a.id ORDER BY id DESC LIMIT 1) AND req5.nombre = 
			'Distrito de procedencia' LIMIT 1), '')
	  ) AS lugar_procedencia,
	  IFNULL((SELECT CONCAT(r.nombre, ' / Cuarto: ', CAST(c.numero AS CHAR)) FROM asignacion_cuartos ac 
		INNER JOIN cuartos c ON c.id = ac.cuarto_id 
		INNER JOIN residencias r ON r.id = c.residencia_id 
		WHERE ac.alumno_id = a.id AND ac.estado = 'activo' ORDER BY ac.fecha_asignacion DESC LIMIT 1), '') AS direccion,
	  a.correo_institucional,
	  a.codigo_estudiante,
	  a.sexo,
	  a.edad,
	  COALESCE(
		CASE 
		  WHEN EXISTS (SELECT 1 FROM articulos ar
			   JOIN faltas_articulos fa ON fa.articulo_id = ar.id
			   WHERE fa.falta_id = f.id AND ar.gravedad = 'grave') THEN 'grave'
		  WHEN EXISTS (SELECT 1 FROM articulos ar
			   JOIN faltas_articulos fa ON fa.articulo_id = ar.id
			   WHERE fa.falta_id = f.id AND ar.gravedad = 'leve') THEN 'leve'
		END, '-') AS gravedad
	FROM faltas f
	INNER JOIN alumnos a ON f.alumno_id = a.id
	ORDER BY f.created_at DESC;
	`
	var faults []*FaultWithStudent
	err := s.DB.Select(&faults, sqlGetAll)
	if err != nil {
		return nil, err
	}
	return faults, nil
}

// ----------------- MÉTODOS PARA LA TABLA FALTAS_ARTICULOS -----------------

// CreateFaultArticulo crea una relación entre falta y artículo
func (s *sqlserver) createFaultArticulo(m *FaultArticulo) error {
	date := time.Now()
	m.CreatedAt = date
	m.UpdatedAt = date

	const sqlInsert = `
INSERT INTO faltas_articulos (id, falta_id, articulo_id, created_at, updated_at)
VALUES (:id, :falta_id, :articulo_id, :created_at, :updated_at)`
	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// ----------------- MÉTODOS PARA LA TABLA FALTAS_INCISOS -----------------

// CreateFaultInciso crea una relación entre falta e inciso
func (s *sqlserver) createFaultInciso(m *FaultInciso) error {
	date := time.Now()
	m.CreatedAt = date
	m.UpdatedAt = date

	const sqlInsert = `
INSERT INTO faltas_incisos (id, falta_id, inciso_id, created_at, updated_at)
VALUES (:id, :falta_id, :inciso_id, :created_at, :updated_at)`
	rs, err := s.DB.NamedExec(sqlInsert, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// ----------------- MÉTODOS ADICIONALES -----------------

func (s *sqlserver) getAlumnoProfileByDni(dni string) (*Alumno, error) {
	const query = `
SELECT id, codigo_estudiante, DNI, nombres, apellido_paterno, apellido_materno,
       sexo, facultad, escuela_profesional, edad, correo_institucional,
       direccion, lugar_procedencia, celular_estudiante
FROM alumnos WHERE DNI = ? LIMIT 1
`
	var alumno Alumno
	err := s.DB.Get(&alumno, query, dni)
	if err != nil {
		return nil, err
	}
	return &alumno, nil
}

func (s *sqlserver) GetDetalleFalta(faltaID string) ([]*FaultDetalle, error) {
	const query = `
SELECT 
	f.id AS falta_id, 
	f.observacion, 
	f.fuente_informacion,
    f.fecha_falta,
    f.servicio_id,
    a.dni, 
    a.nombres,
    a.apellido_paterno, 
    a.apellido_materno, 
    a.sexo, 
    a.facultad,
    a.escuela_profesional, 
    a.edad, 
    a.correo_institucional AS correo_institucional,
    COALESCE(
        CONCAT(r.direccion, ' - ', r.nombre, ' - Cuarto ', cu.numero),
        a.direccion
    ) AS direccion,
    a.lugar_procedencia, 
    a.celular_estudiante,
    ar.id AS articulo_id, 
    ar.descripcion AS articulo_descripcion, 
    ar.gravedad AS articulo_gravedad,
    i.id AS inciso_id, 
    i.nombre AS inciso_nombre, 
    i.descripcion AS inciso_descripcion,
    c.id AS capitulo_id, 
    c.nombre AS capitulo_nombre,
    n.id AS resolucion_id, 
    n.nombre AS resolucion_nombre,
    d.url AS documento_url
FROM faltas f
JOIN alumnos a ON f.alumno_id = a.id
LEFT JOIN asignacion_cuartos ac ON ac.alumno_id = a.id AND ac.estado = 'activo'
LEFT JOIN cuartos cu ON cu.id = ac.cuarto_id
LEFT JOIN residencias r ON r.id = cu.residencia_id
LEFT JOIN faltas_articulos fa ON f.id = fa.falta_id
LEFT JOIN articulos ar ON fa.articulo_id = ar.id
LEFT JOIN capitulos c ON ar.capitulo_id = c.id
LEFT JOIN resoluciones n ON c.resolucion_id = n.id
LEFT JOIN faltas_incisos fi ON f.id = fi.falta_id
LEFT JOIN incisos i ON fi.inciso_id = i.id
LEFT JOIN faltas_documentos d ON f.id = d.falta_id
WHERE f.id = ?;
`
	var detalles []*FaultDetalle
	err := s.DB.Select(&detalles, query, faltaID)
	if err != nil {
		fmt.Printf("ERROR EN QUERY: %v\n", err)
		return nil, err
	}
	return detalles, nil
}

func (s *sqlserver) GetDetalleFaltaAgrupadoJSON(faltaID string) (string, error) {
	detalles, err := s.GetDetalleFalta(faltaID)
	if err != nil {
		return "", err
	}

	if len(detalles) == 0 {
		return "", fmt.Errorf("no existen detalles para la falta")
	}

	nombreServicio, err := s.GetServicioNombreByID(detalles[0].ServicioID)
	if err != nil {
		nombreServicio = ""
	}

	agrupado := AgruparDetalleFalta(detalles, nombreServicio)

	data, err := json.MarshalIndent(agrupado, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al serializar JSON: %v", err)
	}

	return string(data), nil
}

func (s *sqlserver) getNuevoNumeroNotificacion() (string, error) {
	anio := time.Now().Year()

	tx, err := s.DB.Beginx()
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec(`
		MERGE notificaciones_contador AS target
		USING (SELECT ? AS anio) AS source
		ON target.anio = source.anio
		WHEN MATCHED THEN 
			UPDATE SET ultimo_numero = target.ultimo_numero + 1
		WHEN NOT MATCHED THEN 
			INSERT (anio, ultimo_numero) VALUES (source.anio, 1);
	`, anio)
	if err != nil {
		return "", err
	}

	var numero int
	err = tx.Get(&numero, "SELECT ultimo_numero FROM notificaciones_contador WHERE anio = ?", anio)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%04d", numero), nil
}

func (s *sqlserver) GetServicioNombreByID(servicioID int64) (string, error) {
	var nombre string
	err := s.DB.Get(&nombre, "SELECT nombre FROM servicios WHERE id = ?", servicioID)
	if err != nil {
		return "", err
	}
	return nombre, nil
}

func (s *sqlserver) CreateFaultDocumento(doc *FaultDocumento) error {
	const query = `
INSERT INTO faltas_documentos (id, falta_id, url, archivo, created_at)
VALUES (?, ?, ?, ?, ?)
`
	_, err := s.DB.Exec(query, doc.ID, doc.FaultID, doc.URL, doc.Archivo, doc.CreatedAt)
	return err
}

func (s *sqlserver) GetFaultDocumentoByID(id string) (*FaultDocumento, error) {
	const query = `
SELECT id, falta_id, url, archivo, created_at
FROM faltas_documentos
WHERE id = ?
LIMIT 1
`
	var doc FaultDocumento
	err := s.DB.Get(&doc, query, id)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// UpdateEstado actualiza solo el campo estado de una falta
func (s *sqlserver) UpdateEstado(faltaID string, nuevoEstado string) error {
	const sqlUpdateEstado = `
		UPDATE faltas
		SET estado = @estado,
			updated_at = GETDATE()
		WHERE id = @id
	`
	_, err := s.DB.Exec(sqlUpdateEstado, map[string]interface{}{
		"id":     faltaID,
		"estado": nuevoEstado,
	})
	return err
}
