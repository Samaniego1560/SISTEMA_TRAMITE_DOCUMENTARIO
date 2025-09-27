<?php 

namespace App\Exports;

use App\Exceptions\ExceptionGenerate;
use App\Models\Convocatoria;
use App\Models\DetalleSolicitud;
use App\Models\Requisito;
use App\Models\Servicio;
use App\Models\ServicioSolicitado;
use App\Models\Solicitud;
use Illuminate\Support\Facades\DB;
use Illuminate\Database\Eloquent\Collection;

class SolicitudExportDataPorServicioEstadoConvocatoria
{
    private string $servicio;
    private string $estado;
    private int $convocatoriaId;

    public function __construct(string $servicio, string $estado, int $convocatoriaId)
    {
        $this->servicio = $servicio;
        $this->estado = $estado;
        $this->convocatoriaId = $convocatoriaId;
    }

    public function collection(): Collection
    {
        $convocatoria = Convocatoria::find($this->convocatoriaId);

        // Obtener las solicitudes según la convocatoria
        $solicitudesIDsTemporales = Solicitud::select('solicitudes.id as id')
            ->join('alumnos as a', 'solicitudes.alumno_id', '=', 'a.id')
            ->where('solicitudes.convocatoria_id', $convocatoria->id)
            ->orderBy('solicitudes.created_at', 'desc')
            ->get();

        $solicitudesTemporales = Solicitud::select(
            'a.codigo_estudiante',
            'solicitudes.created_at as fecha_solicitud',
            'a.DNI',
            DB::raw('concat(a.apellido_paterno, \' \', a.apellido_materno, \', \', a.nombres) as alumno'),
            'a.sexo',
            'a.facultad',
            'a.escuela_profesional',
            'a.ultimo_semestre',
            'a.estado_matricula',
            'a.creditos_matriculados',
            'a.num_semestres_cursados',
            'a.pps',
            'a.ppa',
            'a.tca',
            'a.fecha_nacimiento',
            'a.edad',
            'a.lugar_nacimiento',
            'a.lugar_procedencia',
            'a.direccion',
            'a.celular_padre',
            'a.celular_estudiante',
            'a.correo_institucional',
            'a.modalidad_ingreso',
            'a.correo_personal'
        )
            ->join('alumnos as a', 'solicitudes.alumno_id', '=', 'a.id')
            ->where('solicitudes.convocatoria_id', $convocatoria->id)
            ->orderBy('solicitudes.created_at', 'desc')
            ->get();

        $solicitudesIDs = new Collection();
        $solicitudes = new Collection();

        // Agrupar las solicitudes únicas
        for ($i = 0; $i < count($solicitudesTemporales); $i++) {
            if (!$solicitudes->contains('codigo_estudiante', $solicitudesTemporales[$i]->codigo_estudiante)) {
                $solicitudes->push($solicitudesTemporales[$i]);
                $solicitudesIDs->push($solicitudesIDsTemporales[$i]);
            }
        }

        // Recorrer las solicitudes
        for ($i = 0; $i < count($solicitudes); $i++) {
            // Obtener detalles de la solicitud
            $detalle_solicitud = DetalleSolicitud::select('detalle_solicitudes.*')
                ->join('solicitudes as s', 'detalle_solicitudes.solicitud_id', '=', 's.id')
                ->join('alumnos as a', 's.alumno_id', '=', 'a.id')
                ->where([['s.id', $solicitudesIDs[$i]->id], ['s.convocatoria_id', $convocatoria->id]])
                ->get();

            // Recorrer los requisitos de la solicitud
            for ($j = 22; $j < count($detalle_solicitud); $j++) {
                $requisito = Requisito::select('tipo_requisito_id', 'nombre', 'id')
                    ->where('id', $detalle_solicitud[$j]->requisito_id)
                    ->first();

                $dato = "";
                if ($requisito->tipo_requisito_id == 1 || $requisito->tipo_requisito_id == 2 || $requisito->tipo_requisito_id == 8) {
                    $dato = $detalle_solicitud[$j]->url_documento;
                }
                if ($requisito->tipo_requisito_id == 3) {
                    $dato = $detalle_solicitud[$j]->respuesta_formulario;
                }
                if ($requisito->tipo_requisito_id == 4) {
                    $dato = $detalle_solicitud[$j]->opcion_seleccion;
                }
                if ($requisito->tipo_requisito_id == 5 || $requisito->tipo_requisito_id == 6 || $requisito->tipo_requisito_id == 7) {
                    $dato = $detalle_solicitud[$j]->respuesta_formulario;
                }

                // Asignar el dato al requisito correspondiente
                $solicitudes[$i]->{$requisito->nombre} = $dato;
            }

            // Obtener los servicios solicitados
            $servicios = Servicio::all();
            for ($k = 0; $k < count($servicios); $k++) {
                try {
                    $servicio_solicitado = ServicioSolicitado::select(
                        'servicio_solicitado.id',
                        'servicio_solicitado.estado'
                    )
                        ->join('solicitudes as s', 'servicio_solicitado.solicitud_id', '=', 's.id')
                        ->join('alumnos as a', 's.alumno_id', '=', 'a.id')
                        ->where([
                            ['servicio_solicitado.servicio_id', $servicios[$k]->id],
                            ['servicio_solicitado.solicitud_id', $solicitudesIDs[$i]->id]
                        ])
                        ->first();

                    $estado = "";
                    if ($servicio_solicitado && $servicio_solicitado->estado) {
                        $estado = $servicio_solicitado->estado;
                    }

                    // Comparación en minúsculas para que no importe si el servicio se pasa como "Comedor" o "comedor"
                    if (strtolower($servicios[$k]->nombre) == strtolower($this->servicio) && $estado == $this->estado) {
                        $solicitudes[$i]->{$servicios[$k]->nombre} = $estado;
                    }
                } catch (ExceptionGenerate $e) {
                    $solicitudes[$i]->{$i . '-' . $k} = '';
                }
            }
        }

        // Filtrar las solicitudes solo con el servicio específico y el estado "aprobado"
        $solicitudes_filtradas = $solicitudes->filter(function ($solicitud) {
            return isset($solicitud->{$this->servicio}) && $solicitud->{$this->servicio} == $this->estado;
        });

        return $solicitudes_filtradas;
    }
}
