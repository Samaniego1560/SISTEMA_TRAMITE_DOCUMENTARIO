<?php

namespace App\Http\Requests\CalendarSetting;

use Anik\Form\FormRequest;

class CalendarSettingRequest extends FormRequest
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
            '*.start_date' => [
                'required',
                'date',
            ],
            '*.end_date' => [
                'nullable',
                'date',
            ],
            '*.title' => [
                'required',
                'string',
                'max:255',
            ],
            '*.description' => [
                'nullable',
                'string',
                'max:1000',
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
            '*.start_date.required' => 'La fecha de inicio es requerida para cada ajuste.',
            '*.start_date.date' => 'La fecha de inicio debe ser válida.',
            '*.end_date.required' => 'La fecha de finalización es requerida para cada ajuste.',
            '*.end_date.date' => 'La fecha de finalización debe ser válida.',
           '*.title.required' => 'El título es requerido para cada ajuste.',
            '*.title.string' => 'El título debe ser una cadena de texto.',
            '*.title.max' => 'El título no debe superar los 255 caracteres.',
            '*.description.string' => 'La descripción debe ser una cadena de texto.',
            '*.description.max' => 'La descripción no debe superar los 1000 caracteres.',
        ];
    }
}
