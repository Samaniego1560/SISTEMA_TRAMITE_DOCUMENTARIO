<?php

namespace App\Services\Solicitud;

use App\Exceptions\ExceptionGenerate;
use App\Models\Alumno;
use App\Models\Convocatoria;
use App\Models\Requisito;
use App\Models\ServicioSolicitado;
use App\Models\Solicitud;
use DateTime;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class CreateSolicitudService
{
    private array $requisitosCache = [];
    private array $alumnosCache = [];

    public function create(array $data)
    {
        $solicitudes = Solicitud::where('convocatoria_id', $data['convocatoria_id'])
            ->where('alumno_id', $data['alumno_id'])
            ->get();

        foreach ($solicitudes as $solicitud) {
            ServicioSolicitado::where('solicitud_id', $solicitud->id)->update(['estado' => 'rechazado']);
        }

        DB::beginTransaction();
        try {
            $solicitud = Solicitud::create([
                'convocatoria_id' => $data['convocatoria_id'],
                'alumno_id' => $data['alumno_id']
            ]);

            foreach ($data['servicios_solicitados'] as $servicioSolicitadoData) {
                $this->servicioSolicitado($servicioSolicitadoData, $solicitud);
            }
            DB::commit();

            $solicitud->load(['servicioSolicitados']);

            foreach ($data['detalle_solicitudes'] as $detalleSolicitudData) {
                $this->detalleSolicitud($detalleSolicitudData, $solicitud);
            }
            DB::commit();

            $solicitud->load(['detalleSolicitudes']);

            return $solicitud;
        } catch (\Exception $e) {
            DB::rollBack();
            Log::error('solicitud create: ' . $e->getMessage() . ', Line: ' . $e->getLine());
        }
    }
    public function uploadFile(array $file): array
    {
        $file_name = $file['id_convocatoria'] . '_' . $file['dni_alumno'] . '_' . $file['name_file'];

        //Guardando Documento
        $this->saveFileLocal($file['file'], $file_name);

        $file['url_file'] = $this->getUrlFile($file);

        return $file;
    }
    private function servicioSolicitado(array $data, Solicitud $solicitud): Model
    {
        return $solicitud->servicioSolicitados()->create([
            'estado' => $data['estado'],
            'servicio_id' => $data['servicio_id'],
        ]);
    }

    private function detalleSolicitud(array $data, Solicitud $solicitud): Model
    {
        $requisito = $this->getRequisitoById((int) $data["requisito_id"]);
        $alumno = $this->getAlumnoById((int) $solicitud->alumno_id);

        if (!$requisito) {
            return $solicitud->detalleSolicitudes()->create([
                "respuesta_formulario" => $data['respuesta_formulario'] ?? null,
                "url_documento" => $data['url_documento'] ?? null,
                "opcion_seleccion" => $data['opcion_seleccion'] ?? null,
                'requisito_id' => $data['requisito_id'],
                'order' => $data['order'],
            ]);
        }

        $nombrePss = mb_strtolower(trim($requisito->nombre));
        if ($nombrePss === 'pss' || $nombrePss === 'ponderado semestral (pss)') {
            $data['respuesta_formulario'] = $alumno ? (string)$alumno->pps : ($data['respuesta_formulario'] ?? null);
        }
        if ($alumno) {
            if ($requisito->nombre == "Departamento de procedencia" || $requisito->nombre == "Provincia de procedencia" || $requisito->nombre == "Distrito de procedencia") {
                if (count(explode("/", (string)$alumno->lugar_procedencia)) >= 4)
                    $alumno->lugar_procedencia = "";
                $alumno->lugar_procedencia = $alumno->lugar_procedencia . strtoupper((string)($data['respuesta_formulario'] ?? '')) . '/';
            }
            if ($requisito->nombre == "Departamento de nacimiento" || $requisito->nombre == "Provincia de nacimiento" || $requisito->nombre == "Distrito de nacimiento") {
                if (count(explode("/", (string)$alumno->lugar_nacimiento)) >= 4)
                    $alumno->lugar_nacimiento = "";
                $alumno->lugar_nacimiento = $alumno->lugar_nacimiento . strtoupper((string)($data['respuesta_formulario'] ?? '')) . '/';
            }
            if ($requisito->nombre == "Celular De Estudiante") {
                $alumno->celular_estudiante = empty($data['respuesta_formulario']) ? "" : $data['respuesta_formulario'];
            }
            if ($requisito->nombre == "Celular Padre") {
                $alumno->celular_padre = empty($data['respuesta_formulario']) ? "" : $data['respuesta_formulario'];
            }
            $alumno->update();
        }
        //recuperar el documento y almacenar en el storage
        return $solicitud->detalleSolicitudes()->create([
            "respuesta_formulario" => $data['respuesta_formulario'],
            "url_documento" => $data['url_documento'],
            "opcion_seleccion" => $data['opcion_seleccion'],
            'requisito_id' => $data['requisito_id'],
            'order' => $data['order'],
        ]);
    }

    private function getRequisitoById(int $id): ?Requisito
    {
        if (!array_key_exists($id, $this->requisitosCache)) {
            $this->requisitosCache[$id] = Requisito::find($id);
        }

        return $this->requisitosCache[$id];
    }

    private function getAlumnoById(int $id): ?Alumno
    {
        if (!array_key_exists($id, $this->alumnosCache)) {
            $this->alumnosCache[$id] = Alumno::find($id);
        }

        return $this->alumnosCache[$id];
    }
    private function standarDates(string $dateString)
    {
        $dateFormats = ['d-m-Y', 'Y-m-d', 'd/m/Y', 'Y/m/d'];
        $date = false;
        foreach ($dateFormats as $dateFormat) {
            $date = DateTime::createFromFormat($dateFormat, $dateString);
            if ($date) {
                break;
            }
        }

        // Check if the DateTime object was created successfully
        if ($date) {
            // Convert the DateTime object to the desired format
            $formattedDate = $date->format('d-m-Y');
            return $formattedDate; // Output: 25-12-2022
        } else {
            throw new ExceptionGenerate('No se soportar el tipo fecha ingresado', 200);
        }
    }
    private function saveFileLocal($file, $file_name)
    {
        //Decodificando Documento para ser guardado
        $filea = base64_decode($file);

        //Guardando Documento
        file_put_contents('storage/app/public/' . $file_name, $filea);
    }

    private function getUrlFile($file)
    {
        //Generando url del Documento almacenado
        return env('APP_URL') . '/storage/app/public/' . $file['id_convocatoria'] . '_' . $file['dni_alumno'] . '_' . $file['name_file'];
    }

    /**
     * Reemplaza un archivo existente
     * 1. Elimina el archivo anterior del storage
     * 2. Guarda el nuevo archivo
     * 3. Actualiza la URL en la base de datos
     */
    public function replaceFile(array $data): array
    {
        DB::beginTransaction();
        try {
            // 1. Buscar el detalle de solicitud
            $detalleSolicitud = \App\Models\DetalleSolicitud::find($data['detalle_solicitud_id']);

            if (!$detalleSolicitud) {
                throw new ExceptionGenerate('Detalle de solicitud no encontrado', 404);
            }

            // 2. Eliminar archivo anterior si existe
            if ($detalleSolicitud->url_documento) {
                $this->deleteFileFromUrl($detalleSolicitud->url_documento);
            }

            // 3. Generar nuevo nombre de archivo
            $file_name = $data['id_convocatoria'] . '_' .
                $data['dni_alumno'] . '_' .
                $data['name_file'];

            // 4. Guardar nuevo archivo
            $this->saveFileLocal($data['file'], $file_name);

            // 5. Generar nueva URL
            $new_url = $this->getUrlFile($data);

            // 6. Actualizar base de datos
            $detalleSolicitud->url_documento = $new_url;
            $detalleSolicitud->save();

            DB::commit();

            Log::info('Archivo reemplazado exitosamente', [
                'detalle_solicitud_id' => $detalleSolicitud->id,
                'old_url' => $detalleSolicitud->url_documento,
                'new_url' => $new_url
            ]);

            return [
                'id' => $detalleSolicitud->id,
                'url_file' => $new_url,
                'message' => 'Archivo reemplazado exitosamente'
            ];
        } catch (ExceptionGenerate $e) {
            DB::rollBack();
            throw $e;
        } catch (\Exception $e) {
            DB::rollBack();
            Log::error('Error al reemplazar archivo: ' . $e->getMessage() . ', Line: ' . $e->getLine());
            throw new ExceptionGenerate('Error al reemplazar el archivo', 500);
        }
    }

    /**
     * Elimina un archivo físico dado su URL
     */
    private function deleteFileFromUrl(string $url): void
    {
        try {
            // Extraer el nombre del archivo de la URL
            // URL format: http://domain/storage/app/public/filename.ext
            $url_parts = explode('/storage/app/public/', $url);

            if (count($url_parts) === 2) {
                $filename = $url_parts[1];
                $file_path = 'storage/app/public/' . $filename;

                // Eliminar archivo si existe
                if (file_exists($file_path)) {
                    unlink($file_path);
                    Log::info('Archivo eliminado correctamente: ' . $file_path);
                } else {
                    Log::warning('Archivo no encontrado para eliminar: ' . $file_path);
                }
            }
        } catch (\Exception $e) {
            Log::warning('No se pudo eliminar el archivo anterior: ' . $e->getMessage());
            // No lanzar excepción, continuar con el proceso
        }
    }
}
