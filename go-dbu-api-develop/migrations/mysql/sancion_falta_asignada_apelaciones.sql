-- Migración para tabla sancion_falta_asignada
CREATE TABLE IF NOT EXISTS sancion_falta_asignada (
    id VARCHAR(36) PRIMARY KEY,
    resolucion_id VARCHAR(36) NOT NULL,
    sancion_id VARCHAR(36) NOT NULL,
    fecha_asignacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    observaciones TEXT,
    usuario_asigna VARCHAR(36),
    estado VARCHAR(20) DEFAULT 'ACTIVA',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (resolucion_id) REFERENCES resoluciones(id),
    FOREIGN KEY (sancion_id) REFERENCES sancionesAFaltas(id)
);

-- Migración para tabla apelaciones
CREATE TABLE IF NOT EXISTS apelaciones (
    id VARCHAR(36) PRIMARY KEY,
    sancion_falta_asignada_id VARCHAR(36) NOT NULL,
    motivo TEXT,
    estado VARCHAR(20) DEFAULT 'PENDIENTE',
    usuario_apela VARCHAR(36),
    fecha_apelacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_resolucion TIMESTAMP,
    observaciones TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sancion_falta_asignada_id) REFERENCES sancion_falta_asignada(id)
);
