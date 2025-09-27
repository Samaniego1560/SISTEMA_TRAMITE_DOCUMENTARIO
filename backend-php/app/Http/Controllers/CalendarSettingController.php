<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Exceptions\ExceptionGenerate;
use App\Http\Response\Response;
use App\Http\Requests\CalendarSetting\CalendarSettingRequest;
use App\Http\Resources\CalendarSetting\CalendarSettingResource;
use App\Services\CalendarSetting\CreateCalendarSettingService;
use App\Services\CalendarSetting\UpdateCalendarSettingService;
use App\Services\CalendarSetting\ListCalendarSettingService;



class CalendarSettingController extends Controller
{
    /**
     * Display a listing of the resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function index(ListCalendarSettingService $listCalendarSettingService)
    {
        return Response::res('Lista de tesistas o practicantes', CalendarSettingResource::collection($listCalendarSettingService->list()), 200);
    }

    /**
     * Show the form for creating a new resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function create(CalendarSettingRequest $request, CreateCalendarSettingService $calendarSettingService)
    {
        try {
            return Response::res('Calendar setting registrada satisfactoriamente', CalendarSettingResource::make($calendarSettingService->create($request->validated())), 200);
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
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
    public function update(CalendarSettingRequest $request, UpdateCalendarSettingService $calendarSettingService)
    {
        try {
            $data = $request->validated(); // Obtiene los datos validados como un array
            $updatedSettings = $calendarSettingService->update($data); // Pasa directamente el array al servicio

            return Response::res(
                'Calendar settings registrados satisfactoriamente',
                CalendarSettingResource::collection($updatedSettings),
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
