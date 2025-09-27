<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class ApplicationNotice extends Model
{
    protected $table = 'application_notice';

    /**
     * Habilitar timestamps (por defecto es true, así que este paso es opcional).
     *
     * @var bool
     */
    public $timestamps = true;

    /**
     * Los atributos que se pueden asignar masivamente.
     *
     * @var array
     */
    protected $fillable = [
        'id',
        'description',
        'status'
    ];
}
