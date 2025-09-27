<?php

namespace App\Http\Controllers;

use App\Exceptions\ExceptionGenerate;
use App\Http\Requests\Solicitud\SolicitudRequest;
use App\Http\Requests\Solicitud\SolicitudServicioRequest;
use App\Http\Requests\Solicitud\UploadDocumentRequest;
use App\Http\Requests\Solicitud\ValidarDatosAlumnoRequest;
use App\Http\Requests\Solicitud\ValidarDatosTesistaPracticanteAlumno;
use App\Http\Resources\ServicioSolicitado\ServicioSolicitadoResource;
use App\Http\Resources\ServicioSolicitado\UpdateServicioSolicitadoResource;
use App\Http\Resources\Solicitud\SolicitudResource;
use App\Http\Response\Response;
use App\Services\Solicitud\CargaSolicitudAlumnoService;
use App\Services\Solicitud\CargaSolicitudConvocatoriaAlumnoService;
use App\Services\Solicitud\ValidacionSolicitudService as SolicitudValidacionSolicitudService;
use App\Services\Solicitud\CreateSolicitudService;
use App\Services\Solicitud\ListSolicitudService;
use App\Services\Solicitud\ServicioSolicitadoSolicitanteService;
use App\Services\Solicitud\ShowSolicitudService;
use App\Services\Solicitud\SolicitudExportService;
use App\Services\Solicitud\SolicitudServicioService;

class SolicitudController extends Controller
{

    public function validacionSolicitud(ValidarDatosAlumnoRequest $request, SolicitudValidacionSolicitudService $validacionSolicitudService)
    {
        try {
            $result = $validacionSolicitudService->validate($request->validated());
            return Response::res('Usted es apto para solicitar los servicios brindados por OBU - UNAS', $result);
        } catch (ExceptionGenerate $e) {
            error_log('ENTRA EN EXEPCION');
            return Response::res($e->getMessage(), $e->getData(), $e->getStatusCode());
        } catch (\Exception $e) {
            error_log('Capturada Excepción estándar: ' . $e->getMessage());
            return Response::res('Error general', [], 500);
        }
    }
    public function validacionSolicitudPermision(ValidarDatosTesistaPracticanteAlumno $request, SolicitudValidacionSolicitudService $validacionSolicitudService)
    {
        try {
            return Response::res('Usted es apto para solicitar los servicos brindados por OBU - UNAS', ($validacionSolicitudService->validatePermision($request->validated())));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), $e->getData(), $e->getStatusCode());
        }
    }
    public function index($id, ListSolicitudService $listSolicitudService)
    {
        try {
            return Response::res('Solicitudes listadas', SolicitudResource::collection($listSolicitudService->list($id)), 200);
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), $e->getData(), $e->getStatusCode());
        }
    }

    public function create(SolicitudRequest $request, CreateSolicitudService $createService)
    {
        try {
            return Response::res('Solicitud registrada', SolicitudResource::make($createService->create($request->validated())));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), $e->getData(), $e->getStatusCode());
        }
    }

    public function show($id, ShowSolicitudService $showSolicitudService)
    {
        try {
            return Response::res('Solicitud filtrada', SolicitudResource::make($showSolicitudService->show($id)));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), $e->getData(), $e->getStatusCode());
        }
    }

    public function updateServicio(SolicitudServicioRequest $request, SolicitudServicioService $updateSolicitudServicioService)
    {
        try {
            return Response::res('Datos de solicitud actualizada', ServicioSolicitadoResource::collection($updateSolicitudServicioService->updateServicio($request->validated())));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }

    public function getNumberOfVacancies($id, SolicitudServicioService $updateSolicitudServicioService)
    {
        try {
            $data = $updateSolicitudServicioService->getNumberOfVacancies($id);
            return Response::res('Datos de solicitud actualizada', $data);
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }


    public function uploadDocument(UploadDocumentRequest $request, CreateSolicitudService $createService)
    {
        try {
            return Response::res('Documento registrado', ($createService->uploadFile($request->validated())));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }
    public function cargaSolicitudAlumno($dni, CargaSolicitudAlumnoService $cargaSolicitudAlumnoService)
    {
        try {
            return Response::res('Documento registrado', ($cargaSolicitudAlumnoService->cargaSolicitudAlumno($dni)));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }
    public function cargaSolicitudConvocatoriaAlumno($convocatoriaId, $dni, CargaSolicitudConvocatoriaAlumnoService $cargaSolicitudConvocatoriaAlumnoService)
    {
        try {
            return Response::res('Documento registrado', ($cargaSolicitudConvocatoriaAlumnoService->cargaSolicitudConvocatoriaAlumno($convocatoriaId, $dni)));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }
    public function solicitudExport(SolicitudExportService $solicitudExportService)
    {
        return $solicitudExportService->export();
    }
    public function solicitudExportById($id, SolicitudExportService $solicitudExportService)
    {
        return $solicitudExportService->exportById($id);
    }

    public function solicitudExportDataById($id, SolicitudExportService $solicitudExportService)
    {
        try {
            return Response::res('Reporte de la convocatoria', ($solicitudExportService->dataExportByConvocatoriaId($id)));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }

    public function solicitudDataServicioEstadoByConvocatoriaId($servicio, $estado, $id, SolicitudExportService $solicitudExportService)
    {
        try {
            return Response::res('Reporte estados', ($solicitudExportService->dataExportByServicioEstadoConvocatoriaId($servicio, $estado, $id)));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }

    public function servicioSolicitadoSolicitante($dni, ServicioSolicitadoSolicitanteService $servicioSolicitadoSolicitanteService)
    {
        try {
            return Response::res('Servicios solicitados del solicitante', ($servicioSolicitadoSolicitanteService->servicioSolicitante($dni)));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }
}
