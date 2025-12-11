-- Tabla para sesiones de estudiantes
CREATE TABLE `sesiones_estudiante`
(
    `id`                  CHAR(36)  NOT NULL PRIMARY KEY,
    `alumno_id`           bigint(20) unsigned NOT NULL,
    `token_jwt`           TEXT      NOT NULL,
    `direccion_ip`        VARCHAR(45),
    `agente_usuario`      VARCHAR(500),
    `fecha_login`         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `fecha_expiracion`    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `fecha_ultimo_acceso` TIMESTAMP NULL,
    `activo`              BOOLEAN            DEFAULT TRUE,
    `created_at`          TIMESTAMP NULL DEFAULT NULL,
    `updated_at`          TIMESTAMP NULL DEFAULT NULL,
    KEY                   `idx_sesion_alumno` (`alumno_id`),
    KEY                   `idx_sesion_activo` (`activo`),
    KEY                   `idx_sesion_id_activo` (`id`, `activo`),
    CONSTRAINT `fk_sesion_alumno` FOREIGN KEY (`alumno_id`) REFERENCES `alumnos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



