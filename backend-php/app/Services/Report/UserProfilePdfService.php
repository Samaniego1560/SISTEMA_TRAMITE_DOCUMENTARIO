<?php

namespace App\Services\Report;

use App\Models\Alumno;
use App\Models\Solicitud;
use App\Models\Convocatoria;
use App\Models\DetalleSolicitud;
use App\Models\Servicio;
use Illuminate\Support\Collection;
use Barryvdh\DomPDF\Facade\Pdf;

class UserProfilePdfService
{

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
                'year' => '2024'
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
                'details_requests' => $this->processDetailsRequests($item['detalle_solicitudes'])
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
                $detalles[] = [
                    'description' => $detalle['descripcion'],
                    'type' => $detalle['type'],
                    'requirements' => $this->processRequirements($detalle['requisitos'])
                ];
            }
        }
        return $detalles;
    }

    private function processRequirements($requisitos): array
    {
        $procesados = [];
        foreach ($requisitos as $requisito) {
            if ($requisito['nombre'] == 'photo-profile') {
                continue;
            }

            $valor_respuesta = $requisito['respuesta']['respuesta_formulario'] ?? 'No proporcionado';
            if ($requisito['tipo_requisito_id'] == 1 || $requisito['tipo_requisito_id'] == 2) {
                $valor_respuesta = $requisito->respuesta->url_documento ?? 'No proporcionado';
            }

            if ($requisito['tipo_requisito_id'] == 4) {
                $valor_respuesta = $requisito->respuesta->opcion_seleccion ?? 'No proporcionado';
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

            // Obtener las solicitudes del alumno para esta convocatoria
            $solicitudAlumno = Solicitud::where('alumno_id', $alumno->id)
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
                $solicitudes[] = $solicitudAlumno;
            }
        }

        return $solicitudes;
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
