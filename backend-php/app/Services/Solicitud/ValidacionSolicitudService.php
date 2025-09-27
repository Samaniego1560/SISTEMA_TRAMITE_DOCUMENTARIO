<?php

namespace App\Services\Solicitud;

use GuzzleHttp\Client;
use App\Exceptions\ExceptionGenerate;
use App\Http\Response\Response;
use App\Models\Alumno;
use App\Models\Convocatoria;
use App\Models\DatosAlumnoAcademico;
use App\Models\Solicitud;
use App\Models\ThesisPracticacing;
use DateTime;
use Exception;

class ValidacionSolicitudService
{
    public function validate(array $data)
    {
        $response = [];
            // Validando convocatoria actual
            $fechaActualConvocatoria = new DateTime();
            $dni = $data['DNI'];
            $convocatoria = Convocatoria::whereDate('fecha_inicio', '<=', $fechaActualConvocatoria)
                ->whereDate('fecha_fin', '>=', $fechaActualConvocatoria)->first();
            if (!$convocatoria) {
                throw new ExceptionGenerate('Actualmente no existe convocatoria en curso', 200);
            }
    
            // Recupera la data del alumno
            $alumno = Alumno::where('DNI', $dni)->first();
    
            // Valida requisitos
		$data["hora-0"] = date("Y-m-d H:i:s");

            $faltas = $this->validateRequisitosSolicitud($data, $convocatoria);
            if (count($faltas) > 0) {
                $search = $this->validateRequisitosSolicitudThesisOrPracticingOrStudent($data, $convocatoria->id);
                $data["hora-1"] = date("Y-m-d H:i:s");

                if ($search === null) {

                    throw new ExceptionGenerate(
                        'Usted no cumple con los siguientes requisitos para postular a los servicios de comedor e internado de la UNAS',
                        500,
                        $faltas
                    );

                }
    		$data["hora-2"] = date("Y-m-d H:i:s");
                $datosAlumnoAcademico = $this->getDatosAlumnoAcademicoBuscar($dni);
                if ($alumno) {
$data["hora-3"] = date("Y-m-d H:i:s");
                    $this->processAlumno($alumno, $datosAlumnoAcademico, $convocatoria->id, $dni);
                } else {
                    $this->createAlumno($search, $dni, $convocatoria->id);
$data["hora-4"] = date("Y-m-d H:i:s");
                }
$data["hora-5"] = date("Y-m-d H:i:s");

            } else {
                $datosAlumnoAcademico = $this->getDatosAlumnoAcademico($dni);
                $this->processAlumno($alumno, $datosAlumnoAcademico, $convocatoria->id, $dni);
            }

$data["hora-6"] = date("Y-m-d H:i:s");

    
            $response = $data;
    
        return $response;
    }

    public function validatePermision(array $data)
    {
        $response = [];
    
            // Validando convocatoria actual
            $convocatoriaId = $data['announcement_id'];
            $dni = $data['DNI'];
    
            // Recupera la data del alumno
            $alumno = Alumno::where('DNI', $dni)->first();
    
            // Valida requisitos
            $faltas = $this->validateRequisitosSolicitudPermision($data);
            if (count($faltas) > 0) {
                $search = $this->validateRequisitosSolicitudThesisOrPracticingOrStudent($data, $convocatoriaId);
                if ($search === null) {
                    throw new ExceptionGenerate(
                        'Usted no cumple con los siguientes requisitos para postular a los servicios de comedor e internado de la UNAS',
                        500,
                        $faltas
                    );
                }
    		
                $datosAlumnoAcademico = $this->getDatosAlumnoAcademicoBuscar($dni);
                if ($datosAlumnoAcademico) {
                    $this->processAlumno($alumno, $datosAlumnoAcademico, $convocatoriaId, $dni);
                } else {
                    $this->createAlumno($search, $dni, $convocatoriaId);
                }
            } else {
                $datosAlumnoAcademico = $this->getDatosAlumnoAcademico($dni);
                $this->processAlumno($alumno, $datosAlumnoAcademico, $convocatoriaId, $dni);
            }
    
            $response = $data;
    
        return $response;
    }

