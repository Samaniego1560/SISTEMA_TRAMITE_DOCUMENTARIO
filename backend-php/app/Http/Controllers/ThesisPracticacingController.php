<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Exceptions\ExceptionGenerate;
use App\Http\Response\Response;
use App\Http\Requests\ThesisPracticacing\ThesisPracticacingRequest;
use App\Http\Resources\ThesisPracticacing\ThesisPracticacingResource;
use App\Services\ThesisPracticacing\DeleteThesisPracticacingService;
use App\Services\ThesisPracticacing\ListThesisPracticacingService;
use App\Services\ThesisPracticacing\RegisterThesisPracticacingService;
use App\Services\ThesisPracticacing\UpdateThesisPracticacingService;


class ThesisPracticacingController extends Controller
{
    /**
     * Display a listing of the resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function index(ListThesisPracticacingService $listThesisPracticacingService)
    {
        return Response::res('Lista de tesistas o practicantes', ThesisPracticacingResource::collection($listThesisPracticacingService->list()), 200);
    }

    /**
     * Show the form for creating a new resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function create(ThesisPracticacingRequest $request, RegisterThesisPracticacingService $registerThesisPracticacingService)
    {
        try {
            return Response::res('Tesista o practicante registrada satisfactoriamente', ThesisPracticacingResource::make($registerThesisPracticacingService->register($request->validated())), 200);
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
    public function update(ThesisPracticacingRequest $request, UpdateThesisPracticacingService $updateThesisPracticacingService, $id)
    {
        try {
            return Response::res('Datos de tesista o practicante actualizado satisfactoriamente', ThesisPracticacingResource::make($updateThesisPracticacingService->update($id, $request->validated())));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return \Illuminate\Http\Response
     */
    public function destroy($id, DeleteThesisPracticacingService $deleteThesisPracticacingService)
    {
        try {
            return Response::res('Tesista o practicante eliminado', ThesisPracticacingResource::make($deleteThesisPracticacingService->delete($id)));
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        }
    }
}
