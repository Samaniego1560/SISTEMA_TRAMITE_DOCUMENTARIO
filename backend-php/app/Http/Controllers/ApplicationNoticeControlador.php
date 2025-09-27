<?php

namespace App\Http\Controllers;

use App\Http\Response\Response;
use Illuminate\Http\Request;
use App\Exceptions\ExceptionGenerate;
use App\Http\Resources\ApplicationNotice\ApplicationNoticeResource;
use App\Http\Requests\ApplicationNotice\ApplicationNoticeRequest;
use App\Services\ApplicationNotice\UpdateApplicationNoticeService;
use App\Services\ApplicationNotice\ListApplicationNoticeService;
use App\Models\ApplicationNotice;

class ApplicationNoticeControlador extends Controller
{
    /**
     * Display a listing of the resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function index(ListApplicationNoticeService $listApplicationNoticeService)
    {
        return Response::res('Lista de tesistas o practicantes', ApplicationNoticeResource::collection($listApplicationNoticeService->list()), 200);
    }

    /**
     * Show the form for creating a new resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function create()
    {
        //
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param  \Illuminate\Http\Request  $request
     * @return \Illuminate\Http\Response
     */
    public function store(Request $request)
    {
        //
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function show($id)
    {
        //
    }

    /**
     * Show the form for editing the specified resource.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function edit($id)
    {
        //
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function update(ApplicationNoticeRequest $request, UpdateApplicationNoticeService $updateApplicationNoticeService)
    {
        try {

            $data = $request->validated(); // Obtiene los datos validados como un array
            $updatedApplicationNotice = $updateApplicationNoticeService->update($data); // Pasa directamente el array al servicio

            return Response::res(
                'Calendar settings registrados satisfactoriamente',
                ApplicationNoticeResource::collection($updatedApplicationNotice),
                200
            );
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        } catch (\Exception $e) {
            return Response::res('Error interno del servidor', null, 500);
        }
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function destroy($id)
    {
        //
    }
}
