<?php

namespace App\Services\Report;

use App\Models\Alumno;
use App\Models\Solicitud;
use App\Models\Convocatoria;
use App\Models\DetalleSolicitud;
use App\Models\Servicio;
use App\Models\Departament;
use App\Models\Province;
use App\Models\District;
use Illuminate\Support\Collection;
use Barryvdh\DomPDF\Facade\Pdf;

class UserProfilePdfService
{

    private function findPssFromSections($detalleSolicitudes)
    {
        foreach ($detalleSolicitudes as $seccion) {
            foreach ($seccion->requisitos as $requisito) {
                $nombre = $requisito->nombre ?? ($requisito['nombre'] ?? '');
                if (mb_strtolower(trim($nombre)) === 'pss' || mb_strtolower(trim($nombre)) === 'ponderado semestral (pss)') {
                    $respuesta = $requisito->respuesta ?? ($requisito['respuesta'] ?? null);
                    if (!$respuesta) {
                        return null;
                    }
                    if (is_array($respuesta)) {
                        return $respuesta['respuesta_formulario'] ?? $respuesta['opcion_seleccion'] ?? $respuesta['url_documento'] ?? null;
                    }
                    return $respuesta->respuesta_formulario ?? $respuesta->opcion_seleccion ?? $respuesta->url_documento ?? null;
                }
            }
        }
        return null;
    }

    public function getDataReport($dni)
    {
        try {
            $data = $this->getData($dni);
            $primera_seccion = $this->processFirstSection($data);
            $linkProfile = $this->getLinkPhotoProfile($data);
            $convocatorias = $this->processAnnouncements($data);

            return [
                'title' => 'Ficha del Estudiante',
                'link_profile' => $linkProfile,
                'announcements' => $convocatorias,
                'first_section' => $primera_seccion,
                'year' => '2025'
            ];
        } catch (\Exception $e) {
            error_log("Error: " . $e->getMessage());
            return response()->json(['error' => 'Error al generar el reporte: ' . $e->getMessage()], 500);
        }
    }

    public function generatePdf($dni)
    {
        try {
            $data = $this->getData($dni);

            $primera_seccion = $this->processFirstSection($data);
            $linkProfile = $this->getLinkPhotoProfile($data);
            $convocatorias = $this->processAnnouncements($data);

            $pdf = PDF::loadView('profile_template', [
                'title' => 'Ficha del Estudiante',
                'link_profile' => $linkProfile,
                'announcements' => $convocatorias,
                'first_section' => $primera_seccion,
                'year' => '2024'
            ]);

            return $pdf->download('ficha_estudiante.pdf');
        } catch (\Exception $e) {
            error_log("Error: " . $e->getMessage());
            return response()->json(['error' => 'Error al generar el PDF: ' . $e->getMessage()], 500);
        }
    }

    private function processFirstSection($data): array
    {
        $primera_seccion = [];
        foreach ($data as $item) {

            foreach ($item['detalle_solicitudes'] as $detalle) {
                if ($detalle['descripcion'] == 'Datos Personales') {
                    $primera_seccion['description'] = $detalle['descripcion'];
                    $primera_seccion['type'] = $detalle['type'];
                    $primera_seccion['requirements'] = $this->processRequirements($detalle['requisitos']);
                    // Sobrescribir con nombres y apellidos completos si están disponibles
                    foreach ($primera_seccion['requirements'] as &$req) {
                        $nombreCampo = isset($req['name']) ? mb_strtolower(trim($req['name'])) : '';
                        if ($nombreCampo === 'nombres' && !empty($item['alumno_nombres'] ?? '')) {
                            $req['value'] = $item['alumno_nombres'];
                        }
                        if ($nombreCampo === 'apellidos' && !empty(trim($item['alumno_apellidos'] ?? ''))) {
                            $req['value'] = $item['alumno_apellidos'];
                        }
                    }
                    unset($req);
                    break;
                }
            }
        }
        return $primera_seccion;
    }

    private function getLinkPhotoProfile($data)
    {
        $primera_seccion = [];
        foreach ($data as $item) {
            foreach ($item['detalle_solicitudes'] as $detalle) {
                if ($detalle['descripcion'] == 'Datos Personales') {
                    foreach ($detalle['requisitos'] as $requisito) {
                        if ($requisito['nombre'] == 'photo-profile') {
                            if ($requisito['tipo_requisito_id'] == 1 || $requisito['tipo_requisito_id'] == 2) {
                                return $requisito->respuesta->url_documento ?? '';
                            }
                        }
                    }
                }
            }
        }
        return '';
    }

    private function processAnnouncements($data): array
    {
        $convocatorias = [];
        foreach ($data as $item) {
            $convocatoria = [
                'name' => $item['announcement_name'],
                'details_requests' => $this->processDetailsRequests($item['detalle_solicitudes']),
                'pps' => $item['pps'] ?? null,
                'solicitud_status' => $item['solicitud_status'] ?? null,
            ];
            $convocatorias[] = $convocatoria;
        }
        return $convocatorias;
    }

