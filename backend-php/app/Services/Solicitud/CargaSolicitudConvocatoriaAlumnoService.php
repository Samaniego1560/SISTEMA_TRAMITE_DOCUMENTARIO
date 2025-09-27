<?php

namespace App\Services\Solicitud;

use App\Exceptions\ExceptionGenerate;
use App\Models\Alumno;
use App\Models\Convocatoria;
use Illuminate\Support\Facades\DB;

class CargaSolicitudConvocatoriaAlumnoService
{
    public function cargaSolicitudConvocatoriaAlumno(int $convocatoriaId, int $dni): Convocatoria
    {
        $alumno = Alumno::where('DNI', $dni)->first();

        if (!$alumno) {
            throw new ExceptionGenerate('No existe registros del alumno', 200);
        }
        $convocatoria = Convocatoria::find($convocatoriaId);
        $convocatoria->user_id = $alumno->id;
        if (!$convocatoria) {
            throw new ExceptionGenerate('No existe registros del alumno en la actual convocatoria', 200);
        }


        $convocatoria->user_id = $alumno->id;

        // Mapeo de requisitos y valores del alumno
        $alumnoData = [
            'Código estudiante' => $alumno->codigo_estudiante,
            'DNI' => $alumno->DNI,
            'Nombres' => $alumno->nombres,
            'Apellidos' => $alumno->apellido_paterno . ' ' . $alumno->apellido_materno,
            'Sexo' => $alumno->sexo,
            'Facultad' => $alumno->facultad,
            'Escuela profesional' => $alumno->escuela_profesional,
            'Modalidad de ingreso' => $alumno->modalidad_ingreso,
            'Edad' => $alumno->edad,
            'Correo institucional' => $alumno->correo_institucional,
            'Fecha de nacimiento' => $alumno->fecha_nacimiento,
            'Correo personal' => $alumno->correo_personal,
        ];

        $recoverableDetailRequests = [];
        // Asignar valores por defecto a los requisitos del alumno
        foreach ($convocatoria->secciones as $seccion) {
            foreach ($seccion->requisitos as $requisito) {
                if (isset($alumnoData[$requisito->nombre])) {
                    $requisito->default = $alumnoData[$requisito->nombre];
                }

                if ($requisito->is_recoverable) {
                    $recoverableDetailRequest = DB::table('detalle_solicitudes as ds')
                        ->join('solicitudes as sol', 'sol.id', '=', 'ds.solicitud_id')
                        ->join('requisitos as re', 'ds.requisito_id', '=', 're.id')
                        ->join('secciones as se', 're.seccion_id', '=', 'se.id')
                        ->where('sol.alumno_id', $alumno->id)
                        ->where('re.nombre', $requisito->nombre)
                        //->where('re.is_recoverable', true)
                        ->whereIn('ds.solicitud_id', function ($query) {
                            $query->selectRaw('MAX(ds_inner.solicitud_id)')
                                ->from('detalle_solicitudes as ds_inner')
                                ->whereColumn('ds_inner.requisito_id', 'ds.requisito_id')
                                ->groupBy('ds_inner.requisito_id');
                        })
                        ->select(
                            DB::raw("
                     CASE 
                         WHEN re.tipo_requisito_id = 4 THEN ds.opcion_seleccion 
                         WHEN re.tipo_requisito_id IN (1, 2) THEN ds.url_documento 
                         ELSE ds.respuesta_formulario 
                     END as respuesta_formulario
                 "),
                            'ds.requisito_id',
                            'ds.solicitud_id',
                            'ds.order'
                        )
                        ->orderBy('ds.order')
                        ->get();

                    $recoverableDetailRequest->transform(function ($item) use ($requisito) {
                        $item->requisito_id = $requisito->id; // Sobrescribir o agregar el ID del requisito actual
                        return $item;
                    });

                    $recoverableDetailRequests = array_merge($recoverableDetailRequests, $recoverableDetailRequest->toArray());
                }
            }
        }
        $convocatoria->recoverable_detail_requests = $recoverableDetailRequests;


        return $convocatoria;
    }
}
