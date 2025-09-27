<?php

namespace App\Services\Convocatoria;

use App\Models\Convocatoria;
use App\Models\Alumno;
use App\Models\DatosAlumnoAcademico;
use App\Models\ServicioSolicitado;
use App\Models\Solicitud;
use App\Models\Departament;


class ReporteConvocatoriaService
{
    public function reporte($id)
    {
        $convocatoria = Convocatoria::findOrFail($id);

        $cantidadHombres = Alumno::where('sexo', 'M')
            ->where('convocatoria_id', $convocatoria->id)
            ->count();

        $cantidadMujeres = Alumno::where('sexo', 'F')
            ->where('convocatoria_id', $convocatoria->id)
            ->count();

        $facultades = $this->obtenerFacultades();
        $cantidadFacultades = [];

        $solicitudes = Solicitud::where('convocatoria_id', $convocatoria->id)
            ->distinct()
            ->select('alumno_id') // Seleccionamos solo los IDs de los alumnos
            ->pluck('alumno_id');

        foreach ($facultades as $facultad) {
            $cantidadFacultades[$facultad] = Alumno::where('facultad', $facultad)
                ->whereIn('id', $solicitudes)
                ->count();
        }

        $escuelas = $this->obtenerEscuelasProfecionales();
        $cantidadEscuelas = [];

        foreach ($escuelas as $escuela) {
            $cantidadEscuelas[$escuela] = Alumno::where('escuela_profesional', $escuela)
                ->whereIn('id', $solicitudes)
                ->count();
        }

        $cantidadEscuelaGenero = [];
        $cantidadEscuelaMasculino = [];
        $cantidadEscuelaFemenino = [];
        foreach ($escuelas as $escuela) {
            $cantidadEscuelaMasculino[$escuela] = Alumno::where('escuela_profesional', $escuela)
                ->where('sexo', 'M')
                ->whereIn('id', $solicitudes)
                ->count();

            $cantidadEscuelaFemenino[$escuela] = Alumno::where('escuela_profesional', $escuela)
                ->where('sexo', 'F')
                ->whereIn('id', $solicitudes)
                ->count();
        }
        $cantidadEscuelaGenero['Masculino'] =  $cantidadEscuelaMasculino;
        $cantidadEscuelaGenero['Femenino'] =  $cantidadEscuelaFemenino;


        


        $cantidadPendientes = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'pendiente')
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();

        $cantidadRechazados = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'rechazado')
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();

        $cantidadAceptados = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'aceptado')
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();

        $cantidadAprobados = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'aprobado')
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();

        $reporte = [
            'sexo' => [
                'num_hombres' => $cantidadHombres,
                'num_mujeres' => $cantidadMujeres
            ],
            'facultades' => $cantidadFacultades,
            'escuelas_profesionales' => $cantidadEscuelas,
            'estados_solicitud' => [
                'pendiente' => count($cantidadPendientes),
                'rechazado' => count($cantidadRechazados),
                'aceptado' => count($cantidadAceptados),
                'aprobado' => count($cantidadAprobados)
            ],
            'escuela_profesionales_genero' => $cantidadEscuelaGenero,
        ];
        return $reporte;
    }

    public function reportePorServicio($id, $tipoServicio)
    {

        $convocatoria = Convocatoria::findOrFail($id);

        $cantidadHombres = Alumno::where('sexo', 'M')
            ->where('convocatoria_id', $convocatoria->id)
            ->count();

        $cantidadMujeres = Alumno::where('sexo', 'F')
            ->where('convocatoria_id', $convocatoria->id)
            ->count();

        $facultades = $this->obtenerFacultades();
        $cantidadFacultades = [];

        $solicitudes = Solicitud::where('convocatoria_id', $convocatoria->id)
            ->join('servicio_solicitado as ss', 'ss.solicitud_id', '=', 'solicitudes.id')
            ->where('ss.servicio_id', $tipoServicio)
            ->distinct()
            ->select('solicitudes.alumno_id')
            ->pluck('alumno_id');


        foreach ($facultades as $facultad) {
            $cantidadFacultades[$facultad] = Alumno::where('facultad', $facultad)
                ->whereIn('id', $solicitudes)
                ->count();
        }

        $escuelas = $this->obtenerEscuelasProfecionales();
        $cantidadEscuelas = [];

        foreach ($escuelas as $escuela) {
            $cantidadEscuelas[$escuela] = Alumno::where('escuela_profesional', $escuela)
                ->whereIn('id', $solicitudes)
                ->count();
        }

        $cantidadPendientes = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria, $tipoServicio) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'pendiente')
                ->where('servicio_id', $tipoServicio)
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();

        $cantidadRechazados = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria, $tipoServicio) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'rechazado')
                ->where('servicio_id', $tipoServicio)
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();


        $cantidadAceptados = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria, $tipoServicio) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'aceptado')
                ->where('servicio_id', $tipoServicio)
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();

        $cantidadAprobados = ServicioSolicitado::whereHas('solicitud', function ($query) use ($convocatoria, $tipoServicio) {
            $query->where('convocatoria_id', $convocatoria->id)
                ->where('estado', 'aprobado')
                ->where('servicio_id', $tipoServicio)
                ->whereIn('id', function ($subQuery) {
                    $subQuery->selectRaw('MAX(id)')
                        ->from('solicitudes')
                        ->groupBy('alumno_id'); // Obtener el último id por alumno
                });
        })->get();



        $reporte = [
            'sexo' => [
                'num_hombres' => $cantidadHombres,
                'num_mujeres' => $cantidadMujeres
            ],
            'facultades' => $cantidadFacultades,
            'escuelas_profesionales' => $cantidadEscuelas,
            'estados_solicitud' => [
                'pendiente' => count($cantidadPendientes),
                'rechazado' => count($cantidadRechazados),
                'aceptado' => count($cantidadAceptados),
                'aprobado' => count($cantidadAprobados)
            ]
        ];
        return $reporte;
    }

    private function obtenerFacultades()
    {
        $facultades = DatosAlumnoAcademico::groupBy('nomfac')
            ->pluck('nomfac')
            ->unique()
            ->values()
            ->toArray();
        unset($facultades[0]);
        return $facultades;
    }

    private function obtenerEscuelasProfecionales()
    {
        $escuelas = DatosAlumnoAcademico::groupBy('nomesp')
            ->pluck('nomesp')
            ->unique()
            ->values()
            ->toArray();
        unset($escuelas[4]);
        return $escuelas;
    }

    private function obtenerCantidadPorDepartamentoPorConvocatoriaYServicio($convocatoriaId, $servicioId) {
        // cantidad por departamento comedor
        $cantidad_departamento_servicio = [];
        $departamentos = Departament::get();

        $solicitudesComedor = Solicitud::where('convocatoria_id', $convocatoriaId)
            ->join('servicio_solicitado as ss', 'ss.solicitud_id', '=', 'solicitudes.id')
            ->where('ss.servicio_id', $servicioId)
            ->distinct()
            ->select('solicitudes.alumno_id')
            ->pluck('alumno_id');

        foreach($departamentos as $departamento) {
            $cantidad_departamento_servicio[$departamento->name] = Alumno::where('escuela_profesional',  $departamento)
           // ->whereIn('id', $solicitudes)
            ->count();
        }

        return $cantidad_departamento_servicio;
    }
}
