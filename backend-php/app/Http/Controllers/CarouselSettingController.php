<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Exceptions\ExceptionGenerate;
use App\Http\Response\Response;
use App\Http\Resources\CarouselSetting\CarouselSettingResource;
use App\Services\CarouselSetting\CreateCarouselSettingService;
use App\Services\CarouselSetting\UpdateCarouselSettingService;
use App\Services\CarouselSetting\ListCarouselSettingService;
use App\Services\CarouselSetting\DeleteCarouselSettingService;
use Illuminate\Support\Facades\Validator;

class CarouselSettingController extends Controller
{
    /**
     * Display a listing of all carousel settings (admin).
     */
    public function index(ListCarouselSettingService $listCarouselSettingService)
    {
        return Response::res('Lista de configuraciones de carrusel', CarouselSettingResource::collection($listCarouselSettingService->list()), 200);
    }

    /**
     * Display a listing of enabled carousel settings (public).
     */
    public function indexPublic(ListCarouselSettingService $listCarouselSettingService)
    {
        return Response::res('Lista de carruseles activos', CarouselSettingResource::collection($listCarouselSettingService->listEnabled()), 200);
    }

    /**
     * Store a newly created carousel setting.
     */
    public function create(Request $request, CreateCarouselSettingService $carouselSettingService)
    {
        try {
            // Manual validation
            $validator = Validator::make($request->all(), [
                'image' => 'required|image|mimes:jpeg,png,jpg,webp|max:10240',
                'title' => 'nullable|string|max:255',
                'description' => 'nullable|string|max:1000',
                'button_text' => 'nullable|string|max:100',
                'button_link' => 'nullable|url|max:500',
                'is_enabled' => 'boolean',
                'order' => 'integer|min:0'
            ]);

            if ($validator->fails()) {
                return Response::res('Errores de validación', $validator->errors(), 422);
            }

            $data = $validator->validated();

            // Add the uploaded file to the data array
            if ($request->hasFile('image')) {
                $data['image'] = $request->file('image');
            }

            $carouselSetting = $carouselSettingService->create($data);

            return Response::res('Carrusel creado satisfactoriamente', CarouselSettingResource::make($carouselSetting), 200);
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        } catch (\Exception $e) {
            return Response::res('Error al crear carrusel: ' . $e->getMessage(), null, 500);
        }
    }

    /**
     * Update carousel settings (bulk update).
     */
    public function update(Request $request, UpdateCarouselSettingService $carouselSettingService)
    {
        try {
            $validator = Validator::make($request->all(), [
                '*.id' => 'required|integer|exists:carousel_settings,id',
                '*.title' => 'nullable|string|max:255',
                '*.description' => 'nullable|string|max:1000',
                '*.button_text' => 'nullable|string|max:100',
                '*.button_link' => 'nullable|url|max:500',
                '*.is_enabled' => 'boolean',
                '*.order' => 'integer|min:0'
            ]);

            if ($validator->fails()) {
                return Response::res('Errores de validación', $validator->errors(), 422);
            }

            $data = $validator->validated();
            $updatedSettings = $carouselSettingService->update($data);

            return Response::res(
                'Configuraciones de carrusel actualizadas satisfactoriamente',
                CarouselSettingResource::collection($updatedSettings),
                200
            );
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        } catch (\Exception $e) {
            return Response::res('Error interno del servidor: ' . $e->getMessage(), null, 500);
        }
    }

    /**
     * Update a single carousel setting.
     */
    public function updateSingle($id, Request $request, UpdateCarouselSettingService $carouselSettingService)
    {
        try {
            $validator = Validator::make($request->all(), [
                'image' => 'nullable|image|mimes:jpeg,png,jpg,webp|max:10240',
                'title' => 'nullable|string|max:255',
                'description' => 'nullable|string|max:1000',
                'button_text' => 'nullable|string|max:100',
                'button_link' => 'nullable|url|max:500',
                'is_enabled' => 'boolean',
                'order' => 'integer|min:0'
            ]);

            if ($validator->fails()) {
                return Response::res('Errores de validación', $validator->errors(), 422);
            }

            $data = $validator->validated();

            // Add the uploaded file to the data array if present
            if ($request->hasFile('image')) {
                $data['image'] = $request->file('image');
            }

            $carouselSetting = $carouselSettingService->updateSingle($id, $data);

            return Response::res(
                'Carrusel actualizado satisfactoriamente',
                CarouselSettingResource::make($carouselSetting),
                200
            );
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        } catch (\Exception $e) {
            return Response::res('Error al actualizar carrusel: ' . $e->getMessage(), null, 500);
        }
    }

    /**
     * Remove the specified carousel setting.
     */
    public function destroy($id, DeleteCarouselSettingService $deleteCarouselSettingService)
    {
        try {
            $deleteCarouselSettingService->delete($id);
            return Response::res('Carrusel eliminado satisfactoriamente', null, 200);
        } catch (ExceptionGenerate $e) {
            return Response::res($e->getMessage(), null, $e->getStatusCode());
        } catch (\Exception $e) {
            return Response::res('Error al eliminar carrusel: ' . $e->getMessage(), null, 500);
        }
    }
}
