-- ============================================
-- MÓDULO DE FARMACIA - MIGRACIÓN FASE 1
-- Gestión de Medicamentos y Lotes
-- ============================================

-- 1. Tabla de medicamentos (catálogo maestro)
CREATE TABLE `farmacia_medicamentos` (
  `id` CHAR(36) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP(),
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `is_deleted` TINYINT(1) DEFAULT 0,
  `user_deleted` VARCHAR(255) DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  `user_creator` VARCHAR(255) DEFAULT NULL,
  
  -- Información del medicamento
  `codigo` VARCHAR(50) NOT NULL UNIQUE COMMENT 'Código interno del medicamento',
  `nombre_generico` VARCHAR(255) NOT NULL,
  `nombre_comercial` VARCHAR(255) DEFAULT NULL,
  `forma_farmaceutica` ENUM('TABLETA', 'CAPSULA', 'JARABE', 'SUSPENSION', 'INYECTABLE', 'CREMA', 'POMADA', 'GEL', 'SOLUCION', 'OTRO') NOT NULL,
  `concentracion` VARCHAR(100) NOT NULL COMMENT 'Ej: 500mg, 250mg/5ml - Solo informativo',
  `unidad_base` VARCHAR(20) DEFAULT 'UNIDAD' COMMENT 'Siempre UNIDAD',
  `via_administracion` VARCHAR(100) DEFAULT NULL COMMENT 'Oral, IV, IM, Tópica, etc.',
  `requiere_receta` TINYINT(1) DEFAULT 0 COMMENT '1=requiere receta médica',
  `controlado` TINYINT(1) DEFAULT 0 COMMENT '1=medicamento controlado',
  `descripcion` TEXT DEFAULT NULL,
  `estado` ENUM('ACTIVO', 'INACTIVO') DEFAULT 'ACTIVO',
  
  PRIMARY KEY (`id`),
  INDEX `idx_codigo` (`codigo`),
  INDEX `idx_nombre_generico` (`nombre_generico`),
  INDEX `idx_estado` (`estado`),
  INDEX `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Catálogo maestro de medicamentos - NO almacena stock';

-- 2. Tabla de lotes (control real de stock)
CREATE TABLE `farmacia_lotes` (
  `id` CHAR(36) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP(),
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `is_deleted` TINYINT(1) DEFAULT 0,
  `user_deleted` VARCHAR(255) DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  `user_creator` VARCHAR(255) DEFAULT NULL,
  
  -- Relación con medicamento
  `medicamento_id` CHAR(36) NOT NULL,
  
  -- Información del lote
  `lote` VARCHAR(100) NOT NULL COMMENT 'Número de lote del fabricante',
  `fecha_vencimiento` DATE NOT NULL,
  `fecha_ingreso` DATE NOT NULL,
  
  -- Control de stock (en unidades)
  `cantidad_total_unidades` INT NOT NULL COMMENT 'Cantidad inicial del lote',
  `cantidad_disponible` INT NOT NULL COMMENT 'Stock actual disponible',
  
  -- Origen del lote
  `origen` ENUM('COMPRA', 'DONACION', 'MINSA', 'TRANSFERENCIA', 'OTRO') NOT NULL,
  `proveedor` VARCHAR(255) DEFAULT NULL,
  `numero_factura` VARCHAR(100) DEFAULT NULL,
  `observaciones` TEXT DEFAULT NULL,
  
  PRIMARY KEY (`id`),
  FOREIGN KEY (`medicamento_id`) REFERENCES `farmacia_medicamentos`(`id`) ON DELETE RESTRICT,
  INDEX `idx_medicamento` (`medicamento_id`),
  INDEX `idx_fecha_vencimiento` (`fecha_vencimiento`),
  INDEX `idx_disponible` (`cantidad_disponible`),
  INDEX `idx_is_deleted` (`is_deleted`),
  UNIQUE KEY `unique_lote_medicamento` (`medicamento_id`, `lote`, `is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Control de stock por lotes con fechas de vencimiento';

-- 3. Tabla de movimientos (trazabilidad completa)
CREATE TABLE `farmacia_movimientos` (
  `id` CHAR(36) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP(),
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `is_deleted` TINYINT(1) DEFAULT 0 COMMENT 'Para anulaciones',
  `user_deleted` VARCHAR(255) DEFAULT NULL,
  `deleted_at` DATETIME DEFAULT NULL,
  `user_creator` VARCHAR(255) DEFAULT NULL,
  
  -- Tipo de movimiento
  `tipo_movimiento` ENUM('ENTRADA', 'SALIDA', 'AJUSTE_POSITIVO', 'AJUSTE_NEGATIVO') NOT NULL,
  
  -- Relaciones
  `medicamento_id` CHAR(36) NOT NULL,
  `lote_id` CHAR(36) NOT NULL,
  `paciente_id` CHAR(36) DEFAULT NULL COMMENT 'NULL para entradas, obligatorio para salidas',
  `consulta_id` CHAR(36) DEFAULT NULL COMMENT 'Futuro: receta digital',
  
  -- Cantidad
  `cantidad_unidades` INT NOT NULL COMMENT 'Siempre en UNIDAD',
  
  -- Autorización (crítico para salidas)
  `autorizado_por_id` VARCHAR(255) DEFAULT NULL COMMENT 'ID del usuario que autoriza',
  `autorizado_por_nombre` VARCHAR(255) DEFAULT NULL COMMENT 'Nombre completo del autorizador',
  `rol_autorizador` ENUM('MEDICO', 'ODONTOLOGIA', 'FARMACIA', 'SISTEMA') DEFAULT NULL,
  `area_origen` ENUM('MEDICINA', 'ODONTOLOGIA', 'FARMACIA', 'PROCEDIMIENTOS', 'SISTEMA') NOT NULL,
  
  -- Información adicional
  `fecha_movimiento` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  `observacion` TEXT DEFAULT NULL,
  `motivo_ajuste` TEXT DEFAULT NULL COMMENT 'Obligatorio para ajustes',
  
  PRIMARY KEY (`id`),
  FOREIGN KEY (`medicamento_id`) REFERENCES `farmacia_medicamentos`(`id`) ON DELETE RESTRICT,
  FOREIGN KEY (`lote_id`) REFERENCES `farmacia_lotes`(`id`) ON DELETE RESTRICT,
  FOREIGN KEY (`paciente_id`) REFERENCES `pacientes`(`id`) ON DELETE RESTRICT,
  INDEX `idx_tipo_movimiento` (`tipo_movimiento`),
  INDEX `idx_paciente` (`paciente_id`),
  INDEX `idx_fecha` (`fecha_movimiento`),
  INDEX `idx_medicamento` (`medicamento_id`),
  INDEX `idx_lote` (`lote_id`),
  INDEX `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Registro de todos los movimientos de medicamentos (entradas/salidas)';

-- 4. Tabla de configuración de alertas
CREATE TABLE `farmacia_alertas_config` (
  `id` CHAR(36) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP(),
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
  `medicamento_id` CHAR(36) NOT NULL,
  
  -- Umbrales
  `stock_minimo` INT NOT NULL DEFAULT 10 COMMENT 'Alerta cuando stock total < este valor',
  `dias_vencimiento_alerta` INT NOT NULL DEFAULT 30 COMMENT 'Alerta cuando faltan X días para vencer',
  
  -- Estado
  `activo` TINYINT(1) DEFAULT 1,
  
  PRIMARY KEY (`id`),
  FOREIGN KEY (`medicamento_id`) REFERENCES `farmacia_medicamentos`(`id`) ON DELETE CASCADE,
  UNIQUE KEY `unique_medicamento_config` (`medicamento_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Configuración de umbrales de alerta por medicamento';

-- ============================================
-- DATOS DE EJEMPLO (OPCIONAL - PARA TESTING)
-- ============================================

-- Insertar medicamentos de ejemplo
INSERT INTO farmacia_medicamentos (id, codigo, nombre_generico, nombre_comercial, forma_farmaceutica, concentracion, unidad_base, via_administracion, requiere_receta, controlado, estado, user_creator)
VALUES 
  (UUID(), 'MED001', 'Paracetamol', 'Tylenol', 'TABLETA', '500mg', 'UNIDAD', 'Oral', 0, 0, 'ACTIVO', 'SISTEMA'),
  (UUID(), 'MED002', 'Ibuprofeno', 'Advil', 'TABLETA', '400mg', 'UNIDAD', 'Oral', 0, 0, 'ACTIVO', 'SISTEMA'),
  (UUID(), 'MED003', 'Amoxicilina', 'Amoxil', 'CAPSULA', '500mg', 'UNIDAD', 'Oral', 1, 0, 'ACTIVO', 'SISTEMA'),
  (UUID(), 'MED004', 'Diclofenaco', NULL, 'INYECTABLE', '75mg/3ml', 'UNIDAD', 'IM', 0, 0, 'ACTIVO', 'SISTEMA'),
  (UUID(), 'MED005', 'Alcohol', NULL, 'SOLUCION', '70%', 'UNIDAD', 'Tópica', 0, 0, 'ACTIVO', 'SISTEMA');

-- ============================================
-- VERIFICACIÓN DE INSTALACIÓN
-- ============================================

-- Verificar tablas creadas
SELECT 
    TABLE_NAME, 
    TABLE_ROWS, 
    CREATE_TIME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME LIKE 'farmacia_%'
ORDER BY TABLE_NAME;

-- Verificar medicamentos insertados
SELECT 
    codigo,
    nombre_generico,
    forma_farmaceutica,
    concentracion,
    estado
FROM farmacia_medicamentos
WHERE is_deleted = 0;
