<?php

namespace App\Services\ThesisPracticacing;

use App\Models\ThesisPracticacing;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\Hash;
use App\Exceptions\ExceptionGenerate;

class UpdateThesisPracticacingService
{
    public function update($id, array $data): Model
    {
        // Lógica de actualización
        $thesis_practicacing = ThesisPracticacing::where([
            ['dni', '=', $data['dni']],
            ['convocatoria_id', '=', $data['convocatoria_id']],
            ['email', '=', $data['email']],
        ])->first();
        if ($thesis_practicacing && ($thesis_practicacing['id'] != $id)) {
            throw new ExceptionGenerate('Ya existe un tesista o practicante con el dni email en la convocatoria', 404);
        }

        $thesis_practicacing = ThesisPracticacing::whereId($id)->first();
        if (!$thesis_practicacing)
            throw new ExceptionGenerate('Usuario no enontrado', 404);

        $thesis_practicacing->update($data);

        return $thesis_practicacing;
    }
}
