<?php

namespace App\Http\Resources\ThesisPracticacing;

use Illuminate\Http\Resources\Json\JsonResource;

class ThesisPracticacingResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @param  \Illuminate\Http\Request  $request
     * @return array
     */
    public function toArray($request)
    {
        return [
            'id' => $this->id,
            'type_student' => $this->type_student,
            'cod_student' => $this->cod_student,
            'name' => $this->name,
            'appaterno' => $this->appaterno,
            'apmaterno' => $this->apmaterno,
            'sexo' => $this->sexo,
            'age' => $this->age,
            'facultad' => $this->facultad,
            'escuela_profesional' => $this->escuela_profesional,
            'mod_ingreso' => $this->mod_ingreso,
            'dni' => $this->dni,
            'email' => $this->email,
            'convocatoria_id' => $this->convocatoria_id,
            'status_id' => $this->status_id,
            'created_at' => $this->created_at,
            'updated_at' => $this->updated_at,
        ];
    }
}