    private function validateRequisitosSolicitud(array $data, $convocatoria)
    {
        $faltas = [];
        if ($data['type_student'] == 'Tesista' || $data['type_student'] == 'Practicante') {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "No cumple con el semestre necesario",
            ]);
            return $faltas;
        }


        $datosAlumnoAcademico = $this->getDatosAlumnoAcademico($data['DNI']);
        if ($datosAlumnoAcademico == null) {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "No se encontraron registros de matricula",
            ]);
            return $faltas;
        }

        //Validacion de promedio ponderado semestral aprovado
        if (floatval($datosAlumnoAcademico['pps']) < $convocatoria->nota_minima && intval($datosAlumnoAcademico['nume_sem_cur']) != 0) {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "No cumple con el promedio semestal minimo de" .$convocatoria->nota_minima .", de usted es: " . $datosAlumnoAcademico['pps'],
            ]);
        }

        //Validacion de semestre matriculado
        if (strtoupper($datosAlumnoAcademico['est_mat_act']) == 'R') {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "No se encuentra matriculado en el semeste actual",
            ]);
        }

        //Validación de numero de ciclos cursado
        if (intval($datosAlumnoAcademico['nume_sem_cur']) > 12) {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "Ya no puede solicitar por exceso de semestre, de usted es: " . $datosAlumnoAcademico['nume_sem_cur'] . " el maximo se ciclos es 12",
            ]);
        }

        //Creditos matriculados (12 minimo)
        if (intval($datosAlumnoAcademico['credmat']) < $convocatoria->credito_minimo) {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "No cumple con los creditos minimos matriculas, de usted es: " . $datosAlumnoAcademico['nume_sem_cur'] . " el minimo es " .$convocatoria->credito_minimo,
            ]);
        }

        //Validar articulo incurso
        if (strlen($datosAlumnoAcademico['artincurso']) > 0) {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "Usted no puede postular por que se encuentra incurso en los articulos " . $datosAlumnoAcademico['artincurso'],
            ]);
        }
        return $faltas;
    }

    private function validateRequisitosSolicitudPermision(array $data)
    {
        $faltas = [];

        $datosAlumnoAcademico = $this->getDatosAlumnoAcademico($data['DNI']);
        if ($datosAlumnoAcademico == null) {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "No se encontraron registros de matricula",
            ]);
            return $faltas;
        }

        //Validacion de semestre matriculado
        if (strtoupper($datosAlumnoAcademico['est_mat_act']) == 'R') {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "No se encuentra matriculado en el semeste actual",
            ]);
        }

        //Validar articulo incurso
        if (strlen($datosAlumnoAcademico['artincurso']) > 0) {
            array_push($faltas,  [
                "tipo" => "academicos",
                "msg" => "Usted no puede postular por que se encuentra incurso en los articulos " . $datosAlumnoAcademico['artincurso'],
            ]);
        }
        return $faltas;
    }
    
    private function getDatosAlumnoAcademico(String $dni)
    {
        try {
            $datosAlumnoAcademico = DatosAlumnoAcademico::where('tdocumento', 'DNI' . $dni)->get();
            $datosAlumnoAcademico = $datosAlumnoAcademico[(count($datosAlumnoAcademico) - 1)];
            if (!$datosAlumnoAcademico) {
                throw new ExceptionGenerate('No existe registro academico del alumno, una de la razones es que supere los 12 semestres academicos o no estes matriculado en el semestre 2025-I', 200);
            }
                return $datosAlumnoAcademico;
        } catch (Exception $e) {
            throw new ExceptionGenerate('No existe registro academico del alumno, una de la razones es que supere los 12 semestres academicos o no estes matriculado en el semestre 2025-I', 200);
        }
    }

    private function getDatosAlumnoAcademicoBuscar(String $dni)
    {
        try {
            $datosAlumnoAcademico = DatosAlumnoAcademico::where('tdocumento', 'DNI' . $dni)->get();
            $datosAlumnoAcademico = $datosAlumnoAcademico[(count($datosAlumnoAcademico) - 1)];
            if (!$datosAlumnoAcademico){
                return null;
            }
            return $datosAlumnoAcademico;
        } catch (Exception $e) {
            return null;
        }
    }

    private function getYearsInDates(DateTime $firsData, DateTime $secondData)
    {
        $diferencia = $firsData->diff($secondData);
        return $diferencia->y;
    }
    
    // Función auxiliar para procesar un alumno existente o registrar sus datos
    private function processAlumno($alumno, $datosAlumnoAcademico, $convocatoriaId, $dni)
    {
        if ($alumno) {
            $alumno->pps = $datosAlumnoAcademico['pps'] ?? 0;
            $alumno->ppa = $datosAlumnoAcademico['ppa'] ?? 0;
            $alumno->tca = $datosAlumnoAcademico['tca'] ?? 0;
            $alumno->convocatoria_id = $convocatoriaId;
            $alumno->facultad = $datosAlumnoAcademico['nomfac'] ?? $alumno->facultad;
            $alumno->escuela_profesional = $datosAlumnoAcademico['nomesp'] ?? $alumno->escuela_profesional;
            $alumno->save();
        } else {
            $this->createAlumno($datosAlumnoAcademico, $dni, $convocatoriaId);
        }
    }
    
    // Función auxiliar para crear un nuevo alumno
    private function createAlumno($data, $dni, $convocatoriaId)
    {
        Alumno::create([
            "codigo_estudiante" => $data['codalumno'] ?? $data['cod_student'],
            "DNI" => $dni,
            "nombres" => $data['nombre'] ?? $data['name'],
            "apellido_paterno" => $data['appaterno'],
            "apellido_materno" => $data['apmaterno'],
            "sexo" => $data['sexo'],
            "facultad" => $data['nomfac'] ?? $data['facultad'],
            "escuela_profesional" => $data['nomesp'] ?? $data['escuela_profesional'],
            "ultimo_semestre" => $data['codsem'] ?? "",
            "modalidad_ingreso" => $data['mod_ingreso'],
            "edad" => $data['age'] ?? $this->getYearsInDates(new DateTime($data['fecnac']), new DateTime()),
            "correo_institucional" => $data['emailinst'] ?? $data['email'] ?? "",
            "direccion" => $data['direccion'] ?? "",
            "fecha_nacimiento" => $data['fecnac'] ?? new DateTime(),
            "correo_personal" => $data['email'] ?? "",
            "celular_estudiante" => $data['telcelular'] ?? "",
            "celular_padre" => $data['tel_ref'] ?? "",
            "estado_matricula" => $data['est_mat_act'] ?? "N",
            "creditos_matriculados" => $data['credmat'] ?? 0,
            "num_semestres_cursados" => $data['nume_sem_cur'] ?? 0,
            "pps" => $data['pps'] ?? 0,
            "ppa" => $data['ppa'] ?? 0,
            "tca" => $data['tca'] ?? 0,
            "convocatoria_id" => $convocatoriaId,
        ]);
    }
    
    private function validateRequisitosSolicitudThesisOrPracticingOrStudent(array $data, $convocatoria_id)
    {

        $thesis_practicacing = ThesisPracticacing::where([
            ['type_student', '=', $data['type_student']],
            ['dni', '=', $data['DNI']],
            ['email', '=', $data['correo']],
            ['convocatoria_id', '=', $convocatoria_id],
            ['status_id', '=', 1],

        ])->first();
        return $thesis_practicacing;
    }
}
