<?php
namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class CalendarSetting extends Model
{
    protected $table = 'calendar_setting';

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
        'start_date',
        'end_date',
        'title',
        'description',
    ];
}
