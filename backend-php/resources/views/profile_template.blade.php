<!DOCTYPE html>
<html lang="es">

<head>
    <meta charset="UTF-8">
    <title>{{ $title }}</title>
    <style>
        /* Agrega aquí los estilos para el PDF */
        body {
            font-family: Arial, sans-serif;
        }

        h1 {
            color: #333;
        }

        .title-image img {
            width: 100px;
            height: 100px;
            border-radius: 50%;
        }
    </style>
</head>

<body>
    <div class="header">
        <table style="width: 100%; border-spacing: 0; border-collapse: collapse;">
            <tr>
                <!-- Primera Celda -->
                <td style="padding: 0;">
                    <table style="width: 100%; border-spacing: 0; border-collapse: collapse;">
                        <tr>
                            <td style="padding: 0; vertical-align: middle;">
                                <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQR1OpH7wrnsCBwE55kpuJJKgQcTD4mSPHxLg&s"
                                    alt="" style="width: 100px; height: 100px; display: block;">
                            </td>
                            <td style="padding: 0; vertical-align: middle;">
                                <div>
                                    <span style="font-weight: bold; font-size: large; white-space: nowrap;">UNIVERSIDAD
                                        NACIONAL</span><br>
                                    <span style="font-weight: bold; font-size: large; white-space: nowrap;">AGRARIA DE
                                        LA SELVA</span>
                                </div>
                            </td>
                        </tr>
                    </table>
                </td>
                <!-- Segunda Celda -->
                <td style="padding: 0;">
                    <table style="width: 100%; border-spacing: 0; border-collapse: collapse;">
                        <tr>
                            <td style="padding: 0; vertical-align: middle;">
                                <div>
                                    <span
                                        style="font-weight: bold; font-size: large; white-space: nowrap;">DBU</span><br>
                                    <span style="font-weight: bold; font-size: large; white-space: nowrap;">Dirección de
                                        bienestar</span><br>
                                    <span
                                        style="font-weight: bold; font-size: large; white-space: nowrap;">Universitario</span>
                                </div>
                            </td>
                            <td style="padding: 0; vertical-align: middle;">
                                @if($link_profile == '')
                                <img src="https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
                                    alt="Imagen de perfil por defecto"
                                    style="width: 100px; height: 100px; border-radius: 50%; display: block;">
                                @else
                                <img src="{{ $link_profile }}"
                                    alt="Imagen de perfil"
                                    style="width: 100px; height: 100px; border-radius: 50%; display: block;">
                                @endif
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
            <!-- Título Centrado -->
            <tr>
                <td colspan="2" style="text-align: center; padding: 0;">
                    <h1 style="white-space: nowrap; margin: 0;">FICHA DEL ESTUDIANTE {{$year}}</h1>
                </td>
            </tr>
        </table>
    </div>

    <!-- Primera Sección -->
    <div class="section" style="margin-top: 20px;">
        <table style="width: 100%; border-spacing: 0; border-collapse: collapse; font-family: Arial, sans-serif;">
            <tr>
                <td colspan="2"
                    style="background-color: #f5f5f5; padding: 10px; font-weight: bold; text-align: left; border: 1px solid #ccc;">
                    {{ $first_section['description'] }}
                </td>
            </tr>
            @foreach ($first_section['requirements'] as $index => $requirement)
            @if ($index % 2 == 0) <!-- Si es un índice par -->
            <tr>
                @endif

                <td style="padding: 10px; border: 1px solid #ccc; width: 50%; vertical-align: top;">
                    <span style="font-weight: bold;">{{ $requirement['name'] }}:</span>
                    <div style="border-bottom: 1px solid #ccc; padding-bottom: 5px;">
                        @if ($requirement['type'] == 1)
                        <span style="color: blue; font-size: 10px; text-decoration: underline;">
                            {{ $requirement['value'] }}
                        </span>
                        @else
                        {{ $requirement['value'] }}
                        @endif
                    </div>
                </td>

                @if ($index % 2 == 1 || $index == count($first_section['requirements']) - 1)
                <!-- Si es un índice impar o el último -->
            </tr>
            @endif
            @endforeach


        </table>
    </div>
    <!-- Convocatorias  -->
    @foreach ($announcements as $announcement)
    <div class="convocatoria" style="margin-top: 20px;">
        <h2 style="background-color: #f5f5f5; padding: 10px; border: 1px solid gray;">{{ $announcement['name'] }}</h2>
        @foreach ($announcement['details_requests'] as $detail)
        <div class="section" style="margin-top: 20px;">
            <table style="width: 100%; border-spacing: 0; border-collapse: collapse; font-family: Arial, sans-serif;">
                <tr>
                    <td colspan="2"
                        style="background-color: #f5f5f5; padding: 10px; font-weight: bold; text-align: left; border: 1px solid #ccc;">
                        {{ $detail['description'] }}
                    </td>
                </tr>
                @foreach ($detail['requirements'] as $index => $requirement)
                @if ($index % 2 == 0) <!-- Si es un índice par -->
                <tr style="width: 100%;">
                    @endif
                    <td
                        style="padding: 10px; border: 1px solid #ccc; width: 50%; vertical-align: top; word-break: break-word;">
                        <span style="font-weight: bold;">{{ $requirement['name'] }}:</span>
                        <div style="border-bottom: 1px solid #ccc; padding-bottom: 5px;">
                            @if ($requirement['type'] == 1)
                            <span style="color: blue; font-size: 10px; text-decoration: underline;">
                                {{ $requirement['value'] }}
                            </span>
                            @else
                            {{ $requirement['value'] }}
                            @endif
                        </div>
                    </td>
                    @if ($index % 2 == 1 || $index == count($detail['requirements']) - 1)
                    <!-- Si es un índice impar o el último -->
                </tr>
                @endif
                @endforeach
            </table>
        </div>
        @endforeach

    </div>
    @endforeach


    <div style="text-align: right; margin-top: 200px;">
        <p style="border-top:1px solid black; display: inline-block; padding-top: 5px;">Firma de Dirección de Bienestar
            Universitario</p>
    </div>

</body>

</html>