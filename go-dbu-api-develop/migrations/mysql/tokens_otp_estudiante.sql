-- Tabla para tokens OTP de estudiantes
CREATE TABLE `tokens_otp_estudiante`
(
    `id`                CHAR(36)     NOT NULL PRIMARY KEY,
    `alumno_id`         bigint(20) unsigned NOT NULL,
    `dni`               VARCHAR(255) NOT NULL,
    `codigo_otp`        VARCHAR(6)   NOT NULL,
    `correo_destino`    VARCHAR(255) NOT NULL,
    `fecha_generacion`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `fecha_expiracion`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `intentos_fallidos` INT                   DEFAULT 0,
    `estado`            ENUM('pendiente', 'usado', 'expirado', 'bloqueado') DEFAULT 'pendiente',
    `direccion_ip`      VARCHAR(45),
    `created_at`        TIMESTAMP NULL DEFAULT NULL,
    `updated_at`        TIMESTAMP NULL DEFAULT NULL,
    KEY                 `idx_otp_dni` (`dni`),
    KEY                 `idx_otp_codigo` (`codigo_otp`),
    KEY                 `idx_otp_estado` (`estado`),
    KEY                 `idx_otp_alumno` (`alumno_id`),
    CONSTRAINT `fk_otp_alumno` FOREIGN KEY (`alumno_id`) REFERENCES `alumnos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;