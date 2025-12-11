CREATE TABLE faltas
(
    id          CHAR(36)                                                                   NOT NULL PRIMARY KEY,
    alumno_id   BIGINT(20)                                                                  UNSIGNED,
    convocatoria_id BIGINT(20)                                                                   UNSIGNED,
    servicio_id BIGINT(20)                                                                       UNSIGNED,
    observacion TEXT                                                              NOT NULL,
    fuente_informacion      VARCHAR (255)                                           NOT NULL,
    fecha_falta      TIMESTAMP                             NOT NULL,
<<<<<<< HEAD
    estado      ENUM ('pendiente', 'en revision', 'sancionado', 'apelado' ) DEFAULT 'pendiente' NOT NULL,
=======
    estado      ENUM ('pendiente', 'notificado', 'sancionado', 'apelado' ) DEFAULT 'pendiente' NOT NULL,
    url_documento TEXT             NOT NULL
>>>>>>> sanciones
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (alumno_id) REFERENCES alumnos (id) ON DELETE CASCADE,
    FOREIGN KEY (convocatoria_id) REFERENCES convocatorias (id) ON DELETE CASCADE,
    FOREIGN KEY (servicio_id) REFERENCES servicios (id) ON DELETE CASCADE
);
<<<<<<< HEAD


CREATE TABLE articulos_faltas(
    id          CHAR(36)                                                                   NOT NULL PRIMARY KEY,
    articulo_id CHAR(36)                                                                    NOt NULL,
    falta_id CHAR(36)                                                                    NOt NULL,
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (articulo_id) REFERENCES articulos (id) ON DELETE CASCADE,
    FOREIGN KEY (falta_id) REFERENCES faltas (id) ON DELETE CASCADE
);



CREATE TABLE incisos_faltas(
    id          CHAR(36)                                                                   NOT NULL PRIMARY KEY,
    inciso_id CHAR(36)                                                                    NOT NULL,
    falta_id CHAR(36)                                                                    NOT NULL,
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inciso_id) REFERENCES incisos (id) ON DELETE CASCADE,
    FOREIGN KEY (falta_id) REFERENCES faltas (id) ON DELETE CASCADE
);

CREATE TABLE archivos_faltas(
    id          CHAR(36)                                                                     NOT NULL,
    falta_id CHAR(36)                                                                    NOT NULL,
    url_documento      VARCHAR (255)                                           NOT NULL,
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (falta_id) REFERENCES faltas (id) ON DELETE CASCADE
);
=======
CREATE TABLE articulos(
    id          CHAR(36)                                                                   NOT NULL PRIMARY KEY,
    descripcion TEXT NOT NULL,
    capitulo_id CHAR(36)                                                                    NOt NULL,
    gravedad ENUM('leve', 'grave'),
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (capitulo_id) REFERENCES capitulos (id) ON DELETE CASCADE,
    UNIQUE (capitulo_id, descripcion)
);
CREATE TABLE incisos(
    id          CHAR(36)                                                                   NOT NULL PRIMARY KEY,
    nombre      VARCHAR(1),
    descripcion TEXT,
    articulo_id CHAR(36)                                                                    NOT NULL,
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (articulo_id) REFERENCES articulos (id) ON DELETE CASCADE,
    UNIQUE (articulo_id, nombre)
);
CREATE TABLE resoluciones(
    id  CHAR(36)        NOT NULL PRIMARY KEY,
    nombre VARCHAR(255)      NOT NULL,
    descripcion TEXT     NOT NULL,
    estado INT(11),
    servicio_id BIGINT(20) UNSIGNED,
    ruta_archivo VARCHAR(255),
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (servicio_id) REFERENCES servicios (id) ON DELETE CASCADE
);
CREATE TABLE capitulos(
    id  CHAR(36)        NOT NULL PRIMARY KEY,
    resolucion_id CHAR(36) NOT NULL,
    nombre VARCHAR(255)      NOT NULL,
    descripcion TEXT     NOT NULL,
    created_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resolucion_id) REFERENCES resoluciones (id) ON DELETE CASCADE,
    UNIQUE (resolucion_id, nombre)
);
CREATE TABLE faltas_articulos (
    id CHAR(36) NOT NULL PRIMARY KEY,
    falta_id CHAR(36) NOT NULL,
    articulo_id CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE,
    FOREIGN KEY (articulo_id) REFERENCES articulos(id) ON DELETE CASCADE
);
CREATE TABLE faltas_incisos (
    id CHAR(36) NOT NULL PRIMARY KEY,
    falta_id CHAR(36) NOT NULL,
    inciso_id CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE,
    FOREIGN KEY (inciso_id) REFERENCES incisos(id) ON DELETE CASCADE
);
-- Crear tabla servicios
CREATE TABLE servicios (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    descripcion TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE faltas_documentos (
    id CHAR(36) NOT NULL PRIMARY KEY,
    falta_id CHAR(36) NOT NULL,
    url TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (falta_id) REFERENCES faltas(id) ON DELETE CASCADE
);


-- Insertar los dos servicios: Comedor y Residencia Estudiantil
INSERT INTO servicios (nombre, descripcion) VALUES
    ('Comedor', 'Servicio de alimentación para estudiantes.'),
    ('Residencia Estudiantil', 'Servicio de residencia para estudiantes.');

CREATE TABLE sancionesAFaltas (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    resolucion_id VARCHAR(36) NOT NULL,
    articulo_id VARCHAR(36) NOT NULL,
    inciso_sancion VARCHAR(1) NOT NULL,
    detalle_sancion TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (resolucion_id) REFERENCES resoluciones(id) ON DELETE CASCADE,
    FOREIGN KEY (articulo_id) REFERENCES articulos(id) ON DELETE CASCADE
);
>>>>>>> sanciones
