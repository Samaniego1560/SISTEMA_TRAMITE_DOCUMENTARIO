-- Migración para la tabla de notificaciones
CREATE TABLE IF NOT EXISTS notificacion (
    id INT AUTO_INCREMENT PRIMARY KEY,
    falta_id VARCHAR(64) NOT NULL,
    numero_notificacion VARCHAR(10) NOT NULL,
    anio INT NOT NULL,
    fecha TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY (falta_id)
);
