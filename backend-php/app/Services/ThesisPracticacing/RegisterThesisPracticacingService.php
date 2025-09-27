<?php

namespace App\Services\ThesisPracticacing;

use App\Exceptions\ExceptionGenerate;
use App\Models\ThesisPracticacing;
use Illuminate\Support\Facades\Hash;

class RegisterThesisPracticacingService
{
    /**
     * Register a new thesis or practicum entry.
     *
     * @param array $data
     * @return ThesisPracticacing
     * @throws ExceptionGenerate
     */
    public function register(array $data)
    {
        // Check if a record already exists with the same DNI, email, and convocatoria_id but with a different status_id
        $thesis_practicacing = ThesisPracticacing::where([
            ['dni', '=', $data['dni']],
            ['convocatoria_id', '=', $data['convocatoria_id']],
            ['email', '=', $data['email']],
            ['status_id', '!=', 1],
        ])->first();

        // Throw an exception if a record exists
        if ($thesis_practicacing) {
            throw new ExceptionGenerate('Ya existe un tesista o practicante con el DNI o correo electrónico para la convocatoria.', 400);
        }

        // Create and return a new record
        $thesis_practicacing = ThesisPracticacing::create($data);

        return $thesis_practicacing;
    }
}
