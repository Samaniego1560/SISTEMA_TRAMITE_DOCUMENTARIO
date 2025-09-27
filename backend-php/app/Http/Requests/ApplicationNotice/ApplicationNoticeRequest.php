<?php

namespace App\Http\Requests\ApplicationNotice;

use Anik\Form\FormRequest;

class ApplicationNoticeRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize(): bool
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array
     */
    public function rules(): array
    {
        return [
            '*.description' => [
                'required',
                'string',
                'max:1000', // Máximo de 1000 caracteres para cadenas
            ],
            '*.status' => [
                'required',
                'integer', // Validación correcta para enteros
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
            '*.description.required' => 'La descripción es requerida.',
            '*.status.required' => 'El estado es requerido.', // Clave corregida
        ];
    }
}
