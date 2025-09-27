<?php

namespace App\Http\Requests\ThesisPracticacing;

use Anik\Form\FormRequest;

class ThesisPracticacingRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    protected function authorize(): bool
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array
     */
    protected function rules(): array
    {
        return [
            'type_student' => [
                'required',
                'string',
                'max:255',
            ],
            'cod_student' => [
                'required',
                'string',
                'max:255',
            ],
            'name' => [
                'required',
                'string',
                'max:255',
            ],

            'appaterno' => [
                'required',
                'string',
                'max:255',
            ],
            'apmaterno' => [
                'required',
                'string',
                'max:255',
            ],
            'sexo' => [
                'required',
                'string',
                'max:255',
            ],
            'age' => [
                'required',
                'string',
                'max:255',
            ],
            'facultad' => [
                'required',
                'string',
                'max:255',
            ],
            'escuela_profesional' => [
                'required',
                'string',
                'max:255',
            ],
            'mod_ingreso' => [
                'required',
                'string',
                'max:255',
            ],
            'dni' => [
                'required',
                'numeric',
                'digits:8', // Asumiendo que el DNI siempre tiene 8 dígitos
            ],
            'email' => [
                'required',
                'email', // Validación básica para correos electrónicos
                'max:255',
            ],
            'convocatoria_id' => [
                'required',
                'integer',
                'exists:convocatorias,id', // Asegúrate de que el ID exista en la tabla `convocatorias`
            ],
            'status_id' => [
                'required',
                'integer',
            ],
        ];
    }

    /**
     * Get custom validation messages.
     *
     * @return array
     */
    public function messages(): array
    {
        return [
            'type_student.required' => 'El tipo estudiante es requerido',
            'cod_student.required' => 'El Codigo estudiamte es requerido',
            'name.required' => 'El nombre es requerido',
            'apmaterno.required' => 'El Apellido es requerido',
            'apmaterno.required' => 'El Apellido es requerido',
            'sexo.required' => 'El Sexo es requerido',
            'age.required' => 'El Sexo es requerido',
            'facultad.required' => 'El facultad es requerido',
            'escuela_profesional.required' => 'La escuela es requerido',
            'mod_ingreso.required' => 'La modalidad es requerido',
            'dni.required' => 'El DNI es requerido',
            'dni.numeric' => 'El DNI debe ser un número',
            'dni.digits' => 'El DNI debe tener 8 dígitos',
            'email.required' => 'El email es requerido',
            'email.email' => 'El email debe ser una dirección de correo válida',
            'convocatoria_id.required' => 'El ID de convocatoria es requerido',
            'convocatoria_id.integer' => 'El ID de convocatoria debe ser un número entero',
            'convocatoria_id.exists' => 'El ID de convocatoria no existe',
            'status_id.required' => 'El ID de estado es requerido',
            'status_id.integer' => 'El ID de estado debe ser un número entero',
            'status_id.exists' => 'El ID de estado no existe',
        ];
    }
}
