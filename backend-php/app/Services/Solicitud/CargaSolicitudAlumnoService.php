<?php

namespace App\Services\Solicitud;

use App\Exceptions\ExceptionGenerate;
use App\Models\Alumno;
use App\Models\Convocatoria;
use Illuminate\Support\Facades\DB;

class CargaSolicitudAlumnoService
{
    public function cargaSolicitudAlumno(int $dni): Convocatoria
{
    $alumno = Alumno::where('DNI', $dni)->first();

    if (!$alumno) {
        throw new ExceptionGenerate('No existe registros del alumno', 200);
    }

    $convocatoria = Convocatoria::find($alumno->convocatoria_id);
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

    // Asignar valores por defecto a los requisitos del alumno
    foreach ($convocatoria->secciones as $seccion) {
        foreach ($seccion->requisitos as $requisito) {
            if (isset($alumnoData[$requisito->nombre])) {
                $requisito->default = $alumnoData[$requisito->nombre];
            }
        }
    }

    // Obtener los detalles de las solicitudes recuperables en una única consulta
$recoverableDetailRequests = DB::table('detalle_solicitudes as ds')
->join('solicitudes as sol', 'sol.id', '=', 'ds.solicitud_id')
->join('requisitos as re', 'ds.requisito_id', '=', 're.id')
->join('secciones as se', 're.seccion_id', '=', 'se.id')
->where('sol.alumno_id', $alumno->id)
->where('re.is_recoverable', true)
->where('se.convocatoria_id', $convocatoria->id)
->select(
    'ds.respuesta_formulario',
    'ds.requisito_id',
    'ds.solicitud_id',
    'ds.order'
)
->groupBy('ds.solicitud_id', 'ds.requisito_id', 'ds.order', 'ds.respuesta_formulario') // Agrupamos los datos relevantes
->orderBy('ds.order') // Ordenamos si es necesario
->get();

// Asignar los datos recuperables al objeto convocatoria
$convocatoria->recoverable_detail_requests = $recoverableDetailRequests->toArray();


    return $convocatoria;
}





    public function cargaSolicitudAlumno_v1(int $dni): Convocatoria
    {
        $alumno = Alumno::where('DNI', $dni)->first();
        if (!$alumno)
            throw new ExceptionGenerate('No existe registros del alumno', 200);
        $convocatoria = Convocatoria::find($alumno->convocatoria_id);
        $convocatoria->user_id = $alumno->id;
        if (!$convocatoria)
            throw new ExceptionGenerate('No existe registros del alumno en la actual convocatoria', 200);
        //Cargando datos obtenidos del alumno en la solicitud
        //Nombre
        for ($i = 0; $i < count($convocatoria->secciones[0]->requisitos); $i++) {
            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Código estudiante')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->codigo_estudiante;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'DNI')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->DNI;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Nombres')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->nombres;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Apellidos')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->apellido_paterno . ' ' . $alumno->apellido_materno;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Sexo')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->sexo;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Facultad')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->facultad;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Escuela profesional')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->escuela_profesional;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Modalidad de ingreso')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->modalidad_ingreso;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Edad')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->edad;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Correo institucional')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->correo_institucional;

            /* if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Dirección')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->direccion; */

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Feha de nacimiento')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->fecha_nacimiento;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Correo personal')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->correo_personal;

            /* if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'celular de estudiante')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->celular_estudiante;

            if ($convocatoria->secciones[0]->requisitos[$i]->nombre == 'Celular padre')
                $convocatoria->secciones[0]->requisitos[$i]->default = $alumno->celular_padre; */
        }

        $requirementsRecoverable = DB::table('requisitos as re')
    ->join('secciones as se', 're.seccion_id', '=', 'se.id')
    ->where('se.convocatoria_id', $convocatoria->id)
    ->select('re.nombre') // Solo seleccionamos lo necesario
    ->get();

$recoverableDetailRequests = [];

foreach ($requirementsRecoverable as $requirement) {
    // Obtener los IDs de los requisitos con el mismo nombre
    $requirementsIds = DB::table('requisitos')
        ->where('nombre', $requirement->nombre)
        ->pluck('id'); // Usamos `pluck` para obtener un array de IDs directamente

    if ($requirementsIds->isEmpty()) {
        continue; // Si no hay IDs, pasamos al siguiente
    }

    // Obtener el ID máximo de solicitud
    $detalleSolicitudIdMax = DB::table('detalle_solicitudes as ds')
        ->join('solicitudes as sol', 'sol.id', '=', 'ds.solicitud_id')
        ->where('sol.alumno_id', $alumno->id)
        ->whereIn('ds.requisito_id', $requirementsIds)
        ->max('ds.solicitud_id'); // Usamos `max` directamente para obtener el valor máximo

    if (!$detalleSolicitudIdMax) {
        continue; // Si no hay solicitudes, pasamos al siguiente
    }

    // Obtener los detalles de la solicitud recuperable
    $detalleSolicitud = DB::table('detalle_solicitudes as ds')
        ->where('ds.solicitud_id', $detalleSolicitudIdMax)
        ->whereIn('ds.requisito_id', $requirementsIds)
        ->select('ds.respuesta_formulario', 'ds.requisito_id', 'ds.solicitud_id', 'ds.order')
        ->get();

    // Agregar los detalles al array de resultados
    $recoverableDetailRequests = array_merge($recoverableDetailRequests, $detalleSolicitud->toArray());
}

// Asignar los datos al objeto convocatoria
$convocatoria->recoverable_detail_requests = $recoverableDetailRequests;

        

        /*  if ($alumno->lugar_nacimiento != null && strlen($alumno->lugar_nacimiento) > 0) {
            $datosNacimiento = explode("/", $alumno->lugar_nacimiento);
            if ($convocatoria->secciones[1]->requisitos[0]->nombre == 'Departamento de nacimiento')
                $convocatoria->secciones[1]->requisitos[0]->default = $datosNacimiento[0];
            if ($convocatoria->secciones[1]->requisitos[1]->nombre == 'Provincia de nacimiento')
                $convocatoria->secciones[1]->requisitos[1]->default = $datosNacimiento[1];
            if ($convocatoria->secciones[1]->requisitos[2]->nombre == 'Distrito de nacimiento')
                $convocatoria->secciones[1]->requisitos[2]->default = $datosNacimiento[2];
        }

        if ($alumno->lugar_procedencia != null && strlen($alumno->lugar_procedencia) > 0) {
            $datosProcedencia = explode("/", $alumno->lugar_procedencia);
            if ($convocatoria->secciones[2]->requisitos[0]->nombre == 'Departamento de procedencia')
                $convocatoria->secciones[2]->requisitos[0]->default = $datosProcedencia[0];
            if ($convocatoria->secciones[2]->requisitos[1]->nombre == 'Provincia de procedencia')
                $convocatoria->secciones[2]->requisitos[1]->default = $datosProcedencia[1];
            if ($convocatoria->secciones[2]->requisitos[2]->nombre == 'Distrito de procedencia')
                $convocatoria->secciones[2]->requisitos[2]->default = $datosProcedencia[2];
        } */
        for ($i = 0; $i < count($convocatoria->secciones); $i++) {
            $convocatoria->secciones[$i]->requisitos[1];
        }
    
        return $convocatoria;
    }

}
