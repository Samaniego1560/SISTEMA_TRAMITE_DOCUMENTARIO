<?php

namespace App\Services\ThesisPracticacing;

use App\Models\ThesisPracticacing;
use Illuminate\Support\Facades\Hash;
use App\Exceptions\ExceptionGenerate;
use Illuminate\Database\Eloquent\Collection; 

class ListThesisPracticacingService
{
    public function list(): Collection
    {
        return ThesisPracticacing::where('status_id', 1)->get();
    }
}
