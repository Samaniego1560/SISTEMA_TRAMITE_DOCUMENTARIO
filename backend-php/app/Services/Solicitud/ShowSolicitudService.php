<?php

namespace App\Services\Solicitud;

use App\Models\Convocatoria;
use App\Models\DetalleSolicitud;
use App\Models\Solicitud;
use Illuminate\Database\Eloquent\Model;

class ShowSolicitudService
{
    public function show($id)/* : ?Model */
    {
        $solicitud = Solicitud::find($id);

        if (!$solicitud) {
            return null; // En caso de que no exista la solicitud
        }

        $convocatoria = Convocatoria::find($solicitud->convocatoria_id);

        if (!$convocatoria) {
            return $solicitud; // Retorna solo la solicitud si no encuentra convocatoria
        }

        // Aseguramos que se carguen los requisitos de cada sección dinámicamente
        foreach ($convocatoria->secciones as $seccion) {
            $seccion->requisitos;
        }

        // Guardamos las secciones dentro de la solicitud
        $solicitud->detalle_solicitudes = $convocatoria->secciones;

        // Obtenemos los detalles de la solicitud
        $solicitud_detalle = DetalleSolicitud::where('solicitud_id', $solicitud->id)->get();

        // Vinculamos requisitos con sus respuestas
        foreach ($solicitud->detalle_solicitudes as $seccion) {
            foreach ($seccion->requisitos as $requisito) {
                foreach ($solicitud_detalle as $detalle) {
                    if ($requisito->id == $detalle->requisito_id) {
                        $requisito->respuesta = $detalle;
                        break; // cuando encuentra la coincidencia, ya no sigue buscando
                    }
                }
            }
        }

        return $solicitud;
    }
}
