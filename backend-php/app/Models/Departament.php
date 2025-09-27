<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Departament extends Model
{
    protected $connection = 'mysql_dbu';
    protected $table = 'departaments';
    protected $primaryKey = 'id';
    public $incrementing = false; // Si 'id' no es autoincremental
    protected $keyType = 'string'; // Si 'id' es una cadena
    public $timestamps = true; // Asegúrate de que los timestamps están habilitados

    protected $fillable = [
        'name',
    ];
}