    private function processDetailsRequests($detalleSolicitudes): array
    {
        $detalles = [];
        foreach ($detalleSolicitudes as $detalle) {
            if ($detalle['descripcion'] != 'Datos Personales') {
                $reqs = $this->processRequirements($detalle['requisitos']);
                if (($detalle['type'] ?? '') === 'table') {
                    // Agrupar por 'order' para secciones tipo tabla (p.ej., Composición familiar)
                    $grouped = [];
                    foreach ($reqs as $r) {
                        $ord = $r['order'] ?? 1;
                        if (!isset($grouped[$ord])) $grouped[$ord] = [];
                        $grouped[$ord][] = $r;
                    }
                    ksort($grouped);
                    $rows = array_values($grouped);
                    $detalles[] = [
                        'description' => $detalle['descripcion'],
                        'type' => $detalle['type'],
                        // Compatibilidad: mantener requirements plano para la vista existente
                        'requirements' => $reqs,
                        // Para el PDF: filas agrupadas por integrante
                        'is_table' => true,
                        'rows' => $rows,
                    ];
                } else {
                    $detalles[] = [
                        'description' => $detalle['descripcion'],
                        'type' => $detalle['type'],
                        'requirements' => $reqs,
                    ];
                }
            }
        }
        return $detalles;
    }

    private function processRequirements($requisitos): array
    {
        $procesados = [];
        foreach ($requisitos as $requisito) {
            $nombre = $requisito['nombre'] ?? ($requisito->nombre ?? '');
            $nombreNormalizado = mb_strtolower(trim($nombre));
            if ($nombreNormalizado === 'pss' || $nombreNormalizado === 'ponderado semestral (pss)') {
                $respuesta = $requisito['respuesta'] ?? ($requisito->respuesta ?? null);
                $respForm = '';
                $respSel = '';
                $respUrl = '';
                if (is_array($respuesta)) {
                    $respForm = $respuesta['respuesta_formulario'] ?? '';
                    $respSel = $respuesta['opcion_seleccion'] ?? '';
                    $respUrl = $respuesta['url_documento'] ?? '';
                } elseif (is_object($respuesta)) {
                    $respForm = $respuesta->respuesta_formulario ?? '';
                    $respSel = $respuesta->opcion_seleccion ?? '';
                    $respUrl = $respuesta->url_documento ?? '';
                }
                if (trim((string)$respForm) === '' && trim((string)$respSel) === '' && trim((string)$respUrl) === '') {
                    continue;
                }
            }
            if ($requisito['nombre'] == 'photo-profile') {
                continue;
            }

            $valor_respuesta = $requisito['respuesta']['respuesta_formulario'] ?? 'No proporcionado';
            if ($requisito['tipo_requisito_id'] == 1 || $requisito['tipo_requisito_id'] == 2) {
                $valor_respuesta = $requisito->respuesta->url_documento ?? 'No proporcionado';
            }

            if ($requisito['tipo_requisito_id'] == 4) {
                $valor_respuesta = $requisito->respuesta->opcion_seleccion ?? 'No proporcionado';
                // Resolver IDs a descripciones para catálogos geográficos
                $opciones = $requisito['opciones'] ?? ($requisito->opciones ?? null);
                $seleccion = $requisito->respuesta->opcion_seleccion ?? null;
                if ($seleccion && is_string($opciones)) {
                    if ($opciones === 'department') {
                        $valor_respuesta = Departament::where('id', $seleccion)->value('name') ?? $seleccion;
                    } elseif ($opciones === 'province') {
                        $valor_respuesta = Province::where('id', $seleccion)->value('name') ?? $seleccion;
                    } elseif ($opciones === 'district') {
                        $valor_respuesta = District::where('id', $seleccion)->value('name') ?? $seleccion;
                    }
                }
            }
            $order_respuesta = $requisito['respuesta']['order'] ?? 1;
            $procesados[] = [
                'name' => $requisito['nombre'],
                'type' => $requisito['tipo_requisito_id'],
                'value' => $valor_respuesta,
                'order' => $order_respuesta,
            ];
        }
        return $procesados;
    }

