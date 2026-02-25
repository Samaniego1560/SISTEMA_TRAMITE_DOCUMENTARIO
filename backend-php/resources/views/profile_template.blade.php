<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="utf-8">
  <title>{{ $title }}</title>
  <style>
    /* Estilos para PDF */
    * { font-family: "DejaVu Sans", sans-serif; }
    .break-word { word-break: break-word; }

    /* M&aacute;rgenes de p&aacute;gina: 3 cm abajo */
    @page { margin: 15mm 15mm 30mm 15mm; }

    /* Evitar que una tabla se parta entre p&aacute;ginas (para secciones que s&iacute; queremos juntas) */
    .table-block {
      page-break-inside: avoid;
      break-inside: avoid;
      break-inside: avoid-page;
      page-break-after: auto;
      page-break-before: auto;
      margin-bottom: 12px;
      display: block;
    }

    .keep-together {
      page-break-inside: avoid;
      break-inside: avoid-page;
      display: block;
      width: 100%;
    }

    thead { display: table-header-group; }
    tfoot { display: table-footer-group; }
    table, tr, td, th { page-break-inside: avoid; }
    caption { caption-side: top; text-align: left; }

    h1 { color: #333; margin: 0; white-space: nowrap; }

    .title-image img {
      width: 100px; height: 100px; border-radius: 50%;
    }
    /* Footer via DomPDF (dibujado solo en la &uacute;ltima p&aacute;gina) */
  </style>
</head>

<body>
  <div class="header">
    <table style="width:100%; border-spacing:0; border-collapse:collapse; page-break-inside: avoid; break-inside: avoid-page;">
      <tr>
        <!-- Logo izquierdo -->
        <td style="padding:0; vertical-align:middle; text-align:left; width:15%;">
          <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQR1OpH7wrnsCBwE55kpuJJKgQcTD4mSPHxLg&s"
               alt="Escudo UNAS"
               style="width:90px; height:90px; display:block;">
        </td>

        <!-- T&iacute;tulos centrados -->
        <td style="padding:0; vertical-align:middle; text-align:center; width:70%;">
          <div>
            <span style="font-weight:bold; font-size:large; white-space:nowrap;">
              UNIVERSIDAD NACIONAL AGRARIA DE LA SELVA
            </span><br>
            <span style="font-weight:bold; font-size:medium; white-space:nowrap;">
              Direcci&oacute;n de Bienestar Universitario
            </span>
          </div>
        </td>

        <!-- Foto de perfil (o imagen por defecto) a la derecha -->
        <td style="padding:0; vertical-align:middle; text-align:center; width:15%;">
          <div style="width:90px; height:90px; border-radius:50%; overflow:hidden; display:inline-block;">
            @if (empty($link_profile))
              <img
                src="https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
                alt="Imagen de perfil por defecto"
                style="width:100%; height:100%; object-fit:cover; display:block;"
              >
            @else
              <img
                src="{{ $link_profile }}"
                alt="Imagen de perfil"
                style="width:100%; height:100%; object-fit:cover; display:block;"
              >
            @endif
          </div>
        </td>
      </tr>

      <!-- T&iacute;tulo centrado -->
      <tr>
        <td colspan="3" style="text-align:center; padding:0; padding-top:10px;">
          <h1>FICHA DEL ESTUDIANTE</h1>
        </td>
      </tr>
    </table>
  </div>

  <!-- Primera Secci&oacute;n
       OJO: SIN .table-block ni .keep-together para que pueda empezar en la primera p&aacute;gina
       aunque luego se corte en la siguiente -->
  <div class="section" style="margin-top:20px;">
    <table style="width:100%; border-spacing:0; border-collapse:collapse;">
      <thead>
        <tr>
          <th colspan="2" style="background-color:#f5f5f5; padding:10px; font-weight:bold; text-align:left; border:1px solid #ccc;">
            {{ $first_section['description'] }}
          </th>
        </tr>
      </thead>

      @php
        $firstReqs = array_values(array_filter($first_section['requirements'] ?? [], function ($r) {
          $val = isset($r['value']) ? trim((string)$r['value']) : '';
          return $val !== '' && strcasecmp($val, 'No proporcionado') !== 0;
        }));
      @endphp
      @foreach ($firstReqs as $index => $requirement)
        @if ($index % 2 == 0)
          <tr>
        @endif

        <td style="padding:10px; border:1px solid #ccc; width:50%; vertical-align:top;">
          <span style="font-weight:bold;">{{ $requirement['name'] }}:</span>
          <div style="border-bottom:1px solid #ccc; padding-bottom:5px;">
            @if ($requirement['type'] == 1)
              <span style="font-size:10px; text-decoration:underline;">
                {{ $requirement['value'] }}
              </span>
            @else
              {{ $requirement['value'] }}
            @endif
          </div>
        </td>

        @php $isLast = ($index == count($firstReqs) - 1); @endphp
        @if ($index % 2 == 0 && $isLast)
          <td style="padding:10px; border:1px solid #ccc; width:50%; vertical-align:top;">&nbsp;</td>
        @endif

        @if ($index % 2 == 1 || $isLast)
          </tr>
        @endif
      @endforeach
    </table>
  </div>

  <!-- Convocatorias -->
  @foreach ($announcements as $announcement)
    <div class="convocatoria keep-together" style="margin-top:20px; page-break-inside: avoid;">
      <h2 style="background-color:#f5f5f5; padding:10px; border:1px solid gray; margin:0; page-break-after: avoid; page-break-inside: avoid;">
        {{ $announcement['name'] }}
      </h2>
      <div class="section table-block keep-together" style="margin-top:10px;">
        <table style="width: 100%; border-spacing: 0; border-collapse: collapse; font-family: Arial, sans-serif; page-break-inside: avoid; break-inside: avoid-page;">
          <tr>
            <td style="padding: 10px; border: 1px solid #ccc; width: 50%; vertical-align: top;">
              <span style="font-weight: bold;">Ponderado Semestral (PSS):</span>
              <div style="border-bottom: 1px solid #ccc; padding-bottom: 5px;">
                  {{ $announcement['pps'] ?? '' }}
              </div>
            </td>
            <td style="padding: 10px; border: 1px solid #ccc; width: 50%; vertical-align: top;">
              <span style="font-weight: bold;">Estado de Solicitud:</span>
              <div style="border-bottom: 1px solid #ccc; padding-bottom: 5px;">
                  {{ $announcement['solicitud_status'] ?? '' }}
              </div>
            </td>
          </tr>
        </table>
      </div>
      @foreach ($announcement['details_requests'] as $detail)
        <div class="section table-block keep-together" style="margin-top:20px;">
          <table style="width:100%; border-spacing:0; border-collapse:collapse; page-break-inside: avoid; break-inside: avoid-page;">
            <thead>
              <tr>
                <th colspan="2" style="background-color:#f5f5f5; padding:10px; font-weight:bold; text-align:left; border:1px solid #ccc;">
                  {{ $detail['description'] }}
                </th>
              </tr>
            </thead>

            @if(isset($detail['is_table']) && $detail['is_table'] && isset($detail['rows']))
              @foreach ($detail['rows'] as $row)
                @php
                  $rowFiltered = array_values(array_filter($row ?? [], function ($r) {
                    $val = isset($r['value']) ? trim((string)$r['value']) : '';
                    return $val !== '' && strcasecmp($val, 'No proporcionado') !== 0;
                  }));
                @endphp
                @if (count($rowFiltered) === 0)
                  @continue
                @endif
                @foreach ($rowFiltered as $index => $requirement)
                  @if ($index % 2 == 0)
                    <tr style="width:100%;">
                  @endif

                  <td style="padding:10px; border:1px solid #ccc; width:50%; vertical-align:top; word-break:break-word;">
                    <span style="font-weight:bold;">{{ $requirement['name'] }}:</span>
                    <div style="border-bottom:1px solid #ccc; padding-bottom:5px;">
                      @if ($requirement['type'] == 1)
                        <span style="font-size:10px; text-decoration:underline;">
                          {{ $requirement['value'] }}
                        </span>
                      @else
                        {{ $requirement['value'] }}
                      @endif
                    </div>
                  </td>

                  @php $isLastReq = ($index == count($rowFiltered) - 1); @endphp
                  @if ($index % 2 == 0 && $isLastReq)
                    <td style="padding:10px; border:1px solid #ccc; width:50%; vertical-align:top; word-break:break-word;">&nbsp;</td>
                  @endif

                  @if ($index % 2 == 1 || $isLastReq)
                    </tr>
                  @endif
                @endforeach
              @endforeach
            @else
              @php
                $reqsFlat = array_values(array_filter($detail['requirements'] ?? [], function ($r) {
                  $val = isset($r['value']) ? trim((string)$r['value']) : '';
                  return $val !== '' && strcasecmp($val, 'No proporcionado') !== 0;
                }));
              @endphp
              @foreach ($reqsFlat as $index => $requirement)
                @if ($index % 2 == 0)
                  <tr style="width:100%;">
                @endif

                <td style="padding:10px; border:1px solid #ccc; width:50%; vertical-align:top; word-break:break-word;">
                  <span style="font-weight:bold;">{{ $requirement['name'] }}:</span>
                  <div style="border-bottom:1px solid #ccc; padding-bottom:5px;">
                    @if ($requirement['type'] == 1)
                      <span style="font-size:10px; text-decoration:underline;">
                        {{ $requirement['value'] }}
                      </span>
                    @else
                      {{ $requirement['value'] }}
                    @endif
                  </div>
                </td>

                @php $isLastReq = ($index == count($reqsFlat) - 1); @endphp
                @if ($index % 2 == 0 && $isLastReq)
                  <td style="padding:10px; border:1px solid #ccc; width:50%; vertical-align:top; word-break:break-word;">&nbsp;</td>
                @endif

                @if ($index % 2 == 1 || $isLastReq)
                  </tr>
                @endif
              @endforeach
            @endif
          </table>
        </div>
      @endforeach
    </div>
  @endforeach

  <!-- Firma al final de la &uacute;ltima p&aacute;gina (3 cm desde el borde inferior) -->
  <div style="text-align:right; margin-top:200px;">
    <p style="border-top:1px solid black; display:inline-block; padding-top:5px;">
      Firma de Direcci&oacute;n de Bienestar Universitario
    </p>
  </div>
</body>
</html>
