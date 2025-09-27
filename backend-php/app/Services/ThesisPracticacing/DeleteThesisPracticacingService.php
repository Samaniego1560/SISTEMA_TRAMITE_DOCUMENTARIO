<?php

namespace App\Services\ThesisPracticacing;

use App\Models\ThesisPracticacing;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\Hash;
use App\Exceptions\ExceptionGenerate;

class DeleteThesisPracticacingService
{
    public function delete($id): Model
    {
        $thesisPracticacing = ThesisPracticacing::findDA($id);
        if (!$thesisPracticacing) {
            throw new ExceptionGenerate('Tesista o practicante no encontrado', 404);
        }
        $thesisPracticacing->delete();
        return $thesisPracticacing;
    }
}
