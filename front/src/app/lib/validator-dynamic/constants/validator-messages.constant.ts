import {ValidationErrors} from "@angular/forms";

export const VALIDATOR_MESSAGE_DEFAULT: ValidationErrors = {
  required: 'Este campo es requerido',
  min: 'Este campo no debe ser inferior al valor: ${min}, valor actual: ${current}',
  max: 'Este campo no debe ser superior al valor: ${max}, valor actual: ${current}',
  minLength: 'Este campo debe tener un tamaño mínimo de: ${minLength} tamaño actual:${current}',
  maxLength: 'Este campo debe tener un tamaño máximo de: ${maxLength} tamaño actual:${current}',
  email: 'Ingrese un email válido',
  dateRange: 'Este campo supera el rango de fecha, rango:${dateRange}',
};
