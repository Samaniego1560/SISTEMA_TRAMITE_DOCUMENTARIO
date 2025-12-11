-- Programar evento diario para resolver sanciones vencidas automáticamente.
-- Ejecuta la actualización cada día a las 03:00 (hora del servidor, recomendado America/Lima).
DROP EVENT IF EXISTS resolve_expired_sanctions;

DELIMITER $$
CREATE EVENT resolve_expired_sanctions
    ON SCHEDULE EVERY 1 DAY
    STARTS (
        TIMESTAMP(CURRENT_DATE + INTERVAL 1 DAY, '03:00:00')
    )
    DO
BEGIN
    -- Fecha/hora de referencia para las actualizaciones.
    SET @ahora := NOW();

    -- Actualiza sanciones asignadas y faltas cuyo periodo ya terminó
    -- y que no estén apeladas ni marcadas previamente como resueltas.
    UPDATE sancion_falta_asignada AS sfa
    INNER JOIN faltas AS f ON f.id = sfa.falta_id
    SET
        sfa.estado = 'RESUELTA',
        sfa.updated_at = @ahora,
        f.estado = 'resuelta',
        f.updated_at = @ahora
    WHERE sfa.fecha_fin IS NOT NULL
      AND sfa.fecha_fin < @ahora
      AND UPPER(COALESCE(sfa.estado, '')) NOT IN ('RESUELTA', 'APELADA')
      AND UPPER(COALESCE(f.estado, '')) NOT IN ('RESUELTA', 'APELADA');
END$$
DELIMITER ;
