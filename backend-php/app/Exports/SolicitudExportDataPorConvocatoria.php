<?php

namespace App\Exports;
use App\Exceptions\ExceptionGenerate;
use App\Models\Convocatoria;
use App\Models\DetalleSolicitud;
use App\Models\Requisito;
use App\Models\Servicio;
use App\Models\ServicioSolicitado;
use App\Models\Solicitud;
use App\Models\Departament;
use App\Models\Province;
use App\Models\District;
use Illuminate\Support\Facades\DB;
use Illuminate\Database\Eloquent\Collection;

class SolicitudExportDataPorConvocatoria
{
    private int $convocatoriaId;

    public function __construct(int $convocatoriaId)
    {
        $this->convocatoriaId = $convocatoriaId;
    }

    public function collection(): Collection
    {
        $convocatoria = Convocatoria::find($this->convocatoriaId);
        //$convocatoria = $convocatoria[(count($convocatoria) - 1)];
        $solicitudesIDsTemporales = Solicitud::select(
            'solicitudes.id as id',
        )
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
            'a.direccion',
            'a.celular_padre',
            'a.celular_estudiante',
            'a.correo_institucional',
            'a.modalidad_ingreso',
            'a.correo_personal',
        )
            ->join('alumnos as a', 'solicitudes.alumno_id', '=', 'a.id')
            ->where('solicitudes.convocatoria_id', $convocatoria->id)
            ->orderBy('solicitudes.created_at', 'desc')
            ->get();


        $solicitudesIDs = new Collection();
        $solicitudes = new Collection();
        for ($i = 0; $i < count($solicitudesTemporales); $i++) {
            if (!$solicitudes->contains('codigo_estudiante', $solicitudesTemporales[$i]->codigo_estudiante)) {
                $solicitudes->push($solicitudesTemporales[$i]);
                $solicitudesIDs->push($solicitudesIDsTemporales[$i]);
            }
        }

        for ($i = 0; $i < count($solicitudes); $i++) {
            //recorrer los requisitos de tipo 3
            $detalle_solicitud = DetalleSolicitud::select(
                'detalle_solicitudes.*',
            )
                ->join('solicitudes as s', 'detalle_solicitudes.solicitud_id', '=', 's.id')
                ->join('alumnos as a', 's.alumno_id', '=', 'a.id')
                ->where([
                    ['s.id', $solicitudesIDs[$i]->id],
                    ['s.convocatoria_id', $convocatoria->id]
                ])
                ->get();
            for ($j = 22; $j < count($detalle_solicitud); $j++) {
                $requisito = Requisito::select(
                    'tipo_requisito_id',
		    'opciones',
                    'nombre',
                    'id'
                )
                    ->where('id', $detalle_solicitud[$j]->requisito_id)->first();
                $dato = "";
                if ($requisito->tipo_requisito_id == 1 || $requisito->tipo_requisito_id == 2 || $requisito->tipo_requisito_id == 8) {
                    $dato = $detalle_solicitud[$j]->url_documento;
                }
                if ($requisito->tipo_requisito_id == 3) {
                    $dato = $detalle_solicitud[$j]->respuesta_formulario;
                }

                if ($requisito->tipo_requisito_id == 4) {
                    $opciones = $requisito->opciones;

                    if ($opciones == 'department') {
                        $dato = Departament::where('id', $detalle_solicitud[$j]->opcion_seleccion)->value('name');
                    }
                
                    if ($opciones == 'province') {
                        $dato = Province::where('id', $detalle_solicitud[$j]->opcion_seleccion)->value('name');
                    }
                
                    if ($opciones == 'district') {
                        $dato = District::where('id', $detalle_solicitud[$j]->opcion_seleccion)->value('name');
                    }  

                    if ($dato == null) {
                        $dato = $detalle_solicitud[$j]->opcion_seleccion;
                    }
                }

                if ($requisito->tipo_requisito_id == 5) {
                    $dato = $detalle_solicitud[$j]->respuesta_formulario;
                }
                if ($requisito->tipo_requisito_id == 6) {
                    $dato = $detalle_solicitud[$j]->respuesta_formulario;
                }
                if ($requisito->tipo_requisito_id == 7) {
                    $dato = $detalle_solicitud[$j]->respuesta_formulario;
                }
                $solicitudes[$i]->{$requisito->nombre} = $dato;
            }

            
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
                        $solicitudes[$i]->{$servicios[$k]->nombre} = $estado;
                        
                } catch (ExceptionGenerate $e) {
                    $solicitudes[$i]->{$i . '-' . $k} = '';
                }
            }
        }
        return $solicitudes;
    }
}
