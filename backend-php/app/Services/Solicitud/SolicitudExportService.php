<?php

namespace App\Services\Solicitud;

use App\Models\DetalleSolicitud;
use App\Models\Requisito;
use App\Models\Servicio;
use App\Models\ServicioSolicitado;
use App\Models\Solicitud;
use App\Models\Convocatoria;
use App\Exports\SolicitudesExport;
use App\Exports\SolicitudesExportPorConvocatoria;
use App\Exports\SolicitudExportDataPorConvocatoria;
use App\Exports\SolicitudExportDataPorServicioEstadoConvocatoria;
use Maatwebsite\Excel\Facades\Excel;
use Illuminate\Support\Facades\DB;
use Illuminate\Database\Eloquent\Collection;
use App\Services\Convocatoria\UltimaConvocatoriaService;

class SolicitudExportService
{
    public function export()/* : ?Model */
    {
        $convocatoria = Convocatoria::find(1);
        //$convocatoria = $convocatoria[(count($convocatoria) - 1)];
        return Excel::download(new SolicitudesExport, $convocatoria->nombre . '-Solicitantes.xlsx');
    }

    public function exportById($id)/* : ?Model */
    {
        $convocatoria = Convocatoria::find($id);
        //$convocatoria = $convocatoria[(count($convocatoria) - 1)];
        return Excel::download(new SolicitudesExportPorConvocatoria($id), $convocatoria->nombre . '-Solicitantes.xlsx');
    }

    public function dataExportByConvocatoriaId($id): ?Collection
    {
        $exportData = new SolicitudExportDataPorConvocatoria($id);
        return $exportData->collection(); // Llamamos al método collection() que retorna la Collection
    }

    public function dataExportByServicioEstadoConvocatoriaId($servicio, $estado, $id): ?Collection
    {
        $exportData = new SolicitudExportDataPorServicioEstadoConvocatoria($servicio, $estado, $id);
        return $exportData->collection(); // Llamamos al método collection() que retorna la Collection
    }

}
