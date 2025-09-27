<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Services\Report\UserProfilePdfService;

class ReportController extends Controller
{
    /**
     * Display a listing of the resource.
     *
     * @return \Illuminate\Http\Response
     */
    public function index()
    {
        //
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
    public function update(Request $request, $id)
    {
        //
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

    public function generateProfile($dni, UserProfilePdfService $userProfilePdfService)
{
    try {
       return response()->json(['msg' => 'Detalle de perfil', 'detalle' => $userProfilePdfService->getDataReport($dni)]);

    } catch (\Exception $e) {
        error_log("Error: " . $e->getMessage());
        return response()->json(['error' => 'Error al generar el PDF: ' . $e->getMessage()], 500);
    }
}

    public function downloadProfile($dni, UserProfilePdfService $userProfilePdfService)
{
    try {
       return $userProfilePdfService->generatePdf($dni);

    } catch (\Exception $e) {
        // Manejo de errores general
        error_log("Error: " . $e->getMessage());
        return response()->json(['error' => 'Error al generar el PDF: ' . $e->getMessage()], 500);
    }
}
}
