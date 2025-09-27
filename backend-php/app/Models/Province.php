<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Province extends Model
{
    protected $connection = 'mysql_dbu';
    protected $table = 'provinces';
    protected $primaryKey = 'id';
    protected $keyType = 'string'; // Si 'id' es una cadena
    protected $fillable = [
        'name',
        'departament_id',
    ];
}
