<?php

namespace App\Http\Requests\Solicitud;

use Anik\Form\FormRequest;

class ValidarDatosAlumnoRequest extends FormRequest
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
            'correo' => [
                'required',
                'email' => 'required',
                /* 'unas_email' */
            ],
            'DNI' => [
                'required',
                'numeric',
                'dni_length',
            ],
            'announcement_id' => [
                'required',
                'numeric'
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
            'correo.required' => 'El correo es requerido',
            'correo.email' => 'Ingrese el correo institucional brindado por la oficina OTI - UNAS',
            'DNI.required' => 'Se requiere el DNI',
            'DNI.min' => 'Ingrese correctamente su DNI',
        ];
    }
}
