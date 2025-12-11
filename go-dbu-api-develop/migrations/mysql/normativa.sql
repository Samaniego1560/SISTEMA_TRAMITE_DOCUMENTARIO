
-- =====================================
-- TABLA: RESOLUCIONES
-- =====================================
CREATE TABLE resoluciones (
  id CHAR(36) PRIMARY KEY,
  nombre VARCHAR(255) NOT NULL,
  descripcion TEXT NULL,
  servicio_id BIGINT(20) UNSIGNED NULL,
  ruta_archivo VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  estado INT(11) DEFAULT 0,
  INDEX (servicio_id),
  INDEX (nombre)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: CAPITULOS
-- =====================================
CREATE TABLE capitulos (
  id CHAR(36) PRIMARY KEY,
  nombre VARCHAR(255) NOT NULL,
  descripcion TEXT NULL,
  resolucion_id CHAR(36) NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX (resolucion_id),
  FOREIGN KEY (resolucion_id) REFERENCES resoluciones(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: ARTICULOS
-- =====================================
CREATE TABLE articulos (
  id CHAR(36) PRIMARY KEY,
  descripcion TEXT NULL,
  capitulo_id CHAR(36) NULL,
  gravedad ENUM('leve', 'grave') NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX (capitulo_id),
  FOREIGN KEY (capitulo_id) REFERENCES capitulos(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: INCISOS
-- =====================================
CREATE TABLE incisos (
  id CHAR(36) PRIMARY KEY,
  nombre VARCHAR(1) NOT NULL,
  descripcion TEXT NULL,
  articulo_id CHAR(36) NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX (articulo_id),
  FOREIGN KEY (articulo_id) REFERENCES articulos(id) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: FALTAS
-- =====================================
CREATE TABLE faltas (
  id CHAR(36) PRIMARY KEY,
  alumno_id BIGINT(20) UNSIGNED NULL,
  convocatoria_id BIGINT(20) UNSIGNED NULL,
  servicio_id BIGINT(20) UNSIGNED NULL,
  observacion TEXT NOT NULL,
  fuente_informacion VARCHAR(255) NOT NULL,
  fecha_falta TIMESTAMP NOT NULL,
  estado ENUM('registrada','sancionada','notificada','apelada','cerrada') NOT NULL DEFAULT 'registrada',
  apelable TINYINT(1) NOT NULL DEFAULT 1,
  apelacion_documento VARCHAR(255) NULL,
  motivo_resolucion TEXT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX (alumno_id),
  INDEX (convocatoria_id),
  INDEX (servicio_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: FALTAS_ARTICULOS
-- =====================================
CREATE TABLE faltas_articulos (
  id CHAR(36) PRIMARY KEY,
  articulo_id CHAR(36) NOT NULL,
  falta_id CHAR(36) NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (articulo_id) REFERENCES articulos(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: FALTAS_INCISOS
-- =====================================
CREATE TABLE faltas_incisos (
  id CHAR(36) PRIMARY KEY,
  inciso_id CHAR(36) NOT NULL,
  falta_id CHAR(36) NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (inciso_id) REFERENCES incisos(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: FALTAS_DOCUMENTOS
-- =====================================
CREATE TABLE faltas_documentos (
  id CHAR(36) PRIMARY KEY,
  falta_id CHAR(36) NOT NULL,
  url TEXT NOT NULL,
  archivo LONGBLOB NULL,
  created_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: sancionesAFaltas
-- =====================================
CREATE TABLE sanciones_faltas_normativa (
  id VARCHAR(36) PRIMARY KEY,
  resolucion_id VARCHAR(36) NOT NULL,
  articulo_id VARCHAR(36) NOT NULL,
  capitulo_sancion VARCHAR(20) NOT NULL,
  articulo_sancion VARCHAR(20) NOT NULL,
  inciso_sancion VARCHAR(10) NOT NULL,
  detalle_sancion TEXT NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (resolucion_id) REFERENCES resoluciones(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (articulo_id) REFERENCES articulos(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: sancion_falta_asignada
-- =====================================
CREATE TABLE sancion_falta_asignada (
  id VARCHAR(36) PRIMARY KEY,
  resolucion_id VARCHAR(36) NOT NULL,
  sancion_id VARCHAR(36) NOT NULL,
  falta_id VARCHAR(36) NOT NULL,
  fecha_asignacion TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  fecha_inicio DATETIME NULL,
  fecha_fin DATETIME NULL,
  observaciones TEXT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (resolucion_id) REFERENCES resoluciones(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (sancion_id) REFERENCES sancionesAFaltas(id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: apelaciones
-- =====================================
CREATE TABLE apelaciones (
  id VARCHAR(36) PRIMARY KEY,
  sancion_falta_asignada_id VARCHAR(36) NOT NULL,
  motivo TEXT NULL,
  estado VARCHAR(20) DEFAULT 'PENDIENTE',
  usuario_apela VARCHAR(36) NULL,
  fecha_apelacion TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  fecha_resolucion TIMESTAMP NULL,
  observaciones TEXT NULL,
  veredicto_observacion TEXT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (sancion_falta_asignada_id) REFERENCES sancion_falta_asignada(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: apelacion_documentos
-- =====================================
CREATE TABLE apelacion_documentos (
  id VARCHAR(36) PRIMARY KEY,
  apelacion_id VARCHAR(36) NOT NULL,
  documento LONGBLOB NULL,
  nombre VARCHAR(255) NULL,
  tipo VARCHAR(100) NULL,
  created_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (apelacion_id) REFERENCES apelaciones(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: notificacion
-- =====================================
CREATE TABLE notificacion (
  id INT(11) PRIMARY KEY AUTO_INCREMENT,
  falta_id VARCHAR(64) NOT NULL,
  numero_notificacion VARCHAR(10) NOT NULL,
  anio INT(11) NOT NULL,
  fecha TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =====================================
-- TABLA: notificacion_secuencia
-- =====================================
CREATE TABLE notificacion_secuencia (
  anio INT(11) PRIMARY KEY,
  ultimo_numero INT(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
