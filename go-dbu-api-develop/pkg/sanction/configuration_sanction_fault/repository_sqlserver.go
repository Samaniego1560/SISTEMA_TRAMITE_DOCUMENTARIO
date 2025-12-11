package configuration_sanction_fault

import (
	"context"
	"fmt"
	"time"

	internalmodels "dbu-api/internal/models"
	"dbu-api/models"

	"github.com/jmoiron/sqlx"
)

// sqlserver implementa SancionRepository usando SQL Server y sqlx
type sqlserver struct {
	DB   *sqlx.DB
	TxID string
}

// Obtener sanciones asignadas a una falta
func (s *sqlserver) GetSancionesAsignadasPorFalta(faltaID string) ([]*internalmodels.SancionAsignadaDetalle, error) {
	const query = `
		SELECT 
			sfa.id,
			sfa.falta_id,
			sfa.resolucion_id,
			sfa.sancion_id,
			sfa.fecha_asignacion,
			sfa.fecha_inicio,
			sfa.fecha_fin,
			sfa.observaciones,
			sfa.created_at,
			sfa.updated_at,
			sa.capitulo_sancion,
			sa.articulo_sancion,
			sa.inciso_sancion,
			sa.detalle_sancion
		FROM sancion_falta_asignada sfa
		INNER JOIN sanciones_faltas_normativa sa ON sa.id = sfa.sancion_id
		WHERE sfa.falta_id = ?`
	var list []*internalmodels.SancionAsignadaDetalle
	err := s.DB.Select(&list, query, faltaID)
	return list, err
}

func NewSancionSqlServerRepository(db *sqlx.DB, txID string) *sqlserver {
	return &sqlserver{
		DB:   db,
		TxID: txID,
	}
}

// Crear sanción
func (s *sqlserver) Create(m *Sancion) error {
	date := time.Now()
	m.CreatedAt = date
	m.UpdatedAt = date
	const sqlInsert = `
		INSERT INTO sanciones_faltas_normativa 
			(id, resolucion_id, articulo_id, capitulo_sancion, articulo_sancion, inciso_sancion, detalle_sancion, created_at, updated_at)
		VALUES 
			(:id, :resolucion_id, :articulo_id, :capitulo_sancion, :articulo_sancion, :inciso_sancion, :detalle_sancion, :created_at, :updated_at)
	`
	rs, err := s.DB.NamedExec(sqlInsert, m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Actualizar sanción
func (s *sqlserver) Update(m *Sancion) error {
	date := time.Now()
	m.UpdatedAt = date
	const sqlUpdate = `
		UPDATE sanciones_faltas_normativa
		SET resolucion_id = :resolucion_id,
			articulo_id = :articulo_id,
			capitulo_sancion = :capitulo_sancion,
			articulo_sancion = :articulo_sancion,
			inciso_sancion = :inciso_sancion,
			detalle_sancion = :detalle_sancion,
			updated_at = :updated_at
		WHERE id = :id
	`
	rs, err := s.DB.NamedExec(sqlUpdate, m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Eliminar sanción
func (s *sqlserver) Delete(id string) error {
	const sqlDelete = `DELETE FROM sanciones_faltas_normativa WHERE id = :id`
	m := Sancion{ID: id}
	rs, err := s.DB.NamedExec(sqlDelete, &m)
	if err != nil {
		return err
	}
	if i, _ := rs.RowsAffected(); i == 0 {
		return fmt.Errorf("ecatch:108")
	}
	return nil
}

// Obtener sanción por ID
func (s *sqlserver) GetByID(id string) (*Sancion, error) {
	const sqlGetByID = `
		SELECT 
			id, 
			resolucion_id, 
			articulo_id, 
			capitulo_sancion,
			articulo_sancion,
			inciso_sancion, 
			detalle_sancion, 
			created_at, 
			updated_at
		FROM sanciones_faltas_normativa
		WHERE id = ?
	`
	var mdl Sancion
	err := s.DB.Get(&mdl, sqlGetByID, id)
	if err != nil {
		return nil, err
	}
	return &mdl, nil
}

// Obtener todas las sanciones
func (s *sqlserver) GetAll() ([]*Sancion, error) {
	const sqlGetAll = `
		SELECT 
			s.id,
			s.resolucion_id,
			r.nombre AS resolucion_nombre,
			s.articulo_id,
			a.descripcion AS articulo_descripcion,
			a.gravedad,
			s.capitulo_sancion,
			c.nombre AS capitulo_nombre,
			s.articulo_sancion,
			s.inciso_sancion,
			s.detalle_sancion,
			s.created_at,
			s.updated_at
		FROM sanciones_faltas_normativa s
		INNER JOIN resoluciones r ON r.id = s.resolucion_id
		INNER JOIN articulos a ON a.id = s.articulo_id
		INNER JOIN capitulos c ON c.id = a.capitulo_id
		ORDER BY a.gravedad, c.nombre, a.descripcion
	`

	var list []*Sancion
	err := s.DB.Select(&list, sqlGetAll)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// AsignarSancionFalta inserta una nueva asignación de sanción a una falta
// ...existing code...
func (s *sqlserver) AsignarSancionFalta(ctx context.Context, sfa *internalmodels.SancionFaltaAsignada) error {
	query := `INSERT INTO sancion_falta_asignada (
			id, falta_id, resolucion_id, sancion_id, fecha_asignacion, fecha_inicio, fecha_fin, observaciones, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.DB.ExecContext(ctx, query,
		sfa.ID,
		sfa.FaltaID,
		sfa.ResolucionID,
		sfa.SancionID,
		sfa.FechaAsignacion,
		sfa.FechaInicio,
		sfa.FechaFin,
		sfa.Observaciones,
		sfa.CreatedAt,
		sfa.UpdatedAt,
	)
	return err
}

// RegistrarApelacion inserta una nueva apelación sobre una sanción asignada
func (s *sqlserver) RegistrarApelacion(ctx context.Context, ap *models.Apelacion) error {
	query := `INSERT INTO apelaciones (
		id, sancion_falta_asignada_id, motivo, estado, usuario_apela, fecha_apelacion, fecha_resolucion, observaciones, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.DB.ExecContext(ctx, query,
		ap.ID,
		ap.SancionFaltaAsignadaID,
		ap.Motivo,
		ap.Estado,
		ap.UsuarioApela,
		ap.FechaApelacion,
		ap.FechaResolucion,
		ap.Observaciones,
		ap.CreatedAt,
		ap.UpdatedAt,
	)
	return err
}
