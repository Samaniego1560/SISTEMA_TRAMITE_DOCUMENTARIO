<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class ThesisPracticacing extends Model
{
    protected $connection = 'mysql_dbu';
    protected $table = 'thesis_practicacings';
    protected $primaryKey = 'id';
    public $incrementing = true; // Si 'id' no es autoincremental
    public $timestamps = true; // Asegúrate de que los timestamps están habilitados

    protected $fillable = [
        'type_student',
        'cod_student',
        'name',
        'appaterno',
        'apmaterno',
        'sexo',
        'age',
        'facultad',
        'escuela_profesional',
        'mod_ingreso',
        'dni',
        'email',
        'convocatoria_id',
        'status_id',
    ];

    public static function allDA()
    {
        return self::whereIn('status_id', [1])->get();
    }

    public static function findDA($id)
    {
        return self::whereIn('status_id', [1])->find($id);
    }
}
