<?php

use Illuminate\Support\Facades\DB;
/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
|
| Here is where you can register all of the routes for an application.
| It is a breeze. Simply tell Lumen the URIs it should respond to
| and give it the Closure to call when that URI is requested.
|
*/

// Rutas generales, no requiere de inicio de session
$router->group([], function ($router) {
    $router->get('/', function () use ($router) {
        //return ;
        return '<h1>API REST</h1> </br> ';
    });
    $router->post('/login', 'AuthController@login');

    $router->get('/convocatoria/vigente-convocatoria', 'ConvocatoriaController@vigenteConvocatoria');

    $router->post('/solicitud/validacion', 'SolicitudController@validacionSolicitud');
    $router->post('/solicitud/validacion/thesis_practicing_student', 'SolicitudController@validacionSolicitud');
    $router->post('/solicitud/create', 'SolicitudController@create');
    $router->post('/solicitud/uploadDocument', 'SolicitudController@uploadDocument');
    $router->get('/solicitud/alumno/{dni}', 'SolicitudController@cargaSolicitudAlumno');
    $router->get('/solicitud/convocatoria/alumno/{convocatoriaId}/{dni}', 'SolicitudController@cargaSolicitudConvocatoriaAlumno');
    $router->get('/solicitud/servicioSolicitado/{dni}', 'SolicitudController@servicioSolicitadoSolicitante');

    $router->get('/solicitud/export/', 'SolicitudController@solicitudExport');

    $router->get('/solicitud/export/convocatoria/{id}', 'SolicitudController@solicitudExportById');

    // Departaments
    $router->get('/ubigeo/departaments', 'DepartamentController@index');

    // Provinces
    $router->get('/ubigeo/provinces', 'ProvinceController@index');

    // Districts
    $router->get('/ubigeo/districts', 'DistrictController@index');

    // add tesistas
    $router->post('/thesis/register', 'ThesisPracticacingController@create');
    $router->delete('/thesis/destroy/{id}', 'ThesisPracticacingController@destroy');

    //report
    $router->get('/solicitud/export/data/convocatoria/{id}', 'SolicitudController@solicitudExportDataById');
    $router->get('/solicitud/export/data/convocatoria/{servicio}/{estado}/{id}', 'SolicitudController@solicitudDataServicioEstadoByConvocatoriaId');

    $router->get('/report/profile/{dni}', 'ReportController@generateProfile');
    $router->get('/report/profile/download/{dni}', 'ReportController@downloadProfile');

    // calendars
    $router->get('/calendar/setting/list', 'CalendarSettingController@index');

    // application notice
    $router->get('/application/notice/list', 'ApplicationNoticeControlador@index');

    $router->get('/DatosAlumnoAcademico/show/{DNI}', 'DatosAlumnoAcademicoController@show');
});

// Rutas que requieren nivel de acceso 1
$router->group(['middleware' => ['auth', 'restriclevel1']], function ($router) {

    $router->post('/register', 'AuthController@register');

    //Users
    $router->put('/users/update/{id}', 'UserController@update');
    $router->get('/users', 'UserController@index');
    $router->get('/users/show/{id}', 'UserController@show');
    $router->delete('/users/destroy/{id}', 'UserController@destroy');

    //Convocatoria
    $router->post('/convocatoria/create', 'ConvocatoriaController@create');
    $router->get('/convocatoria/show/{id}', 'ConvocatoriaController@show');

    //$router->put('/convocatoria/update/{id}', 'ConvocatoriaController@update');


    $router->get('/DatosAlumnoAcademico', 'DatosAlumnoAcademicoController@index');


    // thesis or practicing
    $router->post('/convocatoria/thesis-practicings/register', 'ThesisPracticacingController@create');
    $router->put('/convocatoria/thesis-practicings/update/{id}', 'ThesisPracticacingController@update');
    $router->get('/convocatoria/thesis-practicings', 'ThesisPracticacingController@index');
    $router->delete('/convocatoria/thesis-practicings/destroy/{id}', 'ThesisPracticacingController@destroy');
});

// Rutas que requieren nivel de acceso 2 y 1
$router->group(['middleware' => ['auth', 'restriclevel2']], function ($router) {

    //Convocatoria
    $router->get('/convocatoria', 'ConvocatoriaController@index');
    $router->get('/convocatoria/reporte/{id}', 'ConvocatoriaController@reporteConvocatoria');
    $router->get('/convocatoria/reporte/{id}/{tipoServicio}', 'ConvocatoriaController@reporteConvocatoriaYServicio');

    $router->get('/servicio', 'ServicioController@index');
});

// Rutas que tienen  acceso todos los niveles de usuarios logeados
$router->group(['middleware' => ['auth', 'restriclevel3']], function ($router) {
    // application-notice
    $router->post('/application/notice/update', 'ApplicationNoticeControlador@update');
    
    //Solicitud
    $router->get('/solicitudes/{id}', 'SolicitudController@index');
    $router->get('/solicitud/show/{id}', 'SolicitudController@show');
    //LevelUsers
    $router->get('/leveluser', 'LevelUserController@index');

    $router->post('/logout', 'AuthController@logout');
    $router->get('/validateToken', 'AuthController@validateToken');

    //Solicitud
    $router->put('/solicitud/servicio', 'SolicitudController@updateServicio');

    //Solicitud
    $router->get('/solicitud/servicio/cantidad/vacantes/{id}', 'SolicitudController@getNumberOfVacancies');

    // Validate for add student by super role
    $router->post('/solicitud/validacion/permision', 'SolicitudController@validacionSolicitudPermision');

    // setting calendar
    $router->post('/calendar/setting/update', 'CalendarSettingController@update');

    // application-notice
    $router->post('/application/notice/update', 'ApplicationNoticeControlador@update');
});