    public function getData($dni): array
    {
        // Obtener el último registro del alumno por DNI
        $alumno = Alumno::where('DNI', $dni)->orderBy('created_at', 'desc')->first();

        if (!$alumno) {
            return []; // Si no se encuentra el alumno, devolver un array vacío
        }

        $solicitudes = [];

        // Obtener todas las convocatorias
        $convocatorias = Convocatoria::all();

        foreach ($convocatorias as $convocatoria) {

            // Buscar el registro del alumno correspondiente a esta convocatoria (si existe)
            $alumnoConv = Alumno::where('DNI', $dni)
                ->where('convocatoria_id', $convocatoria->id)
                ->orderBy('created_at', 'desc')
                ->first();

            $alumnoUsado = $alumnoConv ?: $alumno;

            // Obtener las solicitudes del alumno para esta convocatoria
            $solicitudAlumno = Solicitud::where('alumno_id', $alumnoUsado->id)
                ->where('convocatoria_id', $convocatoria->id)
                ->orderBy('created_at', 'desc') // Ordena por "id" en orden descendente
                ->first();
            $this->cargarSecciones($convocatoria);
            if (!$convocatoria->secciones) {
                continue;
            }


            if ($solicitudAlumno) {

                $solicitudAlumno->detalle_solicitudes = $convocatoria->secciones;

                // Obtener los detalles de la solicitud
                $solicitudDetalle = DetalleSolicitud::where('solicitud_id', $solicitudAlumno->id)->get();
                // Mapear los requisitos al detalle de solicitudes
                $this->mapearRequisitos($solicitudAlumno->detalle_solicitudes, $solicitudDetalle);

                // Agregar el nombre de la convocatoria a la solicitud
                $solicitudAlumno['announcement_name'] = $convocatoria->nombre;
                // Adjuntar el PSS (Promedio Ponderado Semestral) del alumno para esta convocatoria
                $solicitudAlumno['pps'] = $this->findPssFromSections($solicitudAlumno->detalle_solicitudes);
                // Adjuntar nombres y apellidos completos desde Alumno
                $solicitudAlumno['alumno_nombres'] = $alumnoUsado->nombres ?? '';
                $solicitudAlumno['alumno_apellidos'] = trim(($alumnoUsado->apellido_paterno ?? '') . ' ' . ($alumnoUsado->apellido_materno ?? ''));
                // Determinar el estado general de la solicitud para esta convocatoria
                $solicitudAlumno['solicitud_status'] = $this->resolveSolicitudStatus($solicitudAlumno->id);
                $solicitudes[] = $solicitudAlumno;
            }
        }

        return $solicitudes;
    }

    private function resolveSolicitudStatus(int $solicitudId): string
    {
        // Obtiene todos los servicios solicitados y calcula un estado global simple
        $estados = \App\Models\ServicioSolicitado::where('solicitud_id', $solicitudId)
            ->pluck('estado')
            ->filter()
            ->map(function ($e) { return mb_strtolower(trim($e)); })
            ->values();

        if ($estados->isEmpty()) {
            return 'Pendiente';
        }

        // Prioridad: Aceptado > Rechazado > Pendiente
        $aceptados = ['aceptado', 'aprobado', 'aprobada', 'aceptada'];
        $rechazados = ['rechazado', 'denegado', 'denegada', 'rechazada'];
        $pendientes = ['pendiente', 'en proceso', 'proceso', 'revision', 'revisión'];

        if ($estados->first(function ($e) use ($aceptados) { return in_array($e, $aceptados, true); })) {
            return 'Aceptado';
        }
        if ($estados->first(function ($e) use ($rechazados) { return in_array($e, $rechazados, true); })) {
            return 'Rechazado';
        }
        return 'Pendiente';
    }

    private function cargarSecciones($convocatoria)
    {
        if (isset($convocatoria->secciones[0])) {
            $convocatoria->secciones[0]->requisitos;
        }
        if (isset($convocatoria->secciones[1])) {
            $convocatoria->secciones[1]->requisitos;
        }
        if (isset($convocatoria->secciones[2])) {
            $convocatoria->secciones[2]->requisitos;
        }
    }

    private function mapearRequisitos($detalleSolicitudes, $solicitudDetalle)
    {
        foreach ($detalleSolicitudes as $seccion) {
            foreach ($seccion->requisitos as $requisito) {
                $responses = $this->getResponseDetailByRequirementId($seccion->type, $requisito->id, $solicitudDetalle);
                if (count($responses) > 0) {
                    $requisito->respuesta = $responses[0];
                    if (count($responses) > 1) {
                        foreach (array_slice($responses, 1) as $resp) {
                            $nuevo_requisito = clone $requisito; // Clonar el objeto para evitar referencias
                            $nuevo_requisito->respuesta = $resp;
                            $seccion->requisitos[] = $nuevo_requisito; // Añadir al array de requisitos
                        }
                    }
                }
            }
        }
    }

    private function getResponseDetailByRequirementId($typeSection, $requisitoId, $solicitudDetalle)
    {
        $response = [];

        foreach ($solicitudDetalle as $detalle) {
            if ($typeSection === "form" && $requisitoId === $detalle->requisito_id) {
                $response[] = $detalle; // Agregar detalle al array
                break; // Solo interesa el primer resultado en este caso
            }

            if ($typeSection === "table" && $requisitoId === $detalle->requisito_id) {
                $response[] = $detalle; // Agregar detalle al array
            }
        }

        return $response;
    }
}
